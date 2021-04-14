package inits

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/snowflake"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/core/config"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/core/listeners"
	"github.com/zekroTJA/shinpuru/internal/core/middleware"
	"github.com/zekroTJA/shinpuru/internal/core/storage"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/report"
	"github.com/zekroTJA/shinpuru/internal/util/snowflakenodes"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/lctimer"
)

func InitDiscordBotSession(container di.Container) *discordgo.Session {
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

	session, err := discordgo.New()
	if err != nil {
		util.Log.Fatal(err)
	}

	cfg := container.Get(static.DiConfig).(*config.Config)
	db := container.Get(static.DiDatabase).(database.Database)
	storage := container.Get(static.DiObjectStorage).(storage.Storage)
	lct := container.Get(static.DiLifecycleTimer).(*lctimer.LifeCycleTimer)
	pmw := container.Get(static.DiPermissionMiddleware).(*middleware.PermissionsMiddleware)
	gpim := container.Get(static.DiGhostpingIgnoreMiddleware).(*middleware.GhostPingIgnoreMiddleware)

	session.Token = "Bot " + cfg.Discord.Token
	session.StateEnabled = true
	session.Identify.Intents = discordgo.MakeIntent(static.Intents)

	listenerInviteBlock := listeners.NewListenerInviteBlock(db, pmw)
	listenerGhostPing := listeners.NewListenerGhostPing(db, gpim)
	listenerJDoodle := listeners.NewListenerJdoodle(db, pmw)
	listenerColors := listeners.NewColorListener(db, pmw, cfg.WebServer.PublicAddr)

	publicAddr := ""
	if cfg.WebServer != nil {
		publicAddr = cfg.WebServer.PublicAddr
	}
	listenerStarboard := listeners.NewListenerStarboard(db, storage, publicAddr)

	session.AddHandler(listeners.NewListenerReady(cfg, db, lct).Handler)
	session.AddHandler(listeners.NewListenerGuildJoin(cfg).Handler)
	session.AddHandler(listeners.NewListenerMemberAdd(db).Handler)
	session.AddHandler(listeners.NewListenerMemberRemove(db).Handler)
	session.AddHandler(listeners.NewListenerVote(db).Handler)
	session.AddHandler(listeners.NewListenerChannelCreate(db).Handler)
	session.AddHandler(listeners.NewListenerVoiceUpdate(db).Handler)
	session.AddHandler(listeners.NewListenerKarma(db).Handler)
	session.AddHandler(listeners.NewListenerAntiraid(db).HandlerMemberAdd)
	session.AddHandler(listeners.NewListenerBotMention(cfg).Listener)

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

	return session
}
