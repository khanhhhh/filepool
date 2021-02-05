package crypto

import (
	"io"
)

type plainDecryptor struct {
}

func (d *plainDecryptor) Decrypt(cipherText io.Reader) (plainText io.Reader) {
	return cipherText
}

func (d *plainDecryptor) Encrypt(plainText io.Reader) (cipherText io.Reader) {
	return plainText
}

// NewPlainDecryptor :
func NewPlainDecryptor() Decryptor {
	return &plainDecryptor{}
}
