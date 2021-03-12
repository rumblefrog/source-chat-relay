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
    /// Type 0x5 containing Opus payload
    Opus(&'a [u8]),

    /// Type 0x6 containing Opus PLC payload.
    OpusPLC(&'a [u8]),

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
                self.cursor.write_u8(0x5)?;
                self.cursor.write_all(data)?;
            }
            Payload::OpusPLC(data) => {
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
            0x5 => {
                let len = self.cursor.read_u16::<LittleEndian>()?;

                let data = &self.cursor.get_ref()[13..13 + len as usize];

                return Ok(Payload::Opus(data));
            }
            0x6 => {
                let len = self.cursor.read_u16::<LittleEndian>()?;

                let data = &self.cursor.get_ref()[13..13 + len as usize];

                return Ok(Payload::OpusPLC(data));
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
    let data: [u8; 105] = [
        0x0F, 0x0F, 0x33, 0x09, 0x01, 0x00, 0x10, 0x01, 0x0B, 0xC0, 0x5D, 0x06, 0x57, 0x00, 0x53, 0x00,
        0x08, 0x00, 0x68, 0x9B, 0x26, 0x69, 0x8D, 0xEA, 0x93, 0xE1, 0xB9, 0x5E, 0xF1, 0xC6, 0x36, 0x22,
        0x9A, 0x55, 0x5F, 0xE3, 0xB0, 0x04, 0x54, 0x03, 0x11, 0xBC, 0xB8, 0xF7, 0xE2, 0xE8, 0x82, 0xF9,
        0x08, 0xEF, 0x0E, 0x98, 0x04, 0x7B, 0xC7, 0x5E, 0xE0, 0x03, 0xDC, 0x39, 0xB9, 0x4D, 0xEC, 0x66,
        0xFD, 0x1B, 0xF5, 0x49, 0x84, 0xAE, 0xA4, 0x86, 0xAC, 0x5A, 0xE0, 0x83, 0x19, 0x83, 0xCA, 0x8F,
        0x31, 0x2B, 0xFF, 0xBB, 0x6A, 0xB5, 0x1E, 0xF3, 0xCE, 0x37, 0x57, 0x37, 0x4D, 0x95, 0x38, 0x7E,
        0xDB, 0x86, 0xD7, 0x0B, 0xF3, 0xBE, 0xC0, 0x2B, 0xEA
    ];

    let mut p = Packet::from_bytes(&data).unwrap();

    let header = p.header().unwrap();

    if let Payload::Opus(d) = p.payload().unwrap() {
        println!("len {}", d.len());
    }
}
