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
	sync(p.decryptor.Encrypt, p.hasher.Hash, p.client, p.server)
}
func (p *pool) Download() {
	sync(p.decryptor.Decrypt, p.hasher.Hash, p.server, p.client)
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

func sync(
	transform func(dataIn []byte) (dataOut []byte, err error),
	hash func(dataIn []byte) (dataOut []byte),
	fromStorage storage.Storage,
	toStorage storage.Storage,
) {
	for _, filename := range fromStorage.List() {
		if isHashFile(filename) {
			// skip hash file
			continue
		}
		fromText, err := fromStorage.Read(filename)
		if err != nil {
			// skip if read error
			continue
		}
		toText, err := transform(fromText)
		if err != nil {
			// decrypt error
			continue
		}
		wantHashText := hash(toText)
		if storedWantHashText, err := fromStorage.Read(toHashFile(filename)); err != nil || !sameHash(wantHashText, storedWantHashText) {
			err = fromStorage.Write(toHashFile(filename), wantHashText)
			if err != nil {
				continue
			}
		}
		if gotHashText, err := toStorage.Read(toHashFile(filename)); err != nil || !sameHash(wantHashText, gotHashText) {
			err = toStorage.Write(toHashFile(filename), wantHashText)
			if err != nil {
				// write error
				continue
			}
			err = toStorage.Write(filename, toText)
			if err != nil {
				// write error
				continue
			}
		}
	}
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
