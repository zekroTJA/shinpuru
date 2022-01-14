package antiraid

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/guildlog"
	"github.com/zekroTJA/shinpuru/pkg/stringutil"
)

func FlushExpired(db database.Database, gl guildlog.Logger) func() {
	gl = gl.Section("antiraid")
	return func() {
		list, err := db.GetAntiraidJoinList("")
		if err != nil {
			logrus.WithError(err).Error("Failed getting antiraid joinlist")
			return
		}

		now := time.Now()
		var clearedGuilds []string
		for _, e := range list {
			if stringutil.ContainsAny(e.GuildID, clearedGuilds) {
				continue
			}
			if now.After(e.Timestamp.Add(TriggerLifetime)) {
				if err = db.FlushAntiraidJoinList(e.GuildID); err != nil && !database.IsErrDatabaseNotFound(err) {
					gl.Errorf(e.GuildID, "Failed flusing joinlist: %s", err.Error())
					logrus.WithError(err).WithField("gid", e.GuildID).Error("Failed getting antiraid joinlist")
					continue
				}
				clearedGuilds = append(clearedGuilds, e.GuildID)
			}
		}
	}
}
