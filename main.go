package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/khanhhhh/filepool/crypto"
	"github.com/khanhhhh/filepool/pool"
	"github.com/khanhhhh/filepool/storage"
)

func run() {
	server, err := storage.NewFileStorage("./server_data")
	if err != nil {
		log.Fatal(err)
	}
	client, err := storage.NewFileStorage("./client_data")
	if err != nil {
		log.Fatal(err)
	}
	// crypto.NewAESKey("./key")
	// decryptor, err := crypto.NewAESDecryptor("./key")
	// if err != nil {
	// 	log.Fatal(err)
	//}
	decryptor := crypto.NewPlainDecryptor()
	hasher := crypto.NewHasher()
	pool := pool.NewPool(decryptor, hasher, server, client)
	for {
		pool.Upload()
		pool.Download()
		fmt.Print("Press 'Enter' to refresh!")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}
}

func main() {
	run()
}
