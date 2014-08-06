Mycroft-Core [Go]
=================

[![GoDoc](https://godoc.org/github.com/robmcl4/Mycroft-Core-Go?status.png)](https://godoc.org/github.com/robmcl4/Mycroft-Core-Go) [![Build Status](https://travis-ci.org/robmcl4/Mycroft-Core-Go.svg?branch=master)](https://travis-ci.org/robmcl4/Mycroft-Core-Go)

This is a rewrite of [Mycroft-core](https://github.com/rit-sse-mycroft/core/)
in [Go](http://golang.org/). Preliminary tests show large speed improvements
over the c# implementation.


Cloning
--------

This repository is cloned in a manner similar to other go packages:

1. Go code is generally stored in one folder. Create a folder for your
   go code if you do not have one already.
2. Make sure the GOPATH environment variable is set to this folder.
   Use the command `export GOPATH=ABSOLUTE_PATH_HERE`.
3. Get the code using `go get github.com/robmcl4/Mycroft-Core-Go/mycroft`

After this completes the code is available at `$GOPATH/src/github.com/robmcl4/Mycroft-Core-Go/mycroft`.
This is a full git repository, edit the code here and commit changes here.

A binary was also built as a part of `go get`, which is available in `$GOPATH/bin`

Running
-------

After building, run with the command `$GOPATH/bin/mycroft`.
