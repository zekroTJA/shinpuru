package privacy

import (
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/internal/util/vote"
	"github.com/zekroTJA/shinpuru/pkg/multierror"
)

func FlushAllGuildData(
	s Session,
	db Database,
	st Storage,
	state State,
	guildID string,
) (err error) {
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

	for _, v := range vote.VotesRunning {
		if v.GuildID == guildID {
			v.Close(s, vote.VoteStateClosedNC)
		}
	}

	if err = db.FlushGuildData(guildID); err != nil {
		return
	}

	if err = state.RemoveGuild(guildID, true); err != nil {
		return
	}

	mErr := multierror.New()
	for _, b := range backups {
		mErr.Append(st.DeleteObject(static.StorageBucketBackups, b.FileID))
	}
	for _, r := range reports {
		if r.AttachmentURL != "" {
			mErr.Append(st.DeleteObject(static.StorageBucketImages, r.AttachmentURL))
		}
	}

	return mErr.Nillify()
}

func FlushAllUserData(
	db Database,
	state State,
	userID string,
) (res map[string]int, err error) {
	res, err = db.FlushUserData(userID)
	if err != nil {
		return
	}

	guildIDs, err := state.UserGuilds(userID)
	if err != nil {
		return
	}

	for _, gid := range guildIDs {
		if err = state.RemoveMember(gid, userID); err != nil {
			return
		}
	}

	err = state.RemoveUser(userID)
	return
}
