package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	var myInt int64 = 0
	var wg sync.WaitGroup
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			//myInt++ //this line would cause a data race
			atomic.AddInt64(&myInt, 1) //so use this line instead
			wg.Done()
		}()
	}
	wg.Wait()

	//will print 10000 when using atomic
	//will print variable result without atomic
	fmt.Println(myInt) 
}
