package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	discordoauth "github.com/zekroTJA/shinpuru/pkg/discordoauth/v2"
)

type AuthController struct {
	discordOAuth *discordoauth.DiscordOAuth
}

func (c *AuthController) Setup(container di.Container, router fiber.Router) {
	c.discordOAuth = container.Get(static.DiDiscordOAuthModule).(*discordoauth.DiscordOAuth)

	router.Get("/login", c.discordOAuth.HandlerInit)
	router.Get("/oauthcallback", c.discordOAuth.HandlerCallback)
}
