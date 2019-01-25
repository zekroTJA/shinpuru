package commands

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/util"
)

type CmdStats struct {
	PermLvl int
}

func (c *CmdStats) GetInvokes() []string {
	return []string{"stats", "uptime", "numbers"}
}

func (c *CmdStats) GetDescription() string {
	return "dispaly some stats like uptime or guilds/user count"
}

func (c *CmdStats) GetHelp() string {
	return "`stats`"
}

func (c *CmdStats) GetGroup() string {
	return GroupEtc
}

func (c *CmdStats) GetPermission() int {
	return c.PermLvl
}

func (c *CmdStats) SetPermission(permLvl int) {
	c.PermLvl = permLvl
}

func (c *CmdStats) Exec(args *CommandArgs) error {
	uptime := int(time.Since(util.StatsStartupTime).Seconds())
	uptimeDays := int(uptime / (3600 * 24))
	uptimeHours := int(uptime % (3600 * 24) / 3600)
	uptimeMinutes := int(uptime % (3600 * 24) % 3600 / 60)
	uptimeSeconds := uptime % (3600 * 24) % 3600 % 60

	var guildUsers int
	for _, g := range args.Session.State.Guilds {
		guildUsers += g.MemberCount
	}

	emb := &discordgo.MessageEmbed{
		Color: util.ColorEmbedDefault,
		Title: "shinpuru Global Stats",
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name: "Uptime",
				Value: fmt.Sprintf("%d d, %d h, %d min, %d sec",
					uptimeDays, uptimeHours, uptimeMinutes, uptimeSeconds),
			},
			&discordgo.MessageEmbedField{
				Name: "Stats since startup",
				Value: fmt.Sprintf("**%d** Messages analysed for commands\n**%d** commands executed",
					util.StatsMessagesAnalysed, util.StatsCommandsExecuted+1),
			},
			&discordgo.MessageEmbedField{
				Name: "Guilds & Members",
				Value: fmt.Sprintf("Serving **%d** guilds with **%d** members in total.",
					len(args.Session.State.Guilds), guildUsers),
			},
		},
	}
	_, err := args.Session.ChannelMessageSendEmbed(args.Channel.ID, emb)
	return err
}
