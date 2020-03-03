package dbs

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
)

type RedisPools struct {
	pools map[string]*redis.Pool
}

func (this *RedisPools) Init(conf *RedisConf) error {
	idleTimeout, _ := time.ParseDuration(conf.Exts.IdleTimeout)
	maxConnLife, _ := time.ParseDuration(conf.Exts.MaxConnLife)
	connTimeout, _ := time.ParseDuration(conf.Exts.ConnTimeout)
	testOnBorrow := conf.Exts.TestOnBorrow
	maxActive := conf.Exts.MaxActive
	maxIdle := conf.Exts.MaxIdle
	host := conf.Host
	passwd := conf.Passwd
	this.pools = make(map[string]*redis.Pool)
	for _, v := range conf.Dbs {
		pool := newPool(host, passwd, v.Db, maxIdle, maxActive, idleTimeout, maxConnLife, connTimeout, testOnBorrow)
		if pool == nil {
			return fmt.Errorf("invalid redis pool: %v", v.Key)
		}
		this.pools[v.Key] = pool
	}
	return nil
}

func (this *RedisPools) Conn(key string) redis.Conn {
	if p, ok := this.pools[key]; ok {
		return p.Get()
	}
	return nil
}

func (this *RedisPools) PubConn(key string) redis.Conn {
	return this.Conn(key)
}

func (this *RedisPools) SubConn(key string) redis.PubSubConn {
	return redis.PubSubConn{Conn: this.Conn(key)}
}
