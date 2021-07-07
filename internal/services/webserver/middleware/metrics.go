package middleware

import (
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/zekroTJA/shinpuru/internal/services/metrics"
)

func NewMetrics() fiber.Handler {
	// Needing this mutex because otherwise it would result
	// in a hash collission for some reason.
	var mtx sync.Mutex

	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		status := c.Response().StatusCode()

		if err != nil {
			if ferr, ok := err.(*fiber.Error); ok {
				status = ferr.Code
			}
		}

		elapsed := float64(time.Since(start).Nanoseconds()) / 1000000000
		istatus := strconv.Itoa(status)

		mtx.Lock()
		defer mtx.Unlock()

		metrics.RestapiRequests.
			WithLabelValues(c.Method(), istatus).
			Inc()

		metrics.RestapiRequestTimes.
			WithLabelValues(c.Method(), istatus).
			Observe(elapsed)

		return err
	}
}
