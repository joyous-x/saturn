package dbs

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/joyous-x/saturn/common/xlog"
)

// RedisPools redis pools using "github.com/gomodule/redigo/redis"
type RedisPools struct {
	pools map[string]*redis.Pool
}

// Init ...
func (r *RedisPools) Init(conf *RedisConf) error {
	idleTimeout, _ := time.ParseDuration(conf.Exts.IdleTimeout)
	maxConnLife, _ := time.ParseDuration(conf.Exts.MaxConnLife)
	connTimeout, _ := time.ParseDuration(conf.Exts.ConnTimeout)
	testOnBorrow := conf.Exts.TestOnBorrow
	maxActive := conf.Exts.MaxActive
	maxIdle := conf.Exts.MaxIdle
	host := conf.Host
	passwd := conf.Passwd
	r.pools = make(map[string]*redis.Pool)
	for _, v := range conf.Dbs {
		pool := newPool(host, passwd, v.Db, maxIdle, maxActive, idleTimeout, maxConnLife, connTimeout, testOnBorrow)
		if pool == nil {
			return fmt.Errorf("invalid redis pool: %v", v.Key)
		}
		r.pools[v.Key] = pool
	}
	return nil
}

// Conn ...
func (r *RedisPools) Conn(key string) redis.Conn {
	if p, ok := r.pools[key]; ok {
		return p.Get()
	}
	return nil
}

// PubConn ...
func (r *RedisPools) PubConn(key string) redis.Conn {
	return r.Conn(key)
}

// SubConn ...
func (r *RedisPools) SubConn(key string) redis.PubSubConn {
	return redis.PubSubConn{Conn: r.Conn(key)}
}

func newPool(host, passwd string, db, maxIdle, maxActive int, idleTimeout, maxConnLife, connTimeout time.Duration, testOnBorrow string) *redis.Pool {
	pool := &redis.Pool{
		MaxIdle:         maxIdle,
		MaxActive:       maxActive,
		IdleTimeout:     idleTimeout,
		MaxConnLifetime: maxConnLife,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", host, redis.DialDatabase(db), redis.DialConnectTimeout(connTimeout))
			if err != nil {
				xlog.Error("dial redis, host=%v:%v error=%v", host, db, err)
				return nil, err
			}

			if len(passwd) > 0 {
				if _, err := c.Do("AUTH", passwd); err != nil {
					c.Close()
					xlog.Error("auth redis, host=%v:%v error=%v", host, db, err)
					return nil, err
				}
			}

			//> TODO: notify-keyspace-events
			xlog.Info("dial redis ok, host=%v:%v", host, db)
			return c, nil
		},
	}

	if len(testOnBorrow) > 0 {
		pool.TestOnBorrow = func(c redis.Conn, t time.Time) error {
			_, err := c.Do(testOnBorrow)
			return err
		}
	}

	return pool
}
