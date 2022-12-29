package slashcommands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekrotja/dgrs"
	"github.com/zekrotja/ken"
)

type Karma struct {
	ken.EphemeralCommand
}

var (
	_ ken.SlashCommand        = (*Karma)(nil)
	_ permissions.PermCommand = (*Karma)(nil)
)

func (c *Karma) Name() string {
	return "karma"
}

func (c *Karma) Description() string {
	return "Display users karma count or the guilds karma scoreboard."
}

func (c *Karma) Version() string {
	return "1.0.0"
}

func (c *Karma) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *Karma) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionUser,
			Name:        "user",
			Description: "Display karma stats of a specific user.",
		},
	}
}

func (c *Karma) Domain() string {
	return "sp.chat.karma"
}

func (c *Karma) SubDomains() []permissions.SubPermission {
	return nil
}

func (c *Karma) Run(ctx ken.Context) (err error) {
	if err = ctx.Defer(); err != nil {
		return
	}

	if userV, ok := ctx.Options().GetByNameOptional("user"); ok {
		return c.userKarma(ctx, userV.UserValue(ctx))
	}

	db := ctx.Get(static.DiDatabase).(database.Database)
	st := ctx.Get(static.DiState).(*dgrs.State)

	karma, err := db.GetKarma(ctx.User().ID, ctx.GetEvent().GuildID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	karmaSum, err := db.GetKarmaSum(ctx.User().ID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	karmaList, err := db.GetKarmaGuild(ctx.GetEvent().GuildID, 20)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	var karmaListStr string

	var karmaListLn int
	if karmaList != nil {
		karmaListLn = len(karmaList)
	}

	if karmaListLn == 0 {
		karmaListStr = "*No entries for this guild.*"
	}

	var i int
	for _, v := range karmaList {
		m, err := st.Member(v.GuildID, v.UserID)
		if err != nil {
			continue
		}

		i++
		karmaListStr = fmt.Sprintf("%s\n`%d` - %s - **%d**",
			karmaListStr, i, m.User.String(), v.Value)
	}

	emb := &discordgo.MessageEmbed{
		Color: static.ColorEmbedDefault,
		Title: "Karma Scoreboard",
		Description: fmt.Sprintf(
			"Your Karma on this guild: **%d**\n"+
				"Your Global Karma: **%d**",
			karma, karmaSum),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  fmt.Sprintf("Scoreboard (Top %d)", karmaListLn),
				Value: karmaListStr,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Issued by " + ctx.User().String(),
			IconURL: ctx.User().AvatarURL("16x16"),
		},
	}

	return ctx.FollowUpEmbed(emb).Error
}

func (c *Karma) userKarma(ctx ken.Context, user *discordgo.User) error {
	st := ctx.Get(static.DiState).(*dgrs.State)
	db := ctx.Get(static.DiDatabase).(database.Database)

	memb, err := st.Member(ctx.GetEvent().GuildID, user.ID)
	if err != nil {
		return err
	}

	guildKarma, err := db.GetKarma(memb.User.ID, ctx.GetEvent().GuildID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	globalKarma, err := db.GetKarmaSum(memb.User.ID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	return ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Title:       memb.User.String() + "'s Karma Stats",
		Description: fmt.Sprintf("Guild Karma: **`%d`**\nGlobal Karma: **`%d`**", guildKarma, globalKarma),
	}).Error
}
