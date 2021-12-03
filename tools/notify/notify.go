package notify

import (
	"fmt"
)

type Notifier interface {
	SendMessage(message *MessageConfig, targets ...string) error
	SendHostAlert(alert *HostAlertConfig, targets ...string) error
	SendServiceAlert(alert *ServiceAlertConfig, targets ...string) error
}

type MessageConfig struct {
	Subject  string
	Content  string
	From     string
	Markdown bool
}

type Targets struct {
	Tenant string

	Email []string
	Lark  []string

	// TODO: to be supported
	Wework     []string
	Sms        []string
	Pushbullet []string
	Pushover   []string
	Telegram   []string
}

func GetNotifier(name, tenant string) (Notifier, error) {
	switch name {
	case "email":
		return NewEmailNotifier(tenant)
	case "lark":
		return NewLarkNotifier(tenant)
	}

	return nil, fmt.Errorf("notifier %s not supported", name)
}
