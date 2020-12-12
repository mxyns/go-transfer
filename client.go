package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/mxyns/go-tcp/filet"
	"github.com/mxyns/go-transfer/io"
	"os"
)

var path *string

func init() {

	path = flag.String("i", "", "input file path")
}

func startClient(address *string, proto *string, port *uint, timeout *string) {

	client := &filet.Client{
		Address: &filet.Address{
			Proto: *proto,
			Addr:  *address,
			Port:  uint32(*port),
		},
	}
	defer client.Close()
	_, err := client.Start(*timeout)
	if err != nil {
		fmt.Println("Couldn't connect to server.")
		return
	}

	if path != nil {
		io.TrySend(client, *path)
	}
	clientTerminalInput(client)
}

func clientTerminalInput(client *filet.Client) {
	fmt.Print("other path || stop > ")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		line := scanner.Text()
		fmt.Printf("got > %v\n", line)
		if line == "!stop" {
			break
		} else { // client mode
			io.TrySend(client, line)
		}
	}
}
