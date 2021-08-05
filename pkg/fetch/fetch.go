// Package fetch provides functionalities to fetch roles,
// channels, members and users by so called resolavbles.
// That means, these functions try to match a member, role
// or channel by their names, displaynames, IDs or mentions
// as greedy as prossible.
package fetch

import (
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	RoleCheckFuncs = []func(*discordgo.Role, string) bool{
		// 1. ID exact match
		func(r *discordgo.Role, resolvable string) bool {
			return r.ID == resolvable
		},
		// 2. name exact match
		func(r *discordgo.Role, resolvable string) bool {
			return r.Name == resolvable
		},
		// 3. name lowercased exact match
		func(r *discordgo.Role, resolvable string) bool {
			return strings.ToLower(r.Name) == strings.ToLower(resolvable)
		},
		// 4. name lowercased startswith
		func(r *discordgo.Role, resolvable string) bool {
			return strings.HasPrefix(strings.ToLower(r.Name), strings.ToLower(resolvable))
		},
		// 5. name lowercased contains
		func(r *discordgo.Role, resolvable string) bool {
			return strings.Contains(strings.ToLower(r.Name), strings.ToLower(resolvable))
		},
	}

	MemberCheckFuncs = []func(*discordgo.Member, string) bool{
		// 1. ID exact match
		func(r *discordgo.Member, resolvable string) bool {
			return r.User.ID == resolvable
		},
		// 2. username exact match
		func(r *discordgo.Member, resolvable string) bool {
			return r.User.Username == resolvable
		},
		// 3. username lowercased exact match
		func(r *discordgo.Member, resolvable string) bool {
			return strings.ToLower(r.User.Username) == strings.ToLower(resolvable)
		},
		// 4. username lowercased startswith
		func(r *discordgo.Member, resolvable string) bool {
			return strings.HasPrefix(strings.ToLower(r.User.Username), strings.ToLower(resolvable))
		},
		// 5. username lowercased contains
		func(r *discordgo.Member, resolvable string) bool {
			return strings.Contains(strings.ToLower(r.User.Username), strings.ToLower(resolvable))
		},
		// 6. nick exact match
		func(r *discordgo.Member, resolvable string) bool {
			return r.Nick == resolvable
		},
		// 7. nick lowercased exact match
		func(r *discordgo.Member, resolvable string) bool {
			return r.Nick != "" && strings.ToLower(r.Nick) == strings.ToLower(resolvable)
		},
		// 8. nick lowercased starts with
		func(r *discordgo.Member, resolvable string) bool {
			return r.Nick != "" && strings.HasPrefix(strings.ToLower(r.Nick), strings.ToLower(resolvable))
		},
		// 9. nick lowercased contains
		func(r *discordgo.Member, resolvable string) bool {
			return r.Nick != "" && strings.Contains(strings.ToLower(r.Nick), strings.ToLower(resolvable))
		},
	}

	ChannelCheckFuncs = []func(*discordgo.Channel, string) bool{
		// 1. ID exact match
		func(r *discordgo.Channel, resolvable string) bool {
			return r.ID == resolvable
		},
		// 2. mention exact match
		func(r *discordgo.Channel, resolvable string) bool {
			l := len(resolvable)
			return l > 3 && r.ID == resolvable[2:l-1]
		},
		// 3. name exact match
		func(r *discordgo.Channel, resolvable string) bool {
			return r.Name == resolvable
		},
		// 4. name lowercased exact match
		func(r *discordgo.Channel, resolvable string) bool {
			return strings.ToLower(r.Name) == strings.ToLower(resolvable)
		},
		// 5. name lowercased starts with
		func(r *discordgo.Channel, resolvable string) bool {
			return strings.HasPrefix(strings.ToLower(r.Name), strings.ToLower(resolvable))
		},
		// 6. name lowercased contains
		func(r *discordgo.Channel, resolvable string) bool {
			return strings.Contains(strings.ToLower(r.Name), strings.ToLower(resolvable))
		},
	}

	GuildCheckFuncs = []func(*discordgo.Guild, string) bool{
		// 1. ID exact match
		func(r *discordgo.Guild, resolvable string) bool {
			return r.ID == resolvable
		},
		// 2. mention exact match
		func(r *discordgo.Guild, resolvable string) bool {
			l := len(resolvable)
			return l > 3 && r.ID == resolvable[2:l-1]
		},
		// 3. name exact match
		func(r *discordgo.Guild, resolvable string) bool {
			return r.Name == resolvable
		},
		// 4. name lowercased exact match
		func(r *discordgo.Guild, resolvable string) bool {
			return strings.ToLower(r.Name) == strings.ToLower(resolvable)
		},
		// 5. name lowercased starts with
		func(r *discordgo.Guild, resolvable string) bool {
			return strings.HasPrefix(strings.ToLower(r.Name), strings.ToLower(resolvable))
		},
		// 6. name lowercased contains
		func(r *discordgo.Guild, resolvable string) bool {
			return strings.Contains(strings.ToLower(r.Name), strings.ToLower(resolvable))
		},
	}
)

