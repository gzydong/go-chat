package mission

import (
	"embed"
	"fmt"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
	"go-chat/config"
	"gorm.io/gorm"
)

//go:embed sql.sql
var file embed.FS

type MigrateProvider struct {
	Config *config.Config
	DB     *gorm.DB
}

func Migrate(_ *cli.Context, app *MigrateProvider) error {
	fmt.Println("数据库初始化中...")

	content, err := file.ReadFile("sql.sql")
	if err != nil {
		fmt.Println("读取数据库初始化文件失败 Err:", err.Error())
		return err
	}

	for _, sql := range strings.Split(string(content), ";;") {
		if len(sql) > 0 {
			err = app.DB.Exec(strings.TrimSpace(sql)).Error
			if err != nil {
				fmt.Println("执行SQL:", strings.TrimSpace(sql), " Err:", err)
				time.Sleep(5 * time.Second)
			}
		}
	}

	return nil
}
