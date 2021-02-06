package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/khanhhhh/filepool/crypto"
	"github.com/khanhhhh/filepool/pool"
	"github.com/khanhhhh/filepool/storage"
)

const (
	modeUpload   = 0
	modeDownload = 1
)

func upload(keyPath string, serverPath string, clientPath string, mode int, createKey bool) {
	server, err := storage.NewFileStorage(serverPath)
	if err != nil {
		log.Fatal(err)
	}
	client, err := storage.NewFileStorage(clientPath)
	if err != nil {
		log.Fatal(err)
	}
	if createKey {
		err = crypto.NewAESKey(keyPath)
		if err != nil {
			log.Fatal(err)
		}
	}
	decryptor, err := crypto.NewAESDecryptor(keyPath)
	if err != nil {
		log.Fatal(err)
	}
	hasher := crypto.NewHasher()
	p := pool.NewPool(decryptor, hasher, server, client)
	switch mode {
	case modeUpload:
		p.Upload()
	case modeDownload:
		p.Download()
	}
}

func main() {
	var isTest = flag.Bool("test", false, "test")
	var keyPath = flag.String("key-path", "./key", "Path to key")
	var serverPath = flag.String("server-path", "./server_data", "Path to server directory")
	var clientPath = flag.String("client-path", "./client_data", "Path to client directory")
	var modeStr = flag.String("mode", "upload", "Mode: [upload, download]")
	var createKey = flag.Bool("create-key", false, "whether create a new key or not")
	flag.Parse()
	if *isTest {
		test()
		return
	}
	modeMap := map[string]int{
		"upload":   modeUpload,
		"download": modeDownload,
	}
	mode, ok := modeMap[*modeStr]
	if !ok {
		log.Fatal("wrong mode")
	}
	fmt.Println("mode: ", *modeStr)
	fmt.Println("key-path: ", *keyPath)
	fmt.Println("server-path: ", *serverPath)
	fmt.Println("client-path: ", *clientPath)
	fmt.Println("create-key: ", *createKey)
	upload(*keyPath, *serverPath, *clientPath, mode, *createKey)
}

func test() {
	aes1, _ := crypto.NewAESDecryptor("./key")
	aes2, _ := crypto.NewAESDecryptor("./key")
	dataIn := make([]byte, 1000000000)
	dataOut, _ := ioutil.ReadAll(aes1.Encrypt(bytes.NewReader(dataIn)))
	fmt.Println(dataIn[:500])
	fmt.Println(dataOut[:500])
	///
	dataOutOut, _ := ioutil.ReadAll(aes2.Decrypt(bytes.NewReader(dataOut)))
	fmt.Println(dataOutOut[:500])
}
