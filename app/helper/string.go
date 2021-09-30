package helper

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"strings"
	"time"
)

// GetAuthToken 获取登录授权 token
func GetAuthToken(c *gin.Context) string {
	token := c.GetHeader("Authorization")
	token = strings.TrimLeft(token, "Bearer")
	token = strings.TrimSpace(token)

	// Headers 中没有授权信息则读取 url 中的 token
	if len(token) == 0 {
		token = c.DefaultQuery("token", "")
	}

	if len(token) == 0 {
		token = c.DefaultPostForm("token", "")
	}

	return token
}

// GenValidateCode 生成数字验证码
func GenValidateCode(length int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < length; i++ {
		_, _ = fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}

// MtRand 生成指定范围内的随机数
func MtRand(min, max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max-min+1) + min
}

// GetRandomString 生成随机字符串
func GetRandomString(length int) string {
	var result []byte
	bytes := []byte("0123456789abcdefghijklmnopqrstuvwxyz")

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}

	return string(result)
}
