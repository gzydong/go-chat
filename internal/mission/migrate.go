package mission

import (
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
	"go-chat/config"
	"gorm.io/gorm"
)

type MigrateProvider struct {
	Config *config.Config
	DB     *gorm.DB
}

func Migrate(_ *cli.Context, app *MigrateProvider) error {
	fmt.Println("数据库初始化中...")

	content, err := os.ReadFile("./doc/sql/go-chat.sql")
	if err != nil {
		fmt.Println("数据库导入失败", err)
		return err
	}

	for _, sql := range strings.Split(string(content), ";;") {
		if len(sql) > 0 {
			app.DB.Exec(strings.TrimSpace(sql))
		}
	}

	return nil
}
