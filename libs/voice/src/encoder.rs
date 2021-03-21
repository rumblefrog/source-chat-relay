use magnum_opus::{Application, Channels, Encoder as OpusEncoder};

use crate::messages::voice::{OpusDataFrame, OpusSilenceFrame};

use crate::{Result, FRAME_SIZE};

pub struct Encoder {
    encoder: OpusEncoder,

    current_frame: u32,

    samples: u64,
}

impl Encoder {
    pub fn new() -> Result<Self> {
        Ok(Self {
            encoder: OpusEncoder::new(48000, Channels::Stereo, Application::Voip)?,
            current_frame: 0,
            samples: 0
        })
    }

    pub async fn encode(&mut self, data: &[i16]) -> Result<Vec<OpusDataFrame>> {
        let mut out = Vec::new();

        for chunk in data.chunks(FRAME_SIZE) {
            let frame_bytes = self.encode_chunk(chunk).await?;

            // Two bytes per sample
            self.samples += (frame_bytes.len() / 2) as u64;

            out.push(OpusDataFrame {
                index: self.current_frame,
                data: frame_bytes,
            });

            self.current_frame += 1;
        }

        Ok(out)
    }

    pub fn seq_end(&mut self) -> OpusSilenceFrame {
        let frame = OpusSilenceFrame {
            samples: self.samples,
        };

        self.samples = 0;
        self.current_frame = 0;

        frame
    }

    #[inline]
    async fn encode_chunk(&mut self, data: &[i16]) -> Result<Vec<u8>> {
        // Recommended max_data_bytes (https://opus-codec.org/docs/opus_api-1.3.1/group__opus__encoder.html#ga4ae9905859cd241ef4bb5c59cd5e5309)
        let mut out = vec![0u8; 4000];

        let bytes = self.encoder.encode(data, &mut out)?;

        out.truncate(bytes);

        Ok(out)
    }
}
