package util

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/snowflake"
)

type Report struct {
	ID         snowflake.ID
	Type       int
	GuildID    string
	ExecutorID string
	VictimID   string
	Msg        string
}

func (r *Report) AsEmbed() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: "Report " + r.ID.String(),
		Color: ColorEmbedDefault,
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
	}
}
