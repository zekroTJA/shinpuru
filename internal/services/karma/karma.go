package karma

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/util/static"
)

// Service provides functionalities to check karma state,
// karma blocklist and alter karma of a user.
type Service struct {
	s  *discordgo.Session
	db database.Database
}

// NewKarmaService initializes a new Service
// instance.
func NewKarmaService(container di.Container) (k *Service) {
	k = &Service{}

	k.s = container.Get(static.DiDiscordSession).(*discordgo.Session)
	k.db = container.Get(static.DiDatabase).(database.Database)

	return
}

// GetState returns the current enabled/disabled state
// of the karma system on the specified guild.
func (k *Service) GetState(guildID string) (ok bool, err error) {
	ok, err = k.db.GetKarmaState(guildID)
	if database.IsErrDatabaseNotFound(err) {
		err = nil
	}
	return
}

// IsBlockListed returns true if the passed user on the
// specified guild is blocked from gaining or giving
// karma.
func (k *Service) IsBlockListed(guildID, userID string) (isBlocklisted bool, err error) {
	isBlocklisted, err = k.db.IsKarmaBlockListed(guildID, userID)
	if database.IsErrDatabaseNotFound(err) {
		err = nil
	}
	return
}

// Update adds or removes karma of the given value of the
// specified user.
func (k *Service) Update(guildID, userID string, value int) (err error) {
	err = k.db.UpdateKarma(userID, guildID, value)
	return
}

// CheckAndUpdate is shorthand for GetState, IsBlockListed
// and Update in one single pipe.
func (k *Service) CheckAndUpdate(guildID string, object *discordgo.User, value int) (ok bool, err error) {
	if object.Bot {
		return
	}

	enabled, err := k.GetState(guildID)
	if !enabled || err != nil {
		return
	}

	isBLocklisted, err := k.IsBlockListed(guildID, object.ID)
	if isBLocklisted || err != nil {
		return
	}

	err = k.db.UpdateKarma(object.ID, guildID, value)
	ok = err == nil
	return
}
