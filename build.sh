#!/bin/sh

mkdir bin

cd GitGoat
go build -o ../bin/git-goat
GOOS=linux GOARCH=amd64 go build -o ../bin/git-goat-linux-amd64
GOOS=linux GOARCH=arm64 go build -o ../bin/git-goat-linux-arm64
cd ..

cd goat-pusher
cargo build --release
cargo run --release
cd ..
