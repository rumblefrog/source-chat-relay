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

use std::os::raw::c_char;

macro_rules! c_str {
    ($lit:expr) => {
        unsafe {
            std::ffi::CStr::from_ptr(concat!($lit, "\0").as_ptr() as *const std::os::raw::c_char)
                .as_ptr()
        }
    };
}

#[no_mangle]
pub extern "C" fn extension_author() -> *const c_char {
    c_str!("rumblefrog")
}

#[no_mangle]
pub extern "C" fn extension_name() -> *const c_char {
    c_str!("Source Chat Relay")
}

#[no_mangle]
pub extern "C" fn extension_description() -> *const c_char {
    c_str!("Glub Glub")
}

#[no_mangle]
pub extern "C" fn extension_url() -> *const c_char {
    c_str!("https://github.com/rumblefrog/source-chat-relay")
}

#[no_mangle]
pub extern "C" fn extension_license() -> *const c_char {
    c_str!("GPL 3.0")
}

#[no_mangle]
pub extern "C" fn extension_version() -> *const c_char {
    c_str!("3.0.0")
}

#[no_mangle]
pub extern "C" fn extension_date() -> *const c_char {
    c_str!("rusttyy")
}

#[no_mangle]
pub extern "C" fn extension_log_tag() -> *const c_char {
    c_str!("SCR")
}
