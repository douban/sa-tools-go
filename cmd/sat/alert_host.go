package main

import (
	"github.com/douban/sa-tools-go/tools/notify"

	"github.com/spf13/cobra"
)

func sendHostAlert(name, tenant string, alert *notify.HostAlertConfig, targets []string) {
	if notifier, err := notify.GetNotifier(name, tenant); err != nil {
		logger.Errorf("get %s notifier failed: %s", name, err)
	} else {
		if err := notifier.SendHostAlert(alert, targets...); err != nil {
			logger.Errorf("send %s failed: %s", name, err)
		} else {
			logger.Infof("sent %s alert to: %s", name, targets)
		}
	}
}

func cmdAlertHost() *cobra.Command {
	targets := &notifyTargets{}

	cmd := &cobra.Command{
		Use: "host",
		Run: func(cmd *cobra.Command, args []string) {
			alert, err := notify.HostAlertFromEnv()
			if err != nil {
				logger.Fatalln(err)
			}
			if targets.Email != nil {
				sendHostAlert("email", targets.Tenant, alert, targets.Email)
			}
			if targets.Lark != nil {
				sendHostAlert("lark", targets.Tenant, alert, targets.Lark)
			}
		},
	}
	addNotifyTargetFlags(cmd.Flags(), targets)

	return cmd
}
