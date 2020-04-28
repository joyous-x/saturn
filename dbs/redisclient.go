package dbs

import (
	"fmt"
	"time"

	"github.com/joyous-x/saturn/common/xlog"
	"gopkg.in/redis.v5"
)

// RedisClients clients using "gopkg.in/redis.v5"
type RedisClients struct {
	pools map[string]*redis.Client
}

// Init ...
func (r *RedisClients) Init(conf *RedisConf) error {
	idleTimeout, _ := time.ParseDuration(conf.Exts.IdleTimeout)
	maxConnLife, _ := time.ParseDuration(conf.Exts.MaxConnLife)
	connTimeout, _ := time.ParseDuration(conf.Exts.ConnTimeout)
	maxActive := conf.Exts.MaxActive
	maxIdle := conf.Exts.MaxIdle
	host := conf.Host
	passwd := conf.Passwd
	r.pools = make(map[string]*redis.Client)
	for _, v := range conf.Dbs {
		pool := newRedisClient(host, passwd, v.Db, maxIdle+maxActive, idleTimeout, maxConnLife, connTimeout)
		if pool == nil {
			return fmt.Errorf("invalid redis client: %v", v.Key)
		}
		r.pools[v.Key] = pool
	}
	return nil
}

// Client get *redis.Client
func (r *RedisClients) Client(key string) *redis.Client {
	if p, ok := r.pools[key]; ok {
		return p
	}
	return nil
}

func newRedisClient(host, passwd string, db, poolSize int, idleTimeout, maxConnLife, connTimeout time.Duration) *redis.Client {
	name := fmt.Sprintf("DB<%s/%d>", host, db)
	client := redis.NewClient(&redis.Options{
		Addr:        host,
		Password:    passwd,
		DB:          db,
		IdleTimeout: idleTimeout,
		DialTimeout: connTimeout, // default 5s
		PoolSize:    poolSize,    // default 10 connections
	})
	xlog.Debug("newRedisClient %s ok", name)
	return client
}
