package karma

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/services/database"
)

func Alter(db database.Database, guildID string, object *discordgo.User, value int) (ok bool, err error) {
	if object.Bot {
		return
	}

	enabled, err := db.GetKarmaState(guildID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return
	}
	if !enabled {
		return
	}

	isBlacklisted, err := db.IsKarmaBlockListed(guildID, object.ID)
	if err != nil {
		return
	}
	if isBlacklisted {
		return
	}

	err = db.UpdateKarma(object.ID, guildID, value)
	ok = err == nil
	return
}
