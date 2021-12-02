package notify

import (
	"fmt"

	"github.com/pkg/errors"
	gomail "gopkg.in/mail.v2"

	"github.com/douban/sa-tools-go/libs/secrets"
)

type EmailConfig struct {
	Default string `yaml:"default"`

	Tenants map[string]*SMTPConfig `yaml:"tenants"`
}

type SMTPConfig struct {
	From string `yaml:"from"`
	SMTP struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"smtp"`
}

func (n *Notifier) SendEmail(targets ...string) error {
	n.logger.Infof("send email to: %s", targets)

	cfg := EmailConfig{}
	err := secrets.Load("email", &cfg)
	if err != nil {
		return errors.Wrap(err, "load email secrets failed")
	}

	var tenant string
	if n.config.Tenant != "" {
		tenant = n.config.Tenant
	} else {
		tenant = cfg.Default
	}
	tc, ok := cfg.Tenants[tenant]
	if !ok {
		return fmt.Errorf("tenant %s not fond in email secret", tenant)
	}

	if n.config.From != "" {
		tc.From = n.config.From
	}

	if tc.SMTP.Username == "" {
		tc.SMTP.Username = tc.From
	}

	m := gomail.NewMessage()
	m.SetHeader("From", tc.From)
	m.SetHeader("To", targets...)
	m.SetHeader("Subject", n.config.Subject)
	m.SetBody("text/plain", n.config.Content)
	d := gomail.NewDialer(tc.SMTP.Host, tc.SMTP.Port, tc.SMTP.Username, tc.SMTP.Password)
	return d.DialAndSend(m)
}
