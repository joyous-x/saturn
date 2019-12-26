package jredis

import (
	"github.com/joyous-x/saturn/common/xlog"
	"github.com/gomodule/redigo/redis"
	"time"
)

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
