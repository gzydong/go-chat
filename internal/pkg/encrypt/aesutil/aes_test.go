package aesutil

import (
	"testing"
)

func TestEncryptCBCWithIV(t *testing.T) {
	plaintext := "这是测试服时间段卡菲纳上课呢罚款就撒阿萨德按色卡"
	key := []byte("1234567890123458")
	iv := []byte("1234567890123456")

	// 加密
	encrypted, err := EncryptStringWithIV(plaintext, key, iv)
	if err != nil {
		t.Fatalf("加密失败: %v", err)
	}
	t.Logf("加密结果: %s", encrypted)

	// 解密
	decrypted, err := DecryptStringWithIV(encrypted, key, iv)
	if err != nil {
		t.Fatalf("解密失败: %v", err)
	}
	t.Logf("解密结果: %s", decrypted)

	// 验证结果
	if decrypted != plaintext {
		t.Errorf("解密结果不匹配，期望: %s, 实际: %s", plaintext, decrypted)
	}
}

func TestEncryptCBCWithIVMultiple(t *testing.T) {
	testCases := []struct {
		name      string
		plaintext string
		key       []byte
		iv        []byte
	}{
		{
			name:      "test",
			plaintext: "asnfaknajksnfjknasjkfnkj",
			key:       []byte("1234567890123458"),
			iv:        []byte("1234567890123456"),
		},
		{
			name:      "中文加密",
			plaintext: "你好，世界",
			key:       []byte("1234567890123458"),
			iv:        []byte("1234567890123456"),
		},
		{
			name:      "特殊字符加密",
			plaintext: "Hello!@#$%^&*()",
			key:       []byte("1234567890123458"),
			iv:        []byte("1234567890123456"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 加密
			encrypted, err := EncryptStringWithIV(tc.plaintext, tc.key, tc.iv)
			if err != nil {
				t.Fatalf("加密失败: %v", err)
			}
			t.Logf("加密结果: %s", encrypted)

			// 解密
			decrypted, err := DecryptStringWithIV(encrypted, tc.key, tc.iv)
			if err != nil {
				t.Fatalf("解密失败: %v", err)
			}
			t.Logf("解密结果: %s", decrypted)

			// 验证结果
			if decrypted != tc.plaintext {
				t.Errorf("解密结果不匹配，期望: %s, 实际: %s", tc.plaintext, decrypted)
			}
		})
	}
}
