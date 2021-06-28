package commands

import (
	"fmt"
	"runtime"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/bytecount"
	"github.com/zekroTJA/shireikan"
	"github.com/zekrotja/dgrs"
)

type CmdStats struct {
	PermLvl int
}

func (c *CmdStats) GetInvokes() []string {
	return []string{"stats", "uptime", "numbers"}
}

func (c *CmdStats) GetDescription() string {
	return "Display some stats like uptime or guilds/user count."
}

func (c *CmdStats) GetHelp() string {
	return "`stats`"
}

func (c *CmdStats) GetGroup() string {
	return shireikan.GroupEtc
}

func (c *CmdStats) GetDomainName() string {
	return "sp.etc.stats"
}

func (c *CmdStats) GetSubPermissionRules() []shireikan.SubPermission {
	return nil
}

func (c *CmdStats) IsExecutableInDMChannels() bool {
	return true
}

func (c *CmdStats) Exec(ctx shireikan.Context) error {
	uptime := int(time.Since(util.StatsStartupTime).Seconds())
	uptimeDays := int(uptime / (3600 * 24))
	uptimeHours := int(uptime % (3600 * 24) / 3600)
	uptimeMinutes := int(uptime % (3600 * 24) % 3600 / 60)
	uptimeSeconds := uptime % (3600 * 24) % 3600 % 60

	st := ctx.GetObject(static.DiState).(*dgrs.State)
	guilds, err := st.Guilds()
	if err != nil {
		return err
	}

	var guildUsers int
	for _, g := range guilds {
		guildUsers += g.MemberCount
	}

	nGoroutines := runtime.NumGoroutine()
	usedCPUs := runtime.NumCPU()
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	usedHeap := bytecount.Format(memStats.HeapInuse)
	usedStack := bytecount.Format(memStats.StackInuse)

	emb := &discordgo.MessageEmbed{
		Color: static.ColorEmbedDefault,
		Title: "shinpuru Global Stats",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name: "Uptime",
				Value: fmt.Sprintf("%d d, %d h, %d min, %d sec",
					uptimeDays, uptimeHours, uptimeMinutes, uptimeSeconds),
			},
			{
				Name: "Stats since startup",
				Value: fmt.Sprintf("**%d** Messages analysed for commands\n**%d** commands executed",
					util.StatsMessagesAnalysed, util.StatsCommandsExecuted+1),
			},
			{
				Name: "Guilds & Members",
				Value: fmt.Sprintf("Serving **%d** guilds with **%d** members in total.",
					len(guilds), guildUsers),
			},
			{
				Name: "Runtime Stats",
				Value: fmt.Sprintf("Running Go Routines: **%d**\nUsed CPU Threads: **%d**\n"+
					"Used Heap: **%s**\nUsed Stack: **%s**", nGoroutines, usedCPUs, usedHeap, usedStack),
			},
		},
	}
	_, err = ctx.GetSession().ChannelMessageSendEmbed(ctx.GetChannel().ID, emb)
	return err
}
