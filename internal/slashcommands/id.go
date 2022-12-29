package slashcommands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/fetch"
	"github.com/zekrotja/ken"
)

type Id struct {
	ken.EphemeralCommand
}

var (
	_ ken.SlashCommand        = (*Id)(nil)
	_ permissions.PermCommand = (*Id)(nil)
)

func (c *Id) Name() string {
	return "id"
}

func (c *Id) Description() string {
	return "Get the discord ID(s) by resolvable."
}

func (c *Id) Version() string {
	return "1.0.0"
}

func (c *Id) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *Id) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "resolvable",
			Description: "The name of a discord object.",
			Required:    true,
		},
	}
}

func (c *Id) Domain() string {
	return "sp.etc.id"
}

func (c *Id) SubDomains() []permissions.SubPermission {
	return nil
}

func (c *Id) Run(ctx ken.Context) (err error) {
	if err = ctx.Defer(); err != nil {
		return
	}

	resolvable := ctx.Options().GetByName("resolvable").StringValue()

	var user *discordgo.User
	var role *discordgo.Role
	var textChannel *discordgo.Channel
	var voiceChannel *discordgo.Channel

	if u, err := fetch.FetchMember(ctx.GetSession(), ctx.GetEvent().GuildID, resolvable); err == nil {
		user = u.User
	}
	if r, err := fetch.FetchRole(ctx.GetSession(), ctx.GetEvent().GuildID, resolvable); err == nil {
		role = r
	}
	if tc, err := fetch.FetchChannel(ctx.GetSession(), ctx.GetEvent().GuildID, resolvable, func(c *discordgo.Channel) bool {
		return c.Type == discordgo.ChannelTypeGuildText
	}); err == nil {
		textChannel = tc
	}
	if vc, err := fetch.FetchChannel(ctx.GetSession(), ctx.GetEvent().GuildID, resolvable, func(c *discordgo.Channel) bool {
		return c.Type == discordgo.ChannelTypeGuildVoice
	}); err == nil {
		voiceChannel = vc
	}

	if user == nil && role == nil && textChannel == nil && voiceChannel == nil {
		return ctx.FollowUpError(
			"Could not fetch any member, role or channel by this resolvable.", "").
			Error
	}

	emb := &discordgo.MessageEmbed{
		Color:  static.ColorEmbedDefault,
		Fields: make([]*discordgo.MessageEmbedField, 0),
	}

	if user != nil {
		emb.Fields = append(emb.Fields, &discordgo.MessageEmbedField{
			Name:  "Member",
			Value: fmt.Sprintf("<@%s> (%s#%s)\n```\n%s\n```", user.ID, user.Username, user.Discriminator, user.ID),
		})
	}
	if role != nil {
		emb.Fields = append(emb.Fields, &discordgo.MessageEmbedField{
			Name:  "Role",
			Value: fmt.Sprintf("<@&%s> (%s)\n```\n%s\n```", role.ID, role.Name, role.ID),
		})
	}
	if textChannel != nil {
		emb.Fields = append(emb.Fields, &discordgo.MessageEmbedField{
			Name:  "Text Channel",
			Value: fmt.Sprintf("<#%s> (%s)\n```\n%s\n```", textChannel.ID, textChannel.Name, textChannel.ID),
		})
	}
	if voiceChannel != nil {
		emb.Fields = append(emb.Fields, &discordgo.MessageEmbedField{
			Name:  "Voice Channel",
			Value: fmt.Sprintf("%s\n```\n%s\n```", voiceChannel.Name, voiceChannel.ID),
		})
	}
	emb.Fields = append(emb.Fields, &discordgo.MessageEmbedField{
		Name:  "Guild",
		Value: fmt.Sprintf("```\n%s\n```", ctx.GetEvent().GuildID),
	})

	return ctx.FollowUpEmbed(emb).Error
}
