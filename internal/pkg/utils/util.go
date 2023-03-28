package utils

import (
	"bytes"
	"fmt"
	"math/rand"
	"runtime"
)

// MtRand 生成指定范围内的随机数
func MtRand(min, max int) int {
	r := rand.New(rand.NewSource(1))
	return r.Intn(max-min+1) + min
}

func PanicTrace(err interface{}) string {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "%v\n", err)
	for i := 2; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
	}
	return buf.String()
}
