package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekrotja/discordgo"
)

type InviteController struct {
	session *discordgo.Session
}

func (c *InviteController) Setup(container di.Container, router fiber.Router) {
	c.session = container.Get(static.DiDiscordSession).(*discordgo.Session)

	router.Get("", c.getInvite)
}

func (c *InviteController) getInvite(ctx *fiber.Ctx) error {
	return ctx.Redirect(util.GetInviteLink(c.session))
}
