package listeners

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/esimov/stackblur-go"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/models"
	"github.com/zekroTJA/shinpuru/internal/services/config"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/guildlog"
	"github.com/zekroTJA/shinpuru/internal/services/karma"
	"github.com/zekroTJA/shinpuru/internal/services/storage"
	"github.com/zekroTJA/shinpuru/internal/util/imgstore"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
	"github.com/zekroTJA/shinpuru/pkg/embedbuilder"
	"github.com/zekroTJA/shinpuru/pkg/thumbnail"
	"github.com/zekrotja/dgrs"
	"github.com/zekrotja/rogu"
	"github.com/zekrotja/rogu/log"
)

const maxSize = 500.0

var (
	errPublicAddrUnset = errors.New("publicAddr unset")

	rxImageURL = regexp.MustCompile(`https?:\/\/.+\.(?:jpe?g|png|webp|gif|tiff|svg)(?:[\?#&%][^\s]+)?`)
	rxVideoURL = regexp.MustCompile(`https?:\/\/.+\.(?:mp4|m4a|mkv|webm|mpg|mp2|mpeg|mpv|m4p|avi|wmv|mov|flv)(?:[\?#&%][^\s]+)?`)
)

type ListenerStarboard struct {
	publicAddr string

	db    database.Database
	gl    guildlog.Logger
	st    storage.Storage
	karma *karma.Service
	state *dgrs.State
	log   rogu.Logger
}

func NewListenerStarboard(container di.Container) *ListenerStarboard {
	cfg := container.Get(static.DiConfig).(config.Provider)
	return &ListenerStarboard{
		publicAddr: cfg.Config().WebServer.PublicAddr,
		db:         container.Get(static.DiDatabase).(database.Database),
		gl:         container.Get(static.DiGuildLog).(guildlog.Logger).Section("starboard"),
		st:         container.Get(static.DiObjectStorage).(storage.Storage),
		karma:      container.Get(static.DiKarma).(*karma.Service),
		state:      container.Get(static.DiState).(*dgrs.State),
		log:        log.Tagged("Starboard"),
	}
}

func (l *ListenerStarboard) ListenerReactionAdd(s *discordgo.Session, e *discordgo.MessageReactionAdd) {
	self, err := l.state.SelfUser()
	if err != nil {
		l.log.Error().Err(err).Msg("Failed getting self user")
		l.gl.Errorf(e.GuildID, "Failed getting self user: %s", err.Error())
		return
	}

	if e.UserID == self.ID {
		return
	}

	member, err := l.state.Member(e.GuildID, e.UserID)
	if err != nil {
		l.log.Error().Err(err).Msg("Failed getting user")
		l.gl.Errorf(e.GuildID, "Failed getting user (%s): %s", e.UserID, err.Error())
		return
	}

	if member.User.Bot {
		return
	}

	starboardConfig, err := l.db.GetStarboardConfig(e.GuildID)
	if database.IsErrDatabaseNotFound(err) {
		return
	}
	if err != nil {
		l.log.Error().Err(err).Msg("Failed getting guild config")
		l.gl.Errorf(e.GuildID, "Failed getting guild config: %s", err.Error())
		return
	}
	if starboardConfig.ChannelID == "" {
		return
	}

	if e.Emoji.Name != starboardConfig.EmojiID {
		return
	}

	starboardChannel, err := l.state.Channel(starboardConfig.ChannelID)
	if err != nil {
		starboardConfig.ChannelID = ""
		if err = l.db.SetStarboardConfig(starboardConfig); err != nil {
			l.log.Error().Err(err).Msg("Failed disabling starboard")
			l.gl.Errorf(e.GuildID, "Failed disabling starboard: %s", err.Error())
		}
		return
	}

	msgChannel, err := l.state.Channel(e.ChannelID)
	if err != nil {
		l.log.Error().Err(err).Msg("Failed getting message channel")
		l.gl.Errorf(e.GuildID, "Failed getting message channel (%s): %s", e.ChannelID, err.Error())
		return
	}

	msg, err := l.state.Message(e.ChannelID, e.MessageID)

	if err != nil {
		l.log.Error().Err(err).Msg("Failed getting message")
		l.gl.Errorf(e.GuildID, "Failed getting message (%s): %s", e.MessageID, err.Error())
		return
	}

	if msg.Author == nil {
		l.log.Error().Err(err).Msg("Message author is nil")
		l.gl.Errorf(e.GuildID, "Message author is nil (%s)", e.MessageID)
		return
	}

	ok, score := l.hitsThreshhold(msg, starboardConfig)
	if !ok {
		return
	}

	ok, err = l.db.GetUserStarboardOptout(msg.Author.ID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		l.log.Error().Err(err).Msg("Failed getting starboard user optout")
		l.gl.Errorf(e.GuildID, "Failed getting starboard user optout: %s", err.Error())
		return
	}
	if ok {
		return
	}

	starboardEntry, err := l.db.GetStarboardEntry(msg.ID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		l.log.Error().Err(err).Msg("Failed getting starboard entry")
		l.gl.Errorf(e.GuildID, "Failed getting starboard entry: %s", err.Error())
		return
	}

	extractImage(msg)

	censorMedia := msgChannel.NSFW && !starboardChannel.NSFW

	var giveKarma bool
	if database.IsErrDatabaseNotFound(err) || starboardEntry.Deleted {
		giveKarma = database.IsErrDatabaseNotFound(err) && !starboardEntry.Deleted

		if censorMedia {
			newAttachments := make([]*discordgo.MessageAttachment, len(msg.Attachments))
			i := 0
			if len(starboardEntry.MediaURLs) > 0 {
				for i, murl := range starboardEntry.MediaURLs {
					newAttachments[i] = &discordgo.MessageAttachment{
						URL: murl,
					}
				}
			} else {
				for _, attachment := range msg.Attachments {
					newAttachment := &discordgo.MessageAttachment{
						ID:     attachment.ID,
						Width:  attachment.Width,
						Height: attachment.Height,
					}
					newAttachment.URL, err = l.blurImage(attachment.URL)
					if err != nil {
						l.log.Error().Err(err).Msg("Failed bluring image")
						l.gl.Errorf(e.GuildID, "Failed bluring NSFW image (%s): %s", attachment.URL, err.Error())
						continue
					}
					newAttachments[i] = newAttachment
					i++
				}
			}
			msg.Attachments = newAttachments[:i]
		}

		sbMsg, err := s.ChannelMessageSendEmbed(starboardConfig.ChannelID, l.getEmbed(msg, e.GuildID, score))
		if err != nil {
			l.log.Error().Err(err).Msg("Failed sending starboard message")
			l.gl.Errorf(e.GuildID, "Failed sending starboard message: %s", err.Error())
			return
		}

		starboardEntry = models.StarboardEntry{
			MessageID:   msg.ID,
			StarboardID: sbMsg.ID,
			GuildID:     e.GuildID,
			ChannelID:   msg.ChannelID,
			AuthorID:    msg.Author.ID,
			Content:     msg.Content,
			MediaURLs:   make([]string, len(msg.Attachments)),
			Score:       score,
			Deleted:     false,
		}
		for i, a := range msg.Attachments {
			starboardEntry.MediaURLs[i] = a.URL
		}

		if len(msg.Embeds) != 0 && msg.Embeds[0].Video != nil && msg.Embeds[0].Video.URL != "" {
			starboardEntry.MediaURLs = append(starboardEntry.MediaURLs, msg.Embeds[0].Video.URL)
		}
	} else {
		_, err = s.ChannelMessageEditEmbed(starboardConfig.ChannelID, starboardEntry.StarboardID, l.getEmbed(msg, e.GuildID, score))
		if err != nil {
			l.log.Error().Err(err).Msg("Failed updating starboard message")
			l.gl.Errorf(e.GuildID, "Failed updating starboard message: %s", err.Error())
			return
		}

		starboardEntry.Score = score
	}

	err = l.db.SetStarboardEntry(starboardEntry)
	if err != nil {
		l.log.Error().Err(err).Msg("Failed setting starboard entry")
		l.gl.Errorf(e.GuildID, "Failed getting starboard entry: %s", err.Error())
		return
	}

	if giveKarma {
		if _, err = l.karma.CheckAndUpdate(e.GuildID, "", msg.Author, starboardConfig.KarmaGain); err != nil {
			l.log.Error().Err(err).Msg("Failed updating karma")
			l.gl.Errorf(e.GuildID, "Failed updating karma (%s): %s", msg.Author.ID, err.Error())
		}
	}
}

func (l *ListenerStarboard) ListenerReactionRemove(s *discordgo.Session, e *discordgo.MessageReactionRemove) {
	self, err := l.state.SelfUser()
	if err != nil {
		return
	}

	if e.UserID == self.ID {
		return
	}

	member, err := l.state.Member(e.GuildID, e.UserID)
	if err != nil {
		l.log.Error().Err(err).Msg("Failed getting user")
		l.gl.Errorf(e.GuildID, "Failed getting user (%s): %s", e.UserID, err.Error())
		return
	}

	if member.User.Bot {
		return
	}

	starboardConfig, err := l.db.GetStarboardConfig(e.GuildID)
	if database.IsErrDatabaseNotFound(err) {
		return
	}
	if err != nil {
		l.log.Error().Err(err).Msg("Failed getting guild config")
		l.gl.Errorf(e.GuildID, "Failed getting guild config: %s", err.Error())
		return
	}
	if starboardConfig.ChannelID == "" {
		return
	}

	if e.Emoji.Name != starboardConfig.EmojiID {
		return
	}

	// This is to ensure, that the event, which updates the message and its reaction
	// count went through before this event is called and the message is received
	// from cache even before it has been updated by the otehr event listener.
	// Kinda janky, but yeah, async stuff 🎉
	time.Sleep(500 * time.Millisecond)
	msg, err := l.state.Message(e.ChannelID, e.MessageID)
	if err != nil {
		l.log.Error().Err(err).Msg("Failed getting message")
		l.gl.Errorf(e.GuildID, "Failed getting message (%s): %s", e.MessageID, err.Error())
		return
	}

	starboardEntry, err := l.db.GetStarboardEntry(msg.ID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		l.log.Error().Err(err).Msg("Failed getting entry")
		l.gl.Errorf(e.GuildID, "Failed getting entry (%s): %s", msg.ID, err.Error())
		return
	}

	if database.IsErrDatabaseNotFound(err) {
		return
	} else {
		ok, score := l.hitsThreshhold(msg, starboardConfig)
		if !ok {
			starboardEntry.Deleted = true
			if err = s.ChannelMessageDelete(starboardConfig.ChannelID, starboardEntry.StarboardID); err != nil {
				l.log.Error().Err(err).Msg("Failed removing starboard message")
				l.gl.Errorf(e.GuildID, "Failed removing starboard message: %s", err.Error())
			}
		} else {
			_, err = s.ChannelMessageEditEmbed(starboardConfig.ChannelID, starboardEntry.StarboardID, l.getEmbed(msg, e.GuildID, score))
			if err != nil {
				l.log.Error().Err(err).Msg("Failed updating starboard message")
				l.gl.Errorf(e.GuildID, "Failed updating starboard message: %s", err.Error())
			}
		}

		starboardEntry.Score = score
	}

	err = l.db.SetStarboardEntry(starboardEntry)
	if err != nil {
		l.log.Error().Err(err).Msg("Failed setting entry")
		l.gl.Errorf(e.GuildID, "Failed setting entry: %s", err.Error())
		return
	}
}

func (l *ListenerStarboard) hitsThreshhold(msg *discordgo.Message, starboardConfig models.StarboardConfig) (ok bool, count int) {
	for _, r := range msg.Reactions {
		count = r.Count
		ok = r.Emoji.Name == starboardConfig.EmojiID && count >= starboardConfig.Threshold
		if ok {
			return
		}
	}
	return
}

func (l *ListenerStarboard) getEmbed(
	msg *discordgo.Message,
	guildID string,
	count int,
) *discordgo.MessageEmbed {
	var videoURL string
	extractImage(msg)
	videoURL, msg.Content = extractRegex(msg.Content, rxVideoURL)

	emb := embedbuilder.New().
		WithAuthor(msg.Author.String(), "", msg.Author.AvatarURL("16x16"), "").
		WithDescription(fmt.Sprintf("%s\n\n[jump to message](%s)",
			msg.Content, discordutil.GetMessageLink(msg, guildID))).
		WithTimestamp(msg.Timestamp).
		WithColor(static.ColorEmbedDefault).
		WithFooter(fmt.Sprintf("%d ⭐", count), "", "")

	if videoURL != "" {
		emb.WithVideo(videoURL, 0, 0)
	}

	if len(msg.Attachments) > 0 {
		att := msg.Attachments[0]
		emb.WithImage(att.URL, att.ProxyURL, att.Width, att.Height)
	}

	return emb.Build()
}

func (l *ListenerStarboard) blurImage(sourceURL string) (targetURL string, err error) {
	if l.publicAddr == "" {
		err = errPublicAddrUnset
		return
	}

	img, err := imgstore.DownloadFromURL(sourceURL)
	if err != nil {
		return
	}

	iimg, _, err := image.Decode(bytes.NewBuffer(img.Data))
	if err != nil {
		return
	}

	iimg = thumbnail.Make(iimg, int(maxSize))

	iimg, err = stackblur.Process(iimg, 50)
	fmt.Println(err)
	if err != nil {
		return
	}

	newImgData := bytes.NewBuffer([]byte{})
	err = jpeg.Encode(newImgData, iimg, &jpeg.Options{
		Quality: 90,
	})
	if err != nil {
		return
	}

	err = l.st.PutObject(static.StorageBucketImages, img.ID.String(),
		newImgData, int64(newImgData.Len()), "image/jpeg")
	if err != nil {
		return
	}

	targetURL = fmt.Sprintf("%s/imagestore/%s.jpeg", l.publicAddr, img.ID.String())

	return
}

func extractRegex(content string, rx *regexp.Regexp) (url string, rest string) {
	url = rx.FindString(content)
	rest = strings.Replace(content, url, "", 1)
	return
}

func extractImage(msg *discordgo.Message) {
	url, rest := extractRegex(msg.Content, rxImageURL)
	if url == "" {
		return
	}
	msg.Content = rest
	msg.Attachments = append(msg.Attachments, &discordgo.MessageAttachment{
		URL: url,
	})
}
