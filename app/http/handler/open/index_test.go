package open

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"go-chat/testutil"
)

func TestIndex_Index(t *testing.T) {
	index := Index{}
	req := testutil.NewTestRequest("/open", index.Index)
	resp, err := req.Get()
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)
}
