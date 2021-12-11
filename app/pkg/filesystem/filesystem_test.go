package filesystem

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"go-chat/testutil"
	"testing"
)

func testNewFilesystem() *Filesystem {
	conf := testutil.GetConfig()

	return NewFilesystem(conf)
}

func TestFilesystem_Write(t *testing.T) {
	filesystem := testNewFilesystem()

	_ = filesystem.Write([]byte("www"), "public/test/file")
}

func TestFilesystem_WriteLocal(t *testing.T) {
	filesystem := testNewFilesystem()

	localFile := "/Users/yuandong/www/gowork/go-chat/README.md"

	assert.NoError(t, filesystem.WriteLocal(localFile, "private/README.md"))
}

func TestFilesystem_Copy(t *testing.T) {
	filesystem := testNewFilesystem()

	_ = filesystem.Copy("private/README.md", "private/README2.md")
}

func TestFilesystem_Delete(t *testing.T) {
	filesystem := testNewFilesystem()

	assert.NoError(t, filesystem.Delete("private/README2.md"))
}

func TestFilesystem_CreateDir(t *testing.T) {
	filesystem := testNewFilesystem()
	assert.NoError(t, filesystem.CreateDir("public/tmp/test"))
}

func TestFilesystem_DeleteDir(t *testing.T) {
	filesystem := testNewFilesystem()
	assert.NoError(t, filesystem.DeleteDir("public/tmp/test"))
}

func TestCosFilesystem_PublicUrl(t *testing.T) {
	filesystem := testNewFilesystem()

	t.Log(filesystem.PublicUrl("private/README.md"))
}

func TestFilesystem_PrivateUrl(t *testing.T) {
	filesystem := testNewFilesystem()

	t.Log(filesystem.PrivateUrl("private/README.md", 120))
}

func TestFilesystem_Stat(t *testing.T) {
	filesystem := testNewFilesystem()

	info, err := filesystem.Stat("private/README.md")

	assert.NoError(t, err)
	fmt.Printf("%#v\n", info)
}
