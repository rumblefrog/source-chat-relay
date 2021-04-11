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

use thiserror::Error;

#[derive(Debug, Error)]
pub enum Error {
    #[error("IO error {0}")]
    IoError(#[from] std::io::Error),

    #[error("Unable to create opus decoder {0}")]
    Decoder(#[from] magnum_opus::Error),

    #[error("Unable to decode data {0:?}")]
    Decode(magnum_opus::ErrorCode),

    #[error("Unable to resample audio {0}")]
    ResampleError(#[from] samplerate::Error),

    #[error("Invalid packet CRC32 (expected: {0:x}) (actual: {1:x}")]
    InvalidPacketChecksum(u32, u32),

    #[error("Invalid payload type (expected: {0:x}) (actual: {1:x}")]
    InvalidPayloadType(u8, u8),

    #[error("Insufficient data length {0}")]
    InsufficientData(usize),
}

pub type Result<T> = std::result::Result<T, Error>;
