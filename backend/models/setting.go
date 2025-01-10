package models

import (
	"time"
)

type Setting struct {
	ID        uint   `gorm:"primarykey"`
	Key       string `gorm:"size:50;uniqueIndex"`
	Value     string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// 系统设置的键名常量
const (
	SettingKeyTLS           = "tls_config"     // TLS配置
	SettingKeyEmailSMTP     = "email_smtp"     // 邮件服务配置
	SettingKeyDefaultQuota  = "default_quota"  // 默认流量配额
	SettingKeyDefaultExpire = "default_expire" // 默认过期时间
	SettingKeyAnnouncement  = "announcement"   // 系统公告
	SettingKeyMaintenance   = "maintenance"    // 维护模式
)

// TLS配置结构
type TLSConfig struct {
	AutoCert    bool   `json:"auto_cert"`    // 是否自动申请证书
	Domain      string `json:"domain"`       // 域名
	Email       string `json:"email"`        // 邮箱
	Provider    string `json:"provider"`     // 证书提供商
	DNSProvider string `json:"dns_provider"` // DNS提供商
	DNSToken    string `json:"dns_token"`    // DNS API Token
}

// 邮件服务配置
type SMTPConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	From     string `json:"from"`
}
