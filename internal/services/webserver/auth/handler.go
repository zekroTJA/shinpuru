package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zekroTJA/shinpuru/pkg/discordoauth/v2"
)

// RequestHandler provides fiber endpoints and handlers
// to authenticate users via an OAuth2 interface.
type RequestHandler interface {

	// LoginFailedHandler is called when either the
	// user authentication failed or something went
	// wrong during the authentication process.
	LoginFailedHandler(ctx *fiber.Ctx, status int, msg string) error

	// BindRefreshToken generates a refresh token
	// and binds it as cookie to the passed
	// context.
	BindRefreshToken(ctx *fiber.Ctx, uid string) error

	// LoginSuccessHandler is called when the
	// authentication process was successful.
	//
	// The function is getting passed the ident of
	// the authenticated user.
	LoginSuccessHandler(ctx *fiber.Ctx, res discordoauth.SuccessResult) error

	// LogoutHandler is called when the user
	// wants to log out.
	LogoutHandler(ctx *fiber.Ctx) error
}
