package gaap

import (
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/mxpaul/gaap/pkg/loggy"
	"github.com/mxpaul/gaap/pkg/tooling"
)

type Config[AppConfigType any] struct {
	Daemon      DaemonConfig  `yaml:"daemon"`
	Application AppConfigType `yaml:",inline"`
}

type DaemonConfig struct {
	Log          loggy.Config   `yaml:"log,omitempty"`
	Tooling      tooling.Config `yaml:"tooling,omitempty"`
	GracefulWait time.Duration  `yaml:"graceful_wait,omitempty"`
}

func LoadConfigFileYAML(configPath string, dst interface{}) error {
	path, err := filepath.Abs(configPath)
	if err != nil {
		return errors.Wrapf(err, "make %q absolute path failed", configPath)
	}
	configBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.Wrapf(err, "ReadFile(%q) error", path)
	}
	err = yaml.Unmarshal(configBytes, dst)
	if err != nil {
		return errors.Wrapf(err, "yaml unmarshal %q", path)
	}

	return nil
}
