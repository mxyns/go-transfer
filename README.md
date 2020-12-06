# go-transfer
Simple Go peer-to-peer file sharing progam

## Building project
`git clone https://github.com/mxyns/go-transfer`
`cd go-transfer`
`go build`

## Launching
If you want to launch the receiver mode you have to use the `-r` (or `-r=true`) flag, otherwise the sender mode is launched by default.

Common flags :
* `-a`  (string, default = 127.0.0.1) : the address of the receiver or address to listen on (for the receiver)
* `-P`  (string, default = "tcp") : the protocol you want to use
* `-p`  (int, default = 8887) : the port of the server
* `-t`  (string, default = "10s") : the time after client stops trying to connect
* `-l`  (string, default = "panic") : the level of debug, [logrus](https://godoc.org/github.com/sirupsen/logrus) has seven log levels: Trace, Debug, Info, Warn, Error, Fatal, and Panic
