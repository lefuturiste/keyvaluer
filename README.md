# Keyvaluer

Keyvaluer is a simple and light Go implementation of a redis server.

## But why don't use redis direclty ?

The purpose of developing a alternative implementation of the redis server is to have a production ready binary which can be used very easily on any kind of servers (because sometime you can't install via a package manager or build redis on a server).

If you have the possibility to install redis with a package manager or whatever, if you can use the original redis server implementation please do.

## Build

To build simply use this command: `go build -o keyvaluer`.
You can now run keyvaluer: `./keyvaluer`.

## Run unit test

There is *1* test file for now: `main_test.go`.

Use `go test` to launch tests.
