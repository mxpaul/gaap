package gaap

import (
	"github.com/spf13/pflag"
)

type CommandLineOptions struct {
	ConfigPath string
}

func ParseCommandLineOdDie() CommandLineOptions {
	opt := CommandLineOptions{}
	pflag.StringVarP(&opt.ConfigPath, "config", "c", "config.yaml", "path to config file in yaml format")
	pflag.Parse()
	return opt
}
