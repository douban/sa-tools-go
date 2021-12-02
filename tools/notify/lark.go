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

func (n *Notifier) SendLark(targets ...string) error {
	n.logger.Infof("send lark to: %s", targets)

	cfg := LarkTenantConfig{}
	err := secrets.Load("lark", &cfg)
	if err != nil {
		return errors.Wrap(err, "load lark secrets failed")
	}

	var tenant string
	if n.config.Tenant != "" {
		tenant = n.config.Tenant
	} else {
		tenant = cfg.Default
	}
	tc, ok := cfg.Tenants[tenant]
	if !ok {
		return fmt.Errorf("tenant %s not fond in lark secret", tenant)
	}

	bot := lark.NewChatBot(tc.AppID, tc.AppSecret)
	_, err = bot.GetTenantAccessTokenInternal(true)
	if err != nil {
		return errors.Wrap(err, "get lark access token failed")
	}

	b := lark.NewCardBuilder()
	card := b.Card(
		b.Div().Text(b.Text(n.config.Content).LarkMd()),
	).Indigo().Title(n.config.Subject)

	for _, target := range targets {
		var om lark.OutcomingMessage
		if n.config.Markdown {
			msg := lark.NewMsgBuffer(lark.MsgInteractive)
			om = msg.BindEmail(target).Card(card.String()).Build()
		} else {
			msg := lark.NewMsgBuffer(lark.MsgText)
			om = msg.BindEmail(target).Text(n.config.Content).Build()
		}
		if _, serr := bot.PostMessage(om); serr != nil {
			err = serr
			n.logger.Errorf("send lark to %s failed: %s", target, serr)
		}
	}
	return err
}
