package util

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	PresenceSeperator = "|||"
	validStatus       = "dnd online idle invisible"
)

type Presence struct {
	Game   string `json:"game"`
	Status string `json:"status"`
}

func UnmarshalPresence(raw string) (*Presence, error) {
	split := strings.Split(raw, PresenceSeperator)
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

func (p *Presence) Validate() error {
	if strings.Contains(p.Game, PresenceSeperator) {
		return fmt.Errorf("`%s` is used as seperator for the settings saving so it can not be contained in the actual message.",
			PresenceSeperator)
	}

	if !strings.Contains(validStatus, p.Status) {
		return fmt.Errorf("invalid status")
	}

	return nil
}
