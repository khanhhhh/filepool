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
	server := storage.NewMemStorage()
	client, err := storage.NewFileStorage("./data")
	if err != nil {
		log.Fatal(err)
	}
	decryptor := crypto.NewPlainDecryptor()
	hasher := crypto.NewHasher()
	pool := pool.NewPool(decryptor, hasher, server, client)
	for {
		fmt.Print("Press 'Enter' to refresh!")
		pool.Upload()
		pool.Download()
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}
}
