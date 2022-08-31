package slashcommands

import (
	"fmt"
	"runtime"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/bytecount"
	"github.com/zekrotja/dgrs"
	"github.com/zekrotja/ken"
)

type Stats struct{}

var (
	_ ken.SlashCommand        = (*Stats)(nil)
	_ permissions.PermCommand = (*Stats)(nil)
	_ ken.DmCapable           = (*Stats)(nil)
)

func (c *Stats) Name() string {
	return "stats"
}

func (c *Stats) Description() string {
	return "Display some stats like uptime or guilds/user count."
}

func (c *Stats) Version() string {
	return "1.0.0"
}

func (c *Stats) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *Stats) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{}
}

func (c *Stats) Domain() string {
	return "sp.etc.stats"
}

func (c *Stats) SubDomains() []permissions.SubPermission {
	return nil
}

func (c *Stats) IsDmCapable() bool {
	return true
}

func (c *Stats) Run(ctx ken.Context) (err error) {
	if err = ctx.Defer(); err != nil {
		return
	}

	st := ctx.Get(static.DiState).(*dgrs.State)

	uptime := time.Since(util.StatsStartupTime)

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
				Name:  "Uptime",
				Value: uptime.Round(time.Second).String(),
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

	return ctx.FollowUpEmbed(emb).Error
}
