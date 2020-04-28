package crontab

import (
	"context"
	"fmt"
	"reflect"
	"runtime"
	"time"

	"github.com/joyous-x/saturn/common/xlog"
	"github.com/robfig/cron/v3"
	"gopkg.in/redis.v5"
)

// Crontab the warpper for cron.Cron, we can use it simply
type Crontab struct {
	cron        *cron.Cron
	redisClient *redis.Client
}

// Init init crontab
func (c *Crontab) Init(redisClient *redis.Client, opts ...cron.Option) error {
	optLocation := cron.WithLocation(time.Local)
	opts = append(opts, optLocation)
	c.cron = cron.New(opts...)
	c.redisClient = redisClient
	return nil
}

// Start start the cron
func (c *Crontab) Start() {
	c.cron.Start()
}

// Stop stop the cron
func (c *Crontab) Stop() context.Context {
	return c.cron.Stop()
}

// Cron get the raw cron.Cron
func (c *Crontab) Cron() *cron.Cron {
	return c.cron
}

// AddFuncOneInstance add function which only run in an (vm) environment when triggered
//                if lockSec <= 0, we will set lockSec to default value: 60
func (c *Crontab) AddFuncOneInstance(spec string, lockSec int, cmd func()) (cron.EntryID, error) {
	return c.cron.AddFunc(spec, c.OneInstanceCmd(spec, lockSec, cmd))
}

// OneInstanceCmd construct a new cmd which only run in an (vm) environment when triggered
//                if lockSec <= 0, we will set lockSec to default value: 60
func (c *Crontab) OneInstanceCmd(name string, lockSec int, cmd func()) func() {
	cmdName := runtime.FuncForPC(reflect.ValueOf(cmd).Pointer()).Name()
	if lockSec <= 0 {
		lockSec = 60
	}
	newCmd := func() {
		redisKey := fmt.Sprintf("cron_lock_%s_%s", name, cmdName)
		if succ, _ := c.tryLock(redisKey, "1", lockSec); !succ {
			return
		}
		cmd()
	}
	return newCmd
}

func (c *Crontab) tryLock(key, value string, lockSec int) (bool, error) {
	return c.redisClient.SetNX(key, value, time.Second*time.Duration(lockSec)).Result()
}

func (c *Crontab) unlock(key string) bool {
	delCnt, err := c.redisClient.Del(key).Result()
	if err != nil {
		xlog.Error("Crontab unlock %v error=%v", key, err)
		return false
	}
	xlog.Info("Crontab unlock %v: cnt=%v", key, delCnt)
	return true
}
