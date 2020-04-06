package config

import (
	"github.com/joyous-x/saturn/common/xlog"
	"github.com/joyous-x/saturn/dbs"
	"github.com/joyous-x/saturn/common/gins"
)

const (
	CfgKeyDbs  = "ckey_dbs"
	CfgKeyLog  = "ckey_log"
	CfgKeyProj = "ckey_proj"
)

type ConfProj struct {
	HttpConfs []*gins.ServerConfig     `yaml:"httpserver"`
	WxApps    map[string]WxMiniAppInfo `yaml:"wxminiapp"`
}

type WxMiniAppInfo struct {
	AppID          string `yaml:"app_id"`
	AppName        string `yaml:"app_name"`
	AppSecret      string `yaml:"app_secret"`
	EncodingAESKey string `yaml:"app_aeskey"`
	Token          string `yaml:"token"`
}

type ConfDbs struct {
	Redis dbs.RedisConf   `yaml:"redis"`
	Mysql []dbs.MysqlConf `yaml:"mysql"`
}

func (this *ConfigMgr) CfgProj() *ConfProj {
	v, ok := this.Configs[CfgKeyProj]
	if !ok {
		return &ConfProj{}
	}
	return v.Data.(*ConfProj)
}

func (this *ConfigMgr) CfgDbs() *ConfDbs {
	v, ok := this.Configs[CfgKeyDbs]
	if !ok {
		return &ConfDbs{}
	}
	return v.Data.(*ConfDbs)
}

func (this *ConfigMgr) CfgLog() *xlog.XLogConf {
	v, ok := this.Configs[CfgKeyLog]
	if !ok {
		return &xlog.XLogConf{}
	}
	return v.Data.(*xlog.XLogConf)
}
