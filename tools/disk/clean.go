package disk

import (
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

func Clean(dataPath string, dryRun bool, logger *logrus.Logger) error {
	var cleanCMD []string
	if dryRun {
		cleanCMD = []string{"-L", dataPath, "-maxdepth", "1", "-type", "f", "-name", "ncdu-export-*", "-mtime", "+1"}
	} else {
		cleanCMD = []string{"-L", dataPath, "-maxdepth", "1", "-type", "f", "-name", "ncdu-export-*", "-mtime", "+1", "-delete"}
	}
	cmd := exec.Command("find", cleanCMD...)
	logger.Debugf("executing: %s\n", strings.Join(cmd.Args, " "))
	return cmd.Run()
}
