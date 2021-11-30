package disk

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/douban/sa-tools-go/libs/humanize"
	"github.com/gofrs/flock"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type UsageConfig struct {
	NcduDataPath string
	DataFile     string
	TopN         uint64
	MaxDepth     uint64
	CheckPath    string
	ForceCheck   bool
	ForceRead    bool
	Block        bool
}

type DiskUsageChecker struct {
	config *UsageConfig
	logger *logrus.Logger
}

func NewDiskUsageChecker(config *UsageConfig, logger *logrus.Logger) (*DiskUsageChecker, error) {
	checker := &DiskUsageChecker{
		config: config,
		logger: logger,
	}

	p, err := filepath.Abs(config.CheckPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get check path")
	}
	checker.config.CheckPath = p
	return checker, nil
}

func (c *DiskUsageChecker) HasDataFile() bool {
	return c.config.DataFile != ""
}

func (c *DiskUsageChecker) IsForceRead() bool {
	return c.config.ForceRead
}

func (c *DiskUsageChecker) FindLatestExportedData() error {
	c.logger.Infof("Finding lastest ncdu exported data file from %s\n", c.config.NcduDataPath)
	findLatestCMD := fmt.Sprintf("ls -lt %s | grep ncdu-export-%s-", c.config.NcduDataPath, c.getEscapedCheckPath())
	out, err := exec.Command("sh", "-c", findLatestCMD).Output()
	if err != nil {
		return nil
	}
	latestLine := strings.Fields(strings.Split(string(out), "\n")[0])
	latestFileName := latestLine[len(latestLine)-1]
	dataFile := filepath.Join(c.config.NcduDataPath, latestFileName)
	c.logger.Infof("Found ncdu exported data file %s\n", dataFile)

	dataFileTimeStr := RE_TIMESTRING.FindAllString(dataFile, -1)[0]
	dataFileTime, err := time.Parse("20060102150405", dataFileTimeStr)
	if err != nil {
		return errors.Wrap(err, "parse ncdu exported data file time error")
	}
	if !c.IsForceRead() && time.Since(dataFileTime) > DATA_STALE_THRESHOLD {
		c.logger.Warning("But data file is out of date.")
		return nil
	}
	c.config.DataFile = dataFile
	return nil
}

func (c *DiskUsageChecker) getEscapedCheckPath() string {
	return strings.ReplaceAll(c.config.CheckPath, "/", "%")
}

func (c *DiskUsageChecker) getLockFile() string {
	return NCDU_JOB_LOCK_PREFIX + "-" + c.getEscapedCheckPath()
}

func (c *DiskUsageChecker) getDataFilePath() error {
	timeNowStr := time.Now().Format("20060102150405")
	dataFileName := fmt.Sprintf("ncdu-export-%s-%s.gz", c.getEscapedCheckPath(), timeNowStr)
	dataFilePath := filepath.Join(c.config.NcduDataPath, dataFileName)
	_, err := os.OpenFile(dataFilePath, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return errors.Wrap(err, "failed to open data file")
	}
	c.config.DataFile = dataFilePath
	return nil
}

func (c *DiskUsageChecker) checkDataPathAvailable() error {
	out, err := exec.Command("df", "-BM", c.config.NcduDataPath).Output()
	if err != nil {
		return errors.Wrapf(err, "failed to get data path available")
	}
	availStr := strings.TrimRight(strings.Fields(string(out))[10], "M")
	avail, err := strconv.Atoi(availStr)
	if err != nil {
		return errors.Wrap(err, "failed to parse data path available")
	}
	if avail < 50 {
		return fmt.Errorf("not enough space(<50M) for storing datafile on %s, disk usage scanning stopped", c.config.NcduDataPath)
	}
	return nil
}

func (c *DiskUsageChecker) runNcdu() error {
	ncduCmd := fmt.Sprintf("ncdu -0xo- %s | gzip > %s", c.config.CheckPath, c.config.DataFile)
	p := exec.Command("sh", "-c", ncduCmd)
	if err := p.Start(); err != nil {
		return errors.Wrap(err, "failed to start ncdu")
	}
	c.logger.Infof("Checking %s using ncdu, and exporting datafile %s on background", c.config.CheckPath, c.config.DataFile)

	if c.config.Block {
		if err := p.Wait(); err != nil {
			return errors.Wrap(err, "failed to wait for ncdu process")
		}
		if err := c.ReadData(); err != nil {
			return errors.Wrap(err, "read ncdu result error")
		}
	}
	return nil
}

func (c *DiskUsageChecker) ReadData() error {
	readDataCmd := fmt.Sprintf("zcat %s | ncdu -f-", c.config.DataFile)
	if c.config.TopN == 0 {
		c.logger.Infof("About to open ncdu ui, read datafile %s\n", c.config.DataFile)
		c.logger.Info("Press the Enter Key to continue...")
		fmt.Scanln()
		if err := exec.Command("sh", "-c", readDataCmd).Run(); err != nil {
			return errors.Wrap(err, "failed to read ncdu data")
		}
	} else {
		c.logger.Infof("About to calc top %d huge items in depth %d.\n", c.config.TopN, c.config.MaxDepth)
		out, err := exec.Command("zcat", c.config.DataFile).Output()
		if err != nil {
			return errors.Wrap(err, "failed to extract ncdu data")
		}
		ncduDetail := []json.RawMessage{}
		err = json.Unmarshal(out, &ncduDetail)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal ncdu data")
		}
		c.logger.Debug("calc top huge dirs from ncdu...")
		ret, _, err := TopHugeDirsFromNcdu(c.config.TopN, ncduDetail[3], "", c.config.MaxDepth, 0)
		if err != nil {
			return errors.Wrap(err, "failed to get top huge dirs")
		}
		c.logger.Debug("done top huge dirs from ncdu")
		fmt.Printf("Top %d huge items in depth %d:\n", c.config.TopN, c.config.MaxDepth)
		for k, v := range ret {
			fmt.Printf("%s\t%s\n", humanize.ByteCountIEC(v), k)
		}
	}
	return nil
}

func (c *DiskUsageChecker) Check() error {
	if err := c.checkDataPathAvailable(); err != nil {
		return errors.Wrap(err, "error checking data path available")
	}

	fileLock := flock.New(c.getLockFile())
	locked, err := fileLock.TryLock()
	if err != nil || !locked {
		return errors.New("sat disk check job already in running")
	}
	defer fileLock.Unlock()

	if err := c.getDataFilePath(); err != nil {
		return errors.Wrap(err, "error getting data file path")
	}

	if err := c.runNcdu(); err != nil {
		return errors.Wrap(err, "error running ncdu")
	}

	return nil
}
