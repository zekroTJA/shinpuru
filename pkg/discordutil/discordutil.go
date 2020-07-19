// Package discordutil provides general purpose extensuion
// functionalities for discordgo.
package discordutil

import (
	"fmt"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

// GetMessageLink assembles and returns a message link by
// passed msg object and guildID.
func GetMessageLink(msg *discordgo.Message, guildID string) string {
	return fmt.Sprintf("https://discord.com/channels/%s/%s/%s", guildID, msg.ChannelID, msg.ID)
}

// GetDiscordSnowflakeCreationTime returns the time.Time
// of creation of the passed snowflake string.
//
// Returns an error when the passed snowflake string could
// not be parsed to an integer.
func GetDiscordSnowflakeCreationTime(snowflake string) (time.Time, error) {
	sfI, err := strconv.ParseInt(snowflake, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	timestamp := (sfI >> 22) + 1420070400000
	return time.Unix(timestamp/1000, timestamp), nil
}

// IsAdmin returns true if one of the members roles has
// admin (0x8) permissions on the passed guild.
func IsAdmin(g *discordgo.Guild, m *discordgo.Member) bool {
	if m == nil || g == nil {
		return false
	}

	for _, r := range g.Roles {
		if r.Permissions&0x8 != 0 {
			for _, mrID := range m.Roles {
				if r.ID == mrID {
					return true
				}
			}
		}
	}

	return false
}

// DeleteMessageLater tries to delete the passed msg after
// the specified duration.
//
// If the message was already removed, the error will be
// ignored.
func DeleteMessageLater(s *discordgo.Session, msg *discordgo.Message, duration time.Duration) {
	if msg == nil {
		return
	}
	time.AfterFunc(duration, func() {
		s.ChannelMessageDelete(msg.ChannelID, msg.ID)
	})
}

// GetGuild first tries to retrieve a guild object by passed
// ID from the discordgo state cache. If there is no value
// available, the guild will be fetched via API.
func GetGuild(s *discordgo.Session, guildID string) (g *discordgo.Guild, err error) {
	if g, err = s.State.Guild(guildID); err != nil {
		g, err = s.Guild(guildID)
	}
	return
}
