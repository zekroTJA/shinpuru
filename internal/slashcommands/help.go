package slashcommands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/services/config"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekrotja/ken"
)

type Help struct{}

var (
	_ ken.Command             = (*Help)(nil)
	_ permissions.PermCommand = (*Help)(nil)
	_ ken.DmCapable           = (*Help)(nil)
)

func (c *Help) Name() string {
	return "help"
}

func (c *Help) Description() string {
	return "Show the shinpuru help center."
}

func (c *Help) Version() string {
	return "1.0.0"
}

func (c *Help) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *Help) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{}
}

func (c *Help) Domain() string {
	return "sp.etc.help"
}

func (c *Help) SubDomains() []permissions.SubPermission {
	return nil
}

func (c *Help) IsDmCapable() bool {
	return true
}

func (c *Help) Run(ctx *ken.Ctx) (err error) {
	cfg := ctx.Get(static.DiConfig).(config.Provider).Config()

	webAddr := cfg.WebServer.PublicAddr
	webEnabled := cfg.WebServer.Enabled && webAddr != ""

	cmdHelp := "[Here](https://github.com/zekroTJA/shinpuru/wiki/Commands) you can find a list of all commands " +
		"and detailed meta and help information."
	if webEnabled {
		cmdHelp += fmt.Sprintf("\nThere is also an interactive [web view](%s/commands) where you can "+
			"discover and look up command information.", webAddr)
	}

	emb := &discordgo.MessageEmbed{
		Title: "Help Center",
		Description: "If you generally need help with the usage or setup of shinpuru, take a look into the " +
			"[**Wiki**](https://github.com/zekroTJA/shinpuru/wiki). There you can find a lot of useful resources " +
			"around shinpuru's features.",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Command Help",
				Value: cmdHelp,
			},
			{
				Name: "Contact",
				Value: "If you need help with shinpuru, feel free to ask on my [dev discord](https://dc.zekro.de). " +
					"You can also contact me directly, either via email (`contact@zekro.de`) or via " +
					"[twitter](https://twitter.com/zekrotja).",
			},
		},
	}

	err = ctx.RespondEmbed(emb)
	return
}
