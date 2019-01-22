package listeners

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/zekroTJA/shinpuru/internal/core"
	"github.com/zekroTJA/shinpuru/internal/util"
)

var voiceStateCashe = map[string]*discordgo.VoiceState{}

type ListenerVoiceUpdate struct {
	db core.Database
}

func NewListenerVoiceUpdate(db core.Database) *ListenerVoiceUpdate {
	return &ListenerVoiceUpdate{
		db: db,
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
	_, err = s.Channel(voiceLogChan)
	if err != nil {
		fmt.Println("reset vlog chan")
		l.db.SetGuildVoiceLog(e.GuildID, "")
		return
	}
	if vsOld == nil || (vsOld != nil && vsOld.ChannelID == "") {
		newChan, err := s.Channel(e.ChannelID)
		if err != nil {
			return
		}
		msgTxt := fmt.Sprintf(":arrow_right:  Joined **`%s`**", newChan.Name)
		l.sendVLCMessage(s, voiceLogChan, e.UserID, msgTxt, util.ColorEmbedGreen)
	} else if vsOld != nil && vsNew.ChannelID != "" && vsOld.ChannelID != vsNew.ChannelID {
		newChan, err := s.Channel(e.ChannelID)
		if err != nil {
			return
		}
		oldChan, err := s.Channel(vsOld.ChannelID)
		if err != nil {
			return
		}
		msgTxt := fmt.Sprintf(":left_right_arrow:  Moved from **`%s`** to **`%s`**", oldChan.Name, newChan.Name)
		l.sendVLCMessage(s, voiceLogChan, e.UserID, msgTxt, util.ColorEmbedCyan)
	} else if vsOld != nil && vsNew.ChannelID == "" {
		oldChan, err := s.Channel(vsOld.ChannelID)
		if err != nil {
			return
		}
		msgTxt := fmt.Sprintf(":arrow_left:  Left **`%s`**", oldChan.Name)
		l.sendVLCMessage(s, voiceLogChan, e.UserID, msgTxt, util.ColorEmbedOrange)
	}
}
