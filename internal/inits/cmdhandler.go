package inits

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/commands"
	"github.com/zekroTJA/shinpuru/internal/core/backup"
	"github.com/zekroTJA/shinpuru/internal/core/config"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/core/middleware"
	"github.com/zekroTJA/shinpuru/internal/core/storage"
	"github.com/zekroTJA/shinpuru/internal/core/twitchnotify"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
	"github.com/zekroTJA/shinpuru/pkg/lctimer"
	"github.com/zekroTJA/shireikan"
)

func InitCommandHandler(s *discordgo.Session, cfg *config.Config, db database.Database, st storage.Storage,
	tnw *twitchnotify.NotifyWorker, lct *lctimer.LifeCycleTimer, pmw *middleware.PermissionsMiddleware,
	gpim *middleware.GhostPingIgnoreMiddleware) shireikan.Handler {

	cmdHandler := shireikan.NewHandler(&shireikan.Config{
		GeneralPrefix:         cfg.Discord.GeneralPrefix,
		AllowBots:             false,
		AllowDM:               true,
		DeleteMessageAfter:    true,
		ExecuteOnEdit:         true,
		InvokeToLower:         true,
		UseDefaultHelpCommand: false,

		OnError:           errorHandler,
		GuildPrefixGetter: db.GetGuildPrefix,
	})

	cmdHandler.SetObject("db", db)
	cmdHandler.SetObject("config", cfg)
	cmdHandler.SetObject("storage", st)
	cmdHandler.SetObject("tnw", tnw)
	cmdHandler.SetObject("lct", lct)
	cmdHandler.SetObject("backup", backup.New(s, db, st))
	cmdHandler.SetObject("pmw", pmw)

	cmdHandler.RegisterMiddleware(pmw)
	cmdHandler.RegisterMiddleware(gpim)
	cmdHandler.RegisterMiddleware(&middleware.CommandStatsMiddleware{})

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

	if util.Release != "TRUE" {
		cmdHandler.RegisterCommand(&commands.CmdTest{})
	}

	util.Log.Infof("%d commands registered", len(cmdHandler.GetCommandInstances()))

	cmdHandler.RegisterHandlers(s)

	return cmdHandler
}

func errorHandler(ctx shireikan.Context, errTyp shireikan.ErrorType, err error) {
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
