package im

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSender(t *testing.T) {
	obj := NewSender("default")

	fmt.Printf("%#v\n", obj)
	fmt.Printf("%#v\n", Manager)

	obj.Send()

	assert.Equal(t, true, true)
}
