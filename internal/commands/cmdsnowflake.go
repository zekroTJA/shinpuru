package commands

import (
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/snowflake"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/snowflakenodes"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/timeutil"
	"github.com/zekroTJA/shireikan"
)

const (
	snowflakeTypeDiscord = iota
	snowflakeTypeShinpuru
)

type CmdSnowflake struct {
	rx *regexp.Regexp
}

func (c *CmdSnowflake) GetInvokes() []string {
	return []string{"snowflake", "sf"}
}

func (c *CmdSnowflake) GetDescription() string {
	return "Calculate information about a Discord or Shinpuru snowflake"
}

func (c *CmdSnowflake) GetHelp() string {
	return "`snowflake <snowflake> (dc/sp)` - get snowflake information\n" +
		"If you attach `dc` (Discord) or `sp` (shinpuru), you will force the calculation " +
		"mode for the snowflake. With nothing given, the mode will be chosen automatically."
}

func (c *CmdSnowflake) GetGroup() string {
	return shireikan.GroupEtc
}

func (c *CmdSnowflake) GetDomainName() string {
	return "sp.etc.snowflake"
}

func (c *CmdSnowflake) GetSubPermissionRules() []shireikan.SubPermission {
	return nil
}

func (c *CmdSnowflake) IsExecutableInDMChannels() bool {
	return true
}

func (c *CmdSnowflake) Exec(ctx shireikan.Context) error {
	if c.rx == nil {
		c.rx = regexp.MustCompile(`(\d+)\s*([a-zA-Z]+)?`)
	}

	matches := c.rx.FindStringSubmatch(strings.Join(ctx.GetArgs(), " "))

	if len(matches) < 2 {
		return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
			"Please enter a Snowflake which should be calculated.").
			DeleteAfter(8 * time.Second).Error()
	}

	sf := matches[1]
	sfTyp := -1

	if len(matches) >= 3 {
		typ := strings.ToLower(matches[2])
		switch typ {
		case "dc":
			sfTyp = snowflakeTypeDiscord
		case "sp":
			sfTyp = snowflakeTypeShinpuru
		}
	}

	sfAsDc, err := snowflakenodes.ParseDiscordSnowflake(sf)
	if err != nil {
		return err
	}

	sfAsSp, err := snowflake.ParseString(sf)
	if err != nil {
		return err
	}

	if sfTyp == -1 {
		if math.Abs(float64(time.Now().Year()-sfAsDc.Time.Year())) > 10 {
			sfTyp = snowflakeTypeShinpuru
		} else {
			sfTyp = snowflakeTypeDiscord
		}
	}

	switch sfTyp {
	case snowflakeTypeDiscord:
		return c.printSfDc(ctx, sfAsDc)
	case snowflakeTypeShinpuru:
		return c.printSfSp(ctx, sfAsSp)
	}

	return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
		"Unknown snowflake type was provided.").
		DeleteAfter(8 * time.Second).Error()
}

func (c *CmdSnowflake) printSfDc(ctx shireikan.Context, sf *snowflakenodes.DiscordSnowflake) error {
	emb := &discordgo.MessageEmbed{
		Title: "Snowflake Info",
		Color: 0x7289DA,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Snowflake",
				Value: fmt.Sprintf("```\n%s\n```", sf.Snowflake),
			},
			{
				Name:  "Type",
				Value: "Discord Snowflake ID",
			},
			{
				Name:  "Timestamp",
				Value: sf.Time.Format(time.RFC1123),
			},
			{
				Name:   "Worker ID",
				Value:  util.EnsureNotEmpty(fmt.Sprintf("%d", sf.WorkerID), "*<empty>*"),
				Inline: true,
			},
			{
				Name:   "Process ID",
				Value:  util.EnsureNotEmpty(fmt.Sprintf("%d", sf.ProcessID), "*<empty>*"),
				Inline: true,
			},
			{
				Name:   "Incremental ID",
				Value:  util.EnsureNotEmpty(fmt.Sprintf("%d", sf.IncrementalID), "*<empty>*"),
				Inline: true,
			},
		},
	}

	_, err := ctx.GetSession().ChannelMessageSendEmbed(ctx.GetChannel().ID, emb)
	return err
}

func (c *CmdSnowflake) printSfSp(ctx shireikan.Context, sf snowflake.ID) error {
	emb := &discordgo.MessageEmbed{
		Title: "Snowflake Info",
		Color: static.ColorEmbedDefault,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Snowflake",
				Value: fmt.Sprintf("```\n%s\n```", sf.String()),
			},
			{
				Name:  "Type",
				Value: "Discord Snowflake ID",
			},
			{
				Name: "Timestamp",
				Value: timeutil.
					FromUnix(int(sf.Time())).
					Format(time.RFC1123),
			},
			{
				Name:  "Node Name",
				Value: util.EnsureNotEmpty(snowflakenodes.GetNodeName(sf.Node()), "*<empty>*"),
			},
			{
				Name:   "Node ID",
				Value:  util.EnsureNotEmpty(fmt.Sprintf("%d", sf.Node()), "*<empty>*"),
				Inline: true,
			},
			{
				Name:   "Incremental ID",
				Value:  util.EnsureNotEmpty(fmt.Sprintf("%d", sf.Step()), "*<empty>*"),
				Inline: true,
			},
		},
	}

	_, err := ctx.GetSession().ChannelMessageSendEmbed(ctx.GetChannel().ID, emb)
	return err
}
