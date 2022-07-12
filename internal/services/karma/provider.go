package karma

import "github.com/bwmarrin/discordgo"

type Provider interface {
	GetState(guildID string) (ok bool, err error)
	IsBlockListed(guildID, userID string) (isBlocklisted bool, err error)
	Update(guildID, userID, executorID string, value int) (err error)
	ApplyPenalty(guildID, userID string) (err error)
	CheckAndUpdate(guildID, executorID string, object *discordgo.User, value int) (ok bool, err error)
}
