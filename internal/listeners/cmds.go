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
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
)

var (
	argsRx = regexp.MustCompile(`(?:[^\s"]+|"[^"]*")+`)
)

// ListenerCmds is the message listener which
// is processing and parsing chat commands to
// execute them using the CmdHandler.
type ListenerCmds struct {
	config     *config.Config
	db         database.Database
	cmdHandler *commands.CmdHandler
}

// NewListenerCmd initializes a new instance of ListenerCmd
// using the passed config, database driver and commandHandler
// instance.
func NewListenerCmd(config *config.Config, db database.Database, cmdHandler *commands.CmdHandler) *ListenerCmds {
	return &ListenerCmds{
		config:     config,
		db:         db,
		cmdHandler: cmdHandler,
	}
}

func (l *ListenerCmds) Handler(s *discordgo.Session, e *discordgo.MessageCreate) {
	// Count up messages analysed stats count.
	util.StatsMessagesAnalysed++

	// Filter messages comming from the bot users account and from other bots.
	if e.Message.Author.ID == s.State.User.ID || e.Message.Author.Bot {
		return
	}

	channel, err := s.Channel(e.ChannelID)
	if err != nil {
		util.Log.Errorf("Failed getting discord channel from ID (%s): %s", e.ChannelID, err.Error())
		return
	}

	// Get the specific guild prefix of the guild where the message was sent,
	// if it was specified in the guild settings.
	guildPrefix, err := l.db.GetGuildPrefix(e.GuildID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		util.Log.Errorf("Failed fetching guild prefix from database: %s", err.Error())
	}

	// Check if the message content starts either with the global prefix or
	// with the guild specific prefix, if specified.
	var pre string
	if strings.HasPrefix(e.Message.Content, l.config.Discord.GeneralPrefix) {
		pre = l.config.Discord.GeneralPrefix
	} else if guildPrefix != "" && strings.HasPrefix(e.Message.Content, guildPrefix) {
		pre = guildPrefix
	} else {
		return
	}

	// Split message content into arguments by spaces.
	// Arguments are not split when contained in quotation
	// marks. The quotation marks are included in the
	// arguments. So, after splitting, strip the
	// quotation marks away from the argument.
	args := argsRx.FindAllString(e.Message.Content, -1)
	for i, k := range args {
		if strings.Contains(k, "\"") {
			args[i] = strings.Replace(k, "\"", "", -1)
		}
	}

	// The first argument is the command invoke, which is
	// taken and then lowercased.
	invoke := strings.ToLower(
		args[0][len(pre):])

	// Try to get the command instance by invoke from
	// the command handler. When the command could not
	// be found, return.
	cmdInstance, ok := l.cmdHandler.GetCommand(invoke)
	if !ok {
		return
	}

	// Mark if the command was sent into a DM channel
	// or not.
	isDM := channel.Type == discordgo.ChannelTypeDM

	// Only continue if the command message was either
	// sent into a guild text channel or into a DM channel
	// when DM execution is allowed by the command instance.
	if !(channel.Type == discordgo.ChannelTypeGuildText ||
		isDM && cmdInstance.IsExecutableInDMChannels()) {
		util.SendEmbedError(s, channel.ID, "This command can not be executed in a DM channel.", "")
		return
	}

	// Now, the first item of the args (the invoke) is
	// stripped away.
	args = args[1:]

	var guild *discordgo.Guild
	if !isDM {
		// Get the guild object where the command was
		// executed from discordgo state cache.
		guild, err = discordutil.GetGuild(s, e.GuildID)
		// If guild object was not found in discordgo
		// state cache, query the guild object from API.
		if err == discordgo.ErrNilState || err == discordgo.ErrStateNotFound {
			guild, err = s.Guild(e.GuildID)
		}
		if err != nil {
			util.Log.Errorf("Failed getting discord guild from ID (%s): %s", e.GuildID, err.Error())
			return
		}
	}

	// Assemble the command args passed to the
	// exec handler of the command instance.
	cmdArgs := &commands.CommandArgs{
		Args:       args,
		Channel:    channel,
		CmdHandler: l.cmdHandler,
		Guild:      guild,
		Message:    e.Message,
		Session:    s,
		User:       e.Author,
		IsDM:       isDM,
	}

	// Check if the command executor has the
	// permissions to execute the command.
	var guildID string
	if guild != nil {
		guildID = guild.ID
	}
	ok, _, err = l.cmdHandler.CheckPermissions(s, guildID, e.Author.ID, cmdInstance.GetDomainName())
	// Return and send error message to channel when
	// something went wrong on getting the permission
	// rules from database.
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		util.SendEmbedError(s, channel.ID, fmt.Sprintf("Failed getting permission from database: ```\n%s\n```", err.Error()), "Permission Error")
		return
	}
	// Return and send not-permitted message to channel
	// when the user has insufficient permission to execute
	// the command.
	if !ok {
		util.SendEmbedError(s, channel.ID, "You are not permitted to use this command!", "Missing permission").
			DeleteAfter(8 * time.Second).Error()
		return
	}

	// If message contains any mentions which are not
	// addressed to a bot, add the message to the
	// notified command message register of the command
	// handler.
	if len(e.Message.Mentions) > 0 {
	mentionsLoop:
		for _, m := range e.Message.Mentions {
			if !m.Bot {
				l.cmdHandler.AddNotifiedCommandMsg(e.Message.ID)
				break mentionsLoop
			}
		}
	}

	// Delete the command initialization message.
	// When the message contains any attachments which
	// may be used by the command, "schedule" the deletion
	// for the end of the handler function.
	if len(e.Message.Attachments) > 0 {
		defer s.ChannelMessageDelete(channel.ID, e.Message.ID)
	} else {
		s.ChannelMessageDelete(channel.ID, e.Message.ID)
	}

	// Execute the command instance handler with the
	// prepared command arguments.
	err = cmdInstance.Exec(cmdArgs)
	// When the command handler failed, send the error
	// as message to the channel with details about the
	// error.
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

	// Count up command executed stats count.
	util.StatsCommandsExecuted++

	// Print command log message when command logging
	// is enabled by config.
	if l.config.Logging.CommandLogging {
		if isDM {
			util.Log.Infof("Executed Command: %s[%s]@DM - %s", e.Author.Username, e.Author.ID, e.Message.Content)
		} else {
			util.Log.Infof("Executed Command: %s[%s]@%s[%s] - %s", e.Author.Username, e.Author.ID, guild.Name, guild.ID, e.Message.Content)
		}
	}
}
