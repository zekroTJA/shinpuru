package commands

import (
	"fmt"
	"image"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/generaltso/vibrant"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
	"github.com/zekroTJA/shinpuru/pkg/embedbuilder"
	"github.com/zekroTJA/shinpuru/pkg/httpreq"
	"github.com/zekroTJA/shireikan"
)

type CmdGuild struct {
}

func (c *CmdGuild) GetInvokes() []string {
	return []string{"guild", "guildinfo", "g", "gi"}
}

func (c *CmdGuild) GetDescription() string {
	return "Outputs info about the current guild."
}

func (c *CmdGuild) GetHelp() string {
	return "`guild` - Prints guild info"
}

func (c *CmdGuild) GetGroup() string {
	return shireikan.GroupChat
}

func (c *CmdGuild) GetDomainName() string {
	return "sp.chat.guild"
}

func (c *CmdGuild) GetSubPermissionRules() []shireikan.SubPermission {
	return nil
}

func (c *CmdGuild) IsExecutableInDMChannels() bool {
	return false
}

func (c *CmdGuild) Exec(ctx shireikan.Context) error {
	g := ctx.GetGuild()
	clr := 0

	if body, err := httpreq.GetFile(g.IconURL()); err == nil {
		if imgData, _, err := image.Decode(body); err == nil {
			if palette, err := vibrant.NewPaletteFromImage(imgData); err == nil {
				for name, swatch := range palette.ExtractAwesome() {
					if name == "Vibrant" {
						clr = int(swatch.Color)
					}
				}
			}
		}
	}

	nTextChans, nVoiceChans, nCategoryChans := 0, 0, 0
	for _, c := range g.Channels {
		switch c.Type {
		case discordgo.ChannelTypeGuildCategory:
			nCategoryChans++
		case discordgo.ChannelTypeGuildVoice:
			nVoiceChans++
		default:
			nTextChans++
		}
	}
	chans := fmt.Sprintf("Category Channels: `%d`\nText Channels: `%d`\nVoice Channels: `%d`",
		nCategoryChans, nTextChans, nVoiceChans)

	roles := make([]string, len(g.Roles))
	for i, r := range g.Roles {
		roles[i] = r.Mention()
	}

	createdTime, err := discordutil.GetDiscordSnowflakeCreationTime(g.ID)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}

	emb := embedbuilder.New().
		WithThumbnail(g.IconURL(), "", 100, 100).
		WithColor(clr).
		AddField("Name", g.Name, false).
		AddField("ID", fmt.Sprintf("```\n%s\n```", g.ID), false).
		AddField("Created", createdTime.Format(time.RFC1123), false).
		AddField("Owner", fmt.Sprintf("<@%s>", g.OwnerID), false).
		AddField("Server Region", g.Region, false).
		AddField("Channels", chans, false).
		AddField("Member Count", fmt.Sprintf("State: %d / Approx.: %d", g.MemberCount, g.ApproximateMemberCount), false).
		AddField("Roles", strings.Join(roles, ", "), false).
		WithFooter(fmt.Sprintf("issued by %s", ctx.GetUser().String()), "", "").
		Build()

	_, err = ctx.GetSession().ChannelMessageSendEmbed(ctx.GetChannel().ID, emb)
	return err
}
