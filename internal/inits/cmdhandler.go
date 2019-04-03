package inits

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/commands"
	"github.com/zekroTJA/shinpuru/internal/core"
	"github.com/zekroTJA/shinpuru/internal/util"
)

func InitCommandHandler(s *discordgo.Session, config *core.Config, database core.Database, twitchNotifyWorker *core.TwitchNotifyWorker) *commands.CmdHandler {
	cmdHandler := commands.NewCmdHandler(s, database, config, twitchNotifyWorker)

	cmdHandler.RegisterCommand(&commands.CmdHelp{PermLvl: 0})
	cmdHandler.RegisterCommand(&commands.CmdPrefix{PermLvl: 10})
	cmdHandler.RegisterCommand(&commands.CmdPerms{PermLvl: 10})
	cmdHandler.RegisterCommand(&commands.CmdClear{PermLvl: 8})
	cmdHandler.RegisterCommand(&commands.CmdMvall{PermLvl: 5})
	cmdHandler.RegisterCommand(&commands.CmdInfo{PermLvl: 0})
	cmdHandler.RegisterCommand(&commands.CmdSay{PermLvl: 3})
	cmdHandler.RegisterCommand(&commands.CmdQuote{PermLvl: 0})
	cmdHandler.RegisterCommand(&commands.CmdGame{PermLvl: 999})
	cmdHandler.RegisterCommand(&commands.CmdAutorole{PermLvl: 9})
	cmdHandler.RegisterCommand(&commands.CmdReport{PermLvl: 5})
	cmdHandler.RegisterCommand(&commands.CmdModlog{PermLvl: 6})
	cmdHandler.RegisterCommand(&commands.CmdKick{PermLvl: 6})
	cmdHandler.RegisterCommand(&commands.CmdBan{PermLvl: 8})
	cmdHandler.RegisterCommand(&commands.CmdVote{PermLvl: 0})
	cmdHandler.RegisterCommand(&commands.CmdProfile{PermLvl: 0})
	cmdHandler.RegisterCommand(&commands.CmdId{PermLvl: 0})
	cmdHandler.RegisterCommand(&commands.CmdMute{PermLvl: 4})
	cmdHandler.RegisterCommand(&commands.CmdMention{PermLvl: 4})
	cmdHandler.RegisterCommand(&commands.CmdNotify{PermLvl: 0})
	cmdHandler.RegisterCommand(&commands.CmdVoicelog{PermLvl: 6})
	cmdHandler.RegisterCommand(&commands.CmdBug{PermLvl: 0})
	cmdHandler.RegisterCommand(&commands.CmdStats{PermLvl: 0})
	cmdHandler.RegisterCommand(&commands.CmdTwitchNotify{PermLvl: 5})
	cmdHandler.RegisterCommand(&commands.CmdGhostping{PermLvl: 3})
	cmdHandler.RegisterCommand(&commands.CmdExec{PermLvl: 5})
	cmdHandler.RegisterCommand(&commands.CmdBackup{PermLvl: 9})
	cmdHandler.RegisterCommand(&commands.CmdInviteBlock{PermLvl: 6})

	if util.Release != "TRUE" {
		cmdHandler.RegisterCommand(&commands.CmdTest{})
	}

	if config.Permissions != nil {
		cmdHandler.UpdateCommandPermissions(config.Permissions.CustomCmdPermissions)
		if config.Permissions.BotOwnerLevel > 0 {
			util.PermLvlBotOwner = config.Permissions.BotOwnerLevel
		}
		if config.Permissions.GuildOwnerLevel > 0 {
			util.PermLvlGuildOwner = config.Permissions.GuildOwnerLevel
		}
	}

	util.Log.Infof("%d commands registered", cmdHandler.GetCommandListLen())

	return cmdHandler
}
