package pool

import (
	"github.com/khanhhhh/filepool/crypto"
	"github.com/khanhhhh/filepool/storage"
)

type pool struct {
	client    storage.Storage
	server    storage.Storage
	decryptor crypto.Decryptor
	hasher    crypto.Hasher
}

func getHashFile(path string) string {
	return path + ".hash"
}

func isHashFile(path string) bool {
	return path[len(path)-5:] == ".hash"
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

func (p *pool) sync(fromStorage storage.Storage, toStorage storage.Storage) {
	toStoragePathSet := make(map[string]struct{})
	for _, path := range toStorage.List() {
		toStoragePathSet[path] = struct{}{}
	}
	for _, path := range fromStorage.List() {
		if isHashFile(path) {
			continue
		}
		var data []byte = nil
		wantHash, err := fromStorage.Read(getHashFile(path))
		if err != nil {
			// write hash to from storage
			data, err = fromStorage.Read(path)
			if err != nil {
				// cannot read file, give up
				continue
			}
			wantHash = p.hasher.Hash(data)
			fromStorage.Write(getHashFile(path), wantHash)
		}
		gotHash, err := toStorage.Read(path)
		if err == nil && sameHash(gotHash, wantHash) {
			continue
		}
		if data == nil {
			// load data
			data, err = fromStorage.Read(path)
			if err != nil {
				continue
			}
		}
		// write
		err = toStorage.Write(path, data)
		// skip error
	}
}
