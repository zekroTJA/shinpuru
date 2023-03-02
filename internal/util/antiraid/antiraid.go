package antiraid

import (
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/guildlog"
	"github.com/zekroTJA/shinpuru/internal/services/timeprovider"
	"github.com/zekroTJA/shinpuru/pkg/stringutil"
	"github.com/zekrotja/rogu/log"
)

var tl = log.Tagged("Antiraid")

func FlushExpired(db database.Database, gl guildlog.Logger, tp timeprovider.Provider) func() {
	gl = gl.Section("antiraid")
	return func() {
		list, err := db.GetAntiraidJoinList("")
		if err != nil {
			tl.Error().Err(err).Msg("Failed getting antiraid joinlist")
			return
		}

		now := tp.Now()
		var clearedGuilds []string
		for _, e := range list {
			if stringutil.ContainsAny(e.GuildID, clearedGuilds) {
				continue
			}
			if now.After(e.Timestamp.Add(TriggerLifetime)) {
				if err = db.FlushAntiraidJoinList(e.GuildID); err != nil && !database.IsErrDatabaseNotFound(err) {
					gl.Errorf(e.GuildID, "Failed flusing joinlist: %s", err.Error())
					tl.Error().Err(err).Field("gid", e.GuildID).Msg("Failed getting antiraid joinlist")
					continue
				}
				clearedGuilds = append(clearedGuilds, e.GuildID)
			}
		}
	}
}
