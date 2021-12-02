package notify

import "github.com/sirupsen/logrus"

type NotifierConfig struct {
	Subject  string
	Content  string
	Tenant   string
	From     string
	Markdown bool
}

type Notifier struct {
	config *NotifierConfig
	logger *logrus.Logger
}

func NewNotifier(config *NotifierConfig, logger *logrus.Logger) *Notifier {
	return &Notifier{
		config: config,
		logger: logger,
	}
}

func (n *Notifier) SendEmail(targets ...string) error {
	n.logger.Debugf("send email to: %s", targets)
	return sendEmail(n.config.Subject, n.config.Content, n.config.From, targets)
}

func (n *Notifier) SendWecom(targets ...string) error {
	n.logger.Debugf("send wecom to: %s", targets)
	return nil
}

func (n *Notifier) SendLark(targets ...string) error {
	n.logger.Debugf("send lark to: %s", targets)
	return nil
}

func (n *Notifier) SendTelegram(targets ...string) error {
	n.logger.Debugf("send telegram to: %s", targets)
	return nil
}
