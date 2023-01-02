package slashcommands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/util/presence"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekrotja/ken"
)

type Presence struct {
	ken.EphemeralCommand
}

var (
	_ ken.SlashCommand        = (*Presence)(nil)
	_ permissions.PermCommand = (*Presence)(nil)
)

func (c *Presence) Name() string {
	return "presence"
}

func (c *Presence) Description() string {
	return "Set the 'now playing' state of the bot."
}

func (c *Presence) Version() string {
	return "1.0.0"
}

func (c *Presence) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *Presence) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "message",
			Description: "The presence message.",
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "status",
			Description: "The presence status.",
			Choices: []*discordgo.ApplicationCommandOptionChoice{
				{Name: string(presence.StatusOnline), Value: presence.StatusOnline},
				{Name: string(presence.StatusIdle), Value: presence.StatusIdle},
				{Name: string(presence.StatusDnD), Value: presence.StatusDnD},
				{Name: string(presence.StatusInvisible), Value: presence.StatusInvisible},
			},
		},
	}
}

func (c *Presence) Domain() string {
	return "sp.presence"
}

func (c *Presence) SubDomains() []permissions.SubPermission {
	return nil
}

func (c *Presence) Run(ctx ken.Context) (err error) {
	if err = ctx.Defer(); err != nil {
		return
	}

	db := ctx.Get(static.DiDatabase).(database.Database)

	rawPre, err := db.GetSetting(static.SettingPresence)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	pre, _ := presence.Unmarshal(rawPre)
	if pre == nil {
		pre = &presence.Presence{
			Game:   "shnp.de",
			Status: presence.StatusOnline,
		}
	}

	if messageV, ok := ctx.Options().GetByNameOptional("message"); ok {
		pre.Game = messageV.StringValue()
	}

	if statusV, ok := ctx.Options().GetByNameOptional("status"); ok {
		pre.Status = presence.Status(statusV.StringValue())
	}

	if err = pre.Validate(); err != nil {
		return ctx.FollowUpError(err.Error(), "").Send().Error
	}

	err = ctx.GetSession().UpdateStatusComplex(pre.ToUpdateStatusData())
	if err != nil {
		return err
	}

	err = db.SetSetting(static.SettingPresence, pre.Marshal())
	if err != nil {
		return err
	}

	return ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: "Presence updated.",
	}).Send().Error
}
