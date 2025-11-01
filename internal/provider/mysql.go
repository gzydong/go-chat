package provider

import (
	"fmt"
	"log"
	"time"

	logger2 "github.com/gzydong/go-chat/internal/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/gzydong/go-chat/config"
)

func NewMySQLClient(conf *config.Config) *gorm.DB {
	file := logger2.CreateFileWriter(conf.Log.LogFilePath("slow-sql.log"))

	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: logger.New(log.New(file, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold:             5 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
		}),
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       conf.MySQL.Dsn(),
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), gormConfig)

	if err != nil {
		panic(fmt.Errorf("mysql connect error :%v", err))
	}

	if db.Error != nil {
		panic(fmt.Errorf("database error :%v", err))
	}

	sqlDB, _ := db.DB()

	sqlDB.SetMaxIdleConns(conf.MySQL.MaxIdleConnNum)
	sqlDB.SetMaxOpenConns(conf.MySQL.MaxOpenConnNum)
	sqlDB.SetConnMaxLifetime(time.Duration(conf.MySQL.ConnMaxLifetime) * time.Second)

	return db
}
