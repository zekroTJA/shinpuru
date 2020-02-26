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
	return GroupEtc
}

func (c *CmdSnowflake) GetDomainName() string {
	return "sp.etc.snowflake"
}

func (c *CmdSnowflake) GetSubPermissionRules() []SubPermission {
	return nil
}

func (c *CmdSnowflake) Exec(args *CommandArgs) error {
	if c.rx == nil {
		c.rx = regexp.MustCompile(`(\d+)\s*([a-zA-Z]+)?`)
	}

	matches := c.rx.FindStringSubmatch(strings.Join(args.Args, " "))

	if len(matches) < 2 {
		msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
			"Please enter a Snowflake which should be calculated.")
		util.DeleteMessageLater(args.Session, msg, 5*time.Second)
		return err
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
		return c.printSfDc(args, sfAsDc)
	case snowflakeTypeShinpuru:
		return c.printSfSp(args, sfAsSp)
	}

	msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
		"Unknown snowflake type was provided.")
	util.DeleteMessageLater(args.Session, msg, 5*time.Second)
	return err
}

func (c *CmdSnowflake) printSfDc(args *CommandArgs, sf *snowflakenodes.DiscordSnowflake) error {
	emb := &discordgo.MessageEmbed{
		Title: "Snowflake Info",
		Color: 0x7289DA,
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:  "Snowflake",
				Value: fmt.Sprintf("```\n%s\n```", sf.Snowflake),
			},
			&discordgo.MessageEmbedField{
				Name:  "Type",
				Value: "Discord Snowflake ID",
			},
			&discordgo.MessageEmbedField{
				Name:  "Timestamp",
				Value: sf.Time.Format(time.RFC1123),
			},
			&discordgo.MessageEmbedField{
				Name:   "Worker ID",
				Value:  util.EnsureNotEmpty(fmt.Sprintf("%d", sf.WorkerID), "*<empty>*"),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Process ID",
				Value:  util.EnsureNotEmpty(fmt.Sprintf("%d", sf.ProcessID), "*<empty>*"),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Incremental ID",
				Value:  util.EnsureNotEmpty(fmt.Sprintf("%d", sf.IncrementalID), "*<empty>*"),
				Inline: true,
			},
		},
	}

	_, err := args.Session.ChannelMessageSendEmbed(args.Channel.ID, emb)
	return err
}

func (c *CmdSnowflake) printSfSp(args *CommandArgs, sf snowflake.ID) error {
	emb := &discordgo.MessageEmbed{
		Title: "Snowflake Info",
		Color: static.ColorEmbedDefault,
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:  "Snowflake",
				Value: fmt.Sprintf("```\n%s\n```", sf.String()),
			},
			&discordgo.MessageEmbedField{
				Name:  "Type",
				Value: "Discord Snowflake ID",
			},
			&discordgo.MessageEmbedField{
				Name: "Timestamp",
				Value: snowflakenodes.
					ParseUnixTime(int(sf.Time())).
					Format(time.RFC1123),
			},
			&discordgo.MessageEmbedField{
				Name:  "Node Name",
				Value: util.EnsureNotEmpty(snowflakenodes.GetNodeName(sf.Node()), "*<empty>*"),
			},
			&discordgo.MessageEmbedField{
				Name:   "Node ID",
				Value:  util.EnsureNotEmpty(fmt.Sprintf("%d", sf.Node()), "*<empty>*"),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Incremental ID",
				Value:  util.EnsureNotEmpty(fmt.Sprintf("%d", sf.Step()), "*<empty>*"),
				Inline: true,
			},
		},
	}

	_, err := args.Session.ChannelMessageSendEmbed(args.Channel.ID, emb)
	return err
}
