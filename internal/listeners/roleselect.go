package listeners

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
	"github.com/zekrotja/dgrs"
	"github.com/zekrotja/ken"
	"github.com/zekrotja/rogu"
	"github.com/zekrotja/rogu/log"
)

type ListenerRoleselect struct {
	db  database.Database
	ken ken.IKen
	st  dgrs.IState
	log *rogu.Logger
}

func NewListenerRoleselect(container di.Container) *ListenerRoleselect {
	return &ListenerRoleselect{
		db:  container.Get(static.DiDatabase).(database.Database),
		ken: container.Get(static.DiCommandHandler).(ken.IKen),
		st:  container.Get(static.DiState).(dgrs.IState),
		log: log.Tagged("RoleSelect"),
	}
}

func (t *ListenerRoleselect) HandlerMessageDelete(s discordutil.ISession, e *discordgo.MessageDelete) {
	t.deleteForMessage(e.GuildID, e.ChannelID, e.ID)
}

func (t *ListenerRoleselect) HandlerMessageBulkDelete(s discordutil.ISession, e *discordgo.MessageDeleteBulk) {
	for _, msg := range e.Messages {
		t.deleteForMessage(e.GuildID, e.ChannelID, msg)
	}
}

func (t *ListenerRoleselect) Ready(s discordutil.ISession, e *discordgo.Ready) {
	roleSelects, err := t.db.GetRoleSelects()
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		t.log.Error().Err(err).Msg("Retrieving stored selects failed")
	}

	type perMessage struct {
		GuildID   string
		ChannelID string
		MessageID string
		RoleIDs   []string
	}

	perMessages := make(map[string]*perMessage)
	for _, rs := range roleSelects {
		key := fmt.Sprintf("%s:%s:%s", rs.GuildID, rs.ChannelID, rs.MessageID)
		if _, ok := perMessages[key]; !ok {
			perMessages[key] = &perMessage{
				GuildID:   rs.GuildID,
				ChannelID: rs.ChannelID,
				MessageID: rs.MessageID,
			}
		}
		perMessages[key].RoleIDs = append(perMessages[key].RoleIDs, rs.RoleID)
	}

	if len(perMessages) > 0 {
		t.log.Info().Field("n-messages", len(perMessages)).Msg("Re-attaching button handlers ...")
	}

	for _, pm := range perMessages {
		roles := make([]*discordgo.Role, 0, len(pm.RoleIDs))
		for _, rid := range pm.RoleIDs {
			role, err := t.st.Role(pm.GuildID, rid)
			if err != nil {
				continue
			}
			roles = append(roles, role)
		}
		b := t.ken.Components().Add(pm.MessageID, pm.ChannelID)
		_, err = util.AttachRoleSelectButtons(b, roles)
		if err != nil {
			if discordutil.IsErrCode(err, discordgo.ErrCodeUnknownMessage) {
				t.log.Info().Fields(
					"guild", pm.GuildID,
					"channel", pm.ChannelID,
					"message", pm.MessageID,
				).Msg("Removing role select entries for deleted message")
				t.db.RemoveRoleSelect(pm.GuildID, pm.ChannelID, pm.MessageID)
				continue
			}
			t.log.Error().Fields(
				"guild", pm.GuildID,
				"channel", pm.ChannelID,
				"message", pm.MessageID,
			).Err(err).Msg("Re-Attaching failed")
		}
	}
}

func (t *ListenerRoleselect) deleteForMessage(guildID, channelID, messageID string) {
	err := t.db.RemoveRoleSelect(guildID, channelID, messageID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		t.log.Error().Err(err).Field("guild", guildID).Msg("Removing etries failed")
	}
}
