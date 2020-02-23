package report

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/snowflake"
	"github.com/zekroTJA/shinpuru/internal/util/imgstore"
	"github.com/zekroTJA/shinpuru/internal/util/static"
)

type Report struct {
	ID            snowflake.ID `json:"id"`
	Type          int          `json:"type"`
	GuildID       string       `json:"guild_id"`
	ExecutorID    string       `json:"executor_id"`
	VictimID      string       `json:"victim_id"`
	Msg           string       `json:"message"`
	AttachmehtURL string       `json:"attachment_url"`
}

func (r *Report) GetTimestamp() time.Time {
	return time.Unix(r.ID.Time()/1000, 0)
}

func (r *Report) AsEmbed(publicAddr string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: "Case " + r.ID.String(),
		Color: static.ReportColors[r.Type],
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
				Value: static.ReportTypes[r.Type],
			},
			&discordgo.MessageEmbedField{
				Name:  "Description",
				Value: r.Msg,
			},
		},
		Timestamp: r.GetTimestamp().Format("2006-01-02T15:04:05.000Z"),
		Image: &discordgo.MessageEmbedImage{
			URL: imgstore.GetImageLink(r.AttachmehtURL, publicAddr),
		},
	}
}

func (r *Report) AsEmbedField(publicAddr string) *discordgo.MessageEmbedField {
	attachmentTxt := ""
	if r.AttachmehtURL != "" {
		attachmentTxt = fmt.Sprintf("Attachment: [[open](%s)]\n", imgstore.GetImageLink(r.AttachmehtURL, publicAddr))
	}

	return &discordgo.MessageEmbedField{
		Name: "Case " + r.ID.String(),
		Value: fmt.Sprintf("Time: %s\nExecutor: <@%s>\nVictim: <@%s>\nType: `%s`\n%s__Reason__:\n%s",
			r.GetTimestamp().Format("2006/01/02 15:04:05"), r.ExecutorID, r.VictimID, static.ReportTypes[r.Type], attachmentTxt, r.Msg),
	}
}
