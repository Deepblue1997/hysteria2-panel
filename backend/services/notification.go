package services

import (
	"fmt"
	"hysteria2-panel/models"
	"time"

	"gorm.io/gorm"
)

type NotificationService struct {
	db          *gorm.DB
	mailService *MailService
}

func NewNotificationService(db *gorm.DB, mailService *MailService) *NotificationService {
	return &NotificationService{
		db:          db,
		mailService: mailService,
	}
}

// 发送订阅到期提醒
func (s *NotificationService) SendExpirationNotices() error {
	var users []models.User
	// 查找3天内即将到期的用户
	err := s.db.Where("expire_at BETWEEN ? AND ?",
		time.Now(),
		time.Now().AddDate(0, 0, 3)).
		Find(&users).Error
	if err != nil {
		return err
	}

	for _, user := range users {
		daysLeft := int(time.Until(user.ExpireAt).Hours() / 24)
		if err := s.mailService.SendExpirationNotice(user.Email, daysLeft); err != nil {
			fmt.Printf("发送到期提醒失败，用户ID: %d, 错误: %v\n", user.ID, err)
		}
	}

	return nil
}

// 发送流量使用提醒
func (s *NotificationService) SendTrafficNotices() error {
	var users []models.User
	// 查找流量使用超过80%的用户
	err := s.db.Where("traffic_limit > 0 AND traffic >= traffic_limit * 0.8").
		Find(&users).Error
	if err != nil {
		return err
	}

	for _, user := range users {
		usagePercent := float64(user.Traffic) / float64(user.TrafficLimit) * 100
		subject := "流量使用提醒"
		body := fmt.Sprintf(`
			<h3>流量使用提醒</h3>
			<p>您的流量使用已达到 %.1f%%。</p>
			<p>为避免服务中断，请及时续费。</p>
		`, usagePercent)

		if err := s.mailService.SendMail([]string{user.Email}, subject, body); err != nil {
			fmt.Printf("发送流量提醒失败，用户ID: %d, 错误: %v\n", user.ID, err)
		}
	}

	return nil
}
