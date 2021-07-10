package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// Logger returns a middleware handler to log incoming
// requests on the debug log channel.
func Logger() fiber.Handler {
	return func(ctx *fiber.Ctx) (err error) {
		start := time.Now()

		err = ctx.Next()

		d := time.Since(start)

		entry := logrus.WithFields(logrus.Fields{
			"code":     ctx.Response().StatusCode(),
			"method":   ctx.Method(),
			"duration": d,
			"ip":       ctx.IP(),
		})

		line := fmt.Sprintf("WS :: %-5s %s", ctx.Method(), ctx.Path())

		if err != nil {
			entry = entry.WithError(err)
		}

		entry.Debug(line)

		return
	}
}
