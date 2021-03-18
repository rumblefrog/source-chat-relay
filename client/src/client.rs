use std::collections::hash_map::Entry;
use std::collections::HashMap;
use std::sync::Arc;

use lazy_static::lazy_static;

use tokio::runtime::Runtime;
use tokio::sync::RwLock;

use crate::packet::{Packet, Payload};
use crate::player::Player;
use crate::Result;

lazy_static! {
    static ref RUNTIME: Runtime = Runtime::new().unwrap();
}

/// Source Chat Relay client.
///
/// This is only constructed from the shim only.
pub struct Client {
    /// Represents all players on the server.
    ///
    /// It is garbage-collected on silence packet.
    players: Arc<RwLock<HashMap<u64, Player>>>,

    // temp: Vec<u8>,
    temp: Arc<RwLock<Vec<u8>>>,
}

// By default *mut T is not safe for send.
// Need to ensure the shim side safety instead.
unsafe impl Send for Client {}

impl Default for Client {
    fn default() -> Self {
        Self {
            players: Arc::new(RwLock::new(HashMap::new())),
            temp: Arc::new(RwLock::new(Vec::new())),
        }
    }
}

impl Client {
    pub async fn receive_audio(
        &mut self,
        _id: u64,
        data: &[u8],
        _force_steam_voice: bool,
    ) -> Result<()> {
        let mut packet = Packet::from_bytes(data)?;

        let header = packet.header()?;

        let payload = packet.payload()?;

        let mut players = self.players.write().await;
        let mut temp = self.temp.write().await;

        let player = match players.entry(header.steam_id) {
            Entry::Occupied(p) => p.into_mut(),
            Entry::Vacant(v) => v.insert(Player::new()?),
        };

        match payload {
            Payload::OpusPLC(data) => {
                println!("!!!!! Opus PLC {}", data.len());

                match player.transcode(data) {
                    Ok(mut d) => {
                        println!("ok transcode {}", d.len());
                        temp.append(&mut d);
                    }
                    Err(e) => println!("{:?}", e),
                }
            }
            Payload::Silence(ns) => {
                println!("!!!!! Silence {}", ns);

                players.retain(|_k, p| !p.is_stale());

                println!("========== written {}", temp.len());

                std::fs::write("out.data", &*temp)?;
                temp.clear();

                // Silence payload should also be sent on the wire.
            }
        }

        Ok(())
    }
}

#[no_mangle]
pub extern "C" fn new_client() -> *mut Client {
    let c = Client::default();
    let b = Box::new(c);

    Box::into_raw(b)
}

// TODO: Remove force_steam_voice
#[no_mangle]
pub unsafe extern "C" fn receive_audio(
    client: *mut Client,
    id: u64,
    bytes: i32,
    data: *const std::os::raw::c_char,
    force_steam_voice: bool,
) {
    if !client.is_null() {
        let d = std::slice::from_raw_parts(data as *const u8, bytes as usize).to_owned();

        let client = &mut *client;

        RUNTIME.spawn(async move {
            let _ = client.receive_audio(id, &d, force_steam_voice).await;
        });
    }
}

#[no_mangle]
pub unsafe extern "C" fn free_client(client: *mut Client) {
    if client.is_null() {
        return;
    }

    Box::from_raw(client);
}
