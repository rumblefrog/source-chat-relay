fn main() {
    prost_build::compile_protos(&["voice.proto", "common.proto"], &["../messages/"])
        .unwrap()
}
