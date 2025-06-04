package main

import (
	"fmt"
	"time"
)

func main() {
	var sum float64 = 0
	start := time.Now()
	for i := float64(0); i < 1_000_000_000; i++ {
		sum += i
	}
	elapsed := time.Since(start)
	fmt.Printf("Go Loop Benchmark: %.5fms\n", float64(elapsed.Nanoseconds())/1_000_000)
	fmt.Printf("Sum: %f\n", sum)
}
