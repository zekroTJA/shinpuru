package util

import (
	"sort"

	"github.com/bwmarrin/discordgo"
)

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

	sortRoleArray(membRoles, reversed)

	return membRoles, nil
}

func GetSortedGuildRoles(s *discordgo.Session, guildID string, reversed bool) ([]*discordgo.Role, error) {
	roles, err := s.GuildRoles(guildID)
	if err != nil {
		return nil, err
	}

	sortRoleArray(roles, reversed)

	return roles, nil
}

func sortRoleArray(r []*discordgo.Role, reversed bool) {
	f := func(i, j int) bool {
		return r[i].Position < r[j].Position
	}

	if reversed {
		f = func(i, j int) bool {
			return r[i].Position > r[j].Position
		}
	}

	sort.Slice(r, f)
}
