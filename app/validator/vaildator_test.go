package validator

import (
	"fmt"
	"regexp"
	"testing"
)

func TestName(t *testing.T) {
	matched, err := regexp.MatchString("^1[3456789][0-9]{9}$", "18798276809")
	fmt.Println(matched, err)
}
