package slashcommands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/services/config"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/embedbuilder"
	"github.com/zekroTJA/shinpuru/pkg/stringutil"
	"github.com/zekrotja/ken"
)

type Help struct{}

var (
	_ ken.SlashCommand        = (*Help)(nil)
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
	return "1.1.0"
}

func (c *Help) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *Help) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "command",
			Description: "A specific command you want to know more about.",
		},
	}
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
	} else {
		webAddr = ""
	}

	if nameV, ok := ctx.Options().GetByNameOptional("command"); ok {
		err = c.cmdHelp(ctx, webAddr, nameV.StringValue())
		return
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

func (c *Help) cmdHelp(ctx *ken.Ctx, webAddr, name string) (err error) {
	name = strings.ToLower(name)

	var info *ken.CommandInfo
	for _, info = range ctx.GetKen().GetCommandInfo() {
		if info.ApplicationCommand.Name == name {
			break
		}
	}

	if info == nil {
		return ctx.RespondError("Could not find any command with this name.", "")
	}

	dmCapable := false
	if imp, ok := info.Implementations["IsDmCapable"]; ok && len(imp) != 0 && imp[0].(bool) {
		dmCapable = true
	}

	domain := info.Implementations["Domain"][0].(string)

	emb := embedbuilder.New().
		WithTitle(fmt.Sprintf("/%s Command Help", name)).
		WithDescription(info.ApplicationCommand.Description).
		WithFooter(fmt.Sprintf("Command Version v%s", info.ApplicationCommand.Version), "", "").
		AddInlineField("Domain", domain).
		AddInlineField("Dm Capable", stringutil.FromBool(dmCapable, "Yes", "No"))

	if imp, ok := info.Implementations["SubDomains"]; ok && len(imp) != 0 {
		if sdns, ok := imp[0].([]permissions.SubPermission); ok && len(sdns) != 0 {
			var sdnsTxt strings.Builder
			for _, sdn := range sdns {
				fmt.Fprintf(&sdnsTxt, "`%s`%s - *%s*\n",
					getTermAssembly(domain, sdn.Term),
					stringutil.FromBool(sdn.Explicit, " [explicit]", ""),
					sdn.Description)
			}
			emb.AddField("Sub Domains", sdnsTxt.String())
		}
	}

	if webAddr != "" {
		emb.WithURL(fmt.Sprintf("%s/commands/#%s", webAddr, name))
	}

	options := info.ApplicationCommand.Options
	hasSubs := len(options) != 0 && options[0].Type == discordgo.ApplicationCommandOptionSubCommand

	if hasSubs {
		var subTxt strings.Builder
		for _, sub := range options {
			fmt.Fprintf(&subTxt, "\n\n**__%s__**\n*%s*",
				sub.Name, sub.Description)
			for _, opt := range sub.Options {
				fmt.Fprintf(&subTxt, "\n%s %s `%s`: *%s*",
					stringutil.FromBool(opt.Required, ":small_orange_diamond:", ":white_small_square:"),
					opt.Name, opt.Type.String(), opt.Description)
			}
		}
		emb.AddField("Sub Commands", subTxt.String())
	} else if len(options) != 0 {
		var optTxt strings.Builder
		for _, opt := range options {
			fmt.Fprintf(&optTxt, "\n%s %s `%s`: *%s*",
				stringutil.FromBool(opt.Required, ":small_orange_diamond:", ":white_small_square:"),
				opt.Name, opt.Type.String(), opt.Description)
		}
		emb.AddField("Arguments", optTxt.String())
	}

	return ctx.RespondEmbed(emb.Build())
}

func getTermAssembly(domain, term string) string {
	if strings.HasPrefix(term, "/") {
		return term[1:]
	}
	return domain + "." + term
}
