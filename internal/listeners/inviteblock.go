package listeners

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/zekroTJA/shinpuru/internal/util"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/commands"
	"github.com/zekroTJA/shinpuru/internal/core"
)

var (
	rxInvLink = regexp.MustCompile(`(?i)(https?:\/\/)?(www\.)?(discord\.gg|discordapp.com\/invite)\/.*`)
	rxGenLink = regexp.MustCompile(`(?i)(https?:\/\/)?(www\.)?([\w-\S]+\.)+\w{1,10}\/?[\S]+`)
)

type ListenerInviteBlock struct {
	config     *core.Config
	db         core.Database
	cmdHandler *commands.CmdHandler
}

func NewListenerInviteBlock(db core.Database, cmdHandler *commands.CmdHandler) *ListenerInviteBlock {
	return &ListenerInviteBlock{
		db:         db,
		cmdHandler: cmdHandler,
	}
}

func (l *ListenerInviteBlock) HandlerMessageSend(s *discordgo.Session, e *discordgo.MessageCreate) {
	l.invokeCheck(s, e.Message)
}

func (l *ListenerInviteBlock) HandlerMessageEdit(s *discordgo.Session, e *discordgo.MessageUpdate) {
	l.invokeCheck(s, e.Message)
}

func (l *ListenerInviteBlock) invokeCheck(s *discordgo.Session, msg *discordgo.Message) {
	cont := msg.Content

	if l.checkForInviteLink(cont) {
		l.detected(s, msg)
		return
	}

	link := rxGenLink.FindString(cont)
	if link != "" {
		match, err := l.followLink(link)
		if err != nil {
			util.Log.Error("Failed following link: ", err)
			return
		}
		if match {
			l.detected(s, msg)
		}
	}
}

func (l *ListenerInviteBlock) checkForInviteLink(cont string) bool {
	return rxInvLink.MatchString(cont)
}

func (l *ListenerInviteBlock) followLink(link string) (bool, error) {
	if !strings.HasPrefix(link, "http://") && !strings.HasPrefix(link, "https://") {
		link = "http://" + link
	}

	resp, err := http.DefaultClient.Get(link)
	if err != nil {
		return false, nil
	}

	return l.checkForInviteLink(resp.Request.URL.String()), nil
}

func (l *ListenerInviteBlock) detected(s *discordgo.Session, e *discordgo.Message) error {
	enabled, err := l.db.GetGuildInviteBlock(e.GuildID)
	if core.IsErrDatabaseNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}
	if enabled == "" {
		return nil
	}

	ok, err := l.cmdHandler.CheckPermissions(s, e.GuildID, e.Author.ID, "sp.guild.mod.inviteblock.send")
	if err != nil {
		return err
	}

	if ok {
		return nil
	}

	return s.ChannelMessageDelete(e.ChannelID, e.ID)
}
