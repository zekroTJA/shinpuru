package controllers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekrotja/dgrs"
)

type InviteController struct {
	session *discordgo.Session
	st      *dgrs.State
}

func (c *InviteController) Setup(container di.Container, router fiber.Router) {
	c.session = container.Get(static.DiDiscordSession).(*discordgo.Session)
	c.st = container.Get(static.DiState).(*dgrs.State)

	router.Get("", c.getInvite)
}

func (c *InviteController) getInvite(ctx *fiber.Ctx) error {
	self, err := c.st.SelfUser()
	if err != nil {
		return err
	}
	return ctx.Redirect(util.GetInviteLink(self.ID))
}
