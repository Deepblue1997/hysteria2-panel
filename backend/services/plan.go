package services

import (
	"errors"
	"fmt"
	"hysteria2-panel/models"
	"time"

	"gorm.io/gorm"
)

type PlanService struct {
	db *gorm.DB
}

func NewPlanService(db *gorm.DB) *PlanService {
	return &PlanService{db: db}
}

// 创建套餐
func (s *PlanService) CreatePlan(plan *models.Plan) error {
	return s.db.Create(plan).Error
}

// 获取套餐列表
func (s *PlanService) GetPlans() ([]models.Plan, error) {
	var plans []models.Plan
	err := s.db.Find(&plans).Error
	return plans, err
}

// 更新套餐
func (s *PlanService) UpdatePlan(id uint, updates map[string]interface{}) error {
	result := s.db.Model(&models.Plan{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("套餐不存在")
	}
	return nil
}

// 订阅套餐
func (s *PlanService) Subscribe(userID, planID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 获取套餐信息
		var plan models.Plan
		if err := tx.First(&plan, planID).Error; err != nil {
			return errors.New("套餐不存在")
		}

		// 检查是否有正在生效的订阅
		var count int64
		if err := tx.Model(&models.Subscription{}).
			Where("user_id = ? AND status = 1 AND end_at > ?", userID, time.Now()).
			Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("已有正在生效的订阅")
		}

		// 创建订阅
		subscription := &models.Subscription{
			UserID:  userID,
			PlanID:  planID,
			StartAt: time.Now(),
			EndAt:   time.Now().AddDate(0, 0, plan.Duration),
			Status:  1,
		}

		if err := tx.Create(subscription).Error; err != nil {
			return err
		}

		// 更新用户流量和到期时间
		updates := map[string]interface{}{
			"traffic_limit": plan.TrafficLimit,
			"expire_at":     subscription.EndAt,
		}
		if err := tx.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
			return err
		}

		return nil
	})
}

// 创建订单
func (s *PlanService) CreateOrder(userID, planID uint) (*models.Order, error) {
	var plan models.Plan
	if err := s.db.First(&plan, planID).Error; err != nil {
		return nil, errors.New("套餐不存在")
	}

	order := &models.Order{
		UserID:        userID,
		PlanID:        planID,
		OrderNo:       fmt.Sprintf("%d%d%d", userID, planID, time.Now().Unix()),
		Amount:        plan.Price,
		PaymentStatus: 0,
	}

	if err := s.db.Create(order).Error; err != nil {
		return nil, err
	}

	return order, nil
}

// 处理支付
func (s *PlanService) HandlePayment(orderNo string, method string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var order models.Order
		if err := tx.Where("order_no = ?", orderNo).First(&order).Error; err != nil {
			return errors.New("订单不存在")
		}

		if order.PaymentStatus != 0 {
			return errors.New("订单状态异常")
		}

		// 更新订单状态
		updates := map[string]interface{}{
			"payment_status": 1,
			"payment_method": method,
			"pay_at":         time.Now(),
		}
		if err := tx.Model(&order).Updates(updates).Error; err != nil {
			return err
		}

		// 创建订阅
		return s.Subscribe(order.UserID, order.PlanID)
	})
}
