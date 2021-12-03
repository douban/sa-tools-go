package main

import (
	"github.com/douban/sa-tools-go/tools/notify"

	"github.com/caarlos0/env/v6"
	"github.com/spf13/cobra"
)

func sendServiceAlert(name, tenant string, alert *notify.ServiceAlertConfig, targets []string) {
	if notifier, err := notify.GetNotifier(name, tenant); err != nil {
		logger.Errorf("get %s notifier failed: %s", name, err)
	} else {
		if err := notifier.SendServiceAlert(alert, targets...); err != nil {
			logger.Errorf("send %s failed: %s", name, err)
		} else {
			logger.Infof("sent %s alert to: %s", name, targets)
		}
	}
}

func cmdAlertService() *cobra.Command {
	targets := &notifyTargets{}

	cmd := &cobra.Command{
		Use: "service",
		Run: func(cmd *cobra.Command, args []string) {
			var alert notify.ServiceAlertConfig
			if err := env.Parse(&alert); err != nil {
				logger.Fatalf("parse icinga service alert from env failed: %s", err)
			}
			if targets.Email != nil {
				sendServiceAlert("email", targets.Tenant, &alert, targets.Email)
			}
			if targets.Lark != nil {
				sendServiceAlert("lark", targets.Tenant, &alert, targets.Lark)
			}

		},
	}
	addNotifyTargetFlags(cmd.Flags(), targets)

	return cmd
}