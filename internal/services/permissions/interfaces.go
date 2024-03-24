package permissions

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/pkg/permissions"
	"github.com/zekroTJA/shinpuru/pkg/roleutil"
)

type Database interface {
	GetGuildPermissions(guildID string) (map[string]permissions.PermissionArray, error)
}

type State interface {
	Guild(id string, hydrate ...bool) (v *discordgo.Guild, err error)
}

type Session interface {
	roleutil.Session

	GuildMember(guildID, userID string, options ...discordgo.RequestOption) (st *discordgo.Member, err error)
	UserChannelPermissions(userID, channelID string, fetchOptions ...discordgo.RequestOption) (apermissions int64, err error)
}
