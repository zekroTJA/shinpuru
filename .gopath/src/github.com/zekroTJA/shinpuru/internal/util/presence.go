package util

import (
	"errors"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Presence struct {
	Game   string
	Status string
}

func UnmarshalPresence(raw string) (*Presence, error) {
	split := strings.Split(raw, "|||")
	if len(split) < 2 {
		return nil, errors.New("invalid format")
	}
	return &Presence{
		Game:   split[0],
		Status: split[1],
	}, nil
}

func (p *Presence) Marshal() string {
	return p.Game + "|||" + p.Status
}

func (p *Presence) ToUpdateStatusData() discordgo.UpdateStatusData {
	return discordgo.UpdateStatusData{
		Game: &discordgo.Game{
			Name: p.Game,
		},
		Status: p.Status,
	}
}
