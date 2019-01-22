package util

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
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

func BoolAsString(cond bool, ifTrue, ifFalse string) string {
	if cond {
		return ifTrue
	}
	return ifFalse
}

func IndexOfStrArray(str string, arr []string) int {
	for i, v := range arr {
		if v == str {
			return i
		}
	}
	return -1
}

func GetMessageLink(msg *discordgo.Message) string {
	return fmt.Sprintf("https://discordapp.com/channels/%s/%s/%s", msg.GuildID, msg.ChannelID, msg.ID)
}

func GetDiscordSnowflakeCreationTime(snowflake string) (time.Time, error) {
	sfI, err := strconv.ParseInt(snowflake, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	timestamp := (sfI >> 22) + 1420070400000
	return time.Unix(timestamp/1000, timestamp), nil
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
	if err != nil {
		return nil, err
	}
	rx := regexp.MustCompile("<@&|>")
	resolvable = rx.ReplaceAllString(resolvable, "")

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

func FetchMember(s *discordgo.Session, guildID, resolvable string) (*discordgo.Member, error) {
	guild, err := s.Guild(guildID)
	if err != nil {
		return nil, err
	}
	rx := regexp.MustCompile("<@|!|>")
	resolvable = rx.ReplaceAllString(resolvable, "")

	checkFuncs := []func(*discordgo.Member, string) bool{
		func(r *discordgo.Member, resolvable string) bool {
			return r.User.ID == resolvable
		},
		func(r *discordgo.Member, resolvable string) bool {
			return r.User.Username == resolvable
		},
		func(r *discordgo.Member, resolvable string) bool {
			return strings.ToLower(r.User.Username) == strings.ToLower(resolvable)
		},
		func(r *discordgo.Member, resolvable string) bool {
			return strings.HasPrefix(strings.ToLower(r.User.Username), strings.ToLower(resolvable))
		},
		func(r *discordgo.Member, resolvable string) bool {
			return strings.Contains(strings.ToLower(r.User.Username), strings.ToLower(resolvable))
		},
		func(r *discordgo.Member, resolvable string) bool {
			return r.Nick == resolvable
		},
		func(r *discordgo.Member, resolvable string) bool {
			return r.Nick != "" && strings.ToLower(r.Nick) == strings.ToLower(resolvable)
		},
		func(r *discordgo.Member, resolvable string) bool {
			return r.Nick != "" && strings.HasPrefix(strings.ToLower(r.Nick), strings.ToLower(resolvable))
		},
		func(r *discordgo.Member, resolvable string) bool {
			return r.Nick != "" && strings.Contains(strings.ToLower(r.Nick), strings.ToLower(resolvable))
		},
	}

	for _, checkFunc := range checkFuncs {
		for _, m := range guild.Members {
			if checkFunc(m, resolvable) {
				return m, nil
			}
		}
	}

	return nil, errors.New("could not be fetched")
}

func FetchChannel(s *discordgo.Session, guildID, resolvable string, condition ...func(*discordgo.Channel) bool) (*discordgo.Channel, error) {
	guild, err := s.Guild(guildID)
	if err != nil {
		return nil, err
	}

	checkFuncs := []func(*discordgo.Channel, string) bool{
		func(r *discordgo.Channel, resolvable string) bool {
			return r.ID == resolvable
		},
		func(r *discordgo.Channel, resolvable string) bool {
			return r.Name == resolvable
		},
		func(r *discordgo.Channel, resolvable string) bool {
			return strings.ToLower(r.Name) == strings.ToLower(resolvable)
		},
		func(r *discordgo.Channel, resolvable string) bool {
			return strings.HasPrefix(strings.ToLower(r.Name), strings.ToLower(resolvable))
		},
		func(r *discordgo.Channel, resolvable string) bool {
			return strings.Contains(strings.ToLower(r.Name), strings.ToLower(resolvable))
		},
	}

	for _, checkFunc := range checkFuncs {
		for _, c := range guild.Channels {
			if len(condition) > 0 && condition[0] != nil {
				if !condition[0](c) {
					continue
				}
			}
			if checkFunc(c, resolvable) {
				return c, nil
			}
		}
	}

	return nil, errors.New("could not be fetched")
}
