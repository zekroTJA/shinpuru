package inits

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/snowflake"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/core/config"
	"github.com/zekroTJA/shinpuru/internal/core/listeners"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/report"
	"github.com/zekroTJA/shinpuru/internal/util/snowflakenodes"
	"github.com/zekroTJA/shinpuru/internal/util/static"
)

func InitDiscordBotSession(container di.Container) {
	snowflake.Epoch = static.DefEpoche
	err := snowflakenodes.Setup()
	if err != nil {
		util.Log.Fatal("Failed setting up snowflake nodes: ", err)
	}

	snowflakenodes.NodesReport = make([]*snowflake.Node, len(report.ReportTypes))
	for i, t := range report.ReportTypes {
		if snowflakenodes.NodesReport[i], err = snowflakenodes.RegisterNode(i, "report."+strings.ToLower(t)); err != nil {
			util.Log.Fatal("Failed setting up snowflake nodes: ", err)
		}
	}

	session := container.Get(static.DiDiscordSession).(*discordgo.Session)
	cfg := container.Get(static.DiConfig).(*config.Config)

	session.Token = "Bot " + cfg.Discord.Token
	session.StateEnabled = true
	session.Identify.Intents = discordgo.MakeIntent(static.Intents)

	listenerInviteBlock := listeners.NewListenerInviteBlock(container)
	listenerGhostPing := listeners.NewListenerGhostPing(container)
	listenerJDoodle := listeners.NewListenerJdoodle(container)
	listenerColors := listeners.NewColorListener(container)

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
		util.StatsMessagesAnalysed++
	})

	if cfg.Metrics != nil && cfg.Metrics.Enable {
		session.AddHandler(listeners.NewListenerMetrics().Listener)
	}

	err = session.Open()
	if err != nil {
		util.Log.Fatal("Failed connecting Discord bot session:", err)
	}
}
