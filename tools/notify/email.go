package notify

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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

type EmailNotifier struct {
	config *SMTPConfig
	logger *logrus.Logger
}

func NewEmailNotifier(tenant string, logger *logrus.Logger) (*EmailNotifier, error) {
	var cfg EmailConfig
	err := secrets.Load("email", &cfg)
	if err != nil {
		return nil, errors.Wrap(err, "load email secrets failed")
	}
	if tenant == "" {
		tenant = cfg.Default
	}
	config, ok := cfg.Tenants[tenant]
	if !ok {
		return nil, fmt.Errorf("tenant %s not fond in email secret", tenant)
	}
	if config.SMTP.Username == "" {
		config.SMTP.Username = config.From
	}

	return &EmailNotifier{
		config: config,
		logger: logger,
	}, nil
}

func (n *EmailNotifier) SendMessage(message *MessageConfig, targets ...string) error {
	n.logger.Infof("send email to: %s", targets)

	// NOTE: modified message itself
	if message.From == "" {
		message.From = n.config.From
	}

	m := gomail.NewMessage()
	m.SetHeader("From", n.config.From)
	m.SetHeader("To", targets...)
	m.SetHeader("Subject", message.Subject)
	m.SetBody("text/plain", message.Content)
	d := gomail.NewDialer(
		n.config.SMTP.Host,
		n.config.SMTP.Port,
		n.config.SMTP.Username,
		n.config.SMTP.Password,
	)
	return d.DialAndSend(m)
}
