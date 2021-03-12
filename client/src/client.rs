use std::{fs::File, io::Write};

use magnum_opus::{Decoder, Channels};

use crate::packet::{Packet, Payload};

// TODO: Transcode payload to match Discord opus requirement
pub struct Client {
    decoder: Decoder,

    file: File,
}

impl Default for Client {
    fn default() -> Self {
        let decoder = Decoder::new(24000, Channels::Stereo).unwrap();

        let file = File::create("audio_out").unwrap();

        Self { decoder, file }
    }
}

impl Client {
    pub fn receive_audio(&mut self, _id: u64, data: &[u8], _force_steam_voice: bool) {
        if let Ok(mut p) = Packet::from_bytes(data) {
            let header = p.header().unwrap();
            let payload = p.payload().unwrap();

            println!("@@@@@ > {} - {}", header.steam_id, header.sample_rate);

            match payload {
                Payload::OpusPLC(d) => {
                    // let mut out = Vec::new();
                    // let frame_size = self.decoder.decode(d, &mut out, false).unwrap();

                    // let mut pcm_bytes: Vec<u8> = Vec::with_capacity(2 * frame_size);

                    // for i in 0..2 * frame_size {
                    //     pcm_bytes[2 * i] = (out[i] & 0xFF) as u8;
                    //     pcm_bytes[2 * i + 1] = ((out[i] >> 8) & 0xFF) as u8;
                    // }

                    // self.file.write(&pcm_bytes).unwrap();

                    println!("!!!!! Opus {}", d.len());
                },
                Payload::Silence(ns) => println!("!!!!! Silence {}", ns),
                _ => { /* not interested */ }
            }
        }
    }
}

#[no_mangle]
pub extern "C" fn new_client() -> *mut Client {
    let c = Client::default();
    let b = Box::new(c);

    Box::into_raw(b)
}

#[no_mangle]
pub extern "C" fn receive_audio(
    client: *mut Client,
    id: u64,
    bytes: i32,
    data: *const std::os::raw::c_char,
    force_steam_voice: bool,
) {
    if !client.is_null() {
        let d = unsafe { std::slice::from_raw_parts(data as *const u8, bytes as usize) };

        unsafe { &mut *client }.receive_audio(id, d, force_steam_voice)
    }
}

#[no_mangle]
pub extern "C" fn free_client(client: *mut Client) {
    if client.is_null() {
        return;
    }

    unsafe { Box::from_raw(client) };
}
