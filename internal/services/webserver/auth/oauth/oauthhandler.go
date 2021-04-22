package oauth

import "github.com/gofiber/fiber/v2"

type OAuthHandler interface {
	LoginFailedHandler(ctx *fiber.Ctx, status int, msg string) error
	LoginSuccessHandler(ctx *fiber.Ctx, uid string) error
	LogoutHandler(ctx *fiber.Ctx) error
}
