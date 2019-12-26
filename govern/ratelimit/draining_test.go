package ratelimit

import (
	"testing"
	"time"
)

func Test_Draining(t *testing.T) {
	limiter := NewRateLimitOnDraining(5)
	limiter.Run()
	time.Sleep(1 * time.Second)
	ticker := time.NewTicker(time.Duration(100) * time.Millisecond)

	totalCnt := 10
	allowCnt := 0
	expected := 5

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
		t.Errorf("Test_Draining allow = %v expected = %d", allowCnt, expected)
	}
}
