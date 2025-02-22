package dao

import (
	"GoGateway/conf"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB *gorm.DB
)

func InitDB(dsn string) {
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		hlog.Fatalf("failed to connect database: %v", err)
	}
}

func DefaultDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&&parseTime=True&loc=Local",
		conf.GetConfig().Mysql.Username,
		conf.GetConfig().Mysql.Password,
		conf.GetConfig().Mysql.Host,
		conf.GetConfig().Mysql.Port,
		"gogateway",
	)
}
