package encrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
)

func RSAEncrypt(data, pemPubKey []byte) (string, error) {
	pk, err := parsePublicKey(pemPubKey)
	if err != nil {
		return "", err
	}

	encryptedData, err := rsa.EncryptPKCS1v15(rand.Reader, pk, data)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encryptedData), nil
}

func RSADecrypt(encryptedData string, pemPriKey []byte) (string, error) {
	encryptedDecodeBytes, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", err
	}

	pk, err := parsePrivateKey(pemPriKey)
	if err != nil {
		return "", err
	}

	originalData, err := rsa.DecryptPKCS1v15(rand.Reader, pk, encryptedDecodeBytes)
	if err != nil {
		return "", err
	}
	return string(originalData), err
}

func parsePublicKey(publicKey []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("publicKey format error")
	}

	var pubInterface interface{}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	if pub, ok := pubInterface.(*rsa.PublicKey); ok {
		return pub, nil
	}

	return nil, errors.New("publicKey error")
}

func parsePrivateKey(privateKey []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("privateKey format error")
	}

	switch block.Type {
	case "RSA PRIVATE KEY", "PRIVATE KEY":
		rsaPrivateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		return rsaPrivateKey, nil
	default:
		return nil, errors.New("privateKey error")
	}
}
