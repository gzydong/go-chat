package other

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
)

type MigrateCommand *cli.Command

func NewMigrateCommand(db *gorm.DB) MigrateCommand {
	return &cli.Command{
		Name:  "migrate",
		Usage: "数据库初始化",
		Action: func(tx *cli.Context) error {
			fmt.Println("数据库初始化中...")

			content, err := ioutil.ReadFile("./doc/sql/go-chat.sql")
			if err != nil {
				fmt.Println("数据库导入失败", err)
			}

			sqls := strings.Split(string(content), ";;")
			for _, sql := range sqls {
				if len(sql) > 0 {
					_ = db.Exec(strings.TrimSpace(sql)).Error
				}
			}

			return nil
		},
	}
}
