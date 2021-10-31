package service

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {

	list := []string{"那就开始拿", "na那三级卡你发能进沙发", "问你啊咖金色发"}

	options := make(map[string]string)
	for i, value := range list {
		options[fmt.Sprintf("%c", 65+i)] = value
	}

	fmt.Printf("%#v", options)
}
