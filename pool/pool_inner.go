package pool

import (
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
		plainText, err := p.client.Read(filename)
		if err != nil {
			continue
		}
		wantHashText := p.hasher.Hash(plainText)
		gotHashText, err := p.server.Read(filename)
		if err == nil && sameHash(wantHashText, gotHashText) {
			continue
		}
		// write
		cipherText, err := p.decryptor.Encrypt(plainText)
		if err != nil {
			continue
		}
		err = p.server.Write(toHashFile(filename), wantHashText)
		err = p.server.Write(filename, cipherText)
	}
}
func (p *pool) Download() {
	for _, filename := range p.server.List() {
		if isHashFile(filename) {
			// skip hash file
			continue
		}
		wantHashText, err := p.server.Read(toHashFile(filename))
		if err != nil {
			continue
		}
		gotPlainText, err := p.client.Read(filename)
		if err == nil {
			gotHashText := p.hasher.Hash(gotPlainText)
			if sameHash(wantHashText, gotHashText) {
				continue
			}
		}
		// write
		cipherText, err := p.server.Read(filename)
		if err != nil {
			continue
		}
		plainText, err := p.decryptor.Decrypt(cipherText)
		if err != nil {
			continue
		}
		err = p.client.Write(filename, plainText)
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
