package slashcommands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekrotja/dgrs"
	"github.com/zekrotja/ken"
)

type Mvall struct {
	ken.EphemeralCommand
}

var (
	_ ken.SlashCommand        = (*Mvall)(nil)
	_ permissions.PermCommand = (*Mvall)(nil)
)

func (c *Mvall) Name() string {
	return "moveall"
}

func (c *Mvall) Description() string {
	return "Move all members of the current voice channel to another one."
}

func (c *Mvall) Version() string {
	return "1.0.0"
}

func (c *Mvall) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *Mvall) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:         discordgo.ApplicationCommandOptionChannel,
			Name:         "channel",
			Description:  "Voice channel to move to.",
			Required:     true,
			ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildVoice},
		},
	}
}

func (c *Mvall) Domain() string {
	return "sp.guild.mod.mvall"
}

func (c *Mvall) SubDomains() []permissions.SubPermission {
	return nil
}

func (c *Mvall) Run(ctx ken.Context) (err error) {
	if err = ctx.Defer(); err != nil {
		return
	}

	st := ctx.Get(static.DiState).(*dgrs.State)

	channel := ctx.Options().GetByName("channel").ChannelValue(ctx)

	vs, err := st.VoiceState(ctx.GetEvent().GuildID, ctx.User().ID)
	if err != nil {
		return
	}
	if vs == nil {
		return ctx.FollowUpError(
			"You need to be in a voice channel to use this command.", "").Error
	}
	if vs.ChannelID == channel.ID {
		return ctx.FollowUpError(
			"You are already in the target voice channel.", "").Error
	}

	vss, err := st.VoiceStates(ctx.GetEvent().GuildID)
	if err != nil {
		return err
	}

	var i int
	for _, vs := range vss {
		if vs.ChannelID == vs.ChannelID {
			err := ctx.GetSession().GuildMemberMove(ctx.GetEvent().GuildID, vs.UserID, &channel.ID)
			if err != nil {
				return err
			}
			i++
		}
	}

	return ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: fmt.Sprintf("Moved %d members to channel %s.",
			i, channel.Name),
	}).Error
}
