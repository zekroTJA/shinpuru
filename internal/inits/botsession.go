package inits

import (
	"strings"
	"sync/atomic"

	"github.com/bwmarrin/snowflake"
	"github.com/sarulabs/di/v2"
	"github.com/sirupsen/logrus"
	"github.com/zekroTJA/shinpuru/internal/config"
	"github.com/zekroTJA/shinpuru/internal/listeners"
	"github.com/zekroTJA/shinpuru/internal/models"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/snowflakenodes"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekrotja/discordgo"
)

func InitDiscordBotSession(container di.Container) {
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
	cfg := container.Get(static.DiConfig).(*config.Config)

	session.Token = "Bot " + cfg.Discord.Token
	session.StateEnabled = true
	session.Identify.Intents = discordgo.MakeIntent(static.Intents)

	listenerInviteBlock := listeners.NewListenerInviteBlock(container)
	listenerGhostPing := listeners.NewListenerGhostPing(container)
	listenerColors := listeners.NewColorListener(container)

	listenerJDoodle, err := listeners.NewListenerJdoodle(container)
	if err != nil {
		logrus.WithError(err).Fatal("Failed setting up code execution listener")
	}

	listenerStarboard := listeners.NewListenerStarboard(container)

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

	session.AddHandler(func(s *discordgo.Session, e *discordgo.MessageCreate) {
		atomic.AddUint64(&util.StatsMessagesAnalysed, 1)
	})

	if cfg.Metrics != nil && cfg.Metrics.Enable {
		session.AddHandler(listeners.NewListenerMetrics().Listener)
	}

	err = session.Open()
	if err != nil {
		logrus.WithError(err).Fatal("Failed connecting Discord bot session")
	}
}
