package main

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/douban/sa-tools-go/tools/notify"
)

func sendMessage(name, tenant string, message *notify.MessageConfig, targets []string) {
	if notifier, err := notify.GetNotifier(name, tenant); err != nil {
		logger.Errorf("get %s notifier failed: %s", name, err)
	} else {
		if err := notifier.SendMessage(message, targets...); err != nil {
			logger.Errorf("send %s failed: %s", name, err)
		} else {
			logger.Infof("sent %s message to: %s", name, targets)
		}
	}
}

type notifyTargets struct {
	Tenant string

	Email    []string
	Lark     []string
	Wecom    []string
	Telegram []string
}

func addNotifyTargetFlags(f *pflag.FlagSet, a *notifyTargets) {
	f.StringVarP(&a.Tenant, "tenant", "t", "", "company user in, used when multiple company or tenant is configured")
	f.StringSliceVarP(&a.Email, "email", "", nil, "")
	f.StringSliceVarP(&a.Lark, "lark", "", nil, "")
	f.StringSliceVarP(&a.Wecom, "wecom", "", nil, "")
	f.StringSliceVarP(&a.Telegram, "telegram", "", nil, "")
}

func cmdNotify() *cobra.Command {
	message := &notify.MessageConfig{}
	targets := &notifyTargets{}

	cmd := &cobra.Command{
		Use: "notify",
		Run: func(cmd *cobra.Command, args []string) {
			if targets.Email != nil {
				sendMessage("email", targets.Tenant, message, targets.Email)
			}
			if targets.Lark != nil {
				sendMessage("lark", targets.Tenant, message, targets.Lark)
			}
		},
	}
	addNotifyTargetFlags(cmd.Flags(), targets)
	cmd.Flags().StringVarP(&message.Subject, "subject", "s", "sent from sa-notify", "")
	cmd.Flags().StringVarP(&message.Content, "content", "c", "", "")
	cmd.Flags().StringVarP(&message.From, "from", "", "", "from address, currently only works for email")
	cmd.Flags().BoolVarP(&message.Markdown, "markdown", "", false, "use markdown rendering, only lark & wework & telegram supported")

	return cmd
}
