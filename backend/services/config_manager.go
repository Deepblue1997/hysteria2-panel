package services

import (
	"encoding/json"
	"hysteria2-panel/models"
	"time"

	"gorm.io/gorm"
)

type ConfigManagerService struct {
	db *gorm.DB
}

func NewConfigManagerService(db *gorm.DB) *ConfigManagerService {
	return &ConfigManagerService{db: db}
}

// 获取用户配置
func (s *ConfigManagerService) GetUserConfig(userID uint) (*models.UserConfig, error) {
	var config models.UserConfig
	err := s.db.Where("user_id = ?", userID).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// 更新用户配置
func (s *ConfigManagerService) UpdateUserConfig(userID uint, config *models.UserConfig) error {
	config.UserID = userID
	config.UpdatedAt = time.Now()
	return s.db.Save(config).Error
}

// 生成客户端配置
func (s *ConfigManagerService) GenerateClientConfig(userID uint) (string, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return "", err
	}

	config := map[string]interface{}{
		"server":     "your-server.com",
		"protocol":   "hysteria2",
		"up_mbps":    100,
		"down_mbps":  100,
		"auth":       user.Auth,
		"server_key": "your-server-key",
	}

	jsonBytes, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}
