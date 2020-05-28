package wechat

import (
	"github.com/gomodule/redigo/redis"
)

// RedisFetcher get redis connection
type RedisFetcher func() redis.Conn

type WxConfig struct {
	AppName        string
	AppID          string
	AppSecret      string
	EncodingAESKey string
	PubAccToken    string
	PayMchID       string //支付 - 商户 ID
	PayNotifyURL   string //支付 - 接受微信支付结果通知的接口地址
	PayKey         string //支付 - 商户后台设置的支付 key
	RedisFetcher   RedisFetcher
}

func (s *WxConfig) SetRedisFetcher(fetcher RedisFetcher) {
	s.RedisFetcher = fetcher
}
