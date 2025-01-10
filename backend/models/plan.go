package models

import (
	"time"
)

type Plan struct {
	ID           uint    `gorm:"primarykey"`
	Name         string  `gorm:"size:50;not null"`
	Price        float64 `gorm:"not null"`           // 价格
	Duration     int     `gorm:"not null"`           // 有效期（天）
	TrafficLimit int64   `gorm:"not null"`           // 流量限制（字节）
	SpeedLimit   int     `gorm:"default:0"`          // 速度限制（Mbps，0表示不限制）
	DeviceLimit  int     `gorm:"default:0"`          // 设备限制（0表示不限制）
	Status       int     `gorm:"default:1;not null"` // 状态：0-禁用，1-启用
	Description  string  `gorm:"type:text"`          // 套餐描述
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Subscription struct {
	ID        uint      `gorm:"primarykey"`
	UserID    uint      `gorm:"not null;index"`
	PlanID    uint      `gorm:"not null;index"`
	StartAt   time.Time // 开始时间
	EndAt     time.Time // 结束时间
	Status    int       `gorm:"default:1;not null"` // 状态：0-已取消，1-生效中
	CreatedAt time.Time
	UpdatedAt time.Time
	Plan      Plan `gorm:"foreignKey:PlanID"`
}

type Order struct {
	ID            uint      `gorm:"primarykey"`
	UserID        uint      `gorm:"not null;index"`
	PlanID        uint      `gorm:"not null;index"`
	OrderNo       string    `gorm:"size:50;uniqueIndex"` // 订单号
	Amount        float64   `gorm:"not null"`            // 订单金额
	PaymentMethod string    `gorm:"size:20"`             // 支付方式
	PaymentStatus int       `gorm:"default:0"`           // 支付状态：0-未支付，1-已支付，2-已取消
	PayAt         time.Time // 支付时间
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
