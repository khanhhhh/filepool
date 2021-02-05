package pool

import (
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
		wantHashTextBuf := crypto.NewTransform(plainTextBuf1, p.hasher.Hash)
		wantHashText, err := ioutil.ReadAll(wantHashTextBuf)
		gotHashTextBuf, err := p.server.Read(filename)
		if err == nil {
			defer gotHashTextBuf.Close()
			gotHashText, err := ioutil.ReadAll(gotHashTextBuf)
			if err == nil && sameHash(wantHashText, gotHashText) {
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
		_, err = io.Copy(cipherTextBuf, crypto.NewTransform(plainTextBuf2, p.decryptor.Encrypt))
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
			gotPlainText, err := ioutil.ReadAll(gotPlainTextBuf)
			if err == nil {
				gotHashText, err := p.hasher.Hash(gotPlainText)
				if err == nil && sameHash(wantHashText, gotHashText) {
					continue
				}
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
		_, err = io.Copy(plainTextBuf, crypto.NewTransform(cipherTextBuf, p.decryptor.Decrypt))
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
