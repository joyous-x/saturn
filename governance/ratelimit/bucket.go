package ratelimit

import (
	"time"
	"fmt"
)

// 限流策略：
//     令牌桶模式
//     当令牌桶的 capacity 为 1 时，退化为 滴漏模式
type RateLimitOnBucket struct {
	capacity    int 
	rate        int
	bucket      chan struct{}
	ticker      *time.Ticker
}

func NewRateLimitOnBucket(capacity, rate int) *RateLimitOnBucket {
	tmp := &RateLimitOnBucket{
		capacity: capacity,
		rate  : rate,
		bucket: make(chan struct{}, capacity),
		ticker: nil,
	}
	return tmp
}

func (this *RateLimitOnBucket) IsAllow() bool {
	ret := false
	select {
	case <- this.bucket:
		ret = true
	default:
	}
	return ret
}

func (this *RateLimitOnBucket) Run() error {
	if this.ticker != nil {
		return fmt.Errorf("already running")
	}
	if this.rate <= 0 {
		return fmt.Errorf("invalid rate(%v)", this.rate)
	}
	intervalNs := int(1000000000 / this.rate)
	this.ticker = time.NewTicker(time.Duration(intervalNs) * time.Nanosecond)
	go func() {
		defer this.Stop()
		for {
			select {
			case <- this.ticker.C:
				select {
				case this.bucket <- struct{}{}:
					// xlog.Debugf("--- current bucket status: %v %v \n", len(this.bucket), cap(this.bucket))
				default:
				}
			}
		}
	}()
	return nil
}

func (this *RateLimitOnBucket) Stop() error {
	if this.ticker != nil {
		this.ticker.Stop()
		this.ticker = nil
	}	
	return nil
}