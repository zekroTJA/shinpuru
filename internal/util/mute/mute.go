package mute

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

// SetupChannels tries to set the permission for
// each text channel of the passed guild for the
// passed role ID to disable permission to write.
func SetupChannels(s *discordgo.Session, guildID, roleID string) error {
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
		if err = s.ChannelPermissionSet(c.ID, roleID, "role", 0, 0x00000800); err != nil {
			return err
		}
	}

	return nil
}
