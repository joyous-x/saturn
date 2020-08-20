package config

import (
	"fmt"
	"io/ioutil"

	"github.com/joyous-x/saturn/common/xlog"
	"gopkg.in/yaml.v2"
)

// ConfObjectItem ...
type ConfObjectItem struct {
	Key  string
	Path string
	Data interface{}
}

// MgrCenter ...
type MgrCenter struct {
	Configs map[string]*ConfObjectItem
}

// AddConfObjectItem ...
func (m *MgrCenter) AddConfObjectItem(configKey, configFilePath string, dataObj interface{}) error {
	if nil == m.Configs {
		m.Configs = make(map[string]*ConfObjectItem)
	}
	if _, ok := m.Configs[configKey]; ok {
		return fmt.Errorf("config key:%v already exists", configKey)
	}

	m.Configs[configKey] = &ConfObjectItem{configKey, configFilePath, dataObj}
	return nil
}

// Reload ...
func (m *MgrCenter) Reload() error {
	err := error(nil)
	for k, v := range m.Configs {
		c, e := ioutil.ReadFile(v.Path)
		if e != nil {
			err = e
			break
		}
		if e := yaml.Unmarshal(c, v.Data); e != nil {
			err = e
			break
		}
		xlog.Debug("reload config: %v, %v, ok", k, v.Path)
	}
	return err
}
