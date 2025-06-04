package main

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

func loopTest() {
	var sum int64 = 0
	for i := int64(1); i <= 1_000_000_000; i++ {
		sum += i
	}
	fmt.Println("Loop test Sum:", sum)
}

func concurrencyTest() {
	ch := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			time.Sleep(100 * time.Millisecond)
			ch <- true
		}()
	}
	for i := 0; i < 10; i++ {
		<-ch
	}
	fmt.Println("Concurrency test done")
}

func main() {
	app := fiber.New()
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	go func() {
		// Uncomment one of the following to test:
		loopTest()
		concurrencyTest()
	}()

	app.Listen(":3001")
}
