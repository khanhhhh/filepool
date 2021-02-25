package pool

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"

	"github.com/khanhhhh/filepool/crypto"
	"github.com/khanhhhh/filepool/storage"
)

type pool struct {
	hasher    crypto.Hasher
	decryptor crypto.Decryptor
	server    storage.Storage
	client    storage.Storage
}

func (p *pool) Upload() {
	log.Println("Uploading")
	for _, filename := range p.client.List() {
		log.Println("Reading client file: ", filename)
		plainTextBuf1, err := p.client.Read(filename)
		if err != nil {
			log.Println("Read error: ", err)
			continue
		}
		log.Println("Hashing client file:", filename)
		wantHashText, err := p.hasher.Hash(plainTextBuf1)
		_ = plainTextBuf1.Close()
		log.Println("Hashing client file: ", filename)
		gotHashTextBuf, err := p.server.Read(toHashFile(filename))
		if err == nil {
			gotHashText, err := ioutil.ReadAll(gotHashTextBuf)
			if err == nil && sameHash(wantHashText, gotHashText) {
				fmt.Println("Skipped upload file: ", filename)
				continue
			}
			_ = gotHashTextBuf.Close()
		}
		// write
		log.Println("Writing server file:", filename)
		writeHashTextBuf, _ := p.server.Write(toHashFile(filename))
		_, _ = writeHashTextBuf.Write(wantHashText)
		_ = writeHashTextBuf.Close()
		cipherTextBuf, _ := p.server.Write(filename)
		plainTextBuf2, _ := p.client.Read(filename)
		_, err = io.Copy(cipherTextBuf, p.decryptor.Encrypt(plainTextBuf2))
		_ = cipherTextBuf.Close()
	}
}
func (p *pool) Download() {
	log.Println("Downloading")
	for _, filename := range p.server.List() {
		if isHashFile(filename) {
			// skip hash file
			continue
		}
		log.Println("Reading server hash file: ", filename)
		wantHashTextBuf, err := p.server.Read(toHashFile(filename))
		if err != nil {
			log.Println("Read error: ", err)
			continue
		}
		wantHashText, err := ioutil.ReadAll(wantHashTextBuf)
		if err != nil {
			log.Println("Read error: ", err)
			continue
		}
		log.Println("Reading client file:", filename)
		gotPlainTextBuf, err := p.client.Read(filename)
		if err == nil {
			log.Println("Hashing client file:", filename)
			gotHashText, err := p.hasher.Hash(gotPlainTextBuf)
			if err == nil && sameHash(wantHashText, gotHashText) {
				fmt.Println("Skip download file: ", filename)
				continue
			}
		}
		// write
		log.Println("Writing client file:", filename)
		cipherTextBuf, _ := p.server.Read(filename)
		plainTextBuf, err := p.client.Write(filename)
		_, err = io.Copy(plainTextBuf, p.decryptor.Decrypt(cipherTextBuf))
		_ = cipherTextBuf.Close()
		_ = plainTextBuf.Close()
	}
}

func toHashFile(filename string) string {
	return filename + ".hash"
}

func isHashFile(filename string) bool {
	return len(filename) >= 5 && filename[len(filename)-5:] == ".hash"
}

func sameHash(hash1 []byte, hash2 []byte) bool {
	if len(hash1) != len(hash2) {
		return false
	}
	for i := 0; i < len(hash1); i++ {
		if hash1[i] != hash2[i] {
			return false
		}
	}
	return true
}

// NewPool :
func NewPool(decryptor crypto.Decryptor, hasher crypto.Hasher, server storage.Storage, client storage.Storage) Pool {
	return &pool{
		decryptor: decryptor,
		hasher:    hasher,
		server:    server,
		client:    client,
	}
}
