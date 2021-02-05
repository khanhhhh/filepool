package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/khanhhhh/filepool/crypto"
	"github.com/khanhhhh/filepool/pool"
	"github.com/khanhhhh/filepool/storage"
)

const (
	modeUpload   = 0
	modeDownload = 1
)

func upload(keyPath string, serverPath string, clientPath string, mode int, createKey bool) {
	server, err := storage.NewFileStorage("./server_data")
	if err != nil {
		log.Fatal(err)
	}
	client, err := storage.NewFileStorage("./client_data")
	if err != nil {
		log.Fatal(err)
	}
	// crypto.NewAESKey("./key")
	decryptor, err := crypto.NewAESDecryptor("./key")
	if err != nil {
		log.Fatal(err)
	}
	hasher := crypto.NewHasher()
	pool := pool.NewPool(decryptor, hasher, server, client)
	switch mode {
	case modeUpload:
		for {
			pool.Upload()
			fmt.Print("Press 'Enter' to upload!")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
		}
	case modeDownload:
		pool.Download()
	}
}

func main() {
	var keyPath = flag.String("key", "./key", "Path to key")
	var serverPath = flag.String("server", "./server_data", "Path to server directory")
	var clientPath = flag.String("client", "./server_data", "Path to client directory")
	var modeStr = flag.String("mode", "upload", "Mode: [upload, download]")
	var createKey = flag.Bool("create-key", false, "whether create a new key or not")
	flag.Parse()
	modeMap := map[string]int{
		"upload":   modeUpload,
		"download": modeDownload,
	}
	mode, ok := modeMap[*modeStr]
	if !ok {
		log.Fatal("wrong mode")
	}
	fmt.Println("mode: ", *modeStr)
	fmt.Println("key: ", *keyPath)
	fmt.Println("client: ", *clientPath)
	fmt.Println("server: ", *serverPath)
	upload(*keyPath, *serverPath, *clientPath, mode, *createKey)
}
