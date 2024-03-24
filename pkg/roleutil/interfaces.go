package roleutil

import "github.com/bwmarrin/discordgo"

type Session interface {
	GuildMember(guildID, userID string, options ...discordgo.RequestOption) (st *discordgo.Member, err error)
	GuildRoles(guildID string, options ...discordgo.RequestOption) (st []*discordgo.Role, err error)
}
