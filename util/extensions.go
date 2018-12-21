package util

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func EnsureNotEmpty(str, def string) string {
	if str == "" {
		return def
	}
	return str
}

func DeleteMessageLater(s *discordgo.Session, msg *discordgo.Message, duration time.Duration) {
	if msg == nil {
		return
	}
	time.AfterFunc(duration, func() {
		s.ChannelMessageDelete(msg.ChannelID, msg.ID)
	})
}

func FetchRole(s *discordgo.Session, guildID, resolvable string) (*discordgo.Role, error) {
	guild, err := s.Guild(guildID)
	rx := regexp.MustCompile("<@&|>")
	resolvable = rx.ReplaceAllString(resolvable, "")
	if err != nil {
		return nil, err
	}

	checkFuncs := []func(*discordgo.Role, string) bool{
		func(r *discordgo.Role, resolvable string) bool {
			return r.ID == resolvable
		},
		func(r *discordgo.Role, resolvable string) bool {
			return r.Name == resolvable
		},
		func(r *discordgo.Role, resolvable string) bool {
			return strings.ToLower(r.Name) == strings.ToLower(resolvable)
		},
		func(r *discordgo.Role, resolvable string) bool {
			return strings.HasPrefix(strings.ToLower(r.Name), strings.ToLower(resolvable))
		},
		func(r *discordgo.Role, resolvable string) bool {
			return strings.Contains(strings.ToLower(r.Name), strings.ToLower(resolvable))
		},
	}

	for _, checkFunc := range checkFuncs {
		for _, r := range guild.Roles {
			if checkFunc(r, resolvable) {
				return r, nil
			}
		}
	}

	return nil, errors.New("could not be fetched")
}
