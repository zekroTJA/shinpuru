package inits

import (
	"math/rand"
	"strings"
	"sync/atomic"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/snowflake"
	"github.com/sarulabs/di/v2"
	"github.com/sirupsen/logrus"
	"github.com/zekroTJA/shinpuru/internal/listeners"
	"github.com/zekroTJA/shinpuru/internal/models"
	"github.com/zekroTJA/shinpuru/internal/services/config"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/snowflakenodes"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekrotja/dgrs"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func InitDiscordBotSession(container di.Container) (release func()) {
	snowflake.Epoch = static.DefEpoche
	err := snowflakenodes.Setup()
	if err != nil {
		logrus.WithError(err).Fatal("Failed setting up snowflake nodes")
	}

	snowflakenodes.NodesReport = make([]*snowflake.Node, len(models.ReportTypes))
	for i, t := range models.ReportTypes {
		if snowflakenodes.NodesReport[i], err = snowflakenodes.RegisterNode(i, "report."+strings.ToLower(t)); err != nil {
			logrus.WithError(err).Fatal("Failed setting up snowflake nodes")
		}
	}

	session := container.Get(static.DiDiscordSession).(*discordgo.Session)
	cfg := container.Get(static.DiConfig).(config.Provider)

	session.Token = "Bot " + cfg.Config().Discord.Token
	session.Identify.Intents = discordgo.MakeIntent(static.Intents)
	session.StateEnabled = false

	if shardCfg := cfg.Config().Discord.Sharding; shardCfg.Total > 1 {
		st := container.Get(static.DiState).(*dgrs.State)

		var id int
		if shardCfg.AutoID {
			d := time.Duration(rand.Int63n(int64(5 * time.Second)))
			logrus.
				WithField("d", d.Round(time.Millisecond).String()).
				Info("Sleeping before retrieving shard ID")
			time.Sleep(d)
			if id, err = st.ReserveShard(); err != nil {
				logrus.WithError(err).Fatal("Failed receiving alive shards from state")
			}
		} else {
			id = shardCfg.ID
			if id < 0 || id >= shardCfg.Total {
				logrus.Fatalf("Shard ID must be in range [0, %d)", shardCfg.Total)
			}
			if _, err = st.ReserveShard(id); err != nil {
				logrus.WithError(err).Fatal("Failed reserving predefined shard")
			}
		}

		release = func() {
			logrus.WithField("id", id).Info("Releasing shard ID")
			if err = st.ReleaseShard(id); err != nil {
				logrus.WithError(err).Error("Failed releasing shard ID")
			}
		}

		logrus.
			WithField("id", id).
			WithField("total", shardCfg.Total).
			Info("Running in sharded mode")
		session.Identify.Shard = &[2]int{id, shardCfg.Total}
	}

	listenerInviteBlock := listeners.NewListenerInviteBlock(container)
	listenerGhostPing := listeners.NewListenerGhostPing(container)
	listenerColors := listeners.NewColorListener(container)

	listenerJDoodle, err := listeners.NewListenerJdoodle(container)
	if err != nil {
		logrus.WithError(err).Fatal("Failed setting up code execution listener")
	}

	listenerStarboard := listeners.NewListenerStarboard(container)
	listenerVerification := listeners.NewListenerVerifications(container)

	session.AddHandler(listeners.NewListenerReady(container).Handler)
	session.AddHandler(listeners.NewListenerMemberAdd(container).Handler)
	session.AddHandler(listeners.NewListenerMemberRemove(container).Handler)
	session.AddHandler(listeners.NewListenerVote(container).Handler)
	session.AddHandler(listeners.NewListenerChannelCreate(container).Handler)
	session.AddHandler(listeners.NewListenerVoiceUpdate(container).Handler)
	session.AddHandler(listeners.NewListenerKarma(container).Handler)
	session.AddHandler(listeners.NewListenerAntiraid(container).HandlerMemberAdd)
	session.AddHandler(listeners.NewListenerBotMention(container).Listener)

	session.AddHandler(listenerGhostPing.HandlerMessageCreate)
	session.AddHandler(listenerGhostPing.HandlerMessageDelete)
	session.AddHandler(listenerInviteBlock.HandlerMessageSend)
	session.AddHandler(listenerInviteBlock.HandlerMessageEdit)

	session.AddHandler(listenerJDoodle.HandlerMessageCreate)
	session.AddHandler(listenerJDoodle.HandlerMessageUpdate)
	session.AddHandler(listenerJDoodle.HandlerReactionAdd)

	session.AddHandler(listenerColors.HandlerMessageCreate)
	session.AddHandler(listenerColors.HandlerMessageEdit)
	session.AddHandler(listenerColors.HandlerMessageReaction)

	session.AddHandler(listenerStarboard.ListenerReactionAdd)
	session.AddHandler(listenerStarboard.ListenerReactionRemove)

	session.AddHandler(listenerVerification.HandlerMemberAdd)
	session.AddHandler(listenerVerification.HandlerMemberRemove)

	session.AddHandler(func(s *discordgo.Session, e *discordgo.MessageCreate) {
		atomic.AddUint64(&util.StatsMessagesAnalysed, 1)
	})

	if cfg.Config().Metrics.Enable {
		session.AddHandler(listeners.NewListenerMetrics().Listener)
	}

	err = session.Open()
	if err != nil {
		logrus.WithError(err).Fatal("Failed connecting Discord bot session")
	}

	return
}
