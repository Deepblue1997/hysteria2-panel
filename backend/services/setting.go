package services

import (
	"encoding/json"
	"errors"
	"hysteria2-panel/models"

	"gorm.io/gorm"
)

type SettingService struct {
	db *gorm.DB
}

func NewSettingService(db *gorm.DB) *SettingService {
	return &SettingService{db: db}
}

// 获取设置值
func (s *SettingService) GetSetting(key string) (*models.Setting, error) {
	var setting models.Setting
	err := s.db.Where("key = ?", key).First(&setting).Error
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

// 更新设置
func (s *SettingService) UpdateSetting(key string, value interface{}) error {
	// 将value转换为JSON字符串
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	result := s.db.Where("key = ?", key).
		Assign(models.Setting{Value: string(jsonValue)}).
		FirstOrCreate(&models.Setting{Key: key})

	return result.Error
}

// 获取TLS配置
func (s *SettingService) GetTLSConfig() (*models.TLSConfig, error) {
	setting, err := s.GetSetting(models.SettingKeyTLS)
	if err != nil {
		return nil, err
	}

	var config models.TLSConfig
	if err := json.Unmarshal([]byte(setting.Value), &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// 更新TLS配置
func (s *SettingService) UpdateTLSConfig(config *models.TLSConfig) error {
	return s.UpdateSetting(models.SettingKeyTLS, config)
}

// 获取SMTP配置
func (s *SettingService) GetSMTPConfig() (*models.SMTPConfig, error) {
	setting, err := s.GetSetting(models.SettingKeyEmailSMTP)
	if err != nil {
		return nil, err
	}

	var config models.SMTPConfig
	if err := json.Unmarshal([]byte(setting.Value), &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// 更新SMTP配置
func (s *SettingService) UpdateSMTPConfig(config *models.SMTPConfig) error {
	return s.UpdateSetting(models.SettingKeyEmailSMTP, config)
}

// 获取系统公告
func (s *SettingService) GetAnnouncement() (string, error) {
	setting, err := s.GetSetting(models.SettingKeyAnnouncement)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil
		}
		return "", err
	}
	return setting.Value, nil
}

// 更新系统公告
func (s *SettingService) UpdateAnnouncement(content string) error {
	return s.UpdateSetting(models.SettingKeyAnnouncement, content)
}
