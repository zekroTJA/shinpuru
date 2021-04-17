package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/core/webserver/auth"
	"github.com/zekroTJA/shinpuru/internal/core/webserver/v1/models"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	discordoauth "github.com/zekroTJA/shinpuru/pkg/discordoauth/v2"
)

type AuthController struct {
	discordOAuth *discordoauth.DiscordOAuth
	rth          auth.RefreshTokenHandler
	ath          auth.AccessTokenHandler
	authMw       auth.Middleware
}

func (c *AuthController) Setup(container di.Container, router fiber.Router) {
	c.discordOAuth = container.Get(static.DiDiscordOAuthModule).(*discordoauth.DiscordOAuth)
	c.rth = container.Get(static.DiAuthRefreshTokenHandler).(auth.RefreshTokenHandler)
	c.ath = container.Get(static.DiAuthAccessTokenHandler).(auth.AccessTokenHandler)
	c.authMw = container.Get(static.DiAuthMiddleware).(auth.Middleware)

	router.Get("/login", c.discordOAuth.HandlerInit)
	router.Get("/oauthcallback", c.discordOAuth.HandlerCallback)
	router.Post("/accesstoken", c.postAccessToken)
	router.Get("/check", c.authMw.Handle, c.getCheck)
	router.Post("/logout", c.authMw.Handle, c.postLogout)
}

func (c *AuthController) postAccessToken(ctx *fiber.Ctx) error {
	refreshToken := ctx.Cookies(static.RefreshTokenCookieName)
	if refreshToken == "" {
		return fiber.ErrUnauthorized
	}

	ident, err := c.rth.ValidateRefreshToken(refreshToken)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		util.Log.Error("WEBSERVER :: failed validating refresh token:", err)
	}
	if ident == "" {
		return fiber.ErrUnauthorized
	}

	token, expires, err := c.ath.GetAccessToken(ident)
	if err != nil {
		return err
	}

	return ctx.JSON(&models.AccessTokenResponse{
		Token:   token,
		Expires: expires,
	})
}

func (c *AuthController) getCheck(ctx *fiber.Ctx) error {
	return ctx.JSON(models.Ok)
}

func (c *AuthController) postLogout(ctx *fiber.Ctx) error {
	uid := ctx.Locals("uid").(string)

	err := c.rth.RevokeToken(uid)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	ctx.ClearCookie(static.RefreshTokenCookieName)

	return ctx.JSON(models.Ok)
}
