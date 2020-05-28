package config

import (
	"fmt"
	"github.com/joyous-x/saturn/common/xlog"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type ConfItem struct {
	Key  string
	Path string
	Data interface{}
}

type ConfigMgr struct {
	Configs map[string]*ConfItem
}

func (this *ConfigMgr) Init(confItems []*ConfItem) error {
	for _, v := range confItems {
		if len(v.Key) == 0 || len(v.Path) == 0 {
			return fmt.Errorf("configs init error: invalid key or path")
		}
		this.addConfigItem(v.Key, v)
	}
	return nil
}

func (this *ConfigMgr) addConfigItem(key string, item *ConfItem) error {
	if nil == this.Configs {
		this.Configs = make(map[string]*ConfItem)
	}
	if _, ok := this.Configs[key]; ok {
		return fmt.Errorf("config key:%v already exists", key)
	} else {
		this.Configs[key] = item
	}
	return nil
}

func (this *ConfigMgr) Load() error {
	err := error(nil)
	for k, v := range this.Configs {
		c, e := ioutil.ReadFile(v.Path)
		if e != nil {
			err = e
			break
		}
		if e := yaml.Unmarshal(c, v.Data); e != nil {
			err = e
			break
		}
		xlog.Debug("load config: %v, %v, ok", k, v.Path)
	}
	return err
}
