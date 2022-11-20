package slashcommands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/services/guildlog"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
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
		return
	}

	gl := ctx.Get(static.DiGuildLog).(guildlog.Logger)

	gl.Debugf(ctx.GetEvent().GuildID, "Some debug message!")
	gl.Errorf(ctx.GetEvent().GuildID, "Some error message!")
	gl.Infof(ctx.GetEvent().GuildID, "Some info message!")
	gl.Warnf(ctx.GetEvent().GuildID, "Some warn message!")
	gl.Fatalf(ctx.GetEvent().GuildID, "Some fatal message!")

	return ctx.FollowUpEmbed(&discordgo.MessageEmbed{Description: "Ok"}).Error
}
