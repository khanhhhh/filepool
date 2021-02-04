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
)

type encryptor struct {
	publ *rsa.PublicKey
}

func (e *encryptor) Encrypt(dataIn []byte) (dataOut []byte, err error) {
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, e.publ, dataIn, defaultLabel)
}

// NewEncryptorFromFile :
func NewEncryptorFromFile(filename string) (Encryptor, error) {
	pemPublData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	publBlock, _ := pem.Decode(pemPublData)
	if publBlock == nil {
		return nil, errors.New("bad key")
	}
	if got, want := publBlock.Type, "PUBLIC KEY"; got != want {
		return nil, errors.New(fmt.Sprintf("unknown key type, got %s, want %s", got, want))
	}
	publ, err := x509.ParsePKCS1PublicKey(publBlock.Bytes)
	if err != nil {
		return nil, err
	}
	return &encryptor{
		publ: publ,
	}, nil
}
