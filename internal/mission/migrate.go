package mission

import (
	"embed"
	"fmt"
	"strings"
	"time"

	"github.com/gzydong/go-chat/config"
	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
)

//go:embed resource/lumenim.sql
var file embed.FS

type MigrateProvider struct {
	Config *config.Config
	DB     *gorm.DB
}

func Migrate(_ *cli.Context, app *MigrateProvider) error {
	fmt.Println("数据库初始化中...")
	defer fmt.Println("数据库初始化完成")

	content, err := file.ReadFile("resource/lumenim.sql")
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
