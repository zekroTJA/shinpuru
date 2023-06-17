package slashcommands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/services/config"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/embedbuilder"
	"github.com/zekroTJA/shinpuru/pkg/md"
	"github.com/zekroTJA/shinpuru/pkg/mody"
	"github.com/zekrotja/dgrs"
	"github.com/zekrotja/ken"
)

type Maintenance struct {
	ken.EphemeralCommand
}

var (
	_ ken.SlashCommand        = (*Maintenance)(nil)
	_ permissions.PermCommand = (*Maintenance)(nil)
	_ ken.DmCapable           = (*Maintenance)(nil)
)

func (c *Maintenance) Name() string {
	return "maintenance"
}

func (c *Maintenance) Description() string {
	return "Maintenance utilities."
}

func (c *Maintenance) Version() string {
	return "1.2.0"
}

func (c *Maintenance) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *Maintenance) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "flush-state",
			Description: "Flush dgrs state.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "reconnect",
					Description: "Disconnect and reconnect session after flush.",
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "subkeys",
					Description: "The cache sub keys (comma seperated).",
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "kill",
			Description: "Kill the bot process.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "exitcode",
					Description: "The exit code.",
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "reconnect",
			Description: "Reconnects the Discord session.",
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "reload-config",
			Description: "Reloads the bots config.",
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "set-config-value",
			Description: "Set a specific config value.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "field",
					Description: "The config fild path and name.",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "jsonvalue",
					Description: "The value as JSON representation.",
					Required:    true,
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "guild-info",
			Description: "Display information about a given guild.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "id",
					Description: "The ID of the guild.",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "json",
					Description: "Dispaly output as JSON.",
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "leave-guild",
			Description: "Let shinpuru leave the specified guild.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "id",
					Description: "The ID of the guild.",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "owner-message",
					Description: "Send the given message to the owner og the guild.",
					Required:    false,
				},
			},
		},
	}
}

func (c *Maintenance) Domain() string {
	return "sp.maintenance"
}

func (c *Maintenance) SubDomains() []permissions.SubPermission {
	return nil
}

func (c *Maintenance) IsDmCapable() bool {
	return true
}

func (c *Maintenance) Run(ctx ken.Context) (err error) {
	if err = ctx.Defer(); err != nil {
		return
	}

	err = ctx.HandleSubCommands(
		ken.SubCommandHandler{"flush-state", c.flushState},
		ken.SubCommandHandler{"kill", c.kill},
		ken.SubCommandHandler{"reconnect", c.reconnect},
		ken.SubCommandHandler{"reload-config", c.reloadConfig},
		ken.SubCommandHandler{"set-config-value", c.setConfigValue},
		ken.SubCommandHandler{"guild-info", c.guildInfo},
		ken.SubCommandHandler{"leave-guild", c.leaveGuild},
	)

	return
}

func (c *Maintenance) flushState(ctx ken.SubCommandContext) (err error) {
	st := ctx.Get(static.DiState).(*dgrs.State)

	subkeys := ([]string)(nil)
	if subkeysV, ok := ctx.Options().GetByNameOptional("subkeys"); ok {
		subkeys = strings.Split(subkeysV.StringValue(), ",")
		for i, sk := range subkeys {
			subkeys[i] = strings.TrimSpace(sk)
		}
	}

	if err = st.Flush(subkeys...); err != nil {
		return
	}

	if reconnectV, ok := ctx.Options().GetByNameOptional("reconnect"); ok && reconnectV.BoolValue() {
		ctx.GetSession().Close()
		ctx.GetSession().Open()
	}

	return ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: "✅ State cache flushed.",
		Color:       static.ColorEmbedGreen,
	}).Send().Error
}

func (c *Maintenance) kill(ctx ken.SubCommandContext) (err error) {
	code := 1

	if exitcodeV, ok := ctx.Options().GetByNameOptional("exitcode"); ok {
		code = int(exitcodeV.IntValue())
	}

	err = ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: "👋 Bye.",
		Color:       static.ColorEmbedOrange,
	}).Send().Error
	if err != nil {
		return
	}

	os.Exit(code)

	return
}

func (c *Maintenance) reconnect(ctx ken.SubCommandContext) (err error) {
	if err = ctx.GetSession().Close(); err != nil {
		return
	}

	ctx.GetSession().Open()

	return ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: "✅ Successfully reconnected.",
		Color:       static.ColorEmbedGreen,
	}).Send().Error
}

