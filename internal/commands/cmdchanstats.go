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

	// Check command argument 0
	// If no type is specified, `cStatsTypeMsgs` stays unchanged
	// and channel will be tried to be fetched by first argument.
	if len(args.Args) == 1 {
		t := c.getTyp(args.Args[0])
		if t == -1 {
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

	// If a second argument is passed and if the first one is
	// not a type specifier, this will be interpreted as error.
	// From the first argument, the channel will be tried to
	// be fetched.
	if len(args.Args) == 2 {
		typ := c.getTyp(args.Args[0])
		if typ == -1 {
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

	// Generate and send a status messgae which shows the current count
	// of collected messages.
	statusMsg, err := args.Session.ChannelMessageSendEmbed(args.Channel.ID, c.getCollectedEmbed(0))
	if err != nil {
		return err
	}

	allMsgs := make([]*discordgo.Message, 0)
	var msgs []*discordgo.Message
	var lastMsgID string

	// Fetch all messages in specified channel.
	// Because only 100 messages can be fetched at once,
	// the request needs to be paginated.
	for {
		// Fetch channel messages.
		msgs, err = args.Session.ChannelMessages(channel.ID, 100, lastMsgID, "", "")
		if err != nil {
			return
		}
		if len(msgs) <= 0 {
			break
		}

		// Append messages to list and set last message ID.
		allMsgs = append(allMsgs, msgs...)
		lastMsgID = msgs[len(msgs)-1].ID
		// Update status message.
		statusMsg, err = args.Session.ChannelMessageEditEmbed(args.Channel.ID, statusMsg.ID, c.getCollectedEmbed(len(allMsgs)))
		if err != nil {
			return
		}
	}

	countPerUser := make(map[string]int)
	var title string

	// Setting title and countPerUser
	// depending on the analysis type.
	switch typ {
	// Type: Messages per user
	case cStatsTypeMsgs:
		title = "Messages per User"
		for _, m := range allMsgs {
			uname := m.Author.Username
			if _, ok := countPerUser[uname]; !ok {
				countPerUser[uname] = 1
			} else {
				countPerUser[uname]++
			}
		}

	// Type: Attachments per user
	case cStatsTypeAtt:
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

	// Transform the `countsPerUser` map to an
	// array of chart.Value.
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

	// Sort the resulting array by value descending.
	sort.Slice(values, func(i, j int) bool {
		return values[i].Value > values[j].Value
	})

	// Figure out at which position the command
	// executor themself is and set the contribution
	// value of them.
	var myIndex int
	var myValue float64
	for i, v := range values {
		if v.Label == args.User.Username {
			myIndex = i + 1
			myValue = v.Value
		}
	}

	// If ammount of users is larger than 10,
	// slice the results by 10.
	if len(values) > 10 {
		title += " (Top 10)"
		values = values[:10]
	}

	// Make an array of strings for the result top list embed
	// and put in all formatted values.
	valuesStr := make([]string, len(values))
	for i, v := range values {
		valuesStr[i] = fmt.Sprintf("%d. %s - **%.0f** *(%.2f%%)*", i+1, v.Label, v.Value, (v.Value/summVals)*100)
		if v.Label == args.User.Username {
			valuesStr[i] = fmt.Sprintf("__%s__", valuesStr[i])
		}
	}

	// Assemble the final result embed and set it to the already
	// sent status embed.
	statusMsg, err = args.Session.ChannelMessageEditEmbed(args.Channel.ID, statusMsg.ID, &discordgo.MessageEmbed{
		Color:       static.ColorEmbedGreen,
		Description: fmt.Sprintf("Finished. Collected %d messages.", len(allMsgs)),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  title,
				Value: strings.Join(valuesStr, "\n"),
			},
			{
				Name:  "Your position",
				Value: fmt.Sprintf("%d. %s - **%.0f** *(%.2f%%)*", myIndex, args.User.Username, myValue, (myValue/summVals)*100),
			},
		},
	})

	// Create and assemble GoChart chart
	// from collected values.
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

	// Render chart to byte buffer.
	imgData := []byte{}
	buff := bytes.NewBuffer(imgData)
	err = ch.Render(chart.PNG, buff)
	if err != nil {
		return
	}

	// Send the rendered chart from buffer into the channel.
	_, err = args.Session.ChannelFileSend(args.Channel.ID,
		"channel_stats_chart.png", buff)

	return
}

// getType returns the type number by
// passed argument string.
func (c *CmdChannelStats) getTyp(arg string) int {
	switch strings.ToLower(arg) {

	case "msg", "msgs", "messages":
		return cStatsTypeMsgs

	case "att", "atts", "attachments":
		return cStatsTypeAtt
	}

	return -1
}

// getCollectedEmbed returns a discordgo.MessageEmbed displaying the
// ammount of processed messages.
func (c *CmdChannelStats) getCollectedEmbed(collected int) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Color:       static.ColorEmbedGray,
		Description: fmt.Sprintf(":stopwatch:  Collected %d messages...", collected),
	}
}
