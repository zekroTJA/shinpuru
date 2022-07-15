package util

import (
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
	"github.com/zekroTJA/shinpuru/pkg/multierror"
	"github.com/zekrotja/dgrs"
)

func UpdateGuildMemberStats(st dgrs.IState, s discordutil.ISession) error {
	guilds, err := st.Guilds()
	if err != nil {
		return err
	}

	mErr := multierror.New()
	for _, g := range guilds {
		gwc, err := s.GuildWithCounts(g.ID)
		if err != nil {
			mErr.Append(err)
			continue
		}
		err = st.SetGuild(gwc)
		mErr.Append(err)
	}

	return mErr.Nillify()
}
