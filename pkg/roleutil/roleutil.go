// Package roleutil provides general purpose
// utilities for discordgo.Role objects and
// arrays.
package roleutil

import (
	"sort"

	"github.com/bwmarrin/discordgo"
)

// SortRoles sorts a given array of discordgo.Role
// object references by position in ascending order.
// If reversed, the order is descending.
func SortRoles(r []*discordgo.Role, reversed bool) {
	var f func(i, j int) bool

	if reversed {
		f = func(i, j int) bool {
			return r[i].Position > r[j].Position
		}
	} else {
		f = func(i, j int) bool {
			return r[i].Position < r[j].Position
		}
	}

	sort.Slice(r, f)
}

// GetSortedMemberRoles tries to fetch the roles of a given
// member on a given guild and returns the role array in
// sorted ascending order by position.
// If any error occurs, the error is returned as well.
// If reversed, the order is descending.
func GetSortedMemberRoles(s *discordgo.Session, guildID, memberID string, reversed bool) ([]*discordgo.Role, error) {
	member, err := s.GuildMember(guildID, memberID)
	if err != nil {
		return nil, err
	}

	roles, err := s.GuildRoles(guildID)
	if err != nil {
		return nil, err
	}

	rolesMap := make(map[string]*discordgo.Role)
	for _, r := range roles {
		rolesMap[r.ID] = r
	}

	membRoles := make([]*discordgo.Role, len(member.Roles))
	applied := 0
	for _, rID := range member.Roles {
		if r, ok := rolesMap[rID]; ok {
			membRoles[applied] = r
			applied++
		}
	}

	membRoles = membRoles[:applied]

	SortRoles(membRoles, reversed)

	return membRoles, nil
}

// GetSortedGuildRoles tries to fetch the roles of a given
// guild and returns the role array in sorted ascending
// order by position.
// If any error occurs, the error is returned as well.
// If reversed, the order is descending.
func GetSortedGuildRoles(s *discordgo.Session, guildID string, reversed bool) ([]*discordgo.Role, error) {
	roles, err := s.GuildRoles(guildID)
	if err != nil {
		return nil, err
	}

	SortRoles(roles, reversed)

	return roles, nil
}
