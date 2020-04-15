# Keyvaluer

Keyvaluer is a simple and light Go implementation of a redis server.

## But why don't use redis direclty ?

The purpose of developing a alternative implementation of the redis server is to have a production ready binary which can be used very easily on any kind of servers (because sometime you can't install via a package manager or build redis on a server).

If you have the possibility to install redis with a package manager or whatever, if you can use the original redis server implementation please do.
