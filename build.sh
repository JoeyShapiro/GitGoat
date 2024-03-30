#!/bin/sh

openssl req -x509 -newkey rsa:4096 \
    -keyout key.pem -out cert.pem \
    -sha256 -days 3650 -nodes \
    -subj "/C=XX/ST=StateName/L=CityName/O=CompanyName/OU=CompanySectionName/CN=git-goat"

cd GitGoat
go build -o git-goat
cd ..

cd goat-pusher
cargo build --release
cargo run --release
cd ..
