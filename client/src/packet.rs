// Source Chat Relay
// Copyright (C) 2021  rumblefrog
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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
    /// Type 0x6 containing Opus PLC payload.
    OpusPLC(&'a [u8]),

    // Type 0x0 containing number of samples.
    Silence(u16),
}

impl Default for Packet<Vec<u8>> {
    fn default() -> Self {
        Packet {
            cursor: Cursor::new(Vec::new()),
        }
    }
}

impl Packet<Vec<u8>> {
    pub fn new() -> Packet<Vec<u8>> {
        Self::default()
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
            0x6 => {
                let len = self.cursor.read_u16::<LittleEndian>()?;

                let pos = self.cursor.position() as usize;

                let data = &self.cursor.get_ref()[pos..pos + len as usize];

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
