package main

import (
	"github.com/spf13/cobra"
)

func cmdAlert() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alert",
		Short: "send icinga alert",
	}
	cmd.AddCommand(
		cmdAlertHost(),
		cmdAlertService(),
	)
	return cmd
}
