mod decoder;
mod encoder;
mod error;

pub use error::Error;
use error::Result;

pub use decoder::Decoder;
pub use encoder::Encoder;

pub const FRAME_SIZE: usize = 960;

pub mod messages {
    pub mod common {
        include!(concat!(env!("OUT_DIR"), "/common.rs"));
    }

    pub mod voice {
        include!(concat!(env!("OUT_DIR"), "/voice.rs"));
    }
}
