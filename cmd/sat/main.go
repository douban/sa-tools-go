package main

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	logger = logrus.New()
	ctx    = context.Background()
)

func main() {
	var (
		flagVerbose bool
		flagQuiet   bool
	)

	subCommands := []*cobra.Command{
		cmdDisk(),
		cmdNotify(),
		cmdAlert(),
	}

	rootCmd := &cobra.Command{
		Use: "sat",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			logger.SetOutput(os.Stdout)
			if flagVerbose {
				logger.SetLevel(logrus.DebugLevel)
			} else if flagQuiet {
				logger.SetLevel(logrus.WarnLevel)
			} else {
				logger.SetLevel(logrus.InfoLevel)
			}
		},
	}
	rootCmd.PersistentFlags().BoolVarP(&flagVerbose, "verbose", "v", false, "log debug")
	rootCmd.PersistentFlags().BoolVarP(&flagQuiet, "quiet", "q", false, "log warning")

	rootCmd.AddCommand(subCommands...)
	rootCmd.Execute()
}
