package auth

import "github.com/gofiber/fiber/v2"

type Middleware interface {
	Handle(ctx *fiber.Ctx) error
}
