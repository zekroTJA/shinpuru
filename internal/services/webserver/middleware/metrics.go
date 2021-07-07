package middleware

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/zekroTJA/shinpuru/internal/services/metrics"
)

func NewMetrics() fiber.Handler {
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
		method := string(c.Context().Method())

		metrics.RestapiRequests.
			WithLabelValues(method, istatus).
			Inc()

		metrics.RestapiRequestTimes.
			WithLabelValues(method, istatus).
			Observe(elapsed)

		return err
	}
}
