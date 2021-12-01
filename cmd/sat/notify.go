package main

import (
	"github.com/spf13/cobra"

	"github.com/douban/sa-tools-go/tools/notify"
)

func cmdNotify() *cobra.Command {
	config := &notify.NotifierConfig{}

	var (
		flagEmail    []string
		flagLark     []string
		flagWecom    []string
		flagTelegram []string
	)

	cmd := &cobra.Command{
		Use: "notify",
		Run: func(cmd *cobra.Command, args []string) {
			notifier := notify.NewNotifier(config, logger)
			if flagEmail != nil {
				notifier.SendEmail(flagEmail...)
			}
			if flagLark != nil {
				notifier.SendLark(flagLark...)
			}
			if flagWecom != nil {
				notifier.SendWecom(flagWecom...)
			}
			if flagTelegram != nil {
				notifier.SendTelegram(flagTelegram...)
			}
		},
	}
	cmd.Flags().StringVarP(&config.Subject, "subject", "s", "sent from sa-notify", "")
	cmd.Flags().StringVarP(&config.Content, "content", "c", "", "")
	cmd.Flags().StringVarP(&config.Tenant, "tenant", "t", "", "company user in, used when multiple company or tenant is configured")
	cmd.Flags().StringVarP(&config.From, "from", "", "", "from address, currently only works for email")
	cmd.Flags().BoolVarP(&config.Markdown, "markdown", "", false, "use markdown rendering, only lark & wework & telegram supported")

	cmd.Flags().StringSliceVarP(&flagEmail, "email", "", nil, "")
	cmd.Flags().StringSliceVarP(&flagLark, "lark", "", nil, "")
	cmd.Flags().StringSliceVarP(&flagWecom, "wecom", "", nil, "")
	cmd.Flags().StringSliceVarP(&flagTelegram, "telegram", "", nil, "")

	return cmd
}
