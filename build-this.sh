#!/bin/bash -e

# build package
go build debbie

# setup dirs for package
mkdir -p build/usr/local/bin/
cp debbie build/usr/local/bin/
mv debbie output
if [ -e build/usr/local/bin/d ];then
    rm build/usr/local/bin/d
fi
(cd build/usr/local/bin/; ln -s debbie d)

./build/usr/local/bin/debbie -name debbie -version 0.0.2 -path ./build -output-dir ./output

ls -l build/usr/local/bin
# clean-up files
rm    build/usr/local/bin/debbie
rm    build/usr/local/bin/d
rmdir build/usr/local/bin
rmdir build/usr/local/
rmdir build/usr/
rmdir build/
