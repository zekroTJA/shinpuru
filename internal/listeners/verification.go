package listeners

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/guildlog"
	"github.com/zekroTJA/shinpuru/internal/services/verification"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekrotja/rogu"
	"github.com/zekrotja/rogu/log"
)

type ListenerVerifications struct {
	db  database.Database
	vs  verification.Provider
	gl  guildlog.Logger
	log rogu.Logger
}

func NewListenerVerifications(container di.Container) *ListenerVerifications {
	return &ListenerVerifications{
		db:  container.Get(static.DiDatabase).(database.Database),
		vs:  container.Get(static.DiVerification).(verification.Provider),
		gl:  container.Get(static.DiGuildLog).(guildlog.Logger).Section("verification"),
		log: log.Tagged("Verification"),
	}
}

func (l *ListenerVerifications) HandlerMemberAdd(s *discordgo.Session, e *discordgo.GuildMemberAdd) {
	if !l.enabled(e.GuildID) {
		return
	}

	err := l.vs.EnqueueVerification(*e.Member)
	if err != nil {
		l.log.Error().Err(err).Field("gid", e.GuildID).Msg("Failed enqueueing user to verification queue")
		l.gl.Errorf(e.GuildID, "Failed enqueueing user to verification queue: %s", err.Error())
	}
}

func (l *ListenerVerifications) HandlerMemberRemove(s *discordgo.Session, e *discordgo.GuildMemberRemove) {

}

func (l *ListenerVerifications) enabled(guildID string) (ok bool) {
	ok, err := l.db.GetGuildVerificationRequired(guildID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		l.log.Error().Err(err).Field("gid", guildID).Msg("Failed getting guild verification state from database")
		l.gl.Errorf(guildID, "Failed getting guild verification state from database: %s", err.Error())
	}
	return
}
