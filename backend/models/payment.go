package models

// 支付接口定义
type PaymentProvider interface {
	// 创建支付
	CreatePayment(order *Order) (string, error)
	// 查询支付状态
	QueryPayment(orderNo string) (bool, error)
	// 验证支付回调
	VerifyCallback(params map[string]string) (string, bool)
}

// 支付宝配置
type AlipayConfig struct {
	AppID      string `json:"app_id"`
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
	NotifyURL  string `json:"notify_url"`
	ReturnURL  string `json:"return_url"`
}

// 微信支付配置
type WechatPayConfig struct {
	AppID     string `json:"app_id"`
	MchID     string `json:"mch_id"`
	Key       string `json:"key"`
	NotifyURL string `json:"notify_url"`
}
