package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type HashResponse struct {
	Hash      string `json:"hash"`
	Timestamp int64  `json:"timestamp"`
	Source    string `json:"source"`
}

var hashSeed = []byte("benchmark-test-data")

var (
	goroutineCount = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "go_goroutines",
		Help: "Number of goroutines that currently exist.",
	})
)

func init() {
	prometheus.MustRegister(goroutineCount)
}

func main() {
	// Start Prometheus metrics server on port 2112
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		// Update goroutine count every second
		go func() {
			for {
				goroutineCount.Set(float64(runtime.NumGoroutine()))
				time.Sleep(1 * time.Second)
			}
		}()
		fmt.Println("Prometheus metrics server starting on :2112")
		if err := http.ListenAndServe(":2112", nil); err != nil {
			panic(err)
		}
	}()

	app := fiber.New()

	app.Get("/hash", hashHandler)
	app.Get("/health", healthHandler)

	fmt.Println("Fiber server starting on :3000")
	if err := app.Listen(":3000"); err != nil {
		panic(err)
	}
}

func hashHandler(c fiber.Ctx) error {
	// Keep the hash loop allocation-free:
	// - sha256.Sum256 returns a fixed-size array
	// - input is reused without building new slices per iteration
	sum := sha256.Sum256(hashSeed)
	for i := 1; i < 100; i++ {
		sum = sha256.Sum256(sum[:])
	}

	response := HashResponse{
		Hash:      hex.EncodeToString(sum[:]),
		Timestamp: time.Now().UnixMilli(),
		Source:    "go-fiber",
	}

	return c.JSON(response)
}

func healthHandler(c fiber.Ctx) error {
	return c.SendStatus(http.StatusOK)
}
