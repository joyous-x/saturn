package ginbox

import (
)

type GinServerConf struct {
	Key      string `yaml:"key"`
	Port     int    `yaml:"port"`
	CertFile string `yaml:"certFile"`
	KeyFile  string `yaml:"keyFile"`
}
