package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/fetch"
	"github.com/zekroTJA/shireikan"
	"github.com/zekrotja/dgrs"
)

type CmdKarma struct {
}

func (c *CmdKarma) GetInvokes() []string {
	return []string{"karma", "scoreboard", "leaderboard", "lb", "sb", "top"}
}

func (c *CmdKarma) GetDescription() string {
	return "Display users karma count or the guilds karma scoreboard."
}

func (c *CmdKarma) GetHelp() string {
	return "`karma` - Display karma scoreboard\n" +
		"`karma <userResolvable>` - Display karma count of this user\n"
}

func (c *CmdKarma) GetGroup() string {
	return shireikan.GroupChat
}

func (c *CmdKarma) GetDomainName() string {
	return "sp.chat.karma"
}

func (c *CmdKarma) GetSubPermissionRules() []shireikan.SubPermission {
	return nil
}

func (c *CmdKarma) IsExecutableInDMChannels() bool {
	return false
}

func (c *CmdKarma) Exec(ctx shireikan.Context) error {
	db, _ := ctx.GetObject(static.DiDatabase).(database.Database)

	userRes := ctx.GetArgs().Get(0).AsString()
	if userRes != "" {
		return c.userKarma(ctx, db, userRes)
	}

	karma, err := db.GetKarma(ctx.GetUser().ID, ctx.GetGuild().ID)
	if err != nil && err != database.ErrDatabaseNotFound {
		return err
	}

	karmaSum, err := db.GetKarmaSum(ctx.GetUser().ID)
	if err != nil && err != database.ErrDatabaseNotFound {
		return err
	}

	karmaList, err := db.GetKarmaGuild(ctx.GetGuild().ID, 20)
	if err != nil && err != database.ErrDatabaseNotFound {
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

	st := ctx.GetObject(static.DiState).(*dgrs.State)

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
			Text:    "Issued by " + ctx.GetUser().String(),
			IconURL: ctx.GetUser().AvatarURL("16x16"),
		},
	}

	return util.SendEmbedRaw(ctx.GetSession(), ctx.GetChannel().ID, emb).Error()
}

func (c *CmdKarma) userKarma(ctx shireikan.Context, db database.Database, userRes string) error {
	memb, err := fetch.FetchMember(ctx.GetSession(), ctx.GetGuild().ID, userRes)
	if err != nil {
		return err
	}

	guildKarma, err := db.GetKarma(memb.User.ID, ctx.GetGuild().ID)
	if err != nil {
		return err
	}

	globalKarma, err := db.GetKarmaSum(memb.User.ID)
	if err != nil {
		return err
	}

	return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
		fmt.Sprintf("Guild Karma: **`%d`**\nGlobal Karma: **`%d`**", guildKarma, globalKarma),
		memb.User.String()+"'s Karma Stats", 0).
		Error()
}
