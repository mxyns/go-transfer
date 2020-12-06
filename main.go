package main

import (
	"flag"
	"fmt"
	"github.com/mxyns/go-transfer/io"
	log "github.com/sirupsen/logrus"
)

func main() {

	// define general flags (some are defined in other packages' init functions, e.g: routines)
	runServer := flag.Bool("r", false, "run in server mode")
	address := flag.String("a", "127.0.0.1", "address to host on / connect to")
	proto := flag.String("P", "tcp", "protocol")
	port := flag.Uint("p", 8887, "port")
	timeout := flag.String("t", "10s", "client connection timeout")
	debugLevel := flag.String("l", "panic", "debug level")
	customFormatter := flag.Bool("f", true, "use custom formatter for network log")

	flag.Parse()

	if *customFormatter {
		log.SetFormatter(&io.Formatter{
			TimestampFormat: "2006-01-02 15:04:05",
			LogFormat:       "[go-tcp][%lvl%]: %time% - %msg% %fields%\n",
		})
	}

	// apply custom debug level for go-tcp logs
	level, err := log.ParseLevel(*debugLevel)
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		log.SetLevel(level)
	}

	if *runServer {
		println("Receiver mode")
		startServer(address, proto, port)
	} else {
		println("Sender mode")
		startClient(address, proto, port, timeout)
	}
}