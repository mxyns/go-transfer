package io

import (
	"fmt"
	"github.com/mxyns/go-tcp/filet"
	"github.com/mxyns/go-tcp/filet/requests"
	"github.com/mxyns/go-tcp/filet/requests/defaultRequests"
	"os"
	"path/filepath"
	"strings"
)

func MoveFile(oldPath string, newPath string) error {

	newPath = strings.ReplaceAll(newPath, "\\", "/")
	fileDir, _ := filepath.Split(newPath)

	if info, err := os.Stat(fileDir); err != nil || !info.IsDir() {
		if err2 := os.MkdirAll(fileDir, os.ModeDir); err2 != nil {
			fmt.Printf("Mkdir All %v\n", err)
			return err2
		}
	}
	return os.Rename(oldPath, newPath)
}

func TrySend(client *filet.Client, path string) {
	if info, err := os.Stat(path); len(path) == 0 || err != nil {
		fmt.Println("Wrong file path, use -i <path>")
		return
	} else {
		if !info.IsDir() {
			sendFile(client, path)
		} else {
			err = sendDir(client, path)
			if err != nil {
				fmt.Printf("%v\n", err)
			}
		}
	}
}

func sendFile(client *filet.Client, path string) {

	fmt.Printf("Sending file : %v\n", path)
	resp := (*client.Send(requests.MakeGenericPack(
		defaultRequests.MakeFileRequest(path, true),
		defaultRequests.MakeTextRequest(path),
	))).(*defaultRequests.TextRequest)

	fmt.Printf("receiver => %v\n", resp.GetText())
}

func sendDir(client *filet.Client, root string) error {

	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	rel_root, err := filepath.Rel(wd, filepath.Join(wd, root))
	if err != nil {
		return err
	}

	sendFileToClient := func(path string, info os.FileInfo, err error) error {

		relative, err := filepath.Rel(root, path)
		shards := strings.Split(root, "/")
		relative = filepath.Join(shards[len(shards)-1], relative)
		if err != nil {
			return err
		}

		if info.IsDir() {
			fmt.Printf("Sending dir : %v\n", relative)
		} else {
			sendFile(client, relative)
		}

		return nil
	}

	err = filepath.Walk(root, sendFileToClient)
	if err != nil {
		fmt.Printf("Error while trying to send dir %v : %v\n", rel_root, err)
	}
	return err
}
