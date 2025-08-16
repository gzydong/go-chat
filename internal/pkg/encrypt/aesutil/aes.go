package aesutil

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

// 定义常见错误
var (
	ErrInvalidPKCS7Data = errors.New("aes: 无效的PKCS7数据")
	ErrInvalidIVSize    = errors.New("aes: 无效的IV大小")
)

// PKCS7Padding 添加PKCS7填充
func PKCS7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// PKCS7UnPadding 移除PKCS7填充
func PKCS7UnPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, ErrInvalidPKCS7Data
	}

	padding := int(data[length-1])
	if padding > length {
		return nil, ErrInvalidPKCS7Data
	}

	return data[:length-padding], nil
}

// EncryptCBCWithIV 使用CBC模式和指定的IV加密数据
func EncryptCBCWithIV(plaintext []byte, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 验证IV大小
	if len(iv) != block.BlockSize() {
		return nil, ErrInvalidIVSize
	}

	// 创建填充后的数据
	plaintext = PKCS7Padding(plaintext, block.BlockSize())

	// 创建密文（不包含IV，因为IV是单独传入的）
	ciphertext := make([]byte, len(plaintext))

	// 加密
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plaintext)

	return ciphertext, nil
}

// DecryptCBCWithIV 使用CBC模式和指定的IV解密数据
func DecryptCBCWithIV(ciphertext []byte, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 验证IV大小
	if len(iv) != block.BlockSize() {
		return nil, ErrInvalidIVSize
	}

	// 创建解密后的明文
	plaintext := make([]byte, len(ciphertext))
	copy(plaintext, ciphertext)

	// 解密
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, plaintext)

	// 移除填充
	return PKCS7UnPadding(plaintext)
}

// EncryptStringWithIV 加密字符串（使用CBC模式和指定的IV，并返回base64编码）
func EncryptStringWithIV(plaintext string, key []byte, iv []byte) (string, error) {
	ciphertext, err := EncryptCBCWithIV([]byte(plaintext), key, iv)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptStringWithIV 解密字符串（使用CBC模式和指定的IV，并处理base64编码）
func DecryptStringWithIV(ciphertext string, key []byte, iv []byte) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	plaintext, err := DecryptCBCWithIV(data, key, iv)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
