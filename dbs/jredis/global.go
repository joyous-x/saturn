package jredis

import (
	"github.com/joyous-x/saturn/common/xlog"
	"github.com/joyous-x/saturn/dbs"
	"sync"
)

var g_redisPools *RedisPools
var g_redis_once sync.Once

func GlobalInst(redisConfs ...dbs.RedisConf) *RedisPools {
	if len(redisConfs) > 0 {
		g_redis_once.Do(func() {
			redisConfData := redisConfs[0]
			redisPools := &RedisPools{}
			if err := redisPools.Init(&redisConfData); err != nil {
				xlog.Error("init redispools err:%v", err)
			} else {
				g_redisPools = redisPools
			}
		})
	}
	return g_redisPools
}
