package listeners

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/services/config"
	"github.com/zekroTJA/shinpuru/internal/services/timeprovider"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
	"github.com/zekrotja/dgrs"
	"github.com/zekrotja/rogu/log"
)

type ListenerGuilds struct {
	cfg config.Provider
	st  *dgrs.State
	tp  timeprovider.Provider

	lockUntil *time.Time
}

func NewListenerGuildAdd(container di.Container) *ListenerGuilds {
	return &ListenerGuilds{
		cfg: container.Get(static.DiConfig).(config.Provider),
		st:  container.Get(static.DiState).(*dgrs.State),
		tp:  container.Get(static.DiTimeProvider).(timeprovider.Provider),
	}
}

func (l *ListenerGuilds) HandlerReady(s *discordgo.Session, e *discordgo.Ready) {
	now := l.tp.Now().Add(10 * time.Second)
	l.lockUntil = &now
}

func (l *ListenerGuilds) HandlerCreate(s *discordgo.Session, e *discordgo.GuildCreate) {
	limit := l.cfg.Config().Discord.GuildsLimit
	if limit < 1 {
		return
	}

	if l.lockUntil == nil || l.tp.Now().Before(*l.lockUntil) {
		return
	}

	time.Sleep(2 * time.Second)
	g, err := l.st.Guilds()
	if err != nil {
		log.Error().Tag("GuildLimit").Err(err).Msg("Failed getting guild list")
		return
	}

	log.Debug().Tag("GuildLimit").Field("ng", len(g)).Msg("Status")

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
		log.Error().Tag("GuildLimit").Err(err).Msg("Failed leaving guild")
		return
	}
}
