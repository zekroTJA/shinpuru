package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
)

type CmdKarma struct {
}

func (c *CmdKarma) GetInvokes() []string {
	return []string{"karma", "scoreboard"}
}

func (c *CmdKarma) GetDescription() string {
	return "Display users karma count or the guilds karma scoreboard."
}

func (c *CmdKarma) GetHelp() string {
	return "`karma` - Display karma scoreboard"
}

func (c *CmdKarma) GetGroup() string {
	return GroupChat
}

func (c *CmdKarma) GetDomainName() string {
	return "sp.chat.karma"
}

func (c *CmdKarma) GetSubPermissionRules() []SubPermission {
	return nil
}

func (c *CmdKarma) IsExecutableInDMChannels() bool {
	return false
}

func (c *CmdKarma) Exec(args *CommandArgs) error {

	karma, err := args.CmdHandler.db.GetKarma(args.User.ID, args.Guild.ID)
	if err != nil && err != database.ErrDatabaseNotFound {
		return err
	}

	karmaSum, err := args.CmdHandler.db.GetKarmaSum(args.User.ID)
	if err != nil && err != database.ErrDatabaseNotFound {
		return err
	}

	karmaList, err := args.CmdHandler.db.GetKarmaGuild(args.Guild.ID, 20)
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

	for i, v := range karmaList {
		m, err := discordutil.GetMember(args.Session, v.GuildID, v.UserID)
		if err != nil {
			continue
		}

		karmaListStr = fmt.Sprintf("%s\n`%d` - %s - **%d**",
			karmaListStr, i+1, m.User.String(), v.Value)
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
	}

	return util.SendEmbedRaw(args.Session, args.Channel.ID, emb).Error()
}
