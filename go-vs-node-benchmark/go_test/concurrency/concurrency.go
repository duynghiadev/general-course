package main

import (
	"fmt"
	"sync"
	"time"
)

func task(wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(100 * time.Millisecond)
}

func main() {
	var wg sync.WaitGroup
	start := time.Now()
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go task(&wg)
	}
	wg.Wait()
	elapsed := time.Since(start)
	fmt.Printf("Go Concurrency Time: %.4fms\n", float64(elapsed.Nanoseconds())/1_000_000)
}
