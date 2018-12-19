package listeners

import (
	"fmt"
	"strings"

	"../commands"
	"../core"
	"../util"
	"github.com/bwmarrin/discordgo"
)

type ListenerCmds struct {
	config     *core.Config
	db         core.Database
	cmdHandler *commands.CmdHandler
}

func NewListenerCmd(config *core.Config, db core.Database, cmdHandler *commands.CmdHandler) *ListenerCmds {
	return &ListenerCmds{
		config:     config,
		db:         db,
		cmdHandler: cmdHandler,
	}
}

func (l *ListenerCmds) Handler(s *discordgo.Session, e *discordgo.MessageCreate) {
	if e.Message.Author.ID == s.State.User.ID {
		return
	}
	channel, err := s.Channel(e.ChannelID)
	if err != nil {
		util.Log.Errorf("Failed getting discord channel from ID (%s): %s", e.ChannelID, err.Error())
		return
	}
	if channel.Type != discordgo.ChannelTypeGuildText {
		return
	}
	guildPrefix, err := l.db.GetGuildPrefix(e.GuildID)
	if err != nil && !core.IsErrDatabaseNotFound(err) {
		util.Log.Errorf("Failed fetching guild prefix from database: %s", err.Error())
	}

	var pre string
	if strings.HasPrefix(e.Message.Content, l.config.Discord.GeneralPrefix) {
		pre = l.config.Discord.GeneralPrefix
	} else if guildPrefix != "" && strings.HasPrefix(e.Message.Content, guildPrefix) {
		pre = guildPrefix
	} else {
		return
	}

	contSplit := strings.Fields(e.Message.Content)
	invoke := contSplit[0][len(pre):]
	invoke = strings.ToLower(invoke)

	// UNFINISHED
	if cmdInstance, ok := l.cmdHandler.GetCommand(invoke); ok {
		guild, _ := s.Guild(e.GuildID)
		cmdArgs := &commands.CommandArgs{
			Args:       contSplit[1:],
			Channel:    channel,
			CmdHandler: l.cmdHandler,
			Guild:      guild,
			Message:    e.Message,
			Session:    s,
			User:       e.Author,
		}

		var permLvl = 0
		if e.Author.ID == l.config.Discord.OwnerID {
			permLvl = 1000
		} else if e.Author.ID == guild.OwnerID {
			permLvl = 10
		} else {
			permLvl, err = l.db.GetMemberPermissionLevel(e.GuildID, e.Author.ID)
		}

		if err != nil && !core.IsErrDatabaseNotFound(err) {
			util.SendEmbedError(s, channel.ID, fmt.Sprintf("Failed getting permission from database: ```\n%s\n```", err.Error()), "Permission Error")
			return
		}
		if permLvl < cmdInstance.GetPermission() {
			util.SendEmbedError(s, channel.ID, "You are not permitted to use this command!", "Missing permission")
			return
		}
		err = cmdInstance.Exec(cmdArgs)
		if err != nil {
			util.SendEmbedError(s, channel.ID, fmt.Sprintf("Failed executing command: ```\n%s\n```", err.Error()), "Command execution failed")
		}
	}
}
