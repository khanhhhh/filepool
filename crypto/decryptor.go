package crypto

import "io"

// Decryptor :
type Decryptor interface {
	Decrypt(cipherText io.Reader) (plainText io.Reader)
	Encrypt(plainText io.Reader) (cipherText io.Reader)
}
