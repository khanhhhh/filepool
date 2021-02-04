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

var byte2hex = map[[4]byte]string{
	{0, 0, 0, 0}: "0",
	{0, 0, 0, 1}: "1",
	{0, 0, 1, 0}: "2",
	{0, 0, 1, 1}: "3",
	{0, 1, 0, 0}: "4",
	{0, 1, 0, 1}: "5",
	{0, 1, 1, 0}: "6",
	{0, 1, 1, 1}: "7",
	{1, 0, 0, 0}: "8",
	{1, 0, 0, 1}: "9",
	{1, 0, 1, 0}: "A",
	{1, 0, 1, 1}: "B",
	{1, 1, 0, 0}: "C",
	{1, 1, 0, 1}: "D",
	{1, 1, 1, 0}: "E",
	{1, 1, 1, 1}: "F",
}

var hex2byte = map[string][4]byte{
	"0": {0, 0, 0, 0},
	"1": {0, 0, 0, 1},
	"2": {0, 0, 1, 0},
	"3": {0, 0, 1, 1},
	"4": {0, 1, 0, 0},
	"5": {0, 1, 0, 1},
	"6": {0, 1, 1, 0},
	"7": {0, 1, 1, 1},
	"8": {1, 0, 0, 0},
	"9": {1, 0, 0, 1},
	"A": {1, 0, 1, 0},
	"B": {1, 0, 1, 1},
	"C": {1, 1, 0, 0},
	"D": {1, 1, 0, 1},
	"E": {1, 1, 1, 0},
	"F": {1, 1, 1, 1},
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
