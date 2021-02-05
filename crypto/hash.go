package crypto

import (
	"crypto/sha256"
	"io"
)

// Hasher :
type Hasher interface {
	Hash(dataIn []byte) (dataOut []byte, err error)
	HashStream(readerIn io.Reader) (dataOut []byte, err error)
}

type defaultHasher struct{}

func (h *defaultHasher) Hash(dataIn []byte) (dataOut []byte, err error) {
	hashed := sha256.Sum256(dataIn)
	return hashed[:], nil
}

func (h *defaultHasher) HashStream(readerIn io.Reader) ([]byte, error) {
	hash := sha256.New()
	_, err := io.Copy(hash, readerIn)
	if err != nil {
		return nil, err
	}
	return hash.Sum(nil), nil
}

// NewHasher :
func NewHasher() Hasher {
	return &defaultHasher{}
}
