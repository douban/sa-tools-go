package notify

import (
	"fmt"
	"net/smtp"

	"github.com/douban/sa-tools-go/libs/secrets"
	"github.com/pkg/errors"
)

type EmailConfig struct {
	SMTP struct {
		Host     string `yaml:"host"`
		Port     uint16 `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"smtp"`
	From string `yaml:"from"`
}

func (c *EmailConfig) GetSMTPAddr() string {
	return fmt.Sprintf("%s:%d", c.SMTP.Host, c.SMTP.Port)
}

func sendEmail(subject, content, from string, target []string) error {
	cfg := EmailConfig{}
	err := secrets.Load("email", &cfg)
	if err != nil {
		return errors.Wrap(err, "load email secrets failed")
	}
	if from == "" {
		from = cfg.From
	}
	auth := smtp.PlainAuth("", cfg.SMTP.Username, cfg.SMTP.Password, cfg.SMTP.Host)

	return smtp.SendMail(cfg.GetSMTPAddr(), auth, from, target, []byte(content))

}
