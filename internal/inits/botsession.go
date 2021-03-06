package inits

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/snowflake"
	"github.com/zekroTJA/shinpuru/internal/core/config"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/core/listeners"
	"github.com/zekroTJA/shinpuru/internal/core/middleware"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/report"
	"github.com/zekroTJA/shinpuru/internal/util/snowflakenodes"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/lctimer"
)

func InitDiscordBotSession(session *discordgo.Session, config *config.Config, database database.Database,
	lct *lctimer.LifeCycleTimer, pmw *middleware.PermissionsMiddleware, gpim *middleware.GhostPingIgnoreMiddleware) {

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

	session.Token = "Bot " + config.Discord.Token
	session.StateEnabled = true
	session.Identify.Intents = discordgo.MakeIntent(static.Intents)

	listenerInviteBlock := listeners.NewListenerInviteBlock(database, pmw)
	listenerGhostPing := listeners.NewListenerGhostPing(database, gpim)
	listenerJDoodle := listeners.NewListenerJdoodle(database, pmw)
	listenerColors := listeners.NewColorListener(database, pmw, config.WebServer.PublicAddr)

	session.AddHandler(listeners.NewListenerReady(config, database, lct).Handler)
	session.AddHandler(listeners.NewListenerGuildJoin(config).Handler)
	session.AddHandler(listeners.NewListenerMemberAdd(database).Handler)
	session.AddHandler(listeners.NewListenerMemberRemove(database).Handler)
	session.AddHandler(listeners.NewListenerVote(database).Handler)
	session.AddHandler(listeners.NewListenerChannelCreate(database).Handler)
	session.AddHandler(listeners.NewListenerVoiceUpdate(database).Handler)
	session.AddHandler(listeners.NewListenerKarma(database).Handler)
	session.AddHandler(listeners.NewListenerAntiraid(database).HandlerMemberAdd)

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

	session.AddHandler(func(s *discordgo.Session, e *discordgo.MessageCreate) {
		util.StatsMessagesAnalysed++
	})

	if config.Metrics != nil && config.Metrics.Enable {
		session.AddHandler(listeners.NewListenerMetrics().Listener)
	}

	err = session.Open()
	if err != nil {
		util.Log.Fatal("Failed connecting Discord bot session:", err)
	}
}
