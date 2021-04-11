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

use std::collections::hash_map::Entry;
use std::collections::HashMap;
use std::ffi::CStr;
use std::sync::Arc;

use lazy_static::lazy_static;

use tokio::runtime::Runtime;
use tokio::sync::RwLock;

use crate::packet::{Packet, Payload};
use crate::player::Player;
use crate::Result;

lazy_static! {
    static ref CLIENT: Client = Client::default();
    static ref RUNTIME: Runtime = Runtime::new().unwrap();
}

/// Source Chat Relay client.
///
/// Measures need to be taken to ensure ptr lifetime on the shim.
pub struct Client(Arc<RwLock<ClientInner>>);

struct ClientInner {
    /// Map of players on the server.
    /// This will be updated upon client join/leave.
    players: HashMap<u64, Player>,
}

impl Default for Client {
    fn default() -> Self {
        Self(Arc::new(RwLock::new(ClientInner {
            players: HashMap::new(),
        })))
    }
}

impl Client {
    pub async fn receive_audio(&self, data: &[u8]) -> Result<()> {
        let mut packet = Packet::from_bytes(data)?;

        let header = packet.header()?;

        let payload = packet.payload()?;

        let mut inner = self.0.write().await;

        let player = match inner.players.entry(header.steam_id) {
            Entry::Occupied(p) => p.into_mut(),
            Entry::Vacant(v) => v.insert(Player::new()?),
        };

        match payload {
            Payload::OpusPLC(data) => {
                println!("!!!!! Opus PLC {}", data.len());

                match player.transcode(data) {
                    Ok(d) => {
                        println!("ok transcode {}", d.len());
                    }
                    Err(e) => println!("{:?}", e),
                }
            }
            Payload::Silence(ns) => {
                println!("!!!!! Silence {}", ns);

                // Silence payload should also be sent on the wire.
            }
        }

        Ok(())
    }

    pub async fn client_put_in_server(&self, steamid: u64, name: &str) -> Result<()> {
        let mut inner = self.0.write().await;

        // TODO: Handle name

        inner.players.insert(steamid, Player::new()?);

        Ok(())
    }

    pub async fn client_disconnect(&self, steamid: u64) -> Result<()> {
        let mut inner = self.0.write().await;

        inner.players.remove(&steamid);

        Ok(())
    }
}

#[no_mangle]
pub unsafe extern "C" fn receive_audio(
    bytes: i32,
    data: *const std::os::raw::c_char,
) {
    let d = std::slice::from_raw_parts(data as *const u8, bytes as usize).to_owned();

    RUNTIME.spawn(async move {
        let _ = CLIENT.receive_audio(&d).await;
    });
}

#[no_mangle]
pub unsafe extern "C" fn client_put_in_server(
    steamid: u64,
    name: *const std::os::raw::c_char,
) {
    let name = CStr::from_ptr(name).to_string_lossy();

    RUNTIME.spawn(async move {
        let _ = CLIENT.client_put_in_server(steamid, &name).await;
    });
}

#[no_mangle]
pub unsafe extern "C" fn client_disconnect(steamid: u64) {
    RUNTIME.spawn(async move {
        let _ = CLIENT.client_disconnect(steamid).await;
    });
}
