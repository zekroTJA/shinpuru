package slashcommands

import (
	"fmt"
	"math"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/snowflake"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/util/snowflakenodes"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/stringutil"
	"github.com/zekroTJA/shinpuru/pkg/timeutil"
	"github.com/zekrotja/ken"
)

const (
	snowflakeTypeDiscord = iota
	snowflakeTypeShinpuru
)

type Snowflake struct{}

var (
	_ ken.Command             = (*Snowflake)(nil)
	_ permissions.PermCommand = (*Snowflake)(nil)
	_ ken.DmCapable           = (*Snowflake)(nil)
)

func (c *Snowflake) Name() string {
	return "snowflake"
}

func (c *Snowflake) Description() string {
	return "Calculate information about a Discord or Shinpuru snowflake."
}

func (c *Snowflake) Version() string {
	return "1.0.0"
}

func (c *Snowflake) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *Snowflake) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "snowflake",
			Description: "The snowflake ID.",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "type",
			Description: "The type of snowflake (will be determindes if not specified).",
			Choices: []*discordgo.ApplicationCommandOptionChoice{
				{"discord", snowflakeTypeDiscord},
				{"shinpuru", snowflakeTypeShinpuru},
			},
		},
	}
}

func (c *Snowflake) Domain() string {
	return "sp.etc.snowflake"
}

func (c *Snowflake) SubDomains() []permissions.SubPermission {
	return nil
}

func (c *Snowflake) IsDmCapable() bool {
	return true
}

func (c *Snowflake) Run(ctx *ken.Ctx) (err error) {
	if err = ctx.Defer(); err != nil {
		return
	}

	sfid := ctx.Options().GetByName("snowflake").IntValue()

	typ := -1
	if typV, ok := ctx.Options().GetByNameOptional("type"); ok {
		typ = int(typV.IntValue())
	}

	sfAsDc := snowflakenodes.ParseDiscordSnowflake(int(sfid))
	sfAsSp := snowflake.ParseInt64(sfid)
	if err != nil {
		return err
	}

	if typ == -1 {
		if math.Abs(float64(time.Now().Year()-sfAsDc.Time.Year())) > 10 {
			typ = snowflakeTypeShinpuru
		} else {
			typ = snowflakeTypeDiscord
		}
	}

	var emb *discordgo.MessageEmbed
	switch typ {
	case snowflakeTypeDiscord:
		emb = c.embSfDc(sfAsDc)
	case snowflakeTypeShinpuru:
		emb = c.embSfSp(sfAsSp)
	}

	return ctx.FollowUpEmbed(emb).Error
}

func (c *Snowflake) embSfDc(sf *snowflakenodes.DiscordSnowflake) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: "Snowflake Info",
		Color: 0x7289DA,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Snowflake",
				Value: fmt.Sprintf("```\n%d\n```", sf.Snowflake),
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
				Value:  stringutil.EnsureNotEmpty(fmt.Sprintf("%d", sf.WorkerID), "*<empty>*"),
				Inline: true,
			},
			{
				Name:   "Process ID",
				Value:  stringutil.EnsureNotEmpty(fmt.Sprintf("%d", sf.ProcessID), "*<empty>*"),
				Inline: true,
			},
			{
				Name:   "Incremental ID",
				Value:  stringutil.EnsureNotEmpty(fmt.Sprintf("%d", sf.IncrementalID), "*<empty>*"),
				Inline: true,
			},
		},
	}
}

func (c *Snowflake) embSfSp(sf snowflake.ID) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
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
				Value: stringutil.EnsureNotEmpty(snowflakenodes.GetNodeName(sf.Node()), "*<empty>*"),
			},
			{
				Name:   "Node ID",
				Value:  stringutil.EnsureNotEmpty(fmt.Sprintf("%d", sf.Node()), "*<empty>*"),
				Inline: true,
			},
			{
				Name:   "Incremental ID",
				Value:  stringutil.EnsureNotEmpty(fmt.Sprintf("%d", sf.Step()), "*<empty>*"),
				Inline: true,
			},
		},
	}
}
