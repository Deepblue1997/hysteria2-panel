package models

import (
	"time"
)

type User struct {
	ID           uint   `gorm:"primarykey"`
	Username     string `gorm:"unique"`
	Password     string
	Email        string    `gorm:"unique"`
	Traffic      int64     `gorm:"default:0"` // 已使用流量
	TrafficLimit int64     `gorm:"default:0"` // 流量限制，0表示不限制
	ExpireAt     time.Time // 账户过期时间
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UserConfig struct {
	UserID    uint `gorm:"primarykey"`
	Port      int  `gorm:"unique"`
	Password  string
	UpSpeed   int // 上行速度限制
	DownSpeed int // 下行速度限制
}
