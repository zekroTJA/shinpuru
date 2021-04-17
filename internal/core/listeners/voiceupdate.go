package listeners

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di/v2"

	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
)

var voiceStateCashe = map[string]*discordgo.VoiceState{}

type ListenerVoiceUpdate struct {
	db database.Database
}

func NewListenerVoiceUpdate(container di.Container) *ListenerVoiceUpdate {
	return &ListenerVoiceUpdate{
		db: container.Get(static.DiDatabase).(database.Database),
	}
}

func (l *ListenerVoiceUpdate) sendVLCMessage(s *discordgo.Session, channelID, userID, content string, color int) {
	user, err := s.User(userID)
	if err != nil {
		return
	}
	s.ChannelMessageSendEmbed(channelID, &discordgo.MessageEmbed{
		Color:       color,
		Description: content,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    user.Username + "#" + user.Discriminator,
			IconURL: user.AvatarURL("16x16"),
		},
		Timestamp: time.Now().Format(time.RFC3339),
	})
}

func (l *ListenerVoiceUpdate) sendJoinMsg(s *discordgo.Session, voiceLogChan, userID string, newChan *discordgo.Channel) {
	msgTxt := fmt.Sprintf(":arrow_right:  Joined **`%s`**", newChan.Name)
	l.sendVLCMessage(s, voiceLogChan, userID, msgTxt, static.ColorEmbedGreen)
}

func (l *ListenerVoiceUpdate) sendMoveMsg(s *discordgo.Session, voiceLogChan, userID string, oldChan, newChan *discordgo.Channel) {
	msgTxt := fmt.Sprintf(":left_right_arrow:  Moved from **`%s`** to **`%s`**", oldChan.Name, newChan.Name)
	l.sendVLCMessage(s, voiceLogChan, userID, msgTxt, static.ColorEmbedCyan)
}

func (l *ListenerVoiceUpdate) sendLeaveMsg(s *discordgo.Session, voiceLogChan, userID string, oldChan *discordgo.Channel) {
	msgTxt := fmt.Sprintf(":arrow_left:  Left **`%s`**", oldChan.Name)
	l.sendVLCMessage(s, voiceLogChan, userID, msgTxt, static.ColorEmbedOrange)
}

func (l *ListenerVoiceUpdate) isBlocked(guildID, chanID string) (ok bool) {
	ok, err := l.db.IsGuildVoiceLogIgnored(guildID, chanID)
	if err != nil {
		util.Log.Errorf("VOICELOG :: failed getting blocked state: %s", err.Error())
	}
	return
}

func (l *ListenerVoiceUpdate) Handler(s *discordgo.Session, e *discordgo.VoiceStateUpdate) {
	vsOld, _ := voiceStateCashe[e.UserID]
	vsNew := e.VoiceState
	if vsOld != nil && vsOld.ChannelID == vsNew.ChannelID {
		return
	}

	voiceStateCashe[e.UserID] = vsNew

	voiceLogChan, err := l.db.GetGuildVoiceLog(e.GuildID)
	if err != nil || voiceLogChan == "" {
		return
	}

	_, err = discordutil.GetChannel(s, voiceLogChan)
	if err != nil {
		l.db.SetGuildVoiceLog(e.GuildID, "")
		return
	}

	if vsOld == nil || (vsOld != nil && vsOld.ChannelID == "") {
		newChan, err := discordutil.GetChannel(s, e.ChannelID)
		if err != nil {
			return
		}

		if l.isBlocked(newChan.GuildID, newChan.ID) {
			return
		}

		l.sendJoinMsg(s, voiceLogChan, e.UserID, newChan)

	} else if vsOld != nil && vsNew.ChannelID != "" && vsOld.ChannelID != vsNew.ChannelID {
		newChan, err := discordutil.GetChannel(s, e.ChannelID)
		if err != nil {
			return
		}

		oldChan, err := discordutil.GetChannel(s, vsOld.ChannelID)
		if err != nil {
			return
		}

		newChanBlocked := l.isBlocked(vsNew.GuildID, vsNew.ChannelID)
		oldChanBlocked := l.isBlocked(vsOld.GuildID, vsOld.ChannelID)

		if newChanBlocked && oldChanBlocked {
		} else if newChanBlocked {
			l.sendLeaveMsg(s, voiceLogChan, e.UserID, oldChan)
		} else if oldChanBlocked {
			l.sendJoinMsg(s, voiceLogChan, e.UserID, newChan)
		} else {
			l.sendMoveMsg(s, voiceLogChan, e.UserID, oldChan, newChan)
		}

	} else if vsOld != nil && vsNew.ChannelID == "" {
		oldChan, err := discordutil.GetChannel(s, vsOld.ChannelID)
		if err != nil {
			return
		}

		if l.isBlocked(oldChan.GuildID, oldChan.ID) {
			return
		}

		l.sendLeaveMsg(s, voiceLogChan, e.UserID, oldChan)
	}
}
