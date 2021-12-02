package notify

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type Notifier interface {
	SendMessage(message *MessageConfig, targets ...string) error
}

type MessageConfig struct {
	Subject  string
	Content  string
	From     string
	Markdown bool
}

func GetNotifier(name, tenant string, logger *logrus.Logger) (Notifier, error) {
	switch name {
	case "email":
		return NewEmailNotifier(tenant, logger)
	case "lark":
		return NewLarkNotifier(tenant, logger)
	}

	return nil, fmt.Errorf("notifier %s not supported", name)
}
