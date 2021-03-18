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
