package config

import (
	"github.com/joyous-x/saturn/common/gins"
	"github.com/joyous-x/saturn/common/xlog"
	"github.com/joyous-x/saturn/dbs"
)

const (
	// CfgKeyCommon config item key for project
	CfgKeyCommon = "ckey_common"
)

// ComConfig ...
type ComConfig struct {
	ServerConfs []*gins.ServerConfig     `yaml:"httpserver"`
	XLog        xlog.XLogConf            `yaml:"xlog"`
	Redis       []dbs.RedisConf          `yaml:"db.redis"`
	Mysql       []dbs.MysqlConf          `yaml:"db.mysql"`
	WxApps      map[string]WxMiniAppInfo `yaml:"wxminiapp"`
}

// WxMiniAppInfo ...
type WxMiniAppInfo struct {
	AppID          string `yaml:"app_id"`
	AppName        string `yaml:"app_name"`
	AppSecret      string `yaml:"app_secret"`
	EncodingAESKey string `yaml:"app_aeskey"`
	Token          string `yaml:"token"`
}

// GetComConfig ...
func (m *MgrCenter) GetComConfig() *ComConfig {
	v, ok := m.Configs[CfgKeyCommon]
	if !ok {
		return &ComConfig{}
	}
	return v.Data.(*ComConfig)
}
