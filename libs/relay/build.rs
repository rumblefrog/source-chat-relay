fn main() {
    prost_build::compile_protos(&["relay.proto", "voice.proto"], &["../../messages/"])
        .unwrap()
}
