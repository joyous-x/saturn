package crontab

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/assert"
	"gopkg.in/redis.v5"
)

func newRedis() (*miniredis.Miniredis, *redis.Client) {
	mredis, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	client := redis.NewClient(&redis.Options{
		Addr: mredis.Addr(),
	})
	return mredis, client
}

func Test_Crontab(t *testing.T) {
	miniRedis, redisClient := newRedis()
	defer miniRedis.Close()

	var err error
	instCron := &Crontab{}
	instCron.Init(redisClient, cron.WithSeconds())

	cntA := 0
	_, err = instCron.Cron().AddFunc("*/2 * * * * *", func() {
		cntA = cntA + 1
		t.Logf("test for AddFunc: counter = %v", cntA)
	})
	if err != nil {
		t.Errorf("cron add func error: %v", err)
	}

	cntB, specB := 0, "*/2 * * * * *"
	cmdB := func() {
		cntB = cntB + 1
		t.Logf("test for AddFuncOneInstance: counter=%v", cntB)

		//> miniRedis 的 setnx 函数无效，所以这里手动释放
		redisClient.FlushAll()
		time.Sleep(time.Millisecond * 500)
	}
	_, err = instCron.AddFuncOneInstance(specB, 1, cmdB)
	if err != nil {
		t.Errorf("cron add func error: %v", err)
	}

	_, err = instCron.AddFuncOneInstance(specB, 1, cmdB)
	if err != nil {
		t.Errorf("cron add func error: %v", err)
	}

	instCron.Start()
	time.Sleep(time.Second * 10)

	assert.Equal(t, cntA, 5, "assert for AddFunc")
	assert.Equal(t, cntB, 5, "assert for AddFuncOneInstance")

	miniRedis.FastForward(time.Second * 3600)
}
