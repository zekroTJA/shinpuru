package slashcommands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekrotja/ken"
)

type Debug struct{}

var (
	_ ken.SlashCommand        = (*User)(nil)
	_ permissions.PermCommand = (*User)(nil)
)

func (c *Debug) Name() string {
	return "debug"
}

func (c *Debug) Description() string {
	return "Debug command for development."
}

func (c *Debug) Version() string {
	return "1.0.0"
}

func (c *Debug) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *Debug) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{}
}

func (c *Debug) Domain() string {
	return "sp.debug"
}

func (c *Debug) SubDomains() []permissions.SubPermission {
	return nil
}

func (c *Debug) Run(ctx ken.Context) (err error) {
	return fmt.Errorf("test 123")
}
