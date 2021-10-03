package filesystem

import (
	"fmt"
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

func TestFilesystem_Delete(t *testing.T) {
	filesystem := testNewFilesystem()

	assert.Error(t, filesystem.Delete("zifubao.jpeg"))
}

func TestFilesystem_CreateDir(t *testing.T) {
	filesystem := testNewFilesystem()
	assert.NoError(t, filesystem.CreateDir("tmp/test"))
}

func TestFilesystem_DeleteDir(t *testing.T) {
	filesystem := testNewFilesystem()
	assert.NoError(t, filesystem.DeleteDir("tmp/test"))
}

func TestFilesystem_Stat(t *testing.T) {
	filesystem := testNewFilesystem()

	info, err := filesystem.Stat("zifubao.jpeg")

	assert.NoError(t, err)
	fmt.Printf("%#v", info)
}
