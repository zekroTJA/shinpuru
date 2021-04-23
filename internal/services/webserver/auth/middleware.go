package auth

import "github.com/gofiber/fiber/v2"

// Middleware provides an authorization
// fiber middleware
type Middleware interface {
	Handle(ctx *fiber.Ctx) error
}
