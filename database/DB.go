package database

import (
	"Campus-errand-platform/config"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(cfg config.Config) *gorm.DB {
	db, err := gorm.Open(mysql.Open(cfg.DBDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	// 连接池配置
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("获取底层连接失败: %v", err)
	}
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db
}
