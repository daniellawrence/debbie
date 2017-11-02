#!/bin/bash -e
go build main.go
mkdir -p build/usr/local/bin/
mv main build/usr/local/bin/debbie
go run main.go -name debbie -path ./build
