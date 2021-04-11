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

    #[error("Unable to create opus encoder/decoder {0}")]
    OpusError(#[from] magnum_opus::Error),

    #[error("Unable to decode data {0:?}")]
    Decode(magnum_opus::ErrorCode),
}

pub(crate) type Result<T> = std::result::Result<T, Error>;
