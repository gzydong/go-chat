package filesystem

import (
	"github.com/stretchr/testify/assert"
	"go-chat/testutil"
	"os"
	"testing"
)

func testNewFilesystem() *Filesystem {
	conf := testutil.GetConfig()

	return NewFilesystem(conf)
}

func TestFilesystem_Delete(t *testing.T) {
	filesystem := testNewFilesystem()

	assert.Error(t, filesystem.Delete("image/test"))
}

func TestFilesystem_WriteLocal(t *testing.T) {
	filesystem := testNewFilesystem()

	localFile := "/Users/yuandong.rao/www/mytest/go-chat/uploads/image/zifubao.png"

	filesystem.WriteLocal(localFile, "public/images/test/zifubao.png")
}

func TestFilesystem_Write(t *testing.T) {
	filesystem := testNewFilesystem()

	file, err := os.ReadFile("/Users/yuandong.rao/www/mytest/go-chat/README.md")
	if err != nil {
		return
	}

	filesystem.Write(file, "public/images/test/2README.md")
}
