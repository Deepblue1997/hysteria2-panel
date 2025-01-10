package services

import (
	"fmt"
	"hysteria2-panel/models"
	"os"
	"os/exec"
	"path/filepath"
)

type CertService struct {
	settingService *SettingService
	certDir        string
}

func NewCertService(settingService *SettingService, certDir string) *CertService {
	return &CertService{
		settingService: settingService,
		certDir:        certDir,
	}
}

// 申请证书
func (s *CertService) ObtainCert() error {
	config, err := s.settingService.GetTLSConfig()
	if err != nil {
		return err
	}

	if !config.AutoCert {
		return nil
	}

	// 创建证书目录
	if err := os.MkdirAll(s.certDir, 0755); err != nil {
		return err
	}

	// 根据不同的证书提供商使用不同的命令
	var cmd *exec.Cmd
	switch config.Provider {
	case "acme.sh":
		cmd = s.getAcmeShCommand(config)
	default:
		return fmt.Errorf("不支持的证书提供商: %s", config.Provider)
	}

	// 执行命令
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("申请证书失败: %v, 输出: %s", err, string(output))
	}

	return nil
}

// 获取acme.sh命令
func (s *CertService) getAcmeShCommand(config *models.TLSConfig) *exec.Cmd {
	args := []string{
		"--issue",
		"-d", config.Domain,
		"--cert-file", filepath.Join(s.certDir, "cert.pem"),
		"--key-file", filepath.Join(s.certDir, "key.pem"),
		"--fullchain-file", filepath.Join(s.certDir, "fullchain.pem"),
	}

	// 根据DNS提供商添加相应的参数
	if config.DNSProvider != "" {
		args = append(args, "--dns", config.DNSProvider)
		if config.DNSToken != "" {
			args = append(args, fmt.Sprintf("--dnssleep %s", config.DNSToken))
		}
	} else {
		args = append(args, "--standalone")
	}

	return exec.Command("acme.sh", args...)
}

// 更新证书
func (s *CertService) RenewCert() error {
	config, err := s.settingService.GetTLSConfig()
	if err != nil {
		return err
	}

	if !config.AutoCert {
		return nil
	}

	cmd := exec.Command("acme.sh", "--renew", "-d", config.Domain)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("更新证书失败: %v, 输出: %s", err, string(output))
	}

	return nil
}

// 检查证书是否存在
func (s *CertService) CheckCertExists() bool {
	certFile := filepath.Join(s.certDir, "cert.pem")
	keyFile := filepath.Join(s.certDir, "key.pem")
	_, certErr := os.Stat(certFile)
	_, keyErr := os.Stat(keyFile)
	return !os.IsNotExist(certErr) && !os.IsNotExist(keyErr)
}
