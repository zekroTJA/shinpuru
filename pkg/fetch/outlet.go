package fetch

import (
	"github.com/bwmarrin/discordgo"
)

type DataOutlet interface {
	GuildRoles(guildID string) ([]*discordgo.Role, error)
	GuildMembers(guildID string, after string, limit int) (st []*discordgo.Member, err error)
	GuildChannels(guildID string) (st []*discordgo.Channel, err error)
}
