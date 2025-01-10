package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Hysteria2Service struct {
	configDir    string
	certPath     string
	keyPath      string
}

func NewHysteria2Service(configDir, certPath, keyPath string) *Hysteria2Service {
	return &Hysteria2Service{
		configDir:    configDir,
		certPath:     certPath,
		keyPath:      keyPath,
	}
}

// 生成服务器配置
func (s *Hysteria2Service) GenerateServerConfig(nodeID uint) error {
	config := map[string]interface{}{
		"listen":     ":443",
		"cert":       s.certPath,
		"key":        s.keyPath,
		"auth":       map[string]string{"type": "password", "password": "your-password"},
		"masquerade": map[string]string{"type": "proxy", "proxy": map[string]string{"url": "https://www.google.com"}},
	}

	jsonBytes, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	configPath := filepath.Join(s.configDir, fmt.Sprintf("node_%d.json", nodeID))
	return os.WriteFile(configPath, jsonBytes, 0644)
}

// 生成客户端配置
func (s *Hysteria2Service) GenerateClientConfig(serverAddr, auth string) (string, error) {
	config := map[string]interface{}{
		"server":     serverAddr,
		"auth":       auth,
		"tls":        map[string]interface{}{"sni": "your-domain.com"},
		"transport":  map[string]interface{}{"type": "udp"},
		"bandwidth": map[string]interface{}{
			"up":   "100 mbps",
			"down": "100 mbps",
		},
	}

	jsonBytes, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}
