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
	From      string `yaml:"from"`
	AlertFrom string `yaml:"alert_from"`
	SMTP      struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"smtp"`
}

type EmailNotifier struct {
	config *SMTPConfig
}

func NewEmailNotifier(tenant string) (*EmailNotifier, error) {
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
	if config.AlertFrom == "" {
		config.AlertFrom = config.From
	}

	return &EmailNotifier{
		config: config,
	}, nil
}
func (n *EmailNotifier) sendMail(subject, content, from string, to []string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", content)
	d := gomail.NewDialer(
		n.config.SMTP.Host,
		n.config.SMTP.Port,
		n.config.SMTP.Username,
		n.config.SMTP.Password,
	)
	return d.DialAndSend(m)
}

func (n *EmailNotifier) SendMessage(message *MessageConfig, targets ...string) error {
	// NOTE: modified message itself
	if message.From == "" {
		message.From = n.config.From
	}
	return n.sendMail(message.Subject, message.Content, message.From, targets)
}

func (n *EmailNotifier) SendHostAlert(alert *HostAlertConfig, targets ...string) error {
	contentTmpl := `***** Icinga *****

Type: %s

Host: %s
Address: %s
State: %s

Date/Time: %s

Additional: %s

Comment: [%s] %s

Acknowledge: %s
`
	subject := fmt.Sprintf("Host %s alert for %s(%s)!",
		alert.HostState, alert.HostName, alert.HostAddress)
	content := fmt.Sprintf(
		contentTmpl,
		alert.NotificationType,
		alert.HostDisplayName,
		alert.HostAddress,
		alert.HostState,
		alert.LongDateTime,
		alert.HostOutput,
		alert.NotificationAuthorName,
		alert.NotificationComment,
		getAckLink(alert.AckLinkURL, "", alert.HostName, alert.ContactName),
	)

	return n.sendMail(subject, content, n.config.AlertFrom, targets)
}

func (n *EmailNotifier) SendServiceAlert(alert *ServiceAlertConfig, targets ...string) error {
	subject := fmt.Sprintf("%s - %s/%s is %s",
		alert.NotificationType,
		alert.HostDisplayName,
		alert.ServiceDisplayName,
		alert.ServiceState)
	contentTmpl := `***** Icinga *****

Type: %s

Service: %s
Host: %s
Address: %s
State: %s

Date/Time: %s

Additional: %s

Comment: [%s] %s

Link: %s
Wiki: %s
Acknowledge: %s
`
	content := fmt.Sprintf(
		contentTmpl,
		alert.NotificationType,
		alert.ServiceDisplayName,
		alert.HostDisplayName,
		alert.HostAddress,
		alert.ServiceState,
		alert.LongDateTime,
		alert.ServiceOutput,
		alert.NotificationAuthorName,
		alert.NotificationComment,
		getIcingaLink(alert.IcingaWebBaseURL, alert.HostName, alert.ServiceName),
		alert.ServiceWiki,
		getAckLink(alert.AckLinkURL, alert.ServiceName, alert.HostName, alert.ContactName),
	)

	return n.sendMail(subject, content, n.config.AlertFrom, targets)
}
