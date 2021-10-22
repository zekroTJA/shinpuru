package inits

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di/v2"
	"github.com/sirupsen/logrus"
	"github.com/zekroTJA/shinpuru/internal/commands"
	"github.com/zekroTJA/shinpuru/internal/middleware"
	"github.com/zekroTJA/shinpuru/internal/services/config"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/util/embedded"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
	"github.com/zekroTJA/shireikan"
	"github.com/zekroTJA/shireikan/state"
	"github.com/zekrotja/dgrs"
)

func InitLegacyCommandHandler(container di.Container) shireikan.Handler {

	cfg := container.Get(static.DiConfig).(config.Provider)
	session := container.Get(static.DiDiscordSession).(*discordgo.Session)
	config := container.Get(static.DiConfig).(config.Provider)
	db := container.Get(static.DiDatabase).(database.Database)
	pmw := container.Get(static.DiPermissions).(*permissions.Permissions)
	gpim := container.Get(static.DiGhostpingIgnoreMiddleware).(*middleware.GhostPingIgnoreMiddleware)
	st := container.Get(static.DiState).(*dgrs.State)

	cmdHandler := shireikan.New(&shireikan.Config{
		GeneralPrefix:         config.Config().Discord.GeneralPrefix,
		AllowBots:             false,
		AllowDM:               true,
		DeleteMessageAfter:    true,
		ExecuteOnEdit:         true,
		InvokeToLower:         true,
		UseDefaultHelpCommand: false,

		OnError: legacyErrorHandler,
		GuildPrefixGetter: func(guildID string) (prefix string, err error) {
			if prefix, err = db.GetGuildPrefix(guildID); database.IsErrDatabaseNotFound(err) {
				err = nil
			}
			return
		},

		ObjectContainer: container,
		State:           state.NewDgrs(st),
	})

	if c := cfg.Config().Discord.GlobalCommandRateLimit; c.Burst > 0 && c.LimitSeconds > 0 {
		cmdHandler.RegisterMiddleware(
			middleware.NewGlobalRateLimitMiddleware(c.Burst, time.Duration(c.LimitSeconds)*time.Second))
	}

	cmdHandler.RegisterMiddleware(middleware.NewDisableCommandsMiddleware(container))
	cmdHandler.RegisterMiddleware(pmw)
	cmdHandler.RegisterMiddleware(gpim)

	if cfg.Config().Logging.CommandLogging {
		cmdHandler.RegisterMiddleware(&middleware.LoggerMiddlewrae{})
	}

	cmdHandler.RegisterCommand(&commands.CmdHelp{})
	cmdHandler.RegisterCommand(&commands.CmdPrefix{})
	cmdHandler.RegisterCommand(&commands.CmdPerms{})
	cmdHandler.RegisterCommand(&commands.CmdClear{})
	cmdHandler.RegisterCommand(&commands.CmdMvall{})
	cmdHandler.RegisterCommand(&commands.CmdInfo{})
	cmdHandler.RegisterCommand(&commands.CmdSay{})
	cmdHandler.RegisterCommand(&commands.CmdQuote{})
	cmdHandler.RegisterCommand(&commands.CmdGame{})
	cmdHandler.RegisterCommand(&commands.CmdAutorole{})
	cmdHandler.RegisterCommand(&commands.CmdReport{})
	cmdHandler.RegisterCommand(&commands.CmdModlog{})
	cmdHandler.RegisterCommand(&commands.CmdKick{})
	cmdHandler.RegisterCommand(&commands.CmdBan{})
	cmdHandler.RegisterCommand(&commands.CmdVote{})
	cmdHandler.RegisterCommand(&commands.CmdProfile{})
	cmdHandler.RegisterCommand(&commands.CmdId{})
	cmdHandler.RegisterCommand(&commands.CmdMute{})
	cmdHandler.RegisterCommand(&commands.CmdMention{})
	cmdHandler.RegisterCommand(&commands.CmdNotify{})
	cmdHandler.RegisterCommand(&commands.CmdVoicelog{})
	cmdHandler.RegisterCommand(&commands.CmdBug{})
	cmdHandler.RegisterCommand(&commands.CmdStats{})
	cmdHandler.RegisterCommand(&commands.CmdTwitchNotify{})
	cmdHandler.RegisterCommand(&commands.CmdGhostping{})
	cmdHandler.RegisterCommand(&commands.CmdExec{})
	cmdHandler.RegisterCommand(&commands.CmdBackup{})
	cmdHandler.RegisterCommand(&commands.CmdInviteBlock{})
	cmdHandler.RegisterCommand(&commands.CmdTag{})
	cmdHandler.RegisterCommand(&commands.CmdJoinMsg{})
	cmdHandler.RegisterCommand(&commands.CmdLeaveMsg{})
	cmdHandler.RegisterCommand(&commands.CmdSnowflake{})
	cmdHandler.RegisterCommand(&commands.CmdChannelStats{})
	cmdHandler.RegisterCommand(&commands.CmdKarma{})
	cmdHandler.RegisterCommand(&commands.CmdColorReaction{})
	cmdHandler.RegisterCommand(&commands.CmdLock{})
	cmdHandler.RegisterCommand(&commands.CmdGuild{})
	cmdHandler.RegisterCommand(&commands.CmdLogin{})
	cmdHandler.RegisterCommand(&commands.CmdStarboard{})
	cmdHandler.RegisterCommand(&commands.CmdMaintenance{})

	if !embedded.IsRelease() {
		cmdHandler.RegisterCommand(&commands.CmdTest{})
	}

	logrus.WithField("n", len(cmdHandler.GetCommandInstances())).Info("Commands registered")

	cmdHandler.Setup(session)

	return cmdHandler
}

