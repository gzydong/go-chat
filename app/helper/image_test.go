package helper

import (
	"fmt"
	"testing"
)

func TestReadImage(t *testing.T) {
	src := "/Users/yuandong.rao/www/mytest/go-chat/9911696_14.jpeg"

	data := ReadImage(src)

	fmt.Println(data["width"], data["height"])
}

func TestGenImageName(t *testing.T) {
	fmt.Println(GenImageName("jpeg", 100, 180))
}
