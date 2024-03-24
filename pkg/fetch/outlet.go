package fetch

import (
	"github.com/bwmarrin/discordgo"
)

type Session interface {
	GuildRoles(guildID string, options ...discordgo.RequestOption) ([]*discordgo.Role, error)
	GuildMembers(guildID string, after string, limit int, options ...discordgo.RequestOption) (st []*discordgo.Member, err error)
	GuildChannels(guildID string, options ...discordgo.RequestOption) (st []*discordgo.Channel, err error)
}