func legacyErrorHandler(ctx shireikan.Context, errTyp shireikan.ErrorType, err error) {
	switch errTyp {

	// Command execution failed
	case shireikan.ErrTypCommandExec:
		msg, _ := ctx.ReplyEmbedError(
			fmt.Sprintf("Command execution failed unexpectedly: ```\n%s\n```", err.Error()),
			"Command Execution Failed")
		discordutil.DeleteMessageLater(ctx.GetSession(), msg, 60*time.Second)

	// Failed getting channel
	case shireikan.ErrTypGetChannel:
		msg, _ := ctx.ReplyEmbedError(
			fmt.Sprintf("Failed getting channel: ```\n%s\n```", err.Error()),
			"Unexpected Error")
		discordutil.DeleteMessageLater(ctx.GetSession(), msg, 60*time.Second)

	// Failed getting channel
	case shireikan.ErrTypGetGuild:
		msg, _ := ctx.ReplyEmbedError(
			fmt.Sprintf("Failed getting guild: ```\n%s\n```", err.Error()),
			"Unexpected Error")
		discordutil.DeleteMessageLater(ctx.GetSession(), msg, 60*time.Second)

	// Failed getting guild prefix
	case shireikan.ErrTypGuildPrefixGetter:
		if !database.IsErrDatabaseNotFound(err) {
			msg, _ := ctx.ReplyEmbedError(
				fmt.Sprintf("Failed getting guild specific prefix: ```\n%s\n```", err.Error()),
				"Unexpected Error")
			discordutil.DeleteMessageLater(ctx.GetSession(), msg, 60*time.Second)
		}

	// Middleware failed
	case shireikan.ErrTypMiddleware:
		msg, _ := ctx.ReplyEmbedError(
			fmt.Sprintf("Command Handler Middleware failed: ```\n%s\n```", err.Error()),
			"Unexpected Error")
		discordutil.DeleteMessageLater(ctx.GetSession(), msg, 60*time.Second)

	// Middleware failed
	case shireikan.ErrTypNotExecutableInDM:
		msg, _ := ctx.ReplyEmbedError(
			"This command is not executable in DM channels.", "")
		discordutil.DeleteMessageLater(ctx.GetSession(), msg, 8*time.Second)

	// Ignored Errors
	case shireikan.ErrTypCommandNotFound, shireikan.ErrTypDeleteCommandMessage:
		return

	}
}
