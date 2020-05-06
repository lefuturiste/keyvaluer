# Keyvaluer

[![Build Status](https://img.shields.io/endpoint.svg?url=https%3A%2F%2Factions-badge.atrox.dev%2Flefuturiste%2Fkeyvaluer%2Fbadge&style=for-the-badge)](https://actions-badge.atrox.dev/lefuturiste/keyvaluer/goto)

Keyvaluer is a simple and light Go implementation of a redis server.

## But why don't use redis direclty ?

The purpose of developing a alternative implementation of the redis server is to have a production ready binary which can be used very easily on any kind of servers (because sometime you can't install via a package manager or build redis on a server).

If you have the possibility to install redis with a package manager or whatever, if you can use the original redis server implementation please do.

Another advantage of this approach is to have a relativly light and simple implementation which can be easily tuned, forked to make changes.

## Usage

You can download a binary in the release github page. Or directly build the binary using go. Then run the binary and you are ready.
By default we listen on loopback but you can make keyvaluer to listen to all interfaces. If you are doing that make sure to set a password with the `REQUIRED_PASS` environement variable. 

**This server implementation does not require a password when listening to all interfaces, so be careful.*

### Environment variables

- `HOST`: define the listening host of the server, default value is `"localhost"`
- `PORT`: define the listening port of the server, default value is `6379`
- `REQUIRED_PASS`: define the password of the server (with the `AUTH` command), default value is `""` (empty string). When the value is an empty string, no password is required.
- `LOG_LEVEL`: define the log level can be one of those values: `["trace", "debug", "info", "warn", "error", "fatal", "panic"]`, default value is `"info"`

## Build

To build simply use this command: `go build -o keyvaluer`.
You can now run keyvaluer: `./keyvaluer`.

## Run unit test

There is *1* test file for now: `main_test.go`.

Use `go test` to launch tests.

## Contributions

Contributions are welcomed, feel free to open an *Issue* or a *Pull Request* and I will be happy to help you!