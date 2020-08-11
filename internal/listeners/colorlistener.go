package listeners

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/pkg/colors"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
)

var (
	rxColorHex = regexp.MustCompile(`^#?[\dA-Fa-f]{6,8}$`)
)

type ColorListener struct {
	db database.Database
}

func NewColorListener(db database.Database) *ColorListener {
	return &ColorListener{db}
}

func (l *ColorListener) HandlerMessageCreate(s *discordgo.Session, e *discordgo.MessageCreate) {
	l.process(s, e.Message)
}

func (l *ColorListener) HandlerMessageEdit(s *discordgo.Session, e *discordgo.MessageUpdate) {
	l.process(s, e.Message)
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

	active, err := l.db.GetGuildColorReaction(m.GuildID)
	if err != nil {
		util.Log.Error("[ColorListener] could not get setting from database:", err)
		return
	}
	if !active {
		return
	}

	guild, err := discordutil.GetGuild(s, m.GuildID)
	if err != nil {
		util.Log.Error("[ColorListener] could not fetch guild:", err)
		return
	}

	for _, hexClr := range matches {
		l.createReaction(s, m, guild, hexClr)
	}
}

func (l *ColorListener) createReaction(s *discordgo.Session, m *discordgo.Message, guild *discordgo.Guild, hexClr string) {
	if strings.HasPrefix(hexClr, "#") {
		hexClr = hexClr[1:]
	}

	hexClr = strings.ToLower(hexClr)
	hexClr = strings.Trim(hexClr, " ")

	clr, err := colors.FromHex(hexClr)
	if err != nil {
		return
	}

	img := image.NewRGBA(image.Rect(0, 0, 24, 24))
	draw.Draw(img, img.Bounds(), &image.Uniform{*clr}, image.ZP, draw.Src)

	buff := bytes.NewBuffer([]byte{})
	if err = png.Encode(buff, img); err != nil {
		util.Log.Error("[ColorListener] failed generating image data:", err)
		return
	}

	b64Data := base64.StdEncoding.EncodeToString(buff.Bytes())

	dataUri := fmt.Sprintf("data:image/png;base64,%s", b64Data)

	emoji, err := s.GuildEmojiCreate(guild.ID, hexClr, dataUri, nil)
	if err != nil {
		util.Log.Error("[ColorListener] failed uploading emoji:", err)
		return
	}

	err = s.MessageReactionAdd(m.ChannelID, m.ID, url.QueryEscape(":"+emoji.Name+":"+emoji.ID))
	if err != nil {
		util.Log.Error("[ColorListener] failed creating message reaction:", err)
		return
	}

	time.AfterFunc(5*time.Second, func() {
		if err = s.GuildEmojiDelete(guild.ID, emoji.ID); err != nil {
			util.Log.Error("[ColorListener] failed deleting emoji:", err)
		}
	})
}
