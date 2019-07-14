package util

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/snowflake"
)

type Report struct {
	ID            snowflake.ID
	Type          int
	GuildID       string
	ExecutorID    string
	VictimID      string
	Msg           string
	AttachmehtURL string
}

func (r *Report) GetTimestamp() time.Time {
	return time.Unix(r.ID.Time()/1000, 0)
}

func (r *Report) AsEmbed() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: "Case " + r.ID.String(),
		Color: ReportColors[r.Type],
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Inline: true,
				Name:   "Executor",
				Value:  fmt.Sprintf("<@%s>", r.ExecutorID),
			},
			&discordgo.MessageEmbedField{
				Inline: true,
				Name:   "Victim",
				Value:  fmt.Sprintf("<@%s>", r.VictimID),
			},
			&discordgo.MessageEmbedField{
				Name:  "Type",
				Value: ReportTypes[r.Type],
			},
			&discordgo.MessageEmbedField{
				Name:  "Description",
				Value: r.Msg,
			},
		},
		Timestamp: r.GetTimestamp().Format("2006-01-02T15:04:05.000Z"),
		Image: &discordgo.MessageEmbedImage{
			URL: r.AttachmehtURL,
		},
	}
}

func (r *Report) AsEmbedField() *discordgo.MessageEmbedField {
	attachmentTxt := ""
	if r.AttachmehtURL != "" {
		attachmentTxt = fmt.Sprintf("Attachment: [[open](%s)]\n", r.AttachmehtURL)
	}

	return &discordgo.MessageEmbedField{
		Name: "Case " + r.ID.String(),
		Value: fmt.Sprintf("Time: %s\nExecutor: <@%s>\nVictim: <@%s>\nType: `%s`\n%s__Reason__:\n%s",
			r.GetTimestamp().Format("2006/01/02 15:04:05"), r.ExecutorID, r.VictimID, ReportTypes[r.Type], attachmentTxt, r.Msg),
	}
}
