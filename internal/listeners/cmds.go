package listeners

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/commands"
	"github.com/zekroTJA/shinpuru/internal/core"
	"github.com/zekroTJA/shinpuru/internal/util"
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
	util.StatsMessagesAnalysed++

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

	re := regexp.MustCompile(`(?:[^\s"]+|"[^"]*")+`)
	contSplit := re.FindAllString(e.Message.Content, -1)
	for i, k := range contSplit {
		if strings.Contains(k, "\"") {
			contSplit[i] = strings.Replace(k, "\"", "", -1)
		}
	}
	invoke := contSplit[0][len(pre):]
	invoke = strings.ToLower(invoke)

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
			permLvl = util.PermLvlBotOwner
		} else if e.Author.ID == guild.OwnerID {
			permLvl = util.PermLvlGuildOwner
		} else {
			permLvl, err = l.db.GetMemberPermissionLevel(s, e.GuildID, e.Author.ID)
		}

		if err != nil && !core.IsErrDatabaseNotFound(err) {
			util.SendEmbedError(s, channel.ID, fmt.Sprintf("Failed getting permission from database: ```\n%s\n```", err.Error()), "Permission Error")
			return
		}
		if permLvl < cmdInstance.GetPermission() {
			util.SendEmbedError(s, channel.ID, "You are not permitted to use this command!", "Missing permission")
			return
		}
		s.ChannelMessageDelete(channel.ID, e.Message.ID)
		err = cmdInstance.Exec(cmdArgs)
		if err != nil {
			emb := &discordgo.MessageEmbed{
				Color:       util.ColorEmbedError,
				Title:       "Command execution failed",
				Description: fmt.Sprintf("Failed executing command: ```\n%s\n```", err.Error()),
				Footer: &discordgo.MessageEmbedFooter{
					Text: "This is kind of an unexpected error and means that something is not right in order. " +
						"Does the bot has the right permissions? If there is no issue with the permissions, please report this bug. For more info, use the 'bug' command.",
				},
			}
			_, err := s.ChannelMessageSendEmbed(channel.ID, emb)
			if err != nil {
				util.Log.Error("An error occured sending command error message: ", err)
			}
		}

		util.StatsCommandsExecuted++

		if l.config.Logging.CommandLogging {
			util.Log.Infof("Executed Command: %s[%s]@%s[%s] - %s", e.Author.Username, e.Author.ID, guild.Name, guild.ID, e.Message.Content)
		}
	}
}
