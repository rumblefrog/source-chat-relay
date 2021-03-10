use std::io::{Cursor, Seek, SeekFrom, Write};

use crc::crc32;

use byteorder::{LittleEndian, ReadBytesExt, WriteBytesExt};

use crate::{Error, Result};

/// Packet reader and writer.
///
/// A single packet can contain multiple payloads.
///
/// From observation, every packet contains at least steamid & sample rate payload.
/// Following that observation, steam id & sample rate payload will be considered header in packets.
///
/// With the remaining 1 payload being voice data or silence that's vital.
///
/// Inner working of steamclient has not be RE-ed yet.
/// For the time being, based on observed voice packets,
/// preserve read and write order of payloads.
pub struct Packet<T> {
    cursor: Cursor<T>,
}

/// Packet "headers".
pub struct Header {
    /// Steam ID 64
    pub steam_id: u64,

    /// Voice sample rate
    pub sample_rate: u16,
}

/// Supported payloads.
pub enum Payload<'a> {
    /// Type 0xB containing Opus bytes in which length is computed from.
    Opus(&'a [u8]),

    // Type 0x0 containing number of samples.
    Silence(u16),
}

impl Packet<Vec<u8>> {
    pub fn new() -> Packet<Vec<u8>> {
        Packet {
            cursor: Cursor::new(Vec::new()),
        }
    }

    pub fn header(&mut self, header: Header) -> Result<()> {
        self.cursor.seek(SeekFrom::Start(0))?;

        self.cursor.write_u64::<LittleEndian>(header.steam_id)?;

        self.cursor.write_u8(0xB)?;
        self.cursor.write_u16::<LittleEndian>(header.sample_rate)?;

        Ok(())
    }

    /// Consumes packet and return inner written Vec<u8>.
    pub fn payload(mut self, payload: Payload) -> Result<Vec<u8>> {
        match payload {
            Payload::Opus(data) => {
                self.cursor.write_u8(0xB)?;
                self.cursor.write_all(data)?;
            }
            Payload::Silence(ns) => {
                self.cursor.write_u8(0x0)?;
                self.cursor.write_u16::<LittleEndian>(ns)?;
            }
        }

        let crc32 = crc32::checksum_ieee(self.cursor.get_ref());

        self.cursor.write_u32::<LittleEndian>(crc32)?;

        Ok(self.cursor.into_inner())
    }
}

