package controllers

import (
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/services/kvcache"
	"github.com/zekroTJA/shinpuru/internal/services/webserver/v1/models"
	"github.com/zekroTJA/shinpuru/internal/services/webserver/wsutil"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/fetch"
	"github.com/zekrotja/dgrs"
)

type SearchController struct {
	session *discordgo.Session
	st      *dgrs.State
	kv      kvcache.Provider
}

func (c *SearchController) Setup(container di.Container, router fiber.Router) {
	c.session = container.Get(static.DiDiscordSession).(*discordgo.Session)
	c.st = container.Get(static.DiState).(*dgrs.State)
	c.kv = container.Get(static.DiKVCache).(kvcache.Provider)

	router.Get("", c.getSearch)
}

// @Summary Global Search
// @Description Search through guilds and members by ID, name or displayname.
// @Tags Search
// @Accept json
// @Produce json
// @Param query query string true "The search query (either ID, name or displayname)."
// @Param limit query int false "The maximum amount of result items (per group)." default(50) minimum(1) maximum(100)
// @Success 200 {object} models.SearchResult
// @Failure 400 {object} models.Error
// @Failure 401 {object} models.Error
// @Router /search [get]
func (c *SearchController) getSearch(ctx *fiber.Ctx) (err error) {
	uid := ctx.Locals("uid").(string)
	query := strings.ToLower(ctx.Query("query"))
	limit, err := wsutil.GetQueryInt(ctx, "limit", 50, 1, 100)
	if err != nil {
		return
	}

	if query == "" {
		return fiber.NewError(fiber.StatusBadRequest, "query must be set")
	}

	kvKey := "SEARCH:GUILDS:" + uid
	guilds, ok := c.kv.Get(kvKey).([]*discordgo.Guild)
	if !ok {
		var guildIDs []string
		guildIDs, err = c.st.UserGuilds(uid)
		if err != nil {
			return
		}

		guilds = make([]*discordgo.Guild, len(guildIDs))
		for i, id := range guildIDs {
			if guilds[i], err = c.st.Guild(id, true); err != nil {
				return
			}
		}
		c.kv.Set(kvKey, guilds, 5*time.Minute)
	}

	sr := &models.SearchResult{
		Guilds:  make([]*models.GuildReduced, 0),
		Members: make([]*models.Member, 0),
	}

	var iG, iM int
	for _, g := range guilds {
		if iG < limit {
			for _, f := range fetch.GuildCheckFuncs {
				if f(g, query) {
					sr.Guilds = append(sr.Guilds, models.GuildReducedFromGuild(g))
					iG++
					break
				}
			}
		}
		for _, m := range g.Members {
			if iM == limit {
				break
			}
			for _, f := range fetch.MemberCheckFuncs {
				if f(m, query) {
					sr.Members = append(sr.Members, models.MemberFromMember(m))
					iM++
					break
				}
			}
		}
	}

	return ctx.JSON(sr)
}
