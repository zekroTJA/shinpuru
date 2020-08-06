package inits

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/commands"
	"github.com/zekroTJA/shinpuru/internal/core/backup"
	"github.com/zekroTJA/shinpuru/internal/core/config"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/core/storage"
	"github.com/zekroTJA/shinpuru/internal/core/twitchnotify"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/pkg/lctimer"
	"github.com/zekroTJA/shireikan"
)

func InitCommandHandler(s *discordgo.Session, cfg *config.Config, db database.Database, st storage.Storage, tnw *twitchnotify.NotifyWorker, lct *lctimer.LifeCycleTimer) *commands.CmdHandler {
	cmdHandler := shireikan.NewHandler(&shireikan.Config{
		GeneralPrefix:         cfg.Discord.GeneralPrefix,
		AllowBots:             false,
		AllowDM:               true,
		DeleteMessageAfter:    true,
		ExecuteOnEdit:         true,
		InvokeToLower:         true,
		UseDefaultHelpCommand: false,

		GuildPrefixGetter: db.GetGuildPrefix,
	})

	cmdHandler.SetObject("db", db)
	cmdHandler.SetObject("config", cfg)
	cmdHandler.SetObject("storage", st)
	cmdHandler.SetObject("tnw", tnw)
	cmdHandler.SetObject("lct", lct)
	cmdHandler.SetObject("backup", backup.New(s, db, st))

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

	if util.Release != "TRUE" {
		cmdHandler.RegisterCommand(&commands.CmdTest{})
	}

	util.Log.Infof("%d commands registered", len(cmdHandler.GetCommandInstances()))

	return cmdHandler
}
