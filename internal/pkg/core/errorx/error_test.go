package errorx

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsErrorX(t *testing.T) {
	assert.Equal(t, true, IsError(New(403, "test")))
	assert.Equal(t, false, IsError(errors.New("test")))
}
