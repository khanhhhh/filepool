package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"io/ioutil"
)

type aesDecryptor struct {
	gcm cipher.AEAD
}

func (d *aesDecryptor) Decrypt(cipherText []byte) (plainText []byte, err error) {
	nonceSize := d.gcm.NonceSize()
	nonce, ciphertext := cipherText[:nonceSize], cipherText[nonceSize:]
	return d.gcm.Open(nil, nonce, ciphertext, nil)
}

func (d *aesDecryptor) Encrypt(plainText []byte) (cipherText []byte, err error) {
	nonce := make([]byte, d.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return d.gcm.Seal(nonce, nonce, plainText, nil), nil
}

// NewAESDecryptorToFile :
func NewAESDecryptorToFile(filename string) (Decryptor, error) {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, err
	}
	err := ioutil.WriteFile(filename, key, 0777)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return &aesDecryptor{
		gcm: gcm,
	}, nil
}

// NewAESDecryptorFromFile :
func NewAESDecryptorFromFile(filename string) (Decryptor, error) {
	key, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return &aesDecryptor{
		gcm: gcm,
	}, nil
}
