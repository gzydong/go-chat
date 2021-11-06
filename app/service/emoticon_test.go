package service

import (
	"fmt"
	"strings"
	"testing"
)

func TestJoin(t *testing.T) {
	items := []string{"111", "222", "333"}

	fmt.Println(strings.Join(items, ","))
}
