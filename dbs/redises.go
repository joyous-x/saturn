package dbs

import (
	"fmt"
	"time"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/joyous-x/saturn/common/xlog"
	"gopkg.in/redis.v5"
)

// Redises ...
type Redises struct {
	useClient bool
	// RedisPools redis pools using "github.com/gomodule/redigo/redis"
	pools map[string]*redigo.Pool
	// RedisClients clients using "gopkg.in/redis.v5"
	clients map[string]*redis.Client
}

// Init ...
func (r *Redises) Init(confs []RedisConf, useClient bool) error {
	r.useClient = useClient
	r.pools = make(map[string]*redigo.Pool)
	r.clients = make(map[string]*redis.Client)
	for _, conf := range confs {
		idleTimeout, _ := time.ParseDuration(conf.Exts.IdleTimeout)
		maxConnLife, _ := time.ParseDuration(conf.Exts.MaxConnLife)
		connTimeout, _ := time.ParseDuration(conf.Exts.ConnTimeout)
		maxActive := conf.Exts.MaxActive
		maxIdle := conf.Exts.MaxIdle
		host := conf.Host
		passwd := conf.Passwd
		if useClient {
			client := newRedisClient(host, passwd, conf.Db, maxIdle+maxActive, idleTimeout, maxConnLife, connTimeout)
			if client != nil {
				r.clients[conf.Name] = client
			} else {
				return fmt.Errorf("invalid redis client: %v", conf.Name)
			}
		} else {
			pool := newPool(host, passwd, conf.Db, maxIdle, maxActive, idleTimeout, maxConnLife, connTimeout, conf.Exts.TestOnBorrow)
			if pool != nil {
				r.pools[conf.Name] = pool
			} else {
				return fmt.Errorf("invalid redis pool: %v", conf.Name)
			}
		}
	}
	return nil
}

// Ping ...
func (r *Redises) Ping(key string) (string, error) {
	if r.useClient {
		if redisCon := r.Conn(key); redisCon != nil {
			defer redisCon.Close()
			return redigo.String(redisCon.Do("ping"))
		}
	} else {
		if redisClient := r.Client(key); redisClient != nil {
			return redisClient.Ping().Result()
		}
	}

	return "", fmt.Errorf("invalid key(%s)", key)
}

// Client get *redis.Client (redis.v5)
func (r *Redises) Client(key string) *redis.Client {
	if p, ok := r.clients[key]; ok {
		return p
	}
	return nil
}

// Conn get redis.Conn (redigo)
func (r *Redises) Conn(key string) redigo.Conn {
	if p, ok := r.pools[key]; ok {
		return p.Get()
	}
	return nil
}

// PubConn ...(redigo)
func (r *Redises) PubConn(key string) redigo.Conn {
	return r.Conn(key)
}

// SubConn ...(redigo)
func (r *Redises) SubConn(key string) redigo.PubSubConn {
	return redigo.PubSubConn{Conn: r.Conn(key)}
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

func newPool(host, passwd string, db, maxIdle, maxActive int, idleTimeout, maxConnLife, connTimeout time.Duration, testOnBorrow string) *redigo.Pool {
	pool := &redigo.Pool{
		MaxIdle:         maxIdle,
		MaxActive:       maxActive,
		IdleTimeout:     idleTimeout,
		MaxConnLifetime: maxConnLife,
		Dial: func() (redigo.Conn, error) {
			c, err := redigo.Dial("tcp", host, redigo.DialDatabase(db), redigo.DialConnectTimeout(connTimeout))
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
		pool.TestOnBorrow = func(c redigo.Conn, t time.Time) error {
			_, err := c.Do(testOnBorrow)
			return err
		}
	}

	return pool
}
