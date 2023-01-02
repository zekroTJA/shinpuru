package slashcommands

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/services/config"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/onetimeauth/v2"
	"github.com/zekroTJA/shinpuru/pkg/timerstack"
	"github.com/zekrotja/ken"
)

type Login struct {
	ken.EphemeralCommand
}

var (
	_ ken.SlashCommand        = (*Login)(nil)
	_ permissions.PermCommand = (*Login)(nil)
	_ ken.DmCapable           = (*Login)(nil)
)

func (c *Login) Name() string {
	return "login"
}

func (c *Login) Description() string {
	return "Log in to the web interface."
}

func (c *Login) Version() string {
	return "1.1.0"
}

func (c *Login) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *Login) IsDmCapable() bool {
	return true
}

func (c *Login) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{}
}

func (c *Login) Domain() string {
	return "sp.etc.login"
}

func (c *Login) SubDomains() []permissions.SubPermission {
	return nil
}

func (c *Login) Run(ctx ken.Context) (err error) {
	if err = ctx.Defer(); err != nil {
		return
	}

	cfg := ctx.Get(static.DiConfig).(config.Provider)
	ota := ctx.Get(static.DiOneTimeAuth).(onetimeauth.OneTimeAuth)
	db := ctx.Get(static.DiDatabase).(database.Database)

	enabled, err := db.GetUserOTAEnabled(ctx.User().ID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return
	}

	if !enabled {
		enableLink := fmt.Sprintf("%s/usersettings", cfg.Config().WebServer.PublicAddr)
		err = ctx.FollowUpError(
			"One Time Authorization is disabled by default. If you want to use it, you need "+
				"to enable it first in your [**user settings page**]("+enableLink+").", "").
			Send().Error
		return err
	}

	token, _, err := ota.GetKey(ctx.User().ID, "login-via-dm")
	if err != nil {
		return
	}

	link := fmt.Sprintf("%s/api/ota?token=%s", cfg.Config().WebServer.PublicAddr, token)
	emb := &discordgo.MessageEmbed{
		Color: static.ColorEmbedDefault,
		Description: "Click this [**this link**](" + link + ") and you will be automatically logged " +
			"in to the shinpuru web interface.\n\nThis link expires in one minute.",
	}

	fum := ctx.FollowUpEmbed(emb).AddComponents(func(cb *ken.ComponentBuilder) {
		cb.AddActionsRow(func(b ken.ComponentAssembler) {
			b.Add(discordgo.Button{
				Label: "Login to the Web Interface",
				Style: discordgo.LinkButton,
				URL:   link,
			}, nil)
		})
	}).Send()
	if fum.HasError() {
		return fum.Error
	}

	timerstack.New().After(1*time.Minute, func() bool {
		emb := &discordgo.MessageEmbed{
			Color:       static.ColorEmbedGray,
			Description: "The login link has expired.",
		}
		fum.Edit(&discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{emb},
			Components: &[]discordgo.MessageComponent{discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{discordgo.Button{
					Label:    "Login to the Web Interface",
					Style:    discordgo.LinkButton,
					Disabled: true,
					URL:      link,
				}},
			}},
		})
		return true
	}).RunBlocking()

	return
}
