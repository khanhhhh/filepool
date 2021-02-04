package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

var (
	defaultLabel = []byte("label")
)

type rsaDecryptor struct {
	priv *rsa.PrivateKey
}

func (d *rsaDecryptor) Decrypt(cipherText []byte) (plainText []byte, err error) {
	return rsa.DecryptOAEP(sha256.New(), rand.Reader, d.priv, cipherText, defaultLabel)
}

func (d *rsaDecryptor) Encrypt(plainText []byte) (cipherText []byte, err error) {
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, &d.priv.PublicKey, plainText, defaultLabel)
}

// NewRSAKey :
func NewRSAKey(filename string) error {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}
	pemPrivData := x509.MarshalPKCS1PrivateKey(priv)
	privBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: pemPrivData,
	}
	privPem, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer privPem.Close()
	err = pem.Encode(privPem, privBlock)
	if err != nil {
		return err
	}
	pemPublData := x509.MarshalPKCS1PublicKey(&priv.PublicKey)
	publBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pemPublData,
	}
	publPem, err := os.Create(filename + ".pub")
	if err != nil {
		return err
	}
	defer publPem.Close()
	err = pem.Encode(publPem, publBlock)
	if err != nil {
		return err
	}
	return nil
}

// NewRSADecryptor :
func NewRSADecryptor(filename string) (Decryptor, error) {
	pemPrivData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	privBlock, _ := pem.Decode(pemPrivData)
	if privBlock == nil {
		return nil, errors.New("bad key")
	}
	if got, want := privBlock.Type, "RSA PRIVATE KEY"; got != want {
		return nil, errors.New(fmt.Sprintf("unknown key type, got %s, want %s", got, want))
	}
	priv, err := x509.ParsePKCS1PrivateKey(privBlock.Bytes)
	if err != nil {
		return nil, err
	}
	return &rsaDecryptor{
		priv: priv,
	}, nil
}
