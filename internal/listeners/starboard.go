package listeners

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"

	"github.com/bwmarrin/discordgo"
	"github.com/esimov/stackblur-go"
	"github.com/sarulabs/di/v2"
	"github.com/sirupsen/logrus"
	"github.com/zekroTJA/shinpuru/internal/config"
	"github.com/zekroTJA/shinpuru/internal/models"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/karma"
	"github.com/zekroTJA/shinpuru/internal/services/storage"
	"github.com/zekroTJA/shinpuru/internal/util/imgstore"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
	"github.com/zekroTJA/shinpuru/pkg/embedbuilder"
	"github.com/zekroTJA/shinpuru/pkg/thumbnail"
)

const maxSize = 500.0

var (
	errPublicAddrUnset = errors.New("publicAddr unset")
)

type ListenerStarboard struct {
	publicAddr string

	db    database.Database
	st    storage.Storage
	karma *karma.Service
}

func NewListenerStarboard(container di.Container) *ListenerStarboard {
	cfg := container.Get(static.DiConfig).(*config.Config)
	var publicAddr string
	if cfg.WebServer != nil {
		publicAddr = cfg.WebServer.PublicAddr
	}

	return &ListenerStarboard{
		db:         container.Get(static.DiDatabase).(database.Database),
		st:         container.Get(static.DiObjectStorage).(storage.Storage),
		karma:      container.Get(static.DiKarma).(*karma.Service),
		publicAddr: publicAddr,
	}
}

func (l *ListenerStarboard) ListenerReactionAdd(s *discordgo.Session, e *discordgo.MessageReactionAdd) {
	if e.UserID == s.State.User.ID {
		return
	}

	member, err := discordutil.GetMember(s, e.GuildID, e.UserID)
	if err != nil {
		logrus.WithError(err).Fatal("STARBOARD :: failed getting user")
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
		logrus.WithError(err).Fatal("STARBOARD :: failed getting guild config")
		return
	}
	if starboardConfig.ChannelID == "" {
		return
	}

	if e.Emoji.Name != starboardConfig.EmojiID {
		return
	}

	starboardChannel, err := discordutil.GetChannel(s, starboardConfig.ChannelID)
	if err != nil {
		starboardConfig.ChannelID = ""
		if err = l.db.SetStarboardConfig(starboardConfig); err != nil {
			logrus.WithError(err).Fatal("STARBOARD :: failed disabling starboard")
			return
		}
	}

	msgChannel, err := discordutil.GetChannel(s, e.ChannelID)
	if err != nil {
		logrus.WithError(err).Fatal("STARBOARD :: failed getting message channel")
		return
	}

	msg, err := discordutil.GetMessage(s, e.ChannelID, e.MessageID)
	if err != nil {
		logrus.WithError(err).Fatal("STARBOARD :: failed getting message")
		return
	}

	ok, score := l.hitsThreshhold(msg, starboardConfig)
	if !ok {
		return
	}

	starboardEntry, err := l.db.GetStarboardEntry(msg.ID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		logrus.WithError(err).Fatal("STARBOARD :: failed getting starboard entry")
		return
	}

	censorMedia := msgChannel.NSFW && !starboardChannel.NSFW

	var giveKarma bool
	if database.IsErrDatabaseNotFound(err) || starboardEntry == nil || starboardEntry.Deleted {
		giveKarma = database.IsErrDatabaseNotFound(err) && !starboardEntry.Deleted

		if censorMedia {
			newAttachments := make([]*discordgo.MessageAttachment, len(msg.Attachments))
			i := 0
			if starboardEntry != nil && len(starboardEntry.MediaURLs) > 0 {
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
						logrus.WithError(err).Fatal("STARBOARD :: failed bluring image")
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
			logrus.WithError(err).Fatal("STARBOARD :: failed sending starboard message")
			return
		}

		starboardEntry = &models.StarboardEntry{
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
	} else {
		_, err = s.ChannelMessageEditEmbed(starboardConfig.ChannelID, starboardEntry.StarboardID, l.getEmbed(msg, e.GuildID, score))
		if err != nil {
			logrus.WithError(err).Fatal("STARBOARD :: failed updating starboard message")
			return
		}

		starboardEntry.Score = score
	}

	err = l.db.SetStarboardEntry(starboardEntry)
	if err != nil {
		logrus.WithError(err).Fatal("STARBOARD :: failed setting starboard entry")
		return
	}

	if giveKarma {
		if _, err = l.karma.CheckAndUpdate(e.GuildID, msg.Author, starboardConfig.KarmaGain); err != nil {
			logrus.WithError(err).Fatal("STARBOARD :: failed updating karma")
		}
	}
}

func (l *ListenerStarboard) ListenerReactionRemove(s *discordgo.Session, e *discordgo.MessageReactionRemove) {
	if e.UserID == s.State.User.ID {
		return
	}

	member, err := discordutil.GetMember(s, e.GuildID, e.UserID)
	if err != nil {
		logrus.WithError(err).Fatal("STARBOARD :: failed getting user")
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
		logrus.WithError(err).Fatal("STARBOARD :: failed getting guild config")
		return
	}
	if starboardConfig.ChannelID == "" {
		return
	}

	if e.Emoji.Name != starboardConfig.EmojiID {
		return
	}

	msg, err := discordutil.GetMessage(s, e.ChannelID, e.MessageID)
	if err != nil {
		logrus.WithError(err).Fatal("STARBOARD :: failed getting message")
		return
	}

	starboardEntry, err := l.db.GetStarboardEntry(msg.ID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		logrus.WithError(err).Fatal("STARBOARD :: failed getting entry")
		return
	}

	if database.IsErrDatabaseNotFound(err) || starboardEntry == nil {
		return
	} else {
		ok, score := l.hitsThreshhold(msg, starboardConfig)
		if !ok {
			starboardEntry.Deleted = true
			if err = s.ChannelMessageDelete(starboardConfig.ChannelID, starboardEntry.StarboardID); err != nil {
				logrus.WithError(err).Fatal("STARBOARD :: failed removing starboard message")
			}
		} else {
			_, err = s.ChannelMessageEditEmbed(starboardConfig.ChannelID, starboardEntry.StarboardID, l.getEmbed(msg, e.GuildID, score))
			if err != nil {
				logrus.WithError(err).Fatal("STARBOARD :: failed updating starboard message")
			}
		}

		starboardEntry.Score = score
	}

	err = l.db.SetStarboardEntry(starboardEntry)
	if err != nil {
		logrus.WithError(err).Fatal("STARBOARD :: failed setting entry")
		return
	}
}

func (l *ListenerStarboard) hitsThreshhold(msg *discordgo.Message, starboardConfig *models.StarboardConfig) (ok bool, count int) {
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
	emb := embedbuilder.New().
		WithAuthor(msg.Author.String(), "", msg.Author.AvatarURL("16x16"), "").
		WithDescription(fmt.Sprintf("%s\n\n[jump to message](%s)",
			msg.Content, discordutil.GetMessageLink(msg, guildID))).
		WithTimestamp(string(msg.Timestamp)).
		WithColor(static.ColorEmbedDefault).
		WithFooter(fmt.Sprintf("%d ⭐", count), "", "")

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

	cDone := make(chan struct{}, 1)
	iimg = stackblur.Process(iimg, uint32(iimg.Bounds().Dx()), uint32(iimg.Bounds().Dy()), 50, cDone)
	<-cDone

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
