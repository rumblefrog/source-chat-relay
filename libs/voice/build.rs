fn main() {
    prost_build::compile_protos(&["voice.proto"], &["../../messages/"])
        .unwrap()
}
