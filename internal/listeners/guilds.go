package listeners

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di/v2"
	"github.com/sirupsen/logrus"
	"github.com/zekroTJA/shinpuru/internal/services/config"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
	"github.com/zekrotja/dgrs"
)

type ListenerGuilds struct {
	cfg config.Provider
	st  *dgrs.State

	lockUntil *time.Time
}

func NewListenerGuildAdd(container di.Container) *ListenerGuilds {
	return &ListenerGuilds{
		cfg: container.Get(static.DiConfig).(config.Provider),
		st:  container.Get(static.DiState).(*dgrs.State),
	}
}

func (l *ListenerGuilds) HandlerReady(s *discordgo.Session, e *discordgo.Ready) {
	now := time.Now().Add(10 * time.Second)
	l.lockUntil = &now
}

func (l *ListenerGuilds) HandlerCreate(s *discordgo.Session, e *discordgo.GuildCreate) {
	limit := l.cfg.Config().Discord.GuildsLimit
	if limit < 1 {
		return
	}

	if l.lockUntil == nil || time.Now().Before(*l.lockUntil) {
		return
	}

	time.Sleep(2 * time.Second)
	g, err := l.st.Guilds()
	if err != nil {
		logrus.WithError(err).Error("GUILDLIMIT :: failed getting guild list")
		return
	}

	logrus.WithField("ng", len(g)).WithField("limit", limit).Debug("GUILDLIMIT :: status")
	if len(g) <= limit {
		return
	}

	discordutil.SendDMEmbed(s, e.Guild.OwnerID, &discordgo.MessageEmbed{
		Color: static.ColorEmbedOrange,
		Title: "Guild Limit Reached",
		Description: "This instance of shinpuru is limited to a specific count of guilds which is now exceeded. " +
			"Therefore, the bot has been removed from your guild. Sorry for that inconvenience. Please try again later. 😔",
	})

	if err = s.GuildLeave(e.Guild.ID); err != nil {
		logrus.WithError(err).Error("GUILDLIMIT :: failed leaving guild")
		return
	}
}
