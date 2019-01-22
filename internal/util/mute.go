package util

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

func MuteSetupChannels(s *discordgo.Session, guildID, roleID string) error {
	guild, err := s.Guild(guildID)
	if err != nil {
		return err
	}

	var roleExists bool
	for _, r := range guild.Roles {
		if r.ID == roleID && !roleExists {
			roleExists = true
		}
	}
	if !roleExists {
		return errors.New("role does not exist on guild")
	}

	for _, c := range guild.Channels {
		if c.Type != discordgo.ChannelTypeGuildText {
			continue
		}
		err = s.ChannelPermissionSet(c.ID, roleID, "role", 0, 0x00000800)
	}

	return err
}
