package permissions

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/pkg/permissions"
)

type Database interface {
	GetGuildPermissions(guildID string) (map[string]permissions.PermissionArray, error)
}

type State interface {
	Guild(id string, hydrate ...bool) (v *discordgo.Guild, err error)
}
