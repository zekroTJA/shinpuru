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
	"github.com/zekroTJA/shinpuru/internal/core/config"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/core/storage"
	"github.com/zekroTJA/shinpuru/internal/shared/models"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/imgstore"
	"github.com/zekroTJA/shinpuru/internal/util/karma"
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

	db database.Database
	st storage.Storage
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
		publicAddr: publicAddr,
	}
}

func (l *ListenerStarboard) ListenerReactionAdd(s *discordgo.Session, e *discordgo.MessageReactionAdd) {
	if e.UserID == s.State.User.ID {
		return
	}

	member, err := discordutil.GetMember(s, e.GuildID, e.UserID)
	if err != nil {
		util.Log.Errorf("STARBOARD :: failed getting user: %s", err.Error())
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
		util.Log.Errorf("STARBOARD :: failed getting guild config: %s", err.Error())
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
			util.Log.Errorf("STARBOARD :: failed disabling starboard: %s", err.Error())
			return
		}
	}

	msgChannel, err := discordutil.GetChannel(s, e.ChannelID)
	if err != nil {
		util.Log.Errorf("STARBOARD :: failed getting message channel: %s", err.Error())
		return
	}

	msg, err := discordutil.GetMessage(s, e.ChannelID, e.MessageID)
	if err != nil {
		util.Log.Errorf("STARBOARD :: failed getting message: %s", err.Error())
		return
	}

	ok, score := l.hitsThreshhold(msg, starboardConfig)
	if !ok {
		return
	}

	starboardEntry, err := l.db.GetStarboardEntry(msg.ID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		util.Log.Errorf("STARBOARD :: failed getting starboard entry: %s", err.Error())
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
						util.Log.Errorf("STARBOARD :: failed bluring image: %s", err.Error())
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
			util.Log.Errorf("STARBOARD :: failed sending starboard message: %s", err.Error())
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
			util.Log.Errorf("STARBOARD :: failed updating starboard message: %s", err.Error())
			return
		}

		starboardEntry.Score = score
	}

	err = l.db.SetStarboardEntry(starboardEntry)
	if err != nil {
		util.Log.Errorf("STARBOARD :: failed setting starboard entry: %s", err.Error())
		return
	}

	if giveKarma {
		if _, err = karma.Alter(l.db, e.GuildID, msg.Author, starboardConfig.KarmaGain); err != nil {
			util.Log.Errorf("STARBOARD :: failed updating karma: %s", err.Error())
		}
	}
}

func (l *ListenerStarboard) ListenerReactionRemove(s *discordgo.Session, e *discordgo.MessageReactionRemove) {
	if e.UserID == s.State.User.ID {
		return
	}

	member, err := discordutil.GetMember(s, e.GuildID, e.UserID)
	if err != nil {
		util.Log.Errorf("STARBOARD :: failed getting user: %s", err.Error())
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
		util.Log.Errorf("STARBOARD :: failed getting guild config: %s", err.Error())
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
		util.Log.Errorf("STARBOARD :: failed getting message: %s", err.Error())
		return
	}

	starboardEntry, err := l.db.GetStarboardEntry(msg.ID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		util.Log.Errorf("STARBOARD :: failed getting entry: %s", err.Error())
		return
	}

	if database.IsErrDatabaseNotFound(err) || starboardEntry == nil {
		return
	} else {
		ok, score := l.hitsThreshhold(msg, starboardConfig)
		if !ok {
			starboardEntry.Deleted = true
			if err = s.ChannelMessageDelete(starboardConfig.ChannelID, starboardEntry.StarboardID); err != nil {
				util.Log.Errorf("STARBOARD :: failed removing starboard message: %s", err.Error())
			}
		} else {
			_, err = s.ChannelMessageEditEmbed(starboardConfig.ChannelID, starboardEntry.StarboardID, l.getEmbed(msg, e.GuildID, score))
			if err != nil {
				util.Log.Errorf("STARBOARD :: failed updating starboard message: %s", err.Error())
			}
		}

		starboardEntry.Score = score
	}

	err = l.db.SetStarboardEntry(starboardEntry)
	if err != nil {
		util.Log.Errorf("STARBOARD :: failed setting entry: %s", err.Error())
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
