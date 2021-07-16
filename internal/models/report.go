package models

import (
	"fmt"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/snowflake"
	"github.com/zekroTJA/shinpuru/internal/util/imgstore"
)

type Type int

const (
	TypeKick Type = iota
	TypeBan
	TypeMute
	TypeWarn
	TypeAd

	TypeMax       = iota - 1
	TypesReserved = TypeWarn
)

var (
	ReportTypes = []string{
		"KICK", // 0
		"BAN",  // 1
		"MUTE", // 2
		"WARN", // 3
		"AD",   // 4
	}

	ReportColors = []int{
		0xD81B60, // 0
		0xe53935, // 1
		0x009688, // 2
		0xFB8C00, // 3
		0x8E24AA, // 4
	}
)

func TypeFromString(s string) (typ Type, err error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return
	}
	if i < 0 || i > TypeMax {
		err = fmt.Errorf("type out of bounds ([0..%d])", TypeMax)
	} else {
		typ = Type(i)
	}
	return
}

// Report describes a report object.
type Report struct {
	ID            snowflake.ID `json:"id"`
	Type          Type         `json:"type"`
	GuildID       string       `json:"guild_id"`
	ExecutorID    string       `json:"executor_id"`
	VictimID      string       `json:"victim_id"`
	Msg           string       `json:"message"`
	AttachmehtURL string       `json:"attachment_url"`
	Timeout       *time.Time   `json:"timeout"`
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
	emb := &discordgo.MessageEmbed{
		Title: "Case " + r.ID.String(),
		Color: ReportColors[r.Type],
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
				Value: ReportTypes[r.Type],
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

	if r.Timeout != nil {
		emb.Fields = append(emb.Fields, &discordgo.MessageEmbedField{
			Name:  "Expires",
			Value: r.Timeout.Format("2006-01-02T15:04:05.000Z"),
		})
	}

	if r.Type == TypeBan {
		emb.Description = fmt.Sprintf(
			"If you want to submit an unbanrequest, you can do this [here](%s/unbanme).", publicAddr)
	}

	return emb
}

// AsEmbedField creates a discordgo.MessageEmbedField from
// the report. publicAddr is passed to generate a publicly
// vailable link embedded in the embed field.
func (r *Report) AsEmbedField(publicAddr string) *discordgo.MessageEmbedField {
	attachmentTxt := ""
	if r.AttachmehtURL != "" {
		attachmentTxt = fmt.Sprintf("Attachment: [[open](%s)]\n", imgstore.GetLink(r.AttachmehtURL, publicAddr))
	}

	return &discordgo.MessageEmbedField{
		Name: "Case " + r.ID.String(),
		Value: fmt.Sprintf("Time: %s\nExecutor: <@%s>\nVictim: <@%s>\nType: `%s`\n%s__Reason__:\n%s",
			r.GetTimestamp().Format("2006/01/02 15:04:05"), r.ExecutorID, r.VictimID, ReportTypes[r.Type], attachmentTxt, r.Msg),
	}
}
