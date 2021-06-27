package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/pkg/fetch"
	"github.com/zekroTJA/shireikan"
)

type CmdMvall struct {
}

func (c *CmdMvall) GetInvokes() []string {
	return []string{"mvall", "mva"}
}

func (c *CmdMvall) GetDescription() string {
	return "Move all members in your current voice channel into another one."
}

func (c *CmdMvall) GetHelp() string {
	return "`mvall <otherChanResolvable>`"
}

func (c *CmdMvall) GetGroup() string {
	return shireikan.GroupModeration
}

func (c *CmdMvall) GetDomainName() string {
	return "sp.guild.mod.mvall"
}

func (c *CmdMvall) GetSubPermissionRules() []shireikan.SubPermission {
	return nil
}

func (c *CmdMvall) IsExecutableInDMChannels() bool {
	return false
}

func (c *CmdMvall) Exec(ctx shireikan.Context) error {
	if len(ctx.GetArgs()) < 1 {
		return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
			"Please enter a voice channel as argument.").
			DeleteAfter(8 * time.Second).Error()
	}

	var currVC string
	for _, vs := range ctx.GetGuild().VoiceStates {
		if vs.UserID == ctx.GetUser().ID {
			currVC = vs.ChannelID
		}
	}

	if currVC == "" {
		return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
			"You need to be in a voice channel to use this command.").
			DeleteAfter(8 * time.Second).Error()
	}

	toVC, err := fetch.FetchChannel(ctx.GetSession(), ctx.GetGuild().ID, strings.Join(ctx.GetArgs(), " "),
		func(c *discordgo.Channel) bool {
			return c.Type == discordgo.ChannelTypeGuildVoice
		})
	if err != nil {
		return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
			"Could not find any voice channel passing the resolvable.").
			DeleteAfter(8 * time.Second).Error()
	}

	if toVC.Type != discordgo.ChannelTypeGuildVoice {
		return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
			fmt.Sprintf("The target channel *(`%s`)* is not a voice channel.", toVC.Name)).
			DeleteAfter(8 * time.Second).Error()
	}

	for _, vs := range ctx.GetGuild().VoiceStates {
		if vs.ChannelID == currVC {
			err := ctx.GetSession().GuildMemberMove(ctx.GetGuild().ID, vs.UserID, &toVC.ID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
