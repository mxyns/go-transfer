package main

import (
	"bufio"
	"fmt"
	"github.com/mxyns/go-tcp/filet"
	gotcp "github.com/mxyns/go-tcp/filet/requests"
	"github.com/mxyns/go-tcp/filet/requests/defaultRequests"
	"github.com/mxyns/go-transfer/io"
	"net"
	"os"
	"strings"
	"sync"
)

func startServer(address *string, proto *string, port *uint) {

	openConn := &sync.WaitGroup{}

	server := filet.Server{
		Address: &filet.Address{
			Proto: *proto,
			Addr:  *address,
			Port:  uint32(*port),
		},
		Clients:          make([]*net.Conn, 5),
		ConnectionWaiter: openConn,
		RequestHandler: func(client *net.Conn, request *gotcp.Request) {

			switch (*request).(type) {
			case *gotcp.Pack:
				{
					pack := (*request).(*gotcp.Pack)
					fileReq := (*pack.GetRequests()[0]).(*defaultRequests.FileRequest)
					textReq := (*pack.GetRequests()[1]).(*defaultRequests.TextRequest)

					oldPath := fileReq.GetPath()
					shards := strings.Split(oldPath, "/")
					shards[len(shards)-1] = textReq.GetText()
					newPath := strings.Join(shards, "/")

					if err := io.MoveFile(oldPath, newPath); err != nil {
						var resp gotcp.Request = defaultRequests.MakeTextRequest("ERR")
						_, _, _ = gotcp.SendRequestOn(client, &resp)
						fmt.Printf("Couldn't move file %v\n", err)
						return
					} else {
						fmt.Printf("received => %v (%v bytes)\n", newPath, fileReq.GetFileSize())
					}
				}
			default:
				{
					var resp gotcp.Request = defaultRequests.MakeTextRequest("ERR")
					_, _, _ = gotcp.SendRequestOn(client, &resp)
					return
				}
			}

			if (*request).Info().WantsResponse {
				var resp gotcp.Request = defaultRequests.MakeTextRequest("OK")
				_, _, _ = gotcp.SendRequestOn(client, &resp)
			}
		},
	}
	defer server.Close()
	go server.Start()

	openConn.Add(1)

	serverTerminalInput(openConn)
	openConn.Wait()
}

func serverTerminalInput(group *sync.WaitGroup) {
	fmt.Print("type !stop to stop > ")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		line := scanner.Text()
		fmt.Printf("got > %v\n", line)
		if line == "!stop" { // server mode
			if group != nil {
				(*group).Done()
			}
			break
		}
	}
}
