#!/bin/sh

cd GitGoat
go build -o git-goat
cd ..

cd goat-pusher
cargo build --release
cargo run --release
cd ..
