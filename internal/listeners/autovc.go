package listeners

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekrotja/dgrs"
	"strings"
)

var autovcCache = map[string]string{}
var voiceStateCache = map[string]*discordgo.VoiceState{}

type ListenerAutoVoice struct {
	db  database.Database
	st  *dgrs.State
	pmw *permissions.Permissions
}

func NewListenerAutoVoice(container di.Container) *ListenerAutoVoice {
	return &ListenerAutoVoice{
		db:  container.Get(static.DiDatabase).(database.Database),
		st:  container.Get(static.DiState).(*dgrs.State),
		pmw: container.Get(static.DiPermissions).(*permissions.Permissions),
	}
}

func (l *ListenerAutoVoice) Handler(s *discordgo.Session, e *discordgo.VoiceStateUpdate) {

	allowed, _, err := l.pmw.CheckPermissions(s, e.GuildID, e.UserID, "sp.chat.autochannel")
	if err != nil || !allowed {
		return
	}
	vsOld, _ := voiceStateCache[e.UserID]
	vsNew := e.VoiceState

	voiceStateCache[e.UserID] = vsNew

	ids, err := l.db.GetGuildAutoVC(e.GuildID)
	if err != nil {
		return
	}
	idString := strings.Join(ids, ";")

	if vsOld == nil || (vsOld != nil && vsOld.ChannelID == "") {

		if !strings.Contains(idString, vsNew.ChannelID) {
			return
		}

		if err := l.createAutoVC(s, e.UserID, e.GuildID, vsNew.ChannelID); err != nil {
			return
		}

	} else if vsOld != nil && vsNew.ChannelID != "" && vsOld.ChannelID != vsNew.ChannelID {

		// we don't want to delete the channel, if the user get's moved to their auto voicechannel
		if vsNew.ChannelID == autovcCache[e.UserID] {

		} else if strings.Contains(idString, vsNew.ChannelID) && autovcCache[e.UserID] == "" {
			if autovcCache[e.UserID] == "" {
				if err := l.createAutoVC(s, e.UserID, e.GuildID, vsNew.ChannelID); err != nil {
					return
				}
			} else {
				if err := l.deleteAutoVC(s, e.UserID); err != nil {
					return
				}
			}
		} else if autovcCache[e.UserID] != "" {
			if err := l.deleteAutoVC(s, e.UserID); err != nil {
				return
			}
		}

	} else if vsOld != nil && vsNew.ChannelID == "" {
		if autovcCache[e.UserID] != "" {
			if err := l.deleteAutoVC(s, e.UserID); err != nil {
				return
			}
		}

	}
}

func (l *ListenerAutoVoice) createAutoVC(s *discordgo.Session, userID, guildID, parentChannelId string) error {
	parentCh, err := l.st.Channel(parentChannelId)
	if err != nil {
		return err
	}
	user, err := l.st.User(userID)
	if err != nil {
		return err
	}
	ch, err := s.GuildChannelCreate(guildID, user.Username, discordgo.ChannelTypeGuildVoice)
	if err != nil {
		return err
	}
	ch, err = s.ChannelEditComplex(ch.ID, &discordgo.ChannelEdit{
		ParentID: parentCh.ParentID,
		Position: parentCh.Position,
	})
	if err != nil {
		return err
	}
	autovcCache[userID] = ch.ID
	if err := s.GuildMemberMove(guildID, userID, &ch.ID); err != nil {
		return err
	}
	return nil
}

func (l *ListenerAutoVoice) deleteAutoVC(s *discordgo.Session, userID string) error {
	vcID := autovcCache[userID]
	_, err := s.ChannelDelete(vcID)
	if err != nil {
		return err
	}
	delete(autovcCache, userID)
	return nil
}
