package models

import (
	"time"
)

type Node struct {
	ID          uint      `gorm:"primarykey"`
	Name        string    `gorm:"size:50;not null"`
	Host        string    `gorm:"size:255;not null"`
	Port        int       `gorm:"not null"`
	Status      int       `gorm:"default:0"`                   // 0: 离线, 1: 在线, 2: 维护中
	Type        string    `gorm:"size:20;default:'hysteria2'"` // 节点类型
	TotalUpload int64     `gorm:"default:0"`                   // 总上传流量
	TotalDown   int64     `gorm:"default:0"`                   // 总下载流量
	LastPing    time.Time // 最后在线时间
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type NodeStatus struct {
	CPU          float64 `json:"cpu"`          // CPU 使用率
	Memory       float64 `json:"memory"`       // 内存使用率
	Disk         float64 `json:"disk"`         // 磁盘使用率
	Load         float64 `json:"load"`         // 系统负载
	NetworkIn    int64   `json:"network_in"`   // 网络入站流量
	NetworkOut   int64   `json:"network_out"`  // 网络出站流量
	ActiveUsers  int     `json:"active_users"` // 活跃用户数
	TotalUsers   int     `json:"total_users"`  // 总用户数
	LastReportAt int64   `json:"last_report"`  // 最后报告时间
}
