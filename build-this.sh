#!/bin/bash -e

# build package
go build debbie

# setup dirs for package
mkdir -p build/usr/local/bin/
mv debbie build/usr/local/bin/
./build/usr/local/bin/debbie -name debbie -path ./build -output-dir ./output

# clean-up files
rm    build/usr/local/bin/debbie
rmdir build/usr/local/bin
rmdir build/usr/local/
rmdir build/usr/
rmdir build/
