package database

import (
	"fmt"
	"hysteria2-panel/config"
	"hysteria2-panel/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 自动迁移数据库结构
	if err := db.AutoMigrate(
		&models.User{},
		&models.UserConfig{},
		&models.Node{},
		&models.Setting{},
		&models.Plan{},
		&models.Subscription{},
		&models.Order{},
	); err != nil {
		return nil, err
	}

	return db, nil
}
