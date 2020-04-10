package commands

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"

	"github.com/zekroTJA/shinpuru/internal/util/static"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/util"
)

const (
	cStatsTypeMsgs = iota
	cStatsTypeAtt
)

type CmdChannelStats struct {
}

func (c *CmdChannelStats) GetInvokes() []string {
	return []string{"chanstats", "cstats"}
}

func (c *CmdChannelStats) GetDescription() string {
	return "get channel contribution statistics"
}

func (c *CmdChannelStats) GetHelp() string {
	return "`chanstats (<ChannelIdentifier>)` - get channel stats\n" +
		"`chanstats msgs (<ChannelIdentifier>)` - get channel stats by messages\n" +
		"`chanstats att (<ChannelIdentifier>)` - get channel stats by attachments"
}

func (c *CmdChannelStats) GetGroup() string {
	return GroupGuildConfig
}

func (c *CmdChannelStats) GetDomainName() string {
	return "sp.chat.chanstats"
}

func (c *CmdChannelStats) GetSubPermissionRules() []SubPermission {
	return nil
}

func (c *CmdChannelStats) Exec(args *CommandArgs) (err error) {
	channel := args.Channel
	typ := cStatsTypeMsgs

	if len(args.Args) == 1 {
		t := c.getTyp(args.Args[0])
		if t < 0 {
			channel, err = util.FetchChannel(args.Session, args.Guild.ID, args.Args[0], func(c *discordgo.Channel) bool {
				return c.Type == discordgo.ChannelTypeGuildText
			})
			if err != nil {
				return
			}
			if channel == nil {
				msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
					"Invalid command arguments. Please use `help chanstats` to see how to use this command.")
				util.DeleteMessageLater(args.Session, msg, 8*time.Second)
				return err
			}
		} else {
			typ = t
		}
	}

	if len(args.Args) == 2 {
		typ := c.getTyp(args.Args[0])
		if typ < 0 {
			msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
				"Invalid command arguments. Please use `help chanstats` to see how to use this command.")
			util.DeleteMessageLater(args.Session, msg, 8*time.Second)
			return err
		}

		channel, err = util.FetchChannel(args.Session, args.Guild.ID, args.Args[0], func(c *discordgo.Channel) bool {
			return c.Type == discordgo.ChannelTypeGuildText
		})
		if err != nil {
			return
		}
		if channel == nil {
			msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
				"Invalid command arguments. Please use `help chanstats` to see how to use this command.")
			util.DeleteMessageLater(args.Session, msg, 8*time.Second)
			return err
		}
	}

	statusMsg, err := args.Session.ChannelMessageSendEmbed(args.Channel.ID, c.getCollectedEmbed(0))
	if err != nil {
		return err
	}

	allMsgs := make([]*discordgo.Message, 0)
	var msgs []*discordgo.Message
	var lastMsgID string

	for {
		msgs, err = args.Session.ChannelMessages(channel.ID, 100, lastMsgID, "", "")
		if err != nil {
			return
		}
		if len(msgs) <= 0 {
			break
		}

		allMsgs = append(allMsgs, msgs...)
		lastMsgID = msgs[len(msgs)-1].ID
		statusMsg, err = args.Session.ChannelMessageEditEmbed(args.Channel.ID, statusMsg.ID, c.getCollectedEmbed(len(allMsgs)))
		if err != nil {
			return
		}
	}

	countPerUser := make(map[string]int)
	var title string

	// Type: Messages per user
	if typ == cStatsTypeMsgs {
		title = "Messages per User"
		for _, m := range allMsgs {
			uname := m.Author.Username
			if _, ok := countPerUser[uname]; !ok {
				countPerUser[uname] = 1
			} else {
				countPerUser[uname]++
			}
		}
	}

	// Type: Attachments per user
	if typ == cStatsTypeAtt {
		title = "Attachments per User"
		for _, m := range allMsgs {
			uname := m.Author.Username
			natt := len(m.Attachments)
			if _, ok := countPerUser[uname]; !ok {
				countPerUser[uname] = natt
			} else {
				countPerUser[uname] += natt
			}
		}
	}

	values := make([]chart.Value, len(countPerUser))
	var summVals float64
	i := 0
	for uname, c := range countPerUser {
		v := float64(c)
		values[i] = chart.Value{
			Label: uname,
			Value: v,
		}
		summVals += v
		i++
	}

	sort.Slice(values, func(i, j int) bool {
		return values[i].Value > values[j].Value
	})

	valuesStr := make([]string, len(values))
	for i, v := range values {
		valuesStr[i] = fmt.Sprintf("%d. %s - **%.0f** *(%.2f%%)*", i+1, v.Label, v.Value, (v.Value/summVals)*100)
	}

	statusMsg, err = args.Session.ChannelMessageEditEmbed(args.Channel.ID, statusMsg.ID, &discordgo.MessageEmbed{
		Color:       static.ColorEmbedGreen,
		Description: fmt.Sprintf("Finished. Collected %d messages.", len(allMsgs)),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  title,
				Value: strings.Join(valuesStr, "\n"),
			},
		},
	})

	ch := chart.BarChart{
		Title:      title,
		TitleStyle: chart.StyleShow(),
		XAxis:      chart.StyleShow(),
		YAxis: chart.YAxis{
			Style: chart.StyleShow(),
			ValueFormatter: func(v interface{}) string {
				return fmt.Sprintf("%.0f", v)
			},
			GridMajorStyle: chart.StyleShow(),
			GridMinorStyle: chart.StyleShow(),
		},
		Background: chart.Style{
			Padding: chart.Box{
				Top:    40,
				Right:  40,
				Bottom: 30,
				Left:   10,
			},
			FillColor: drawing.ColorTransparent,
		},
		Canvas: chart.Style{
			FontColor: drawing.ColorBlack,
		},
		Height: 512,
		Width:  1024,
		Bars:   values,
	}

	imgData := []byte{}
	buff := bytes.NewBuffer(imgData)
	err = ch.Render(chart.PNG, buff)
	if err != nil {
		return
	}

	_, err = args.Session.ChannelFileSend(args.Channel.ID,
		"channel_stats_chart.png", buff)

	return
}

func (c *CmdChannelStats) getTyp(arg string) int {
	switch strings.ToLower(arg) {

	case "msg", "msgs", "messages":
		return cStatsTypeMsgs

	case "att", "atts", "attachments":
		return cStatsTypeAtt
	}

	return -1
}

func (c *CmdChannelStats) getCollectedEmbed(collected int) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Color:       static.ColorEmbedGray,
		Description: fmt.Sprintf(":stopwatch:  Collected %d messages...", collected),
	}
}
