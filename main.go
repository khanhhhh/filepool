package main

import (
	"fmt"
	"log"

	"github.com/khanhhhh/filepool/storage"
)

func main() {
	s := storage.NewFileStorage("./data")
	fmt.Println(s.List())
	err := s.Write("./data/bar/haha.txt", []byte("haha"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(s.List())
}
