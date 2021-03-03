// Package discordutil provides general purpose extensuion
// functionalities for discordgo.
package discordutil

import (
	"fmt"
	"strconv"
	"strings"
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

// GetMember first tries to retrieve a member object by passed
// ID from the discordgo state cache. If there is no value
// available, the member will be fetched via API.
func GetMember(s *discordgo.Session, guildID, userID string) (m *discordgo.Member, err error) {
	if m, err = s.State.Member(guildID, userID); err != nil {
		m, err = s.GuildMember(guildID, userID)
	}
	return
}

// GetMembers fetches all members of the guild by utilizing
// paged requests until all members are requested.
func GetMembers(s *discordgo.Session, guildID string) ([]*discordgo.Member, error) {
	lastID := ""
	members := make([]*discordgo.Member, 0)

	for {
		membs, err := s.GuildMembers(guildID, lastID, 1000)
		if err != nil {
			return nil, err
		}

		members = append(members, membs...)

		if len(membs) < 1000 {
			return members, nil
		}

		lastID = membs[999].User.ID
	}
}

// IsCanNotOpenDmToUserError returns true if an returned error
// is caused because a DM channel to a user could not be opened.
func IsCanNotOpenDmToUserError(err error) bool {
	return err != nil && strings.Contains(err.Error(), `"Cannot send messages to this user"`)
}
