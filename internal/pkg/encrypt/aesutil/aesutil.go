package aesutil

import (
	"crypto/md5"
	"encoding/hex"
)

type IAesUtil interface {
	Encrypt(plaintext string) (string, error)
	Decrypt(ciphertext string) (string, error)
}

var _ IAesUtil = (*AesUtil)(nil)

type AesUtil struct {
	key []byte
	iv  []byte
}

func NewAesUtil(key string) IAesUtil {
	h := md5.New()
	h.Write([]byte(key))

	return &AesUtil{
		key: []byte(key),
		iv:  []byte(hex.EncodeToString(h.Sum(nil))[:16]),
	}
}

func (a *AesUtil) Encrypt(plaintext string) (string, error) {
	return EncryptStringWithIV(plaintext, a.key, a.iv)
}

func (a *AesUtil) Decrypt(ciphertext string) (string, error) {
	return DecryptStringWithIV(ciphertext, a.key, a.iv)
}

func (a *AesUtil) EncryptByte(data []byte) ([]byte, error) {
	return EncryptCBCWithIV(data, a.key, a.iv)
}

func (a *AesUtil) DecryptByte(data []byte) ([]byte, error) {
	return DecryptCBCWithIV(data, a.key, a.iv)
}
