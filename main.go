package main

import (
	"fmt"
	"log"

	"github.com/khanhhhh/filepool/crypto"
)

func main() {
	d, err := crypto.NewAESDecryptorToFile("./key")
	if err != nil {
		log.Fatal(err)
	}
	e, err := crypto.NewAESDecryptorFromFile("./key")
	cipher, err := d.Encrypt([]byte("haha"))
	if err != nil {
		log.Fatal(err)
	}
	out, err := e.Decrypt(cipher)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(out))
}
