package controllers

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/gofiber/fiber/v2"
	"github.com/narqo/go-badge"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/services/kvcache"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekrotja/dgrs"
)

type InviteController struct {
	session *discordgo.Session
	st      *dgrs.State
	kv      kvcache.Provider
}

func (c *InviteController) Setup(container di.Container, router fiber.Router) {
	c.session = container.Get(static.DiDiscordSession).(*discordgo.Session)
	c.st = container.Get(static.DiState).(*dgrs.State)
	c.kv = container.Get(static.DiKVCache).(kvcache.Provider)

	router.Get("", c.getInvite)
	router.Get("badge.svg", c.getInviteSvg)
}

func (c *InviteController) getInvite(ctx *fiber.Ctx) error {
	self, err := c.st.SelfUser()
	if err != nil {
		return err
	}
	return ctx.Redirect(util.GetInviteLink(self.ID))
}

func (c *InviteController) getInviteSvg(ctx *fiber.Ctx) error {
	title := ctx.Query("title", "invite")
	color := ctx.Query("color", badge.ColorGreen.String())

	nGuilds, err := c.getNGuilds()
	if err != nil {
		return err
	}

	ctx.Response().Header.SetContentType("image/svg+xml")
	err = badge.Render(
		title,
		fmt.Sprintf("%d %s", nGuilds, util.Pluralize(nGuilds, "guild")),
		badge.Color(color), ctx)
	if err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusOK)
}

// --- HELPERS ---

func (c *InviteController) getNGuilds() (n int, err error) {
	const nGuildsKey = "GLOBAL:N_GUILDS"

	n, ok := c.kv.Get(nGuildsKey).(int)
	if !ok {
		var guilds []*discordgo.Guild
		if guilds, err = c.st.Guilds(); err != nil {
			return
		}
		n = len(guilds)
		c.kv.Set(nGuildsKey, n, 10*time.Minute)
	}

	return
}
