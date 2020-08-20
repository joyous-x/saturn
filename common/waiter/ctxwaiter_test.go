package waiter

import (
	"fmt"
	"testing"
	"time"
)

func Test_CtxWaiter(t *testing.T) {
	realDo := func() {
		for i := 0; i < 100; i++ {
			fmt.Printf("currint id = %d \n", i)
			time.Sleep(time.Duration(time.Second))
		}
	}

	tree := NewFlexWaiter()
	tree.SetTimeout(time.Duration(time.Second * 10))
	tree.NewRoutine(realDo)
	tree.NewRoutine(realDo)
	tree.NewRoutine(realDo)
	tree.WaitAll()
}
