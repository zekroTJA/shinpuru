package report

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/snowflake"
	"github.com/zekroTJA/shinpuru/internal/util/imgstore"
	"github.com/zekroTJA/shinpuru/internal/util/static"
)

// Report describes a report object.
type Report struct {
	ID            snowflake.ID `json:"id"`
	Type          int          `json:"type"`
	GuildID       string       `json:"guild_id"`
	ExecutorID    string       `json:"executor_id"`
	VictimID      string       `json:"victim_id"`
	Msg           string       `json:"message"`
	AttachmehtURL string       `json:"attachment_url"`
}

// GetTimestamp returns the time stamp when the
// report was generated from the reports ID
// snowflake.
func (r *Report) GetTimestamp() time.Time {
	return time.Unix(r.ID.Time()/1000, 0)
}

// AsEmbed creates a discordgo.Embed from the
// report. publicAddr is passed to generate a
// public link for a potential report attachment
// to be displayed in the embeds image section.
func (r *Report) AsEmbed(publicAddr string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: "Case " + r.ID.String(),
		Color: static.ReportColors[r.Type],
		Fields: []*discordgo.MessageEmbedField{
			{
				Inline: true,
				Name:   "Executor",
				Value:  fmt.Sprintf("<@%s>", r.ExecutorID),
			},
			{
				Inline: true,
				Name:   "Victim",
				Value:  fmt.Sprintf("<@%s>", r.VictimID),
			},
			{
				Name:  "Type",
				Value: static.ReportTypes[r.Type],
			},
			{
				Name:  "Description",
				Value: r.Msg,
			},
		},
		Timestamp: r.GetTimestamp().Format("2006-01-02T15:04:05.000Z"),
		Image: &discordgo.MessageEmbedImage{
			URL: imgstore.GetLink(r.AttachmehtURL, publicAddr),
		},
	}
}

// AsEmbedField creates a discordgo.MessageEmbedField from
// the report. publicAddr is passed to generate a publicly
// vailable link embeded in the embed field.
func (r *Report) AsEmbedField(publicAddr string) *discordgo.MessageEmbedField {
	attachmentTxt := ""
	if r.AttachmehtURL != "" {
		attachmentTxt = fmt.Sprintf("Attachment: [[open](%s)]\n", imgstore.GetLink(r.AttachmehtURL, publicAddr))
	}

	return &discordgo.MessageEmbedField{
		Name: "Case " + r.ID.String(),
		Value: fmt.Sprintf("Time: %s\nExecutor: <@%s>\nVictim: <@%s>\nType: `%s`\n%s__Reason__:\n%s",
			r.GetTimestamp().Format("2006/01/02 15:04:05"), r.ExecutorID, r.VictimID, static.ReportTypes[r.Type], attachmentTxt, r.Msg),
	}
}
