package crypto

import "io"

// Decryptor :
type Decryptor interface {
	Decrypt(cipherText []byte) (plainText []byte, err error)
	DecryptStream(cipherText io.Reader, plainText io.Writer) error
	Encrypt(plainText []byte) (cipherText []byte, err error)
	EncryptStream(plainText io.Reader, cipherText io.Writer) error
}
