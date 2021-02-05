package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
	"io/ioutil"
)

const (
	// each of 256MB of file is encrypted and padded with a nonce (32B) at the beginning
	chunkSize = 256 * 1024 * 1024 // 256MB
)

type aesDecryptor struct {
	decryptor cipher.Stream
	encryptor cipher.Stream
}

func (d *aesDecryptor) Decrypt(cipherText io.Reader) (plainText io.Reader) {
	return &cipher.StreamReader{
		S: d.decryptor,
		R: cipherText,
	}
}

func (d *aesDecryptor) Encrypt(plainText io.Reader) (cipherText io.Reader) {
	return &cipher.StreamReader{
		S: d.encryptor,
		R: plainText,
	}
}

// NewAESDecryptor :
func NewAESDecryptor(filename string) (Decryptor, error) {
	keyStr, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	key, err := hex.DecodeString(string(keyStr))
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	iv := make([]byte, block.BlockSize())
	return &aesDecryptor{
		decryptor: cipher.NewCFBDecrypter(block, iv),
		encryptor: cipher.NewCFBEncrypter(block, iv),
	}, nil
}

// NewAESKey :
func NewAESKey(filename string) error {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return err
	}
	if err := ioutil.WriteFile(filename, []byte(hex.EncodeToString(key)), 0777); err != nil {
		return err
	}
	return nil
}
