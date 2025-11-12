package main

import (
	"math"
	"runtime"

	"github.com/gofiber/fiber/v2"
)

func heavyComputation() {
	// More CPU-intensive work
	for i := 0; i < 100_000_000; i++ {
		_ = math.Sin(float64(i)) * math.Sqrt(float64(i))
	}
}

func main() {
	// Use a single CPU core
	runtime.GOMAXPROCS(1)

	app := fiber.New()

	app.Get("/api/hello", func(c *fiber.Ctx) error {
		heavyComputation()
		return c.JSON(fiber.Map{
			"message": "Hello from Fiber!",
		})
	})

	app.Listen(":3000")
}
