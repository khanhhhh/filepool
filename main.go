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

func main() {
	client1, err := storage.NewFileStorage("./client1_data")
	if err != nil {
		log.Fatal(err)
	}
	client2, err := storage.NewFileStorage("./client2_data")
	if err != nil {
		log.Fatal(err)
	}
	server, err := storage.NewFileStorage("./server_data")
	if err != nil {
		log.Fatal(err)
	}
	decryptor, err := crypto.NewAESDecryptorToFile("./key")
	if err != nil {
		log.Fatal(err)
	}
	pool := pool.NewPool(decryptor, server, []storage.Storage{client1, client2})
	for {
		fmt.Print("Press 'Enter' to refresh!")
		pool.Upload()
		pool.Download()
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}
}
