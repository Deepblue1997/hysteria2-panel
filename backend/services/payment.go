package services

import (
	"errors"
	"fmt"
	"hysteria2-panel/models"

	"gorm.io/gorm"
)

type PaymentService struct {
	db             *gorm.DB
	settingService *SettingService
	planService    *PlanService
	providers      map[string]models.PaymentProvider
}

func NewPaymentService(db *gorm.DB, settingService *SettingService, planService *PlanService) *PaymentService {
	return &PaymentService{
		db:             db,
		settingService: settingService,
		planService:    planService,
		providers:      make(map[string]models.PaymentProvider),
	}
}

// 注册支付提供商
func (s *PaymentService) RegisterProvider(name string, provider models.PaymentProvider) {
	s.providers[name] = provider
}

// 创建支付
func (s *PaymentService) CreatePayment(orderNo string, method string) (string, error) {
	var order models.Order
	if err := s.db.Where("order_no = ?", orderNo).First(&order).Error; err != nil {
		return "", errors.New("订单不存在")
	}

	if order.PaymentStatus != 0 {
		return "", errors.New("订单状态异常")
	}

	provider, ok := s.providers[method]
	if !ok {
		return "", fmt.Errorf("不支持的支付方式: %s", method)
	}

	// 更新订单支付方式
	if err := s.db.Model(&order).Update("payment_method", method).Error; err != nil {
		return "", err
	}

	return provider.CreatePayment(&order)
}

// 处理支付回调
func (s *PaymentService) HandleCallback(method string, params map[string]string) error {
	provider, ok := s.providers[method]
	if !ok {
		return fmt.Errorf("不支持的支付方式: %s", method)
	}

	orderNo, verified := provider.VerifyCallback(params)
	if !verified {
		return errors.New("回调验证失败")
	}

	return s.planService.HandlePayment(orderNo, method)
}

// 查询支付状态
func (s *PaymentService) QueryPaymentStatus(orderNo string) (bool, error) {
	var order models.Order
	if err := s.db.Where("order_no = ?", orderNo).First(&order).Error; err != nil {
		return false, errors.New("订单不存在")
	}

	if order.PaymentStatus == 1 {
		return true, nil
	}

	provider, ok := s.providers[order.PaymentMethod]
	if !ok {
		return false, fmt.Errorf("不支持的支付方式: %s", order.PaymentMethod)
	}

	return provider.QueryPayment(orderNo)
}
