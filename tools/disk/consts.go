package disk

import (
	"regexp"
	"time"
)

const (
	NCDU_EXPORT_DATA_PATH = "/data"
	NCDU_JOB_LOCK_PREFIX  = "/tmp/sa-disk-ncdu-lock"
	DATA_STALE_THRESHOLD  = time.Hour
)

var (
	RE_TIMESTRING = regexp.MustCompile(`\d{14}`)
)
