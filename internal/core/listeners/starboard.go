package listeners

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/shared/models"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
)

type ListenerStarboard struct {
	db database.Database
}

func NewListenerStarboard(db database.Database) *ListenerStarboard {
	return &ListenerStarboard{db}
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

	msg, err := discordutil.GetMessage(s, e.ChannelID, e.MessageID)
	if err != nil {
		util.Log.Errorf("STARBOARD :: failed getting message: %s", err.Error())
		return
	}
	// if msg.Author.ID == e.UserID {
	// 	return
	// }

	ok, score := l.hitsThreshhold(msg, starboardConfig)
	if !ok {
		return
	}

	starboardEntry, err := l.db.GetStarboardEntry(msg.ID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		util.Log.Errorf("STARBOARD :: failed getting starboard entry: %s", err.Error())
		return
	}

	if database.IsErrDatabaseNotFound(err) || starboardEntry == nil {
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
			if err = l.db.RemoveStarboardEntry(msg.ID); err != nil {
				util.Log.Errorf("STARBOARD :: failed removing entry: %s", err.Error())
			}
			if err = s.ChannelMessageDelete(starboardConfig.ChannelID, starboardEntry.StarboardID); err != nil {
				util.Log.Errorf("STARBOARD :: failed removing starboard message: %s", err.Error())
			}
			return
		}

		_, err = s.ChannelMessageEditEmbed(starboardConfig.ChannelID, starboardEntry.StarboardID, l.getEmbed(msg, e.GuildID, score))
		if err != nil {
			util.Log.Errorf("STARBOARD :: failed updating starboard message: %s", err.Error())
			return
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

func (l *ListenerStarboard) getEmbed(msg *discordgo.Message, guildID string, count int) *discordgo.MessageEmbed {
	emb := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    msg.Author.String(),
			IconURL: msg.Author.AvatarURL("16x16"),
		},
		Description: fmt.Sprintf("%s\n\n[jump to message](%s)",
			msg.Content, discordutil.GetMessageLink(msg, guildID)),
		Timestamp: string(msg.Timestamp),
		Color:     static.ColorEmbedDefault,
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("%d ⭐", count),
		},
	}

	if len(msg.Attachments) > 0 {
		emb.Image = &discordgo.MessageEmbedImage{
			URL:    msg.Attachments[0].URL,
			Width:  msg.Attachments[0].Width,
			Height: msg.Attachments[0].Height,
		}
	}

	return emb
}
