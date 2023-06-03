package slashcommands

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/util/modnot"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekrotja/ken"
)

type Debug struct{}

var (
	_ ken.SlashCommand        = (*User)(nil)
	_ permissions.PermCommand = (*User)(nil)
)

func (c *Debug) Name() string {
	return "debug"
}

func (c *Debug) Description() string {
	return "Debug command for development."
}

func (c *Debug) Version() string {
	return "1.0.0"
}

func (c *Debug) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *Debug) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{}
}

func (c *Debug) Domain() string {
	return "sp.debug"
}

func (c *Debug) SubDomains() []permissions.SubPermission {
	return nil
}

func (c *Debug) Run(ctx ken.Context) (err error) {
	if err = ctx.Defer(); err != nil {
		return err
	}

	db := ctx.Get(static.DiDatabase).(database.Database)

	err = modnot.Send(db, ctx.GetSession(), ctx.GetEvent().GuildID, &discordgo.MessageEmbed{
		Color:       static.ColorEmbedDefault,
		Description: "Just a test.",
	})
	if err != nil {
		return err
	}

	return ctx.
		FollowUpEmbed(&discordgo.MessageEmbed{Description: "Test mod notification sent."}).
		Send().
		DeleteAfter(2 * time.Second).
		Error
}
