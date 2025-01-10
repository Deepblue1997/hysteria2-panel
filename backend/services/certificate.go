package services

import (
	"os/exec"
)

type CertificateService struct {
	Domain string
	Email  string
}

func (c *CertificateService) ObtainCert() error {
	cmd := exec.Command("acme.sh", "--issue", "-d", c.Domain, "--standalone", "--email", c.Email)
	return cmd.Run()
}

func (c *CertificateService) RenewCert() error {
	cmd := exec.Command("acme.sh", "--renew", "-d", c.Domain)
	return cmd.Run()
}
