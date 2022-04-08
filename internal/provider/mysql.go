package provider

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"go-chat/config"
)

type Writer interface {
	Write(p []byte) (n int, err error)
}

func NewMySQLClient(conf *config.Config) *gorm.DB {

	var out io.Writer

	if conf.Debug() {
		out = os.Stdout
	} else {
		out, _ = os.OpenFile(fmt.Sprintf("%s/logs/sql.log", conf.Log.Dir), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	}

	newLogger := logger.New(
		log.New(out, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,   // 慢 SQL 阈值
			LogLevel:      logger.Silent, // Log level
		},
	)

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       conf.MySQL.GetDsn(), // DSN data source name
		DisableDatetimePrecision:  true,                // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,                // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,                // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,               // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名，启用该选项，此时，`Article` 的表名应该是 `it_article`
		},
		Logger: newLogger,
	})

	if err != nil {
		panic(fmt.Errorf("mysql connect error :%v", err))
	}

	if db.Error != nil {
		panic(fmt.Errorf("database error :%v", err))
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if conf.Debug() {
		db = db.Debug()
	}

	return db
}
