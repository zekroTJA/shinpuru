package listeners

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/commands"
	"github.com/zekroTJA/shinpuru/internal/core/config"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
)

type ListenerCmds struct {
	config     *config.Config
	db         database.Database
	cmdHandler *commands.CmdHandler
}

func NewListenerCmd(config *config.Config, db database.Database, cmdHandler *commands.CmdHandler) *ListenerCmds {
	return &ListenerCmds{
		config:     config,
		db:         db,
		cmdHandler: cmdHandler,
	}
}

func (l *ListenerCmds) Handler(s *discordgo.Session, e *discordgo.MessageCreate) {
	util.StatsMessagesAnalysed++

	if e.Message.Author.ID == s.State.User.ID || e.Message.Author.Bot {
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
	if err != nil && !database.IsErrDatabaseNotFound(err) {
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

		ok, err := l.cmdHandler.CheckPermissions(s, guild.ID, e.Author.ID, cmdInstance.GetDomainName())

		if err != nil && !database.IsErrDatabaseNotFound(err) {
			util.SendEmbedError(s, channel.ID, fmt.Sprintf("Failed getting permission from database: ```\n%s\n```", err.Error()), "Permission Error")
			return
		}

		if !ok {
			util.SendEmbedError(s, channel.ID, "You are not permitted to use this command!", "Missing permission").
				DeleteAfter(8 * time.Second).Error()
			return
		}

		if len(e.Message.Mentions) > 0 {
			userMentions := 0
			for _, m := range e.Message.Mentions {
				if !m.Bot {
					userMentions++
				}
			}
			if userMentions > 0 {
				l.cmdHandler.AddNotifiedCommandMsg(e.Message.ID)
			}
		}

		if len(e.Message.Attachments) > 0 {
			defer s.ChannelMessageDelete(channel.ID, e.Message.ID)
		} else {
			s.ChannelMessageDelete(channel.ID, e.Message.ID)
		}
		err = cmdInstance.Exec(cmdArgs)
		if err != nil {
			emb := &discordgo.MessageEmbed{
				Color:       static.ColorEmbedError,
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
