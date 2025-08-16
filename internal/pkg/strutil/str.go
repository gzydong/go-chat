package strutil

import (
	"fmt"
	"math/rand"
	"path"
	"strings"
	"time"

	"github.com/google/uuid"
)

// GenValidateCode 生成数字验证码
func GenValidateCode(length int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	newRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	var sb strings.Builder
	for i := 0; i < length; i++ {
		_, _ = fmt.Fprintf(&sb, "%d", numeric[newRand.Intn(10)])
	}
	return sb.String()
}

// MtSubstr 字符串截取
func MtSubstr(value string, start, end int) string {

	if start > end {
		return ""
	}

	str := []rune(value)

	if length := len(str); end > length {
		end = length
	}

	return string(str[start:end])
}

// FileSuffix 获取文件后缀名
func FileSuffix(filename string) string {
	return strings.TrimPrefix(path.Ext(filename), ".")
}

func NewMsgId() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

func NewUuid() string {
	return uuid.New().String()
}
