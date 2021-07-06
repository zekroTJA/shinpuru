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

	chans, err := s.GuildChannels(guildID)
	if err != nil {
		return err
	}

	for _, c := range chans {
		if c.Type != discordgo.ChannelTypeGuildText {
			continue
		}
		if err = s.ChannelPermissionSet(c.ID, roleID, discordgo.PermissionOverwriteTypeRole, 0, 0x00000800); err != nil {
			return err
		}
	}

	return nil
}
