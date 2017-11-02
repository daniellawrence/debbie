#!/bin/bash
go build main.go
mkdir -p build
mv main build/debbie
./build/debbie -name debbie -path ./build
dpkg-deb --info /tmp/debbie*deb
