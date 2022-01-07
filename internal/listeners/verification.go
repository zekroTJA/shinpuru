package listeners

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di/v2"
	"github.com/sirupsen/logrus"
	"github.com/zekroTJA/shinpuru/internal/models"
	"github.com/zekroTJA/shinpuru/internal/services/config"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/guildlog"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/internal/util/verification"
)

type ListenerVerifications struct {
	db  database.Database
	cfg config.Provider
	gl  guildlog.Logger
}

func NewListenerVerifications(container di.Container) *ListenerVerifications {
	return &ListenerVerifications{
		db:  container.Get(static.DiDatabase).(database.Database),
		cfg: container.Get(static.DiConfig).(config.Provider),
		gl:  container.Get(static.DiGuildLog).(guildlog.Logger).Section("verification"),
	}
}

func (l *ListenerVerifications) HandlerMemberAdd(s *discordgo.Session, e *discordgo.GuildMemberAdd) {
	if !l.enabled(e.GuildID) {
		return
	}

	if verified, erroneous := l.userVerified(e.Member); verified || erroneous {
		return
	}

	err := l.db.AddVerificationQueue(&models.VerificationQueueEntry{
		GuildID:   e.GuildID,
		UserID:    e.User.ID,
		Timestamp: time.Now(),
	})
	if err != nil {
		logrus.WithError(err).WithField("gid", e.GuildID).Error("Failed adding user to verification queue")
		l.gl.Errorf(e.GuildID, "Failed adding user to verification queue: %s", err.Error())
		return
	}

	timeout := time.Now().Add(verification.ValidationTimeout)
	err = s.GuildMemberTimeout(e.GuildID, e.User.ID, &timeout)
	if err != nil {
		logrus.WithError(err).WithField("gid", e.GuildID).Error("Failed timeouting user for verification")
		l.gl.Errorf(e.GuildID, "Failed timeouting user: %s", err.Error())
		return
	}

	msg := fmt.Sprintf(
		"You need to verify your user account before you can communicate on the guild you joined.\n\n"+
			"Please go to the [**verification page**](%s/verify) and complete the captcha to verify your account.",
		l.cfg.Config().WebServer.PublicAddr,
	)
	l.sendDM(s, e.User.ID, msg, "User Verification", func(content, title string) {
		l.sendToJoinMsgChan(s, e.GuildID, e.User.ID, content, title)
	})
}

func (l *ListenerVerifications) HandlerMemberRemove(s *discordgo.Session, e *discordgo.GuildMemberRemove) {
	// if !l.enabled(e.GuildID) {
	// 	return
	// }

	// if verified, erroneous := l.userVerified(e.Member); verified || erroneous {
	// 	return
	// }
}

func (l *ListenerVerifications) enabled(guildID string) (ok bool) {
	ok, err := l.db.GetGuildVerificationRequired(guildID)
	if err != nil {
		logrus.WithError(err).WithField("gid", guildID).Error("Failed getting guild verification state from database")
		l.gl.Errorf(guildID, "Failed getting guild verification state from database: %s", err.Error())
	}
	return
}

func (l *ListenerVerifications) userVerified(e *discordgo.Member) (verified, erroneous bool) {
	verified, err := l.db.GetUserVerified(e.User.ID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		logrus.WithError(err).WithField("gid", e.GuildID).Error("Failed getting user verification state from database")
		l.gl.Errorf(e.GuildID, "Failed getting user verification state from database: %s", err.Error())
		erroneous = true
		return
	}
	return
}

func (l *ListenerVerifications) sendDM(
	s *discordgo.Session,
	userID, content, title string,
	fallback func(content, title string),
) {
	if fallback == nil {
		fallback = func(content, title string) {}
	}

	ch, err := s.UserChannelCreate(userID)
	if err != nil {
		fallback(content, title)
		return
	}
	err = util.SendEmbed(s, ch.ID, content, title, 0).Error()
	if err != nil {
		fallback(content, title)
		return
	}
}

func (l *ListenerVerifications) sendToJoinMsgChan(s *discordgo.Session, guildID, userID, content, title string) {
	chanID, _, err := l.db.GetGuildJoinMsg(guildID)
	if err != nil {
		return
	}

	s.ChannelMessageSendComplex(chanID, &discordgo.MessageSend{
		Content: "<@" + userID + ">",
		Embed: &discordgo.MessageEmbed{
			Color:       static.ColorEmbedDefault,
			Title:       title,
			Description: content,
		},
	})
}
