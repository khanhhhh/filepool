package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
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
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return &aesDecryptor{
		gcm: gcm,
	}, nil
}
