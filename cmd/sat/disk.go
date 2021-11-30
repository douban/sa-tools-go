package main

import (
	"github.com/spf13/cobra"

	"github.com/douban/sa-tools-go/tools/disk"
)

func cmdDisk() *cobra.Command {

	cmd := &cobra.Command{
		Use: "disk",
		Example: `sat disk usage
sat disk usage -n 5 -d 3
sat disk usage -r /data1/ncdu-export-%-20160513142844.gz
sat disk usage -c /tmp
sat disk usage -p /data
sat disk usage --force-check

sat disk clean`,
	}

	cmd.AddCommand(
		cmdDiskUsage(),
		cmdDiskClean(),
	)
	return cmd
}

func cmdDiskUsage() *cobra.Command {
	config := &disk.UsageConfig{}
	cmd := &cobra.Command{
		Use:   "usage",
		Short: "Disk usage operation",
		Long: `
if -p/--ncdu-data-path not specified, a path will be selected automatically, typically /data or /
if -f/--force-check specified, a check will be forced
if -r/--data-file specified, the ncdu exported data file will be read;
    otherwise a latest data will be selected;
    if no data file within 1 hour found, a check will be performed eventually;
    if -R/--force-read specified, check will be aborted`,
		Run: func(cmd *cobra.Command, args []string) {
			checker, err := disk.NewDiskUsageChecker(config, logger)
			if err != nil {
				logger.Fatalf("failed to get disk usage checker: %s", err)
			}
			if config.ForceCheck {
				err = checker.Check()
				if err != nil {
					logger.Fatalf("failed to check disk: %s", err)
				}
			} else {
				if !checker.HasDataFile() {
					err = checker.FindLatestExportedData()
					if err != nil {
						logger.Fatalf("failed to find latest exported data: %s", err)
					}
				}
				if checker.HasDataFile() {
					err = checker.ReadData()
					if err != nil {
						logger.Fatalf("failed to read disk check data: %s", err)
					}
				} else {
					if checker.IsForceRead() {
						logger.Fatalln("no data file found")
					} else {
						logger.Warning("Recent data file not found. Checking disk usage.")
						err = checker.Check()
						if err != nil {
							logger.Fatalf("failed to check disk: %s", err)
						}
					}
				}
			}
		},
	}
	cmd.Flags().StringVarP(&config.NcduDataPath, "ncdu-data-path", "p", "", "Path to store ncdu exported data")
	cmd.Flags().StringVarP(&config.DataFile, "data-file", "r", "", "Read ncdu export data file.")
	cmd.Flags().Uint64VarP(&config.TopN, "topn", "n", 0, "Specify top N huge files for output.")
	cmd.Flags().Uint64VarP(&config.MaxDepth, "max-depth", "d", 3, "Max path depth to calc top N.")
	cmd.Flags().StringVarP(&config.CheckPath, "check-path", "c", "/", "Check disk and export ncdu data file.")
	cmd.Flags().BoolVarP(&config.ForceCheck, "force-check", "f", false, "Force check, do not use existing exported file.")
	cmd.Flags().BoolVarP(&config.ForceRead, "force-read", "R", false, "Force read latest exported data file, do not check.")
	cmd.Flags().BoolVarP(&config.Block, "block", "b", false, "Disk usage check block mode.")

	return cmd
}

func cmdDiskClean() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "clean",
		Short: "Clean disk, including outdated ncdu exported data file",
	}

	return cmd
}
