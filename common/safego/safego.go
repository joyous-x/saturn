package safego

import (
	"fmt"
	"os"
	"runtime"
)

// RecoverHandler when panic happen, the hander will be called
type RecoverHandler func(err interface{})

var defaultRecoverHandler = func(err interface{}) {
	stackData := make([]byte, 2048)
	stackSize := runtime.Stack(stackData, false)
	fmt.Fprintf(os.Stderr, "=> default panic handler: %s \n %s \n", err, string(stackData[:stackSize]))
}

// Go run the f() with a goroutine and keep it away from panic
// it will use defaultRecoverHandler if argument handler is nil
func Go(f func(), handlers ...RecoverHandler) {
	curFunc := func(args ...interface{}) {
		f()
	}
	GoWith(curFunc, handlers...)()
}

// GoWith run the f(...interface{}) with a goroutine and keep it away from panic
// eg.
//     GoWith(func(args ...interface{}) { }, func(err interface{}) { })(1, 0)
func GoWith(f func(args ...interface{}), handlers ...RecoverHandler) func(args ...interface{}) {
	handler := defaultRecoverHandler
	if len(handlers) > 0 {
		handler = handlers[0]
	}

	curFunc := f
	return func(args ...interface{}) {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					handler(r)
				}
			}()

			curFunc(args...)
		}()
	}
}
