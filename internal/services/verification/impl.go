package verification

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
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
	"github.com/zekroTJA/shinpuru/pkg/multierror"
)

const timeout = 48 * time.Hour

type impl struct {
	s   discordutil.ISession
	db  database.Database
	cfg config.Provider
	gl  guildlog.Logger
}

var _ Provider = (*impl)(nil)

func New(ctn di.Container) Provider {
	return &impl{
		s:   ctn.Get(static.DiDiscordSession).(discordutil.ISession),
		db:  ctn.Get(static.DiDatabase).(database.Database),
		cfg: ctn.Get(static.DiConfig).(config.Provider),
		gl:  ctn.Get(static.DiGuildLog).(guildlog.Logger).Section("verification"),
	}
}

func (p *impl) GetEnabled(guildID string) (ok bool, err error) {
	ok, err = p.db.GetGuildVerificationRequired(guildID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		ok = false
		return
	}
	return
}

func (p *impl) SetEnabled(guildID string, enabled bool) (err error) {
	err = p.db.SetGuildVerificationRequired(guildID, enabled)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	if !enabled {
		err = p.purgeQeueue(guildID, "")
	}

	return
}

func (p *impl) IsVerified(userID string) (ok bool, err error) {
	ok, err = p.db.GetUserVerified(userID)
	if database.IsErrDatabaseNotFound(err) {
		err = nil
	}
	return
}

func (p *impl) EnqueueVerification(guildID, userID string) (err error) {
	verified, err := p.IsVerified(userID)
	if err != nil || verified {
		return
	}

	err = p.db.AddVerificationQueue(&models.VerificationQueueEntry{
		GuildID:   guildID,
		UserID:    userID,
		Timestamp: time.Now(),
	})
	if err != nil {
		return
	}

	timeout := time.Now().Add(timeout)
	err = p.s.GuildMemberTimeout(guildID, userID, &timeout)
	if err != nil {
		return
	}

	msg := fmt.Sprintf(
		"You need to verify your user account before you can communicate on the guild you joined.\n\n"+
			"Please go to the [**verification page**](%s/verify) and complete the captcha to verify your account.",
		p.cfg.Config().WebServer.PublicAddr,
	)
	p.sendDM(p.s, userID, msg, "User Verification", func(content, title string) {
		p.sendToJoinMsgChan(p.s, guildID, userID, content, title)
	})

	return
}

func (p *impl) Verify(userID string) (err error) {
	if err := p.db.SetUserVerified(userID, true); err != nil {
		return err
	}

	err = p.purgeQeueue("", userID)
	return
}

func (p *impl) KickRoutine() {
	queue, err := p.db.GetVerificationQueue("", "")
	if err != nil {
		logrus.WithError(err).Error("Failed getting verification queue from database")
		return
	}

	now := time.Now()
	for _, e := range queue {
		if now.Before(e.Timestamp.Add(timeout)) {
			continue
		}

		if err = p.s.GuildMemberTimeout(e.GuildID, e.UserID, nil); err != nil {
			logrus.WithError(err).Error("Failed removing member timeout")
			p.gl.Errorf(e.GuildID, "Failed removing member timeout: %s", err.Error())
		}
		if err = p.s.GuildMemberDelete(e.GuildID, e.UserID); err != nil {
			logrus.WithError(err).Error("Failed kicking member")
			p.gl.Errorf(e.GuildID, "Failed kicking member: %s", err.Error())
		}
		if _, err = p.db.RemoveVerificationQueue(e.GuildID, e.UserID); err != nil {
			logrus.WithError(err).Error("Failed removing member from verification queue")
			p.gl.Errorf(e.GuildID, "Failed removing member from verification queue: %s", err.Error())
		}
	}
}

func (p *impl) purgeQeueue(guildID, userID string) (err error) {
	queue, err := p.db.GetVerificationQueue(guildID, userID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	mErr := multierror.New()
	for _, e := range queue {
		ok, err := p.db.RemoveVerificationQueue(e.GuildID, e.UserID)
		mErr.Append(err)
		if ok {
			mErr.Append(p.s.GuildMemberTimeout(e.GuildID, e.UserID, nil))
		}
	}

	return mErr.Nillify()
}

func (p *impl) sendDM(
	s discordutil.ISession,
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

func (p *impl) sendToJoinMsgChan(s discordutil.ISession, guildID, userID, content, title string) {
	chanID, _, err := p.db.GetGuildJoinMsg(guildID)
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
