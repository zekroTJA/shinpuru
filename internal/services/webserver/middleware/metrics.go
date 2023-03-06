package middleware

import (
	"regexp"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/zekroTJA/shinpuru/internal/services/metrics"
	"github.com/zekrotja/rogu/log"
)

type MetricsOptions struct {
	IgnorePatterns []string
}

func NewMetrics(opts ...MetricsOptions) fiber.Handler {
	var opt MetricsOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	ignore := make([]*regexp.Regexp, 0, len(opt.IgnorePatterns))
	for _, pattern := range opt.IgnorePatterns {
		rx, err := regexp.Compile(pattern)
		if err == nil {
			ignore = append(ignore, rx)
		} else {
			log.Error().Tag("Metrics").
				Err(err).
				Field("pattern", pattern).
				Msg("Failed parsing regex ignore pattern")
		}
	}

	return func(c *fiber.Ctx) error {
		for _, rx := range ignore {
			if rx.MatchString(c.Path()) {
				log.Debug().Tag("Metrics").
					Field("path", c.Path()).
					Msg("Reqeust has been ignored")
				continue
			}
		}

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
