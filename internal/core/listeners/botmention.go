package listeners

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/core/config"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/embedbuilder"
)

type ListenerBotMention struct {
	config *config.Config

	idLen int
}

func NewListenerBotMention(config *config.Config) *ListenerBotMention {
	return &ListenerBotMention{config, 0}
}

func (l *ListenerBotMention) Listener(s *discordgo.Session, e *discordgo.MessageCreate) {
	if l.idLen == 0 {
		l.idLen = len(s.State.User.ID)
	}

	cLen := len(e.Message.Content)
	if cLen < 3+l.idLen ||
		cLen > 5+l.idLen ||
		e.Message.Content[0] != '<' ||
		e.Message.Content[1] != '@' ||
		e.Author.ID == s.State.User.ID {
		return
	}

	cursor := 2
	if e.Message.Content[2] == '!' {
		cursor = 3
	}

	id := e.Message.Content[cursor : cursor+l.idLen]
	if id != s.State.User.ID {
		return
	}

	prefix := l.config.Discord.GeneralPrefix
	emb := embedbuilder.New().
		WithColor(static.ColorEmbedDefault).
		WithThumbnail(s.State.User.AvatarURL("64x64"), "", 64, 64).
		WithDescription(fmt.Sprintf("shinpuru Discord Bot v.%s (%s)", util.AppVersion, util.AppCommit[:6])).
		WithFooter(fmt.Sprintf("© %d Ringo Hoffmann (zekro Development)", time.Now().Year()), "", "").
		AddField("Help", fmt.Sprintf(
			"Type `%shelp` in the chat to get a list of available commands.\n"+
				"You can also use `%shelp <commandInvoke>` to get more details about a command.\n"+
				"[**Here**](https://github.com/zekroTJA/shinpuru/wiki/commands) you can find "+
				"the wiki page with a detailed list of available commands.", prefix, prefix))

	if l.config.WebServer != nil && l.config.WebServer.Enabled {
		emb.AddField("Web Interface", fmt.Sprintf(
			"[**Here**](%s) you can access the web interface.\n"+
				"You can also use the `%slogin` command if you don't want to log in to the web interface via Discord.",
			l.config.WebServer.PublicAddr, prefix))
	}

	emb.AddField("Repository", fmt.Sprintf(
		"[**Here**](https://github.com/zekroTJA/shinpuru) you can find the open source "+
			"repository of shinpuru. Feel free to contribute issues and pull requests, if you want.\n"+
			"You can also use the `%sinfo` command to get more information.", prefix))

	util.SendEmbedRaw(s, e.ChannelID, emb.Build())
}
