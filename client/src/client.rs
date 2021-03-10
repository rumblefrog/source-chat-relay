use crate::packet::{Packet, Payload};

pub struct Client;

impl Default for Client {
    fn default() -> Self {
        Self {}
    }
}

impl Client {
    pub fn receive_audio(&self, _id: u64, data: &[u8], _force_steam_voice: bool) {
        if let Ok(mut p) = Packet::from_bytes(data) {
            let header = p.header().unwrap();
            let payload = p.payload().unwrap();

            println!("@@@@@ > {} - {}", header.steam_id, header.sample_rate);

            match payload {
                Payload::Opus(d) => println!("!!!!! Opus len: {}", d.len()),
                Payload::Silence(ns) => println!("!!!!! Silence {}", ns),
            }
        }
    }
}

#[no_mangle]
pub extern "C" fn new_client() -> *mut Client {
    let c = Client {};
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

        unsafe { &*client }.receive_audio(id, d, force_steam_voice)
    }
}

#[no_mangle]
pub extern "C" fn free_client(client: *mut Client) {
    if client.is_null() {
        return;
    }

    unsafe { Box::from_raw(client) };
}
