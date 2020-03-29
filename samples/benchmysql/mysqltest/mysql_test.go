package mysqltest

import (
	"testing"
	"fmt"
	"sync"
)

func Test_Insert2Cache(t *testing.T){
	count := 1000
	var wg sync.WaitGroup
	wg.Add(count)

	for m := 0; m < count; m++ {
		go func(base int){
			defer wg.Done()
			db, _ := GetCacheDB()
			for i:=0; i<100; i++{
				fmt.Printf("current %v", base * 1000 + i)
				Insert2Cache(base * 1000 + i, db)
			}
		}(m)
	}
	wg.Wait()
}

// go test -bench=.  -run=none