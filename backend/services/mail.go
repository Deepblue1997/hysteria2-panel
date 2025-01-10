package services

import (
	"fmt"
	"net/smtp"
	"strings"
)

type MailService struct {
	settingService *SettingService
}

func NewMailService(settingService *SettingService) *MailService {
	return &MailService{settingService: settingService}
}

// 发送邮件
func (s *MailService) SendMail(to []string, subject, body string) error {
	config, err := s.settingService.GetSMTPConfig()
	if err != nil {
		return fmt.Errorf("获取SMTP配置失败: %v", err)
	}

	auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)

	// 构建邮件头
	headers := make(map[string]string)
	headers["From"] = config.From
	headers["To"] = strings.Join(to, ",")
	headers["Subject"] = subject
	headers["Content-Type"] = "text/html; charset=UTF-8"

	// 构建邮件内容
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// 发送邮件
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	return smtp.SendMail(addr, auth, config.From, to, []byte(message))
}

// 发送注册验证邮件
func (s *MailService) SendRegisterVerification(email, code string) error {
	subject := "注册验证码"
	body := fmt.Sprintf(`
		<h3>欢迎注册</h3>
		<p>您的验证码是: <strong>%s</strong></p>
		<p>验证码有效期为10分钟，请尽快完成注册。</p>
	`, code)

	return s.SendMail([]string{email}, subject, body)
}

// 发送重置密码邮件
func (s *MailService) SendPasswordReset(email, code string) error {
	subject := "重置密码"
	body := fmt.Sprintf(`
		<h3>重置密码</h3>
		<p>您的验证码是: <strong>%s</strong></p>
		<p>验证码有效期为10分钟，请尽快重置密码。</p>
		<p>如果这不是您的操作，请忽略此邮件。</p>
	`, code)

	return s.SendMail([]string{email}, subject, body)
}

// 发送账户到期提醒
func (s *MailService) SendExpirationNotice(email string, daysLeft int) error {
	subject := "账户即将到期"
	body := fmt.Sprintf(`
		<h3>账户到期提醒</h3>
		<p>您的账户将在 %d 天后到期。</p>
		<p>请及时续费以避免服务中断。</p>
	`, daysLeft)

	return s.SendMail([]string{email}, subject, body)
}
