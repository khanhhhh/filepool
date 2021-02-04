package main

import (
	"filepool/crypto"
	"fmt"
	"log"
)

func main() {
	d, err := crypto.NewDecryptorToFile("./key")
	e, err := crypto.NewEncryptorFromFile("./key.pub")
	if err != nil {
		log.Fatal(err)
	}
	cypher, err := e.Encrypt([]byte("haha"))
	if err != nil {
		log.Fatal(err)
	}
	plain, err := d.Decrypt(cypher)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(plain))
	return
}
