#!/bin/bash -e

# build package
go build debbie

# setup dirs for package
mkdir -p build/usr/local/bin/
mv debbie build/usr/local/bin/
./build/usr/local/bin/debbie -name debbie -path ./build
mv /tmp/debbie_0.0.1_all.deb output/

# clean-up files
mv build/usr/local/bin/debbie output/
rmdir build/usr/local/bin
rmdir build/usr/local/
rmdir build/usr/
rmdir build/