func (c *Maintenance) reloadConfig(ctx ken.SubCommandContext) (err error) {
	cfg := ctx.Get(static.DiConfig).(config.Provider)

	if err = cfg.Parse(); err != nil {
		return
	}

	return ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: "Config has been reloaded.\n\nSome config changes will only take effect after a restart!",
	}).Send().Error
}

func (c *Maintenance) setConfigValue(ctx ken.SubCommandContext) (err error) {
	cfg := ctx.Get(static.DiConfig).(config.Provider)

	field := ctx.Options().GetByName("field").StringValue()
	jsonvalue := ctx.Options().GetByName("jsonvalue").StringValue()

	var errInner error
	err = mody.Catch(func() {
		errInner = mody.UpdateJson(cfg.Config(), field, jsonvalue)
	})
	if err != nil {
		return
	}
	if err = errInner; err != nil {
		return
	}

	return ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: fmt.Sprintf("Config value `%s` has been updated to `%s`.\n\n"+
			"Keep in mind that not all config value changes will be effective.",
			field, jsonvalue),
	}).Send().Error
}

func (c *Maintenance) guildInfo(ctx ken.SubCommandContext) (err error) {
	st := ctx.Get(static.DiState).(dgrs.IState)

	id := ctx.Options().GetByName("id").StringValue()

	asJson := false
	if oAsJson, ok := ctx.Options().GetByNameOptional("json"); ok {
		asJson = oAsJson.BoolValue()
	}

	guild, err := st.Guild(id)
	if err != nil {
		return err
	}

	owner, err := st.Member(guild.ID, guild.OwnerID)
	if err != nil {
		return err
	}

	if asJson {
		guildJson, err := json.MarshalIndent(guild, "", "  ")
		if err != nil {
			return err
		}

		ownerJson, err := json.MarshalIndent(owner, "", "  ")
		if err != nil {
			return err
		}

		return ctx.FollowUp(true, &discordgo.WebhookParams{
			Files: []*discordgo.File{
				{
					Name:        fmt.Sprintf("%s.guild.json", guild.ID),
					ContentType: "application/json",
					Reader:      bytes.NewBuffer(guildJson),
				},
				{
					Name:        fmt.Sprintf("%s.%s.owner.json", guild.ID, guild.OwnerID),
					ContentType: "application/json",
					Reader:      bytes.NewBuffer(ownerJson),
				},
			},
		}).Send().Error
	}

	emb := embedbuilder.New().
		WithTitle(guild.Name).
		WithAuthor(owner.User.String(), "", owner.AvatarURL("16"), "").
		AddField("ID", md.CodeBlock(guild.ID)).
		AddField("Onwer ID", md.CodeBlock(guild.OwnerID)).
		AddField("Region", guild.Region).
		AddField("Joined At", guild.JoinedAt.Format(time.RFC1123)).
		AddField("Membre Count", strconv.Itoa(guild.MemberCount)).
		WithThumbnail(guild.IconURL("64"), "", 64, 64).
		Build()

	return ctx.FollowUpEmbed(emb).Send().Error
}

func (c *Maintenance) leaveGuild(ctx ken.SubCommandContext) (err error) {
	id := ctx.Options().GetByName("id").StringValue()

	var ownerMessage string
	if v, ok := ctx.Options().GetByNameOptional("owner-message"); ok {
		ownerMessage = v.StringValue()
	}

	var msgErr error

	if ownerMessage != "" {
		st := ctx.Get(static.DiState).(dgrs.IState)
		guild, err := st.Guild(id)
		if err != nil {
			return err
		}

		msgErr = sendMessageToUser(ctx.GetSession(), guild.OwnerID, ownerMessage)
	}

	err = ctx.GetSession().GuildLeave(id)
	if err != nil {
		return err
	}

	if msgErr != nil {
		return ctx.FollowUpError(
			fmt.Sprintf("Sending the message to the owner of the guild failed:\n```\n%s\n```\n"+
				"The bot has been removed from the guild though.", err.Error()),
			"Failed sending message to owner").Send().Error
	}

	return ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: "The bot has left the guild.",
	}).Send().Error
}

func sendMessageToUser(s *discordgo.Session, userId string, msg string) error {
	ch, err := s.UserChannelCreate(userId)
	if err != nil {
		return err
	}

	_, err = s.ChannelMessageSendEmbed(ch.ID, &discordgo.MessageEmbed{
		Title:       "⚠️ Important Information",
		Description: msg,
	})
	return err
}
