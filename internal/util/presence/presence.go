package presence

import (
	"errors"
	"fmt"
	"strings"

	"github.com/zekrotja/discordgo"
)

const (
	presenceSeperator = "|||"
	validStatus       = "dnd online idle invisible"
)

// Presence represents a presence status with a game
// message and a status string.
type Presence struct {
	Game   string `json:"game"`
	Status string `json:"status"`
}

// Unmarshal deserializes the passed raw string to
// a presence object and returns an error when the
// raw data has the wrong format.
func Unmarshal(raw string) (*Presence, error) {
	split := strings.Split(raw, presenceSeperator)
	if len(split) < 2 {
		return nil, errors.New("invalid format")
	}
	return &Presence{
		Game:   split[0],
		Status: split[1],
	}, nil
}

// Marshal produces a raw string from the presence.
func (p *Presence) Marshal() string {
	return p.Game + "|||" + p.Status
}

// ToUpdateStatusData returns a discordgo.UpdateStatusData
// from the presence object.
func (p *Presence) ToUpdateStatusData() discordgo.UpdateStatusData {
	return discordgo.UpdateStatusData{
		Activities: []*discordgo.Activity{
			{
				Name: p.Game,
				Type: discordgo.ActivityTypeGame,
			},
		},
		Status: p.Status,
	}
}

// Validate returns an error when either an invalid status
// was specified or the presence text contains the seperator
// used for serialization and deserialization.
func (p *Presence) Validate() error {
	if strings.Contains(p.Game, presenceSeperator) {
		return fmt.Errorf("`%s` is used as seperator for the settings saving so it can not be contained in the actual message.",
			presenceSeperator)
	}

	if !strings.Contains(validStatus, p.Status) {
		return fmt.Errorf("invalid status")
	}

	return nil
}
