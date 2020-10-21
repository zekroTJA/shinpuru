package listeners

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/ratelimit"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
)

type ListenerAntiraid struct {
	db database.Database

	limiters map[string]*ratelimit.Limiter
}

func NewListenerAntiraid(db database.Database) *ListenerAntiraid {
	return &ListenerAntiraid{
		db:       db,
		limiters: make(map[string]*ratelimit.Limiter),
	}
}

func (l *ListenerAntiraid) Handler(s *discordgo.Session, e *discordgo.GuildMemberAdd) {
	ok, limit, burst := l.getGuildSettings(e.GuildID)
	if !ok {
		if _, ok := l.limiters[e.GuildID]; ok {
			delete(l.limiters, e.GuildID)
		}
		return
	}

	limitDur := time.Duration(limit) * time.Second

	limiter, ok := l.limiters[e.GuildID]
	if !ok || limiter == nil {
		limiter = ratelimit.NewLimiter(limitDur, burst)
		l.limiters[e.GuildID] = limiter
	} else {
		if limiter.Burst() != burst {
			limiter.SetBurst(burst)
		}
		if limiter.Limit() != limitDur {
			limiter.SetLimit(limitDur)
		}
	}

	if limiter.Allow() {
		return
	}

	verificationLvl := discordgo.VerificationLevelVeryHigh
	_, err := s.GuildEdit(e.GuildID, discordgo.GuildParams{
		VerificationLevel: &verificationLvl,
	})

	guild, err := discordutil.GetGuild(s, e.GuildID)
	if err != nil {
		util.Log.Errorf("failed getting guild (gid: %s): %s", e.GuildID, err.Error())
		return
	}

	alertDescrition := fmt.Sprintf(
		"Following guild you are admin on is currently being raided!\n\n"+
			"**%s (`%s`)**", guild.Name, e.GuildID)
	if err != nil {
		alertDescrition = fmt.Sprintf("%s\n\n"+
			"**Attention:** Failed to raise guilds verification level because "+
			"following error occured:\n```\n%s\n```", alertDescrition, err.Error())
	}

	emb := &discordgo.MessageEmbed{
		Title:       "⚠ GUILD RAID ALERT",
		Description: alertDescrition,
		Color:       static.ColorEmbedOrange,
	}

	members, err := discordutil.GetMembers(s, e.GuildID)
	if err != nil {
		util.Log.Errorf("failed getting guild members (gid: %s): %s", e.GuildID, err.Error())
		return
	}

	for _, m := range members {
		if discordutil.IsAdmin(guild, m) || guild.OwnerID == m.User.ID {
			ch, err := s.UserChannelCreate(m.User.ID)
			if err != nil {
				continue
			}
			s.ChannelMessageSendEmbed(ch.ID, emb)
		}
	}

	if chanID, _ := l.db.GetGuildModLog(e.GuildID); chanID != "" {
		s.ChannelMessageSendEmbed(chanID, &discordgo.MessageEmbed{
			Color: static.ColorEmbedOrange,
			Title: "⚠ GUILD RAID ALERT",
			Description: "Because an atypical burst of members joined the guild, " +
				"the guilds verification level was raised to `verry high` and all admins " +
				"were informed.",
		})
	}
}

func (l *ListenerAntiraid) getGuildSettings(gid string) (ok bool, limit, burst int) {
	var err error
	var state bool

	state, err = l.db.GetAntiraidState(gid)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		util.Log.Errorf("failed getting antiraid state (gid: %s): %s", gid, err.Error())
		return
	}
	if !state {
		return
	}

	limit, err = l.db.GetAntiraidRegeneration(gid)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		util.Log.Errorf("failed getting antiraid regeneration (gid: %s): %s", gid, err.Error())
		return
	}
	if limit < 1 {
		return
	}

	burst, err = l.db.GetAntiraidBurst(gid)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		util.Log.Errorf("failed getting antiraid burst (gid: %s): %s", gid, err.Error())
		return
	}
	if burst < 1 {
		return
	}

	ok = true

	return
}
