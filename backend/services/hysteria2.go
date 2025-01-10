package services

import (
	"encoding/json"
	"fmt"
	"hysteria2-panel/models"
	"os"
	"path/filepath"
)

type Hysteria2Service struct {
	configDir   string
	certPath    string
	keyPath     string
	defaultAuth AuthConfig
}

func NewHysteria2Service(configDir, certPath, keyPath string) *Hysteria2Service {
	return &Hysteria2Service{
		configDir: configDir,
		certPath:  certPath,
		keyPath:   keyPath,
		defaultAuth: AuthConfig{
			Type: "password",
		},
	}
}

// 生成服务器配置
func (s *Hysteria2Service) GenerateServerConfig(userConfig *models.UserConfig) error {
	config := &Hysteria2Config{
		Listen: fmt.Sprintf(":%d", userConfig.Port),
		TLS: TLSConfig{
			Cert: s.certPath,
			Key:  s.keyPath,
		},
		Auth: AuthConfig{
			Type:     s.defaultAuth.Type,
			Password: userConfig.Password,
		},
		Bandwidth: BWConfig{
			Up:   userConfig.UpSpeed,
			Down: userConfig.DownSpeed,
		},
	}

	// 创建配置目录
	if err := os.MkdirAll(s.configDir, 0755); err != nil {
		return err
	}

	// 生成配置文件
	configPath := filepath.Join(s.configDir, fmt.Sprintf("config_%d.json", userConfig.UserID))
	configData, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, configData, 0644)
}

// 生成客户端配置
func (s *Hysteria2Service) GenerateClientConfig(domain string, userConfig *models.UserConfig) (*ClientConfig, error) {
	return &ClientConfig{
		Server: fmt.Sprintf("%s:%d", domain, userConfig.Port),
		Auth:   userConfig.Password,
		TLS: ClientTLSConfig{
			SNI:       domain,
			Insecure:  false,
			Pinned:    "",
			CA:        "",
			SkipVerif: false,
		},
		Bandwidth: BWConfig{
			Up:   userConfig.UpSpeed,
			Down: userConfig.DownSpeed,
		},
		FastOpen:  true,
		Lazy:      false,
		Socks5:    ClientSocks5Config{Listen: "127.0.0.1:1080"},
		HTTP:      ClientHTTPConfig{Listen: "127.0.0.1:8080"},
		Transport: "",
		Obfs:      "",
		ObfsParam: "",
	}, nil
}

type ClientConfig struct {
	Server    string             `json:"server"`
	Auth      string             `json:"auth"`
	TLS       ClientTLSConfig    `json:"tls"`
	Bandwidth BWConfig           `json:"bandwidth"`
	FastOpen  bool               `json:"fastOpen"`
	Lazy      bool               `json:"lazy"`
	Socks5    ClientSocks5Config `json:"socks5"`
	HTTP      ClientHTTPConfig   `json:"http"`
	Transport string             `json:"transport,omitempty"`
	Obfs      string             `json:"obfs,omitempty"`
	ObfsParam string             `json:"obfsParam,omitempty"`
}

type ClientTLSConfig struct {
	SNI       string `json:"sni"`
	Insecure  bool   `json:"insecure"`
	Pinned    string `json:"pinned,omitempty"`
	CA        string `json:"ca,omitempty"`
	SkipVerif bool   `json:"skipVerif,omitempty"`
}

type ClientSocks5Config struct {
	Listen string `json:"listen"`
}

type ClientHTTPConfig struct {
	Listen string `json:"listen"`
}
