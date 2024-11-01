package main

import (
	"sync"
)

func main() {
	n := 10000
	myMap := make(map[string]int, 1)
	myMap["a"] = 0
	var wg sync.WaitGroup
	var mu sync.Mutex
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			mu.Lock() //without the lock we get fatal error: concurrent map writes at runtime
			defer mu.Unlock()
			myMap["a"] = i
			wg.Done()
		}()
	}
	wg.Wait()
}
