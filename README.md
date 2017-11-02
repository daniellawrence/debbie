debbie
--------

Create a .deb package for debian from a path.

Goals
------

* Deb (*dead*) simple to use
* Pure go, Do not shell out for any reason
* Buffers over temp files, do not use any temp files written to the disk
* Learn, Gain a better understanding of golang + deb packages
* Fast, Be the fastest way to make a deb package on the internet

How to package
-----------------

This is the script to turn this into a deb package

	#!/bin/bash
	go build main.go
	mkdir -p build
	mv main build/debbie
	./build/debbie -name debbie -path ./build
	dpkg-deb --info /tmp/debbie*deb
	
It should be installable via dpkg

How to help
--------------

Let me know I am crazy