impl<'p> Packet<&'p [u8]> {
    /// Forms packet from data reference and validates CRC32.
    pub fn from_bytes(data: &'p [u8]) -> Result<Packet<&'p [u8]>> {
        let data_len = {
            let len = data.len();

            // 11 header + 4 CRC
            if len < 15 {
                return Err(Error::InsufficientData(len));
            }

            len - 4
        };

        let mut cursor = Cursor::new(data);

        cursor.seek(SeekFrom::Start(data_len as u64))?;

        let expected_crc32 = cursor.read_u32::<LittleEndian>()?;
        let actual_crc32 = crc32::checksum_ieee(&cursor.get_ref()[0..data_len]);

        if expected_crc32 != actual_crc32 {
            return Err(Error::InvalidPacketChecksum(expected_crc32, actual_crc32));
        }

        Ok(Packet { cursor })
    }

    pub fn header(&mut self) -> Result<Header> {
        self.cursor.seek(SeekFrom::Start(0))?;

        Ok(Header {
            steam_id: self.cursor.read_u64::<LittleEndian>()?,
            sample_rate: {
                let pt = self.cursor.read_u8()?;

                if pt != 0xB {
                    return Err(Error::InvalidPayloadType(0xB, pt));
                }

                self.cursor.read_u16::<LittleEndian>()?
            },
        })
    }

    pub fn payload(&mut self) -> Result<Payload> {
        self.cursor.seek(SeekFrom::Start(8 + 1 + 2))?;

        let pt = self.cursor.read_u8()?;

        match pt {
            0x6 => {
                let len = self.cursor.read_u16::<LittleEndian>()?;

                let data = &self.cursor.get_ref()[13..13 + len as usize];

                return Ok(Payload::Opus(data));
            }
            0x0 => {
                return Ok(Payload::Silence(self.cursor.read_u16::<LittleEndian>()?));
            }
            _ => {}
        }

        Err(Error::InvalidPayloadType(0x6, pt))
    }
}

#[test]
fn v1() {
    let data: [u8; 281] = [
        0x0F, 0x0F, 0x33, 0x09, 0x01, 0x00, 0x10, 0x01, 0x0B, 0xC0, 0x5D, 0x06, 0x07, 0x01, 0x55,
        0x00, 0x08, 0x00, 0x68, 0x9B, 0x9A, 0x2A, 0x4E, 0xE1, 0x49, 0xC7, 0x2B, 0x54, 0xF2, 0xBC,
        0xCB, 0x92, 0x4A, 0xCA, 0xFC, 0xF3, 0x80, 0x9E, 0xBD, 0x7C, 0xB6, 0x26, 0x41, 0x13, 0xFA,
        0x9F, 0x15, 0xBE, 0x3B, 0x97, 0xF8, 0xCF, 0x14, 0x78, 0x32, 0x27, 0xE3, 0x70, 0x53, 0xF2,
        0x58, 0xB1, 0x90, 0x64, 0xA7, 0xA4, 0xC5, 0xFD, 0x3A, 0x6B, 0xDB, 0xF7, 0xA2, 0x04, 0x7A,
        0x9A, 0x65, 0xDE, 0x06, 0x18, 0xC2, 0xD9, 0x6C, 0xB9, 0x37, 0xC4, 0xE5, 0x55, 0xA9, 0x47,
        0xA8, 0xA8, 0xAA, 0xDD, 0xBC, 0x72, 0x21, 0x06, 0xE0, 0x33, 0x4E, 0xF3, 0xBE, 0x49, 0x00,
        0x09, 0x00, 0x68, 0x9B, 0x9E, 0x77, 0x21, 0xE5, 0xFD, 0xF5, 0x7A, 0x5D, 0xE5, 0x41, 0x1F,
        0xBF, 0x0B, 0x51, 0x9C, 0x7A, 0xCD, 0xC5, 0x17, 0x6B, 0xF4, 0x77, 0x6C, 0x80, 0x6A, 0x65,
        0x72, 0x7A, 0x9B, 0xAF, 0x21, 0xF9, 0x33, 0x1F, 0xCB, 0x6D, 0x8C, 0x18, 0x81, 0xD3, 0x32,
        0x94, 0x0C, 0x13, 0x64, 0xBC, 0xA5, 0x13, 0xAB, 0xF4, 0x76, 0x00, 0x0D, 0x66, 0x6A, 0xEB,
        0x76, 0x7D, 0x78, 0xF8, 0x38, 0x3E, 0x6B, 0x08, 0x94, 0x4F, 0xBC, 0x2D, 0x4B, 0x3F, 0x67,
        0x5D, 0x00, 0x0A, 0x00, 0x68, 0x9B, 0xCB, 0xE1, 0x14, 0x01, 0xD1, 0xAC, 0x0F, 0xFC, 0x09,
        0x9E, 0xD7, 0xBA, 0x49, 0x42, 0x56, 0x28, 0x16, 0x3D, 0xE5, 0x67, 0xCD, 0xFD, 0xE5, 0x82,
        0xC8, 0x20, 0x52, 0xB5, 0x9E, 0x06, 0x4B, 0x37, 0x55, 0x0A, 0xF1, 0x3B, 0xB7, 0x59, 0x54,
        0x8D, 0xFD, 0xAC, 0xA8, 0x0D, 0x9D, 0x52, 0x21, 0x1C, 0x30, 0xFD, 0xAF, 0x55, 0xAF, 0x99,
        0xF3, 0xEA, 0x0B, 0xC1, 0x26, 0xAB, 0x04, 0xC1, 0x40, 0xDB, 0x81, 0xE0, 0x87, 0x6A, 0xEA,
        0x36, 0x32, 0xF1, 0xE5, 0x7B, 0xB2, 0xF4, 0x4A, 0xA9, 0x5D, 0x97, 0x0A, 0x02, 0xC5, 0x79,
        0x3C, 0x52, 0x2B, 0xDC, 0x41, 0x9B, 0x2B, 0xD7, 0xF3, 0x69, 0x9A,
    ];

    let mut p = Packet::from_bytes(&data).unwrap();

    let header = p.header().unwrap();

    let payload = p.payload().unwrap();
}
