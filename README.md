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

The following command should create a package with a single file.

    $ mkdir -p /tmp/data/example/
	$ date > /tmp/data/example/date.text
	$ debbie -name example -path /tmp/data
	example_0.0.1-1.deb
	
It should be installable via dpkg

How to help
--------------

Let me know I am crazy
