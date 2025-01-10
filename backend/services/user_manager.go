package services

import (
	"errors"
	"hysteria2-panel/models"
	"strconv"

	"gorm.io/gorm"
)

type UserManagerService struct {
	db *gorm.DB
}

func NewUserManagerService(db *gorm.DB) *UserManagerService {
	return &UserManagerService{db: db}
}

func (s *UserManagerService) GetUsers(page, pageSize int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	if err := s.db.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := s.db.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	// 清除密码字段
	for i := range users {
		users[i].Password = ""
	}

	return users, total, nil
}

func (s *UserManagerService) GetUserByID(id string) (*models.User, error) {
	var user models.User
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, errors.New("无效的用户ID")
	}

	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	user.Password = "" // 清除密码字段
	return &user, nil
}

func (s *UserManagerService) UpdateUser(id string, updates map[string]interface{}) error {
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return errors.New("无效的用户ID")
	}

	// 不允许更新密码和用户名
	delete(updates, "password")
	delete(updates, "username")

	result := s.db.Model(&models.User{}).Where("id = ?", userID).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("用户不存在")
	}

	return nil
}

func (s *UserManagerService) DeleteUser(id string) error {
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return errors.New("无效的用户ID")
	}

	// 开始事务
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 删除用户配置
		if err := tx.Where("user_id = ?", userID).Delete(&models.UserConfig{}).Error; err != nil {
			return err
		}

		// 删除用户
		result := tx.Delete(&models.User{}, userID)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("用户不存在")
		}

		return nil
	})
}
