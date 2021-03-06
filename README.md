[![Build Status](https://travis-ci.org/daniellawrence/debbie.svg?branch=master)](https://travis-ci.org/daniellawrence/debbie)

debbie
--------

The easiest way to create a .deb package for debian from a directory.

Goals
------

* Deb (*dead*) simple to use
* Pure go, Do not shell out for any reason
* Buffers over temp files, do not use any temp files written to the disk
* Learn, Gain a better understanding of golang + deb packages
* Fast, Be the fastest way to make a deb package on the internet

How to install
-----------------

Its a package!

    wget https://raw.githubusercontent.com/daniellawrence/debbie/output/debbie-0.0.1_all.deb
    sudo dpkg -i debbie-0.0.1_all.deb


How to use
------------

Accept most defaults

    debbie -name nginx -path output

Custom versions

    debbie -name nginx -path output -version 1.2.3

Custom install directory

    debbie -name nginx -path output -install-path /usr/local/nginx

Custom Maintainer

    debbie -name nginx -path output -maintainer "Daniel Lawrence"

Custom Maintainer Email

    debbie -name nginx -path output -maintainer-email "dannyla@linux.com"

Custom Description

    debbie -name nginx -path output -description "dir2deb in go"


Goal review
--------------

*dead simple to use*

* Single file
* Can be installed via apt/dpkg
* No ruby/rbenv/pip/python/virtualenv

*Pure go*

No shells
	 
*Buffers over temp files*

No Temp files
	 
*Learning go + deb*

Its getting better
	 
*Fast*

single file 2Mb package created in 340ms

     $ time build/usr/local/bin/debbie -name debbie -path ./build
     2017/11/01 22:13:26 Created file: /tmp/debbie_0.0.1_all.deb
     build/usr/local/bin/debbie -name debbie -path ./build  0.34s user 0.00s system 100% cpu 0.342 total

Example package (this)
-----------------

This is the script to turn this into a deb package

    #!/bin/bash
    mkdir -p /tmp/
    go build main.go
    mkdir -p build
    mv main build/debbie
    ./build/debbie -name debbie -path ./build
    dpkg-deb --info /tmp/debbie*deb

It should be installable via dpkg

    root@b4c5c08042ae:/# dpkg -i /src/debbie_0.0.1_all.deb
    (Reading database ... 6490 files and directories currently installed.)
    Preparing to unpack /src/debbie_0.0.1_all.deb ...
    Unpacking debbie (0.0.1-1) ...
    Setting up debbie (0.0.1-1) ...


How to help
--------------

Let me know I am crazy - <dannyla@linux.com>

TODO / Ideas
---------------

* Validate .deb packages are as good as they can be
* More deb package options
* Investigate rpm

GOPATH?
-------

    export GOPATH=$PWD:$PWD/vendor


Projects using this
------------------------

* 
