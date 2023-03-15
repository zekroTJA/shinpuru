package controllers

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/webserver/v1/models"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/hashutil"
	"github.com/zekrotja/dgrs"
)

type PublicController struct {
	session *discordgo.Session
	db      database.Database
	st      *dgrs.State
}

func (c *PublicController) Setup(container di.Container, router fiber.Router) {
	c.session = container.Get(static.DiDiscordSession).(*discordgo.Session)
	c.db = container.Get(static.DiDatabase).(database.Database)
	c.st = container.Get(static.DiState).(*dgrs.State)

	router.Get("/guilds/:guildid", c.getGuild)
}

// @Summary Get Public Guild
// @Description Returns public guild information, if enabled by guild config.
// @Tags Public
// @Accept json
// @Produce json
// @Param id path string true "The Guild ID."
// @Success 200 {object} models.GuildReduced
// @Router /public/guilds/{id} [get]
func (c *PublicController) getGuild(ctx *fiber.Ctx) error {
	guildID := ctx.Params("guildid")

	state, err := c.db.GetGuildAPI(guildID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	if !state.Enabled {
		return fiber.ErrNotFound
	}

	if state.Hydrate().Protected {
		token := c.obtainToken(ctx)
		if token == "" {
			return fiber.ErrNotFound
		}
		ok, err := hashutil.Compare(token, state.TokenHash)
		if err != nil {
			return err
		}
		if !ok {
			return fiber.ErrNotFound
		}
	}

	if state.AllowedOrigins == "" {
		state.AllowedOrigins = "*"
	}

	guild, err := c.st.Guild(guildID)
	if err != nil {
		return err
	}

	ctx.Set("Access-Control-Allow-Origin", state.AllowedOrigins)
	ctx.Set("Access-Control-Allow-Methods", "GET")
	ctx.Set("Access-Control-Allow-Headers", "*")

	gr := models.GuildReducedFromGuild(guild)

	return ctx.JSON(gr)
}

func (c *PublicController) obtainToken(ctx *fiber.Ctx) (token string) {
	token = ctx.Query("token")

	if token == "" {
		split := strings.SplitN(ctx.Get("Authorization"), " ", 2)
		if len(split) < 2 {
			return
		}
		if strings.ToLower(split[0]) != "bearer" {
			return
		}
		token = split[1]
	}

	return
}
