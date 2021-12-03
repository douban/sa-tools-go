package main

import (
	"github.com/douban/sa-tools-go/tools/notify"

	"github.com/caarlos0/env/v6"
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
			var alert notify.HostAlertConfig
			if err := env.Parse(&alert); err != nil {
				logger.Fatalf("parse icinga host alert from env failed: %s", err)
			}
			if targets.Email != nil {
				sendHostAlert("email", targets.Tenant, &alert, targets.Email)
			}
			if targets.Lark != nil {
				sendHostAlert("lark", targets.Tenant, &alert, targets.Lark)
			}
		},
	}
	addNotifyTargetFlags(cmd.Flags(), targets)

	return cmd
}
