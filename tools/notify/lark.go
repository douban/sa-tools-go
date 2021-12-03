package notify

import (
	"fmt"

	"github.com/go-lark/lark"
	"github.com/pkg/errors"

	"github.com/douban/sa-tools-go/libs/secrets"
)

type LarkTenantConfig struct {
	Default string `yaml:"default"`

	Tenants map[string]*LarkConfig `yaml:"tenants"`
}

type LarkConfig struct {
	AppID             string `yaml:"app_id"`
	AppSecret         string `yaml:"app_secret"`
	EncryptKey        string `yaml:"encrypt_key"`
	VerificationToken string `yaml:"verification_token"`
}

type LarkNotifier struct {
	bot *lark.Bot
}

func NewLarkNotifier(tenant string) (*LarkNotifier, error) {
	var cfg LarkTenantConfig
	err := secrets.Load("lark", &cfg)
	if err != nil {
		return nil, errors.Wrap(err, "load lark secrets failed")
	}
	if tenant == "" {
		tenant = cfg.Default
	}
	config, ok := cfg.Tenants[tenant]
	if !ok {
		return nil, fmt.Errorf("tenant %s not fond in lark secret", tenant)
	}
	bot := lark.NewChatBot(config.AppID, config.AppSecret)
	// NOTE: should refresh token periodically if used in long running program
	_, err = bot.GetTenantAccessTokenInternal(true)
	if err != nil {
		return nil, errors.Wrap(err, "get lark access token failed")
	}

	return &LarkNotifier{
		bot: bot,
	}, nil
}

func (n *LarkNotifier) SendMessage(message *MessageConfig, targets ...string) error {
	b := lark.NewCardBuilder()
	var errs []error
	for _, target := range targets {
		var om lark.OutcomingMessage
		if message.Markdown {
			msg := lark.NewMsgBuffer(lark.MsgInteractive)
			card := b.Card(
				b.Div().Text(b.Text(message.Content).LarkMd()),
			).Indigo().Title(message.Subject)
			om = msg.BindEmail(target).Card(card.String()).Build()
		} else {
			msg := lark.NewMsgBuffer(lark.MsgText)
			om = msg.BindEmail(target).Text(message.Content).Build()
		}
		if _, serr := n.bot.PostMessage(om); serr != nil {
			errs = append(errs, serr)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return fmt.Errorf("%v", errs)
}

func (n *LarkNotifier) SendHostAlert(alert *HostAlertConfig, targets ...string) error {
	return nil
}

func (n *LarkNotifier) SendServiceAlert(alert *ServiceAlertConfig, targets ...string) error {
	return nil
}
