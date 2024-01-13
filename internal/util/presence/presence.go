package presence

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/pkg/stringutil"
)

type Status string

const (
	StatusOnline    Status = "online"
	StatusDnD       Status = "dnd"
	StatusIdle      Status = "idle"
	StatusInvisible Status = "invisible"

	presenceSeperator = "|||"
)

var (
	validStatus = []string{string(StatusDnD), string(StatusIdle), string(StatusInvisible), string(StatusOnline)}
)

// Presence represents a presence status with a game
// message and a status string.
type Presence struct {
	Game   string `json:"game"`
	Status Status `json:"status"`
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
		Status: Status(split[1]),
	}, nil
}

// Marshal produces a raw string from the presence.
func (p *Presence) Marshal() string {
	return p.Game + "|||" + string(p.Status)
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
		Status: string(p.Status),
	}
}

// Validate returns an error when either an invalid status
// was specified or the presence text contains the seperator
// used for serialization and deserialization.
func (p *Presence) Validate() error {
	if strings.Contains(p.Game, presenceSeperator) {
		return fmt.Errorf("`%s` is used as seperator for the settings saving so it can not be contained in the actual message",
			presenceSeperator)
	}

	if !stringutil.ContainsAny(string(p.Status), validStatus) {
		return fmt.Errorf("invalid status")
	}

	return nil
}
