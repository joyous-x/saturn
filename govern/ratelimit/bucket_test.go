package ratelimit

import (
	"testing"
	"time"
)

func Test_Bucket(t *testing.T) {
	limiter := NewRateLimitOnBucket(10, 5)
	limiter.Run()
	time.Sleep(1 * time.Second)
	ticker := time.NewTicker(time.Duration(50) * time.Millisecond)

	totalCnt := 10
	allowCnt := 0
	expected := 7

	for i := 0; i < totalCnt; i++ {
		select {
		case <-ticker.C:
			allow := limiter.IsAllow()
			if allow {
				allowCnt += 1
			}
			t.Logf("---- %v \n", allow)
		}
	}

	if expected != allowCnt {
		t.Errorf("Test_Bucket allow = %v expected = %d", allowCnt, expected)
	}
}
