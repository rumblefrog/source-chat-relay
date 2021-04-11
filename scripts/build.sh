#!/bin/bash
set -e

# Client
cargo build --target i686-unknown-linux-gnu -p client
mkdir -p build && cd build
python3 ../client/shim/configure.py -s sdk2013 --hl2sdk-root ../ --mms_path ../metamod-source/ --client-path ../target/i686-unknown-linux-gnu/debug/ --sm-path ../sourcemod/ --enable-optimize
ambuild

# Rest
cargo build --workspace --exclude client