// FetchRoles tries to fetch a role on the specified guild
// by given resolvable and returns this role, when found.
// You can pass a condition function which ignores the result
// if this functions returns false on the given object.
// If no object was found, ErrNotFound is returned.
// If any other unexpected error occurs during fetching,
// this error is returned as well.
func FetchRole(s DataOutlet, guildID, resolvable string, condition ...func(*discordgo.Role) bool) (*discordgo.Role, error) {
	roles, err := s.GuildRoles(guildID)
	if err != nil {
		return nil, err
	}
	rx := regexp.MustCompile("<@&|>")
	resolvable = rx.ReplaceAllString(resolvable, "")

	for _, checkFunc := range RoleCheckFuncs {
		for _, r := range roles {
			if len(condition) > 0 && condition[0] != nil {
				if !condition[0](r) {
					continue
				}
			}
			if checkFunc(r, resolvable) {
				return r, nil
			}
		}
	}

	return nil, ErrNotFound
}

// FetchMember tries to fetch a member on the specified guild
// by given resolvable and returns this member, when found.
// You can pass a condition function which ignores the result
// if this functions returns false on the given object.
// If no object was found, ErrNotFound is returned.
// If any other unexpected error occurs during fetching,
// this error is returned as well.
func FetchMember(s DataOutlet, guildID, resolvable string, condition ...func(*discordgo.Member) bool) (*discordgo.Member, error) {
	rx := regexp.MustCompile("<@|!|>")
	resolvable = rx.ReplaceAllString(resolvable, "")
	var lastUserID string

	for {
		members, err := s.GuildMembers(guildID, lastUserID, 1000)
		if err != nil {
			return nil, err
		}

		if len(members) < 1 {
			break
		}

		lastUserID = members[len(members)-1].User.ID

		for _, checkFunc := range MemberCheckFuncs {
			for _, m := range members {
				if len(condition) > 0 && condition[0] != nil {
					if !condition[0](m) {
						continue
					}
				}
				if checkFunc(m, resolvable) {
					return m, nil
				}
			}
		}
	}

	return nil, ErrNotFound
}

// FetchChannel tries to fetch a channel on the specified guild
// by given resolvable and returns this channel, when found.
// You can pass a condition function which ignores the result
// if this functions returns false on the given object.
// If no object was found, ErrNotFound is returned.
// If any other unexpected error occurs during fetching,
// this error is returned as well.
func FetchChannel(s DataOutlet, guildID, resolvable string, condition ...func(*discordgo.Channel) bool) (*discordgo.Channel, error) {
	channels, err := s.GuildChannels(guildID)
	if err != nil {
		return nil, err
	}

	for _, checkFunc := range ChannelCheckFuncs {
		for _, c := range channels {
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

	return nil, ErrNotFound
}
