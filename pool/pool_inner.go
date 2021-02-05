package pool

import (
	"fmt"
	"io"
	"io/ioutil"

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
	for _, filename := range p.client.List() {
		plainTextBuf1, err := p.client.Read(filename)
		defer plainTextBuf1.Close()
		if err != nil {
			continue
		}
		wantHashText, err := p.hasher.Hash(plainTextBuf1)
		fmt.Println(wantHashText)
		gotHashTextBuf, err := p.server.Read(toHashFile(filename))
		if err == nil {
			defer gotHashTextBuf.Close()
			gotHashText, err := ioutil.ReadAll(gotHashTextBuf)
			fmt.Println(gotHashText)
			if err == nil && sameHash(wantHashText, gotHashText) {
				fmt.Println("skip upload")
				continue
			}
		}
		// write
		writeHashTextBuf, err := p.server.Write(toHashFile(filename))
		if err != nil {
			continue
		}
		defer writeHashTextBuf.Close()
		writeHashTextBuf.Write(wantHashText)
		cipherTextBuf, err := p.server.Write(filename)
		if err != nil {
			continue
		}
		defer cipherTextBuf.Close()
		plainTextBuf2, err := p.client.Read(filename)
		if err != nil {
			continue
		}
		defer plainTextBuf2.Close()
		_, err = io.Copy(cipherTextBuf, p.decryptor.Encrypt(plainTextBuf2))
	}
}
func (p *pool) Download() {
	for _, filename := range p.server.List() {
		if isHashFile(filename) {
			// skip hash file
			continue
		}
		wantHashTextBuf, err := p.server.Read(toHashFile(filename))
		if err != nil {
			continue
		}
		wantHashText, err := ioutil.ReadAll(wantHashTextBuf)
		if err != nil {
			continue
		}
		gotPlainTextBuf, err := p.client.Read(filename)
		if err == nil {
			gotHashText, err := p.hasher.Hash(gotPlainTextBuf)
			if err == nil && sameHash(wantHashText, gotHashText) {
				fmt.Println("skip download")
				continue
			}
		}
		// write
		cipherTextBuf, err := p.server.Read(filename)
		if err != nil {
			continue
		}
		defer cipherTextBuf.Close()
		plainTextBuf, err := p.client.Write(filename)
		if err != nil {
			continue
		}
		defer plainTextBuf.Close()
		_, err = io.Copy(plainTextBuf, p.decryptor.Decrypt(cipherTextBuf))
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
