/*
	Package waiter provides some interface, with which we can wait some routines at the same time simply.
	We can wait a routine for it's completion in three lines, eg:
		waiter := NewCtxWaiter()
		// waiter.SetTimeout(time.Duration(time.Second * 10))
		waiter.NewRoutine(realDo)
		waiter.WaitAll()
*/
package waiter
