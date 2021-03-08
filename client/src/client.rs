pub struct Client {

}

impl Client {
    #[no_mangle]
    pub extern fn new_client() -> *mut Client {
        let c = Client {};
        let b = Box::new(c);

        Box::into_raw(b)
    }
}
