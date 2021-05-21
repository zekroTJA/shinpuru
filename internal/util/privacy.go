package util

import (
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/storage"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/multierror"
)

func FlushAllGuildData(db database.Database, st storage.Storage, guildID string) (err error) {
	backups, err := db.GetBackups(guildID)
	if err != nil {
		return
	}

	reportsCount, err := db.GetReportsGuildCount(guildID)
	if err != nil {
		return
	}
	reports, err := db.GetReportsGuild(guildID, 0, reportsCount)
	if err != nil {
		return
	}

	if err = db.FlushGuildData(guildID); err != nil {
		return
	}

	mErr := multierror.New()
	for _, b := range backups {
		mErr.Append(st.DeleteObject(static.StorageBucketBackups, b.FileID))
	}
	for _, r := range reports {
		if r.AttachmehtURL != "" {
			mErr.Append(st.DeleteObject(static.StorageBucketImages, r.AttachmehtURL))
		}
	}

	return mErr.Nillify()
}
