package rsautil

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
)

var _ IRsa = (*Rsa)(nil)

type IRsa interface {
	Encrypt(data []byte) (string, error)
	Decrypt(ciphertext string) ([]byte, error)
}

type Rsa struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func (r *Rsa) Encrypt(data []byte) (string, error) {
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, r.publicKey, data)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (r *Rsa) Decrypt(ciphertext string) ([]byte, error) {
	decodeString, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, errors.New("base64 decode error")
	}

	decryptedData, err := rsa.DecryptPKCS1v15(rand.Reader, r.privateKey, decodeString)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %v", err)
	}

	return decryptedData, nil
}

func NewRsa(publicKey []byte, privateKey []byte) IRsa {
	rssPublicKey, err := parsePublicKey(publicKey)
	if err != nil {
		panic(err.Error())
	}

	rsaPrivateKey, err := parsePrivateKey(privateKey)
	if err != nil {
		panic(err.Error())
	}

	return &Rsa{
		publicKey:  rssPublicKey,
		privateKey: rsaPrivateKey,
	}
}

func parsePublicKey(publicKey []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return pubInterface.(*rsa.PublicKey), nil
}

func parsePrivateKey(privateKey []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("failed to decode private key")
	}

	privateKeys, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err == nil {
		return privateKeys, nil
	}

	privateKeys2, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return privateKeys, nil
	}

	rsaPrivateKey, ok := privateKeys2.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("private key is not of type RSA")
	}

	return rsaPrivateKey, nil
}
