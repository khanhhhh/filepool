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
	gcm cipher.AEAD
}

func (d *aesDecryptor) Decrypt(cipherText []byte) (plainText []byte, err error) {
	nonceSize := d.gcm.NonceSize()
	nonce, ciphertext := cipherText[:nonceSize], cipherText[nonceSize:]
	return d.gcm.Open(nil, nonce, ciphertext, nil)
}

func (d *aesDecryptor) DecryptStream(cipherTextBuf io.Reader, plainTextBuf io.Writer) error {
	nonceSize := d.gcm.NonceSize()
	cipherTextChunk := make([]byte, chunkSize)
	for {
		n, err := cipherTextBuf.Read(cipherTextChunk)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		nonce, cipherText := cipherTextChunk[:nonceSize], cipherTextChunk[nonceSize:n]
		plaintext, err := d.gcm.Open(nil, nonce, cipherText, nil)
		if err != nil {
			return err
		}
		_, err = plainTextBuf.Write(plaintext)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *aesDecryptor) Encrypt(plainText []byte) (cipherText []byte, err error) {
	nonce := make([]byte, d.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return d.gcm.Seal(nonce, nonce, plainText, nil), nil
}

func (d *aesDecryptor) EncryptStream(plainTextBuf io.Reader, cipherTextBuf io.Writer) error {
	plainText := make([]byte, chunkSize)
	for {
		n, err := plainTextBuf.Read(plainText)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		nonce := make([]byte, d.gcm.NonceSize())
		if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
			return err
		}
		cipherTextChunk := d.gcm.Seal(nonce, nonce, plainText[:n], nil)
		_, err = cipherTextBuf.Write(cipherTextChunk)
		if err != nil {
			return err
		}
	}
	return nil
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
