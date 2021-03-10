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
