package slashcommands

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekrotja/ken"
)

type analysisMode string

const (
	analyzeMessages    analysisMode = "messages"
	analyzeAttachments analysisMode = "attachments"

	hardMessageLimit = 10000
)

type Chanstats struct{}

var (
	_ ken.SlashCommand        = (*Chanstats)(nil)
	_ permissions.PermCommand = (*Chanstats)(nil)
)

func (c *Chanstats) Name() string {
	return "channelstats"
}

func (c *Chanstats) Description() string {
	return "Get channel contribution statistics."
}

func (c *Chanstats) Version() string {
	return "1.0.0"
}

func (c *Chanstats) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *Chanstats) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "mode",
			Description: "The analysis mode.",
			Required:    true,
			Choices: []*discordgo.ApplicationCommandOptionChoice{
				{
					Name:  string(analyzeMessages),
					Value: analyzeMessages,
				},
				{
					Name:  string(analyzeAttachments),
					Value: analyzeAttachments,
				},
			},
		},
		{
			Type:         discordgo.ApplicationCommandOptionChannel,
			Name:         "channel",
			Description:  "The channel to be analyzed (defaultly current channel).",
			ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildText},
		},
		{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "limit",
			Description: "The maximum amount of messages analyzed.",
		},
	}
}

func (c *Chanstats) Domain() string {
	return "sp.chat.chanstats"
}

func (c *Chanstats) SubDomains() []permissions.SubPermission {
	return nil
}

func (c *Chanstats) Run(ctx ken.Context) (err error) {
	if err = ctx.Defer(); err != nil {
		return
	}

	mode := analysisMode(ctx.Options().GetByName("mode").StringValue())

	channelID := ctx.GetEvent().ChannelID
	if channelV, ok := ctx.Options().GetByNameOptional("channel"); ok {
		ch := channelV.ChannelValue(ctx)
		channelID = ch.ID
	}

	limit := hardMessageLimit
	if limitV, ok := ctx.Options().GetByNameOptional("limit"); ok {
		limit = int(limitV.IntValue())
	}
	if limit < 1 || limit > hardMessageLimit {
		err = ctx.FollowUpError(
			fmt.Sprintf("Message limit must be in range [1, %d]", hardMessageLimit), "").
			Send().
			Error
		return
	}

	// Generate and send a status messgae which shows the current count
	// of collected messages.
	fum := ctx.FollowUpEmbed(c.getCollectedEmbed(0)).Send()
	if err = fum.Error; err != nil {
		return
	}

	allMsgs := make([]*discordgo.Message, 0)
	var msgs []*discordgo.Message
	var lastMsgID string

	// Fetch all messages in specified channel.
	// Because only 100 messages can be fetched at once,
	// the request needs to be paginated.
	for {
		// Fetch channel messages.
		msgs, err = ctx.GetSession().ChannelMessages(channelID, 100, lastMsgID, "", "")
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
		if err = fum.EditEmbed(c.getCollectedEmbed(len(allMsgs))); err != nil {
			return
		}

		// If collected messages are equal ore above limit,
		// break further message collection
		if len(allMsgs) >= limit {
			allMsgs = allMsgs[:limit]
			break
		}
	}

	countPerUser := make(map[string]int)
	var title string

	// Setting title and countPerUser
	// depending on the analysis type.
	switch mode {
	// Type: Messages per user
	case analyzeMessages:
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
	case analyzeAttachments:
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
		if v.Label == ctx.User().Username {
			myIndex = i + 1
			myValue = v.Value
		}
	}

	// If amount of users is larger than 10,
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
		if v.Label == ctx.User().Username {
			valuesStr[i] = fmt.Sprintf("__%s__", valuesStr[i])
		}
	}

	// Assemble the final result embed and set it to the already
	// sent status embed.
	myPositionStr := "*You did not contributed any messages in this channel in the given range.*"
	if myValue > 0 {
		myPositionStr = fmt.Sprintf("%d. %s - **%.0f** *(%.2f%%)*", myIndex, ctx.User().Username, myValue, (myValue/summVals)*100)
	}
	err = fum.EditEmbed(&discordgo.MessageEmbed{
		Color:       static.ColorEmbedGreen,
		Description: fmt.Sprintf("Finished. Collected %d messages.", len(allMsgs)),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  title,
				Value: strings.Join(valuesStr, "\n"),
			},
			{
				Name:  "Your position",
				Value: myPositionStr,
			},
		},
	})
	if err != nil {
		return
	}

	// If `values` has only 1 entry, append another
	// "empty" value to bypass "invalid data range;
	// cannot be zero" error.
	if len(values) == 1 {
		values = append(values, chart.Value{
			Label: "",
			Value: 0,
		})
	}

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
	_, err = ctx.GetSession().ChannelFileSend(ctx.GetEvent().ChannelID,
		"channel_stats_chart.png", buff)

	return
}

func (c *Chanstats) getCollectedEmbed(collected int) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Color:       static.ColorEmbedGray,
		Description: fmt.Sprintf(":stopwatch:  Collected %d messages...", collected),
	}
}
