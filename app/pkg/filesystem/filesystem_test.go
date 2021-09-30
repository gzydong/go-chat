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

func TestFilesystem_Write(t *testing.T) {
	filesystem := testNewFilesystem()

	file, err := os.ReadFile("/Users/yuandong/www/gowork/go-chat/test.jpeg")
	if err != nil {
		return
	}

	_ = filesystem.Write(file, "/images/20201025/test.jpeg")
}

func TestFilesystem_WriteLocal(t *testing.T) {
	filesystem := testNewFilesystem()

	localFile := "/Users/yuandong/www/gowork/go-chat/test.jpeg"

	assert.NoError(t, filesystem.WriteLocal(localFile, "zifubao.jpeg"))
}

func TestFilesystem_Copy(t *testing.T) {
	filesystem := testNewFilesystem()

	_ = filesystem.Copy("public/images/test/2README.md", "public/images/test/6README.md")
}

func TestFilesystem_Stat(t *testing.T) {
	filesystem := testNewFilesystem()

	filesystem.Stat("zifubao.jpeg")
}

func TestFilesystem_Delete(t *testing.T) {
	filesystem := testNewFilesystem()

	assert.Error(t, filesystem.Delete("zifubao.jpeg"))
}
