package waiter

import (
	"sync"
	"time"
)

// IFlexWaiter 一种易用的等待routine完成接口
type IFlexWaiter interface {
	// SetTimeout set timeout for CtxWaiter.
	// NOTE:
	//   This should be called before NewRoutine.
	SetTimeout(duration time.Duration)
	// NewRoutine generate a new go routine which run the function fn
	NewRoutine(runnable func())
	// WaitAll wait all routines for completion or timeout
	WaitAll()
}

// NewFlexWaiter new an instance of IFlexWaiter
func NewFlexWaiter() IFlexWaiter {
	waiter := &CtxWaiter{}
	waiter.waitGroup = &sync.WaitGroup{}
	return waiter
}
