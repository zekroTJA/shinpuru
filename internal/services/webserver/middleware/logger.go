package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/zekrotja/rogu/log"
)

var mwLog = log.Tagged("WebServer")

// Logger returns a middleware handler to log incoming
// requests on the debug log channel.
func Logger() fiber.Handler {
	return func(ctx *fiber.Ctx) (err error) {
		start := time.Now()

		err = ctx.Next()

		d := time.Since(start)

		entry := mwLog.Debug().Fields(
			"code", ctx.Response().StatusCode(),
			"method", ctx.Method(),
			"duration", d,
			"ip", ctx.IP(),
		)

		if err != nil {
			entry = entry.Err(err)
		}

		entry.Msgf("%-5s %s", ctx.Method(), ctx.Path())

		return
	}
}
