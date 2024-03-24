package controllers

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/services/config"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/timeprovider"
	"github.com/zekroTJA/shinpuru/internal/services/webserver/auth"
	_ "github.com/zekroTJA/shinpuru/internal/services/webserver/v1/models" // Import for API documentation
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/discordoauth/v2"
	"github.com/zekroTJA/shinpuru/pkg/onetimeauth/v2"
)

type OTAController struct {
	db      Database
	session Session
	
	cfg          config.Provider
	ota          onetimeauth.OneTimeAuth
	oauthHandler auth.RequestHandler
	tp           timeprovider.Provider
}

func (c *OTAController) Setup(container di.Container, router fiber.Router) {
	c.session = container.Get(static.DiDiscordSession).(*discordgo.Session)
	c.cfg = container.Get(static.DiConfig).(config.Provider)
	c.db = container.Get(static.DiDatabase).(database.Database)
	c.ota = container.Get(static.DiOneTimeAuth).(onetimeauth.OneTimeAuth)
	c.oauthHandler = container.Get(static.DiOAuthHandler).(auth.RequestHandler)
	c.tp = container.Get(static.DiTimeProvider).(timeprovider.Provider)

	router.Get("", c.getOta)
}

// @Summary OTA Login
// @Description Logs in the current browser session by using the passed pre-obtained OTA token.
// @Tags OTA
// @Accept json
// @Produce json
// @Success 200
// @Failure 401 {object} models.Error
// @Router /ota [get]
func (c *OTAController) getOta(ctx *fiber.Ctx) error {
	token := ctx.Query("token")

	if token == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid ota token")
	}

	userID, err := c.ota.ValidateKey(token, "login-via-dm")
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid ota token")
	}

	enabled, err := c.db.GetUserOTAEnabled(userID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	if !enabled {
		return fiber.NewError(fiber.StatusUnauthorized, "ota disabled")
	}

	if ch, err := c.session.UserChannelCreate(userID); err == nil {
		ipaddr := ctx.IP()
		useragent := string(ctx.Context().UserAgent())
		emb := &discordgo.MessageEmbed{
			Color: static.ColorEmbedOrange,
			Description: fmt.Sprintf("Someone logged in to the web interface as you.\n"+
				"\n**Details:**\nIP Address: ||`%s`||\nUser Agent: `%s`\n\n"+
				"If this was not you, consider disabling OTA [**here**](%s/usersettings).",
				ipaddr, useragent, c.cfg.Config().WebServer.PublicAddr),
			Timestamp: c.tp.Now().Format(time.RFC3339),
		}
		c.session.ChannelMessageSendEmbed(ch.ID, emb)
	}

	return c.oauthHandler.LoginSuccessHandler(ctx, discordoauth.SuccessResult{
		UserID: userID,
	})
}
