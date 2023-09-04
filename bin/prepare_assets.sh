#!/bin/bash

set -e

# ensure existence of release folder
if ! [ -d "./release" ]; then
    mkdir ./release
fi

# ensure zip is installed
if [ "$(which zip)" = "" ]; then
    apt-get update && apt-get install -y zip
fi

# add execution permission
chmod 750 ./build/zeedas-cli-freebsd-386
chmod 750 ./build/zeedas-cli-freebsd-amd64
chmod 750 ./build/zeedas-cli-freebsd-arm
chmod 750 ./build/zeedas-cli-linux-386
chmod 750 ./build/zeedas-cli-linux-amd64
chmod 750 ./build/zeedas-cli-linux-arm
chmod 750 ./build/zeedas-cli-linux-arm64
chmod 750 ./build/zeedas-cli-linux-riscv64
chmod 750 ./build/zeedas-cli-netbsd-386
chmod 750 ./build/zeedas-cli-netbsd-amd64
chmod 750 ./build/zeedas-cli-netbsd-arm
chmod 750 ./build/zeedas-cli-openbsd-386
chmod 750 ./build/zeedas-cli-openbsd-amd64
chmod 750 ./build/zeedas-cli-openbsd-arm
chmod 750 ./build/zeedas-cli-openbsd-arm64
chmod 750 ./build/zeedas-cli-windows-386.exe
chmod 750 ./build/zeedas-cli-windows-amd64.exe
chmod 750 ./build/zeedas-cli-windows-arm64.exe

# create archives
zip -j ./release/zeedas-cli-freebsd-386.zip ./build/zeedas-cli-freebsd-386
zip -j ./release/zeedas-cli-freebsd-amd64.zip ./build/zeedas-cli-freebsd-amd64
zip -j ./release/zeedas-cli-freebsd-arm.zip ./build/zeedas-cli-freebsd-arm
zip -j ./release/zeedas-cli-linux-386.zip ./build/zeedas-cli-linux-386
zip -j ./release/zeedas-cli-linux-amd64.zip ./build/zeedas-cli-linux-amd64
zip -j ./release/zeedas-cli-linux-arm.zip ./build/zeedas-cli-linux-arm
zip -j ./release/zeedas-cli-linux-arm64.zip ./build/zeedas-cli-linux-arm64
zip -j ./release/zeedas-cli-linux-riscv64.zip ./build/zeedas-cli-linux-riscv64
zip -j ./release/zeedas-cli-netbsd-386.zip ./build/zeedas-cli-netbsd-386
zip -j ./release/zeedas-cli-netbsd-amd64.zip ./build/zeedas-cli-netbsd-amd64
zip -j ./release/zeedas-cli-netbsd-arm.zip ./build/zeedas-cli-netbsd-arm
zip -j ./release/zeedas-cli-openbsd-386.zip ./build/zeedas-cli-openbsd-386
zip -j ./release/zeedas-cli-openbsd-amd64.zip ./build/zeedas-cli-openbsd-amd64
zip -j ./release/zeedas-cli-openbsd-arm.zip ./build/zeedas-cli-openbsd-arm
zip -j ./release/zeedas-cli-openbsd-arm64.zip ./build/zeedas-cli-openbsd-arm64
zip -j ./release/zeedas-cli-windows-386.zip ./build/zeedas-cli-windows-386.exe
zip -j ./release/zeedas-cli-windows-amd64.zip ./build/zeedas-cli-windows-amd64.exe
zip -j ./release/zeedas-cli-windows-arm64.zip ./build/zeedas-cli-windows-arm64.exe

# handle apple binaries
unzip ./build/zeedas-cli-darwin.zip
chmod 750 ./build/zeedas-cli-darwin-amd64
chmod 750 ./build/zeedas-cli-darwin-arm64
zip -j ./release/zeedas-cli-darwin-amd64.zip ./build/zeedas-cli-darwin-amd64
zip -j ./release/zeedas-cli-darwin-arm64.zip ./build/zeedas-cli-darwin-arm64

# calculate checksums
for file in  ./release/*; do
	checksum=$(sha256sum "${file}" | cut -d' ' -f1)
	filename=$(echo "${file}" | rev | cut -d/ -f1 | rev)
	echo "${checksum} ${filename}" >> ./release/checksums_sha256.txt
done
