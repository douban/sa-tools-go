package secrets

import (
	"io/ioutil"
	"path/filepath"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

const (
	secretDir = "/etc/sa-tools"
)

func Load(name string, target interface{}) error {
	configFile := filepath.Join(secretDir, name) + ".yaml"
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		return errors.Wrapf(err, "failed to read secret file %s", configFile)
	}
	err = yaml.Unmarshal(content, target)
	if err != nil {
		return errors.Wrapf(err, "failed to unmarshal secret file %s", configFile)
	}
	return nil
}
