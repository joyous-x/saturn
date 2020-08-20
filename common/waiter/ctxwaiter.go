package waiter

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// CtxWaiter an implemention of IFlexWaiter with context
type CtxWaiter struct {
	ctx       context.Context
	fnCancel  context.CancelFunc
	waitGroup *sync.WaitGroup
}

// SetTimeout ...
func (c *CtxWaiter) SetTimeout(duration time.Duration) {
	if c.ctx != nil {

	}
	c.ctx, c.fnCancel = context.WithTimeout(context.Background(), duration)
}

// NewRoutine ...
func (c *CtxWaiter) NewRoutine(runnable func()) {
	c.waitGroup.Add(1)
	go c.newRoutine(runnable)
}

// WaitAll ...
func (c *CtxWaiter) WaitAll() {
	c.waitGroup.Wait()
}

func (c *CtxWaiter) run(fn func(), chOut chan int) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("RoutineTree run panic:", err)
		}
	}()
	fn()
	chOut <- 0
}

func (c *CtxWaiter) newRoutine(fn func()) {
	if nil == c.waitGroup {
		go fn()
		return
	}

	if c.ctx == nil {
		c.ctx = context.Background()
	}

	chOut := make(chan int)
	go c.run(fn, chOut)

	select {
	case <-c.ctx.Done():
	case <-chOut:
	}
	c.waitGroup.Done()
}

// NewCtxWaiter ...
func NewCtxWaiter() IFlexWaiter {
	waiter := &CtxWaiter{}
	waiter.waitGroup = &sync.WaitGroup{}
	return waiter
}
