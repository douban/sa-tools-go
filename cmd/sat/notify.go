package main

import (
	"github.com/spf13/cobra"

	"github.com/douban/sa-tools-go/tools/notify"
)

func sendNotification(name, tenant string, message *notify.MessageConfig, target []string) {
	notifier, err := notify.GetNotifier(name, tenant, logger)
	if err != nil {
		logger.Errorf("get %s notifier failed: %s", name, err)
	}
	err = notifier.SendMessage(message, target...)
	if err != nil {
		logger.Errorf("send %s failed: %s", name, err)
	}
}

func cmdNotify() *cobra.Command {
	message := &notify.MessageConfig{}

	var (
		flagEmail    []string
		flagLark     []string
		flagWecom    []string
		flagTelegram []string
		flagTenant   string
	)

	cmd := &cobra.Command{
		Use: "notify",
		Run: func(cmd *cobra.Command, args []string) {
			if flagEmail != nil {
				sendNotification("email", flagTenant, message, flagEmail)
			}
			if flagLark != nil {
				sendNotification("lark", flagTenant, message, flagLark)
			}

		},
	}
	cmd.Flags().StringVarP(&message.Subject, "subject", "s", "sent from sa-notify", "")
	cmd.Flags().StringVarP(&message.Content, "content", "c", "", "")
	cmd.Flags().StringVarP(&message.From, "from", "", "", "from address, currently only works for email")
	cmd.Flags().BoolVarP(&message.Markdown, "markdown", "", false, "use markdown rendering, only lark & wework & telegram supported")

	cmd.Flags().StringVarP(&flagTenant, "tenant", "t", "", "company user in, used when multiple company or tenant is configured")
	cmd.Flags().StringSliceVarP(&flagEmail, "email", "", nil, "")
	cmd.Flags().StringSliceVarP(&flagLark, "lark", "", nil, "")
	cmd.Flags().StringSliceVarP(&flagWecom, "wecom", "", nil, "")
	cmd.Flags().StringSliceVarP(&flagTelegram, "telegram", "", nil, "")

	return cmd
}
