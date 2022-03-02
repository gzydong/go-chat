package filesystem

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"go-chat/testutil"
)

func testNewFilesystem() *Filesystem {
	conf := testutil.GetConfig()

	return NewFilesystem(conf)
}

func TestFilesystem_Write(t *testing.T) {
	filesystem := testNewFilesystem()

	_ = filesystem.Default.Write([]byte("www"), "public/test/file")
}

func TestFilesystem_WriteLocal(t *testing.T) {
	filesystem := testNewFilesystem()

	localFile := "/Users/yuandong/www/gowork/go-chat/README.md"

	assert.NoError(t, filesystem.Default.WriteLocal(localFile, "private/README.md"))
}

func TestFilesystem_Copy(t *testing.T) {
	filesystem := testNewFilesystem()

	_ = filesystem.Default.Copy("private/README.md", "private/README2.md")
}

func TestFilesystem_Delete(t *testing.T) {
	filesystem := testNewFilesystem()

	assert.NoError(t, filesystem.Default.Delete("private/README2.md"))
}

func TestFilesystem_CreateDir(t *testing.T) {
	filesystem := testNewFilesystem()
	assert.NoError(t, filesystem.Default.CreateDir("public/tmp/test"))
}

func TestFilesystem_DeleteDir(t *testing.T) {
	filesystem := testNewFilesystem()
	assert.NoError(t, filesystem.Default.DeleteDir("public/tmp/test"))
}

func TestCosFilesystem_PublicUrl(t *testing.T) {
	filesystem := testNewFilesystem()

	t.Log(filesystem.Default.PublicUrl("private/README.md"))
}

func TestFilesystem_PrivateUrl(t *testing.T) {
	filesystem := testNewFilesystem()

	t.Log(filesystem.Default.PrivateUrl("private/README.md", 120))
}

func TestFilesystem_Stat(t *testing.T) {
	filesystem := testNewFilesystem()

	info, err := filesystem.Default.Stat("private/README.md")

	assert.NoError(t, err)
	fmt.Printf("%#v\n", info)
}

func TestCosFilesystem_ReadContent(t *testing.T) {
	filesystem := testNewFilesystem()

	info, err := filesystem.Default.ReadStream("private/tmp/20211218/ba6680b1da03bbae24081f2f4ba09a4e/3-ba6680b1da03bbae24081f2f4ba09a4e.tmp")

	assert.NoError(t, err)
	fmt.Printf("%#v\n", info)
}
