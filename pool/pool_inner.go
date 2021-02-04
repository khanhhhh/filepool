package pool

import (
	"github.com/khanhhhh/filepool/crypto"
	"github.com/khanhhhh/filepool/storage"
)

type pool struct {
	decryptor  crypto.Decryptor
	server     storage.Storage
	clientList []storage.Storage
}

func (p *pool) Upload() {
	for _, client := range p.clientList {
		sync(p.decryptor, client, []storage.Storage{p.server})
	}
}
func (p *pool) Download() {
	sync(p.decryptor, p.server, p.clientList)
}

func sync(decryptor crypto.Decryptor, fromStorage storage.Storage, toStorageList []storage.Storage) {
	for _, filename := range fromStorage.List() {
		fromStat, err := fromStorage.Stat(filename)
		if err != nil {
			// skip if read error
			continue
		}
		// otherwise, write
		cipherText, err := fromStorage.Read(filename)
		if err != nil {
			// skip if read error
			continue
		}
		plainText, err := decryptor.Decrypt(cipherText)
		if err != nil {
			// decrypt error
			continue
		}
		for _, toStorage := range toStorageList {
			toStat, err := toStorage.Stat(filename)
			if err == nil && toStat.ModTime == fromStat.ModTime {
				// skip if same mod time
				continue
			}
			err = toStorage.Write(filename, plainText)
			if err != nil {
				// write error
				continue
			}
		}
	}
}

// NewPool :
func NewPool(decryptor crypto.Decryptor, server storage.Storage, clientList []storage.Storage) Pool {
	return &pool{
		decryptor:  decryptor,
		server:     server,
		clientList: clientList,
	}
}
