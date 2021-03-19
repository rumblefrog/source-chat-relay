use magnum_opus::{Channels, Decoder as OpusDecoder};

use crate::messages::voice::opus_frame::Frame;
use crate::{Result, FRAME_SIZE};

pub struct Decoder {
    decoder: OpusDecoder,

    current_frame: u32,

    last_samples_count: u64,
}

pub struct DecodedData {
    pub samples: Vec<i16>,

    // Whether this payload of frames contains the silence frame.
    // End user may use this as indication to get last_samples_count.
    pub is_last: bool,
}

impl Decoder {
    pub fn new() -> Result<Self> {
        Ok(Self {
            decoder: OpusDecoder::new(48000, Channels::Stereo)?,
            current_frame: 0,
            last_samples_count: 0,
        })
    }

    pub async fn decode(&mut self, frames: &[Frame]) -> Result<DecodedData> {
        let mut decoded_data = DecodedData {
            // FRAME_SIZE * frames.len() for initial size, if loss is detected, it will grow pcm.
            samples: Vec::with_capacity(FRAME_SIZE * frames.len()),
            is_last: false
        };

        for frame in frames {
            match frame {
                Frame::Data(data) => {
                    let current_frame = data.index;
                    let previous_frame = self.current_frame;

                    if current_frame >= previous_frame {
                        let mut decoded = if data.index == previous_frame {
                            self.current_frame = data.index + 1;

                            self.decode_frame(&data.data).await?
                        } else {
                            self.decode_loss((current_frame - previous_frame) as usize).await?
                        };

                        decoded_data.samples.append(&mut decoded);
                    }

                },
                // Silence frame is always the last in a series
                Frame::Silence(silence) => {
                    self.current_frame = 0;
                    self.last_samples_count = silence.samples;
                    decoded_data.is_last = true;
                }
            }
        }

        Ok(decoded_data)
    }

    /// Last samples count for a series of frames.
    /// This should be used when is_last is true, otherwise the count will be prior series.
    pub fn last_samples_count(&self) -> u64 {
        self.last_samples_count
    }

    #[inline]
    async fn decode_frame(&mut self, data: &[u8]) -> Result<Vec<i16>> {
        let mut out = vec![0; FRAME_SIZE];

        self.decoder.decode(&data, &mut out, false)?;

        Ok(out)
    }

    #[inline]
    async fn decode_loss(&mut self, samples_loss: usize) -> Result<Vec<i16>> {
        let samples_loss = std::cmp::min(samples_loss, 10);

        let mut out = Vec::with_capacity(FRAME_SIZE * samples_loss);

        for _ in 0..samples_loss {
            let mut o = vec![0; FRAME_SIZE];

            self.decoder.decode(&[], &mut o, false)?;

            out.append(&mut o);
        }

        Ok(out)
    }
}
