Mycroft-Core [Go]
=================

This is a rewrite of [Mycroft-core](https://github.com/rit-sse-mycroft/core/)
in [Go](http://golang.org/). Preliminary tests show large speed improvements
over the c# implementation.


Building
--------

1. Clone this repository

2. Run `export GOPATH=...` where `...` is the directory containing 'src' 'pkg'
and 'bin' folders.

3. Compile using `go install github.com/robmcl4/mycroft`


Running
-------

After building, run with the command `$GOPATH/bin/mycroft`. Note that currently
this application does not support c#.
