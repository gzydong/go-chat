package testutil

import (
	"fmt"
	"go-chat/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

func GetDb() *gorm.DB {

	conf := GetConfig()

	dsn := getDsn(conf)

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   conf.MySQL.Prefix, // 表名前缀，`Article` 的表名应该是 `it_articles`
			SingularTable: true,              // 使用单数表名，启用该选项，此时，`Article` 的表名应该是 `it_article`
		},
	})

	if err != nil {
		fmt.Printf("mysql connect error :%v", err)
	}

	if db.Error != nil {
		fmt.Printf("database error :%v", db.Error)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	db.Debug()

	return db
}

func getDsn(conf *config.Config) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		conf.MySQL.UserName,
		conf.MySQL.Password,
		conf.MySQL.Host,
		conf.MySQL.Port,
		conf.MySQL.Database,
		conf.MySQL.Charset,
	)
}
