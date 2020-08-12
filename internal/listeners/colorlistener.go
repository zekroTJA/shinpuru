package listeners

import (
	"encoding/base64"
	"fmt"
	"image/color"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/pkg/colors"
	"github.com/zekroTJA/timedmap"
)

var (
	rxColorHex = regexp.MustCompile(`^#?[\dA-Fa-f]{6,8}$`)
)

type emojiCacheEntry struct {
	clr *color.RGBA
}

type ColorListener struct {
	db         database.Database
	emojiCahce *timedmap.TimedMap
}

func NewColorListener(db database.Database) *ColorListener {
	return &ColorListener{db, timedmap.New(1 * time.Minute)}
}

func (l *ColorListener) HandlerMessageCreate(s *discordgo.Session, e *discordgo.MessageCreate) {
	l.process(s, e.Message)
}

func (l *ColorListener) HandlerMessageEdit(s *discordgo.Session, e *discordgo.MessageUpdate) {
	l.process(s, e.Message)
}

func (l *ColorListener) HandlerMessageReaction(s *discordgo.Session, e *discordgo.MessageReactionAdd) {
	// e.
}

func (l *ColorListener) process(s *discordgo.Session, m *discordgo.Message) {
	if len(m.Content) < 6 {
		return
	}

	matches := make([]string, 0)

	// Find color hex in message content using
	// predefined regex.
	for _, v := range strings.Split(m.Content, " ") {
		if rxColorHex.MatchString(v) {
			matches = append(matches, v)
		}
	}

	// Get color reaction enabled guild setting
	// and return when disabled
	active, err := l.db.GetGuildColorReaction(m.GuildID)
	if err != nil {
		util.Log.Error("[ColorListener] could not get setting from database:", err)
		return
	}
	if !active {
		return
	}

	// Execute reaction for each match
	for _, hexClr := range matches {
		l.createReaction(s, m, hexClr)
	}
}

func (l *ColorListener) createReaction(s *discordgo.Session, m *discordgo.Message, hexClr string) {
	buff, err := colors.CreateImage(hexClr, 24, 24)
	if err != nil {
		util.Log.Error("[ColorListener] failed generating image data:", err)
		return
	}

	// Encode the raw image data to a base64 string
	b64Data := base64.StdEncoding.EncodeToString(buff.Bytes())

	// Envelope the base64 data into data uri format
	dataUri := fmt.Sprintf("data:image/png;base64,%s", b64Data)

	// Upload guild emote
	emoji, err := s.GuildEmojiCreate(m.GuildID, hexClr, dataUri, nil)
	if err != nil {
		util.Log.Error("[ColorListener] failed uploading emoji:", err)
		return
	}

	// Add reaction of the uploaded emote to the message
	err = s.MessageReactionAdd(m.ChannelID, m.ID, url.QueryEscape(":"+emoji.Name+":"+emoji.ID))
	if err != nil {
		util.Log.Error("[ColorListener] failed creating message reaction:", err)
		return
	}

	// Delete the uploaded emote after 5 seconds
	// to give discords caching or whatever some
	// time to save the emoji.
	time.AfterFunc(5*time.Second, func() {
		if err = s.GuildEmojiDelete(m.GuildID, emoji.ID); err != nil {
			util.Log.Error("[ColorListener] failed deleting emoji:", err)
		}
	})
}
