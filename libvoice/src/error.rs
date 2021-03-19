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
