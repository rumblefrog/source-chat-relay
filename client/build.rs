use std::env;

use cbindgen::Builder;

fn main() {
    let crate_dir = env::var("CARGO_MANIFEST_DIR").unwrap();

    Builder::new()
        .with_crate(crate_dir)
        .generate()
        .expect("Unable to generate bindings")
        .write_to_file("shim/bindings.h");
}
