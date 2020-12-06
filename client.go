package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/mxyns/go-tcp/filet"
	"github.com/mxyns/go-tcp/filet/requests"
	"github.com/mxyns/go-tcp/filet/requests/defaultRequests"
	"os"
	"strings"
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
		sendFile(client, *path)
	}
	clientTerminalInput(client)
}

func sendFile(client *filet.Client, path string) {

	if _, err := os.Stat(path); len(path) == 0 || err != nil {
		fmt.Println("Wrong file path, use -i <path>")
		return
	} else {
		fmt.Printf("Sending file : %v\n", path)
	}

	shards := strings.Split(path, "/")
	resp := (*client.Send(requests.MakeGenericPack(
		defaultRequests.MakeFileRequest(path, true),
		defaultRequests.MakeTextRequest(shards[len(shards)-1]),
	))).(*defaultRequests.TextRequest)

	fmt.Printf("receiver => %v\n", resp.GetText())
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
			sendFile(client, line)
		}
	}
}
