package listeners

import (
	"fmt"
	"strings"
	"time"

	"github.com/zekroTJA/shinpuru/internal/util"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/core"
)

const gpDelay = 12 * time.Hour

type ListenerGhostPing struct {
	db core.Database
}

func NewListenerGhostPing(db core.Database) *ListenerGhostPing {
	return &ListenerGhostPing{
		db: db,
	}
}

func (l *ListenerGhostPing) Handler(s *discordgo.Session, e *discordgo.MessageDelete) {
	fmt.Println("TEST", e.Author)
	if e == nil || e.Author == nil || e.Author.Bot || len(e.Mentions) == 0 {
		return
	}

	timeSent, err := e.Timestamp.Parse()
	if err != nil {
		return
	}

	if time.Since(timeSent) <= gpDelay {
		gpMsg, err := l.db.GetGuildGhostpingMsg(e.GuildID)
		if err != nil {
			if !core.IsErrDatabaseNotFound(err) {
				util.Log.Errorf("failed getting ghost ping msg for guild %s: %s\n", e.GuildID, err.Error())
			}
			return
		}

		uPinged := e.Mentions[0]

		gpMsg = strings.Replace(gpMsg, "{pinger}", e.Author.Mention(), -1)
		gpMsg = strings.Replace(gpMsg, "{pinged}", uPinged.Mention(), -1)
		gpMsg = strings.Replace(gpMsg, "{msg}", e.Content, -1)

		s.ChannelMessageSend(e.ChannelID, gpMsg)
	}
}
