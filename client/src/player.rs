use std::time::{Instant, Duration};
use std::io::{Cursor, Seek, SeekFrom};

use byteorder::{LittleEndian, ReadBytesExt};

use samplerate::{Samplerate, ConverterType};

use magnum_opus::{Channels, Decoder, Encoder, Application};

use crate::Result;

const FRAME_SIZE: usize = 480;

/// Player struct represents an unique player on the server.
///
/// Player will have its own encoder & decoder stream.
pub struct Player {
    encoder: Encoder,

    decoder: Decoder,

    last_receive: Instant,

    current_sample: u16,

    resampler: Samplerate,
}

// Requires shim safety
unsafe impl Send for Player {}
unsafe impl Sync for Player {}

impl Player {
    pub fn new() -> Result<Player> {
        Ok(Self {
            encoder: Encoder::new(48000, Channels::Stereo, Application::Voip)?,
            decoder: Decoder::new(24000, Channels::Mono)?,
            last_receive: Instant::now(),
            current_sample: 0,
            resampler: Samplerate::new(ConverterType::SincFastest, 24000, 48000, 2)?,
        })
    }

    pub fn is_stale(&self) -> bool {
        // Considered stale if last received is more than 10 minutes ago
        self.last_receive.elapsed() > Duration::from_secs(60 * 10)
    }

    pub fn transcode(&mut self, data: &[u8]) -> Result<Vec<u8>> {
        let payload_len = data.len() as u64;

        let mut data = Cursor::new(data);
        let mut out: Vec<u8> = Vec::new();

        while data.position() < payload_len {
            let chunk_len = data.read_i16::<LittleEndian>()?;

            // End of packet sequence
            if chunk_len == -1 {
                self.current_sample = 0;
                break;
            }

            let current_sample = data.read_u16::<LittleEndian>()?;
            let prev_sample = self.current_sample;

            let pos = data.position() as usize;
            let chunk_len = chunk_len as usize;

            data.seek(SeekFrom::Current(chunk_len as i64))?;

            if current_sample >= prev_sample {
                let decoded = if current_sample == prev_sample {
                    self.current_sample = current_sample + 1;

                    self.decode_chunk(&data.get_ref()[pos .. pos + &chunk_len])?
                } else {
                    self.decode_loss((current_sample - prev_sample) as usize)?
                };

                for decoded_chunk in decoded.chunks(FRAME_SIZE) {
                    let resampled = self.resample(decoded_chunk)?;
                    let mut chunk = self.encode_chunk(&resampled)?;

                    // let mut rem = vec![0u8; FRAME_SIZE * 2];

                    // for (i, f) in decoded_chunk.iter().enumerate() {
                    //     rem[2 * i] = *f as u8;
                    //     rem[2 * i + 1] = (*f >> 8) as u8;
                    // }

                    // let resampled = self.resample(&decoded_chunk)?;

                    // let mut chunk = self.encode_chunk(&resampled)?;
                    out.append(&mut chunk);
                }
            }
        }

        Ok(out)
    }

    /// Interleave two mono channels to one stereo.
    #[inline]
    fn duplicate_interleave(pcm: &[f32]) -> Vec<f32> {
        let mut out = vec![0f32; FRAME_SIZE * 2];

        for i in 0..FRAME_SIZE {
            out[i * 2] = pcm[i];
            out[i * 2 + 1] = pcm[i];
        }

        out
    }

    #[inline]
    fn resample(&mut self, pcm: &[f32]) -> Result<Vec<f32>> {
        let interleaved = Self::duplicate_interleave(pcm);

        Ok(self.resampler.process(&interleaved)?)
    }

    #[inline]
    fn encode_chunk(&mut self, data: &[f32]) -> Result<Vec<u8>> {
        // Recommended max_data_bytes (https://opus-codec.org/docs/opus_api-1.3.1/group__opus__encoder.html#ga4ae9905859cd241ef4bb5c59cd5e5309)
        let mut out = vec![0u8; 4000];

        let bytes = self.encoder.encode_float(data, &mut out)?;

        out.truncate(bytes);

        Ok(out)
    }

    #[inline]
    fn decode_chunk(&mut self, data: &[u8]) -> Result<Vec<f32>> {
        let mut out = vec![0.0; FRAME_SIZE];

        self.decoder
            .decode_float(&data, &mut out, false)?;

        Ok(out)
    }

    #[inline]
    fn decode_loss(&mut self, samples_loss: usize) -> Result<Vec<f32>> {
        let samples_loss = std::cmp::min(samples_loss, 10);

        let mut out = Vec::with_capacity(FRAME_SIZE * samples_loss);

        for _ in 0..samples_loss {
            let mut o = vec![0.0; FRAME_SIZE];

            self.decoder
                .decode_float(&[], &mut o, false)?;

            out.append(&mut o);
        }

        Ok(out)
    }
}
