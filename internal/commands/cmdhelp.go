package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sahilm/fuzzy"
	"github.com/zekroTJA/shinpuru/internal/core/config"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shireikan"
)

type CmdHelp struct {
}

func (c *CmdHelp) GetInvokes() []string {
	return []string{"help", "h", "?", "man"}
}

func (c *CmdHelp) GetDescription() string {
	return "Display list of command or get help for a specific command."
}

func (c *CmdHelp) GetHelp() string {
	return "`help` - display command list\n" +
		"`help <command>` - display help of specific command"
}

func (c *CmdHelp) GetGroup() string {
	return shireikan.GroupGeneral
}

func (c *CmdHelp) GetDomainName() string {
	return "sp.etc.help"
}

func (c *CmdHelp) GetSubPermissionRules() []shireikan.SubPermission {
	return nil
}

func (c *CmdHelp) IsExecutableInDMChannels() bool {
	return true
}

func (c *CmdHelp) Exec(ctx shireikan.Context) error {
	cmdhandler, _ := ctx.GetObject(shireikan.ObjectMapKeyHandler).(shireikan.Handler)
	cfg, _ := ctx.GetObject("config").(*config.Config)

	emb := &discordgo.MessageEmbed{
		Color:  static.ColorEmbedDefault,
		Fields: make([]*discordgo.MessageEmbedField, 0),
	}

	if len(ctx.GetArgs()) == 0 {
		commandListEmbed(cmdhandler.GetCommandInstances(), cfg, emb)
	} else {
		query := strings.ToLower(ctx.GetArgs().Get(0).AsString())
		cmd, ok := cmdhandler.GetCommand(query)

		if ok {
			commandEmbed(cmd, cfg, emb)
		} else {
			invokes := make([]string, len(cmdhandler.GetCommandMap()))

			i := 0
			for k := range cmdhandler.GetCommandMap() {
				invokes[i] = k
				i++
			}

			matches := fuzzy.Find(query, invokes)

			if matches.Len() == 0 {
				return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
					"Could not find any matches.", "").
					DeleteAfter(8 * time.Second).
					Error()
			}

			if matches.Len() == 1 {
				cmd, _ := cmdhandler.GetCommand(matches[0].Str)
				commandEmbed(cmd, cfg, emb)
			} else {
				cmdInstancesMap := make(map[string]shireikan.Command)
				for _, match := range matches {
					cmd := cmdhandler.GetCommandMap()[match.Str]
					cmdInstancesMap[cmd.GetInvokes()[0]] = cmd
				}

				cmdInstances := make([]shireikan.Command, len(cmdInstancesMap))
				i = 0
				for _, cmd := range cmdInstancesMap {
					cmdInstances[i] = cmd
					i++
				}

				commandListEmbed(cmdInstances, cfg, emb)
			}
		}
	}

	return sendUserOrInChannel(ctx, emb)
}

func commandListEmbed(cmdInstances []shireikan.Command, cfg *config.Config, emb *discordgo.MessageEmbed) {
	cmds := make(map[string][]shireikan.Command)
	for _, c := range cmdInstances {
		group := c.GetGroup()
		if _, ok := cmds[group]; !ok {
			cmds[group] = make([]shireikan.Command, 0)
		}
		cmds[group] = append(cmds[group], c)
	}

	manualURL := static.CommandManualDocument
	if cfg.WebServer != nil && cfg.WebServer.Enabled && cfg.WebServer.PublicAddr != "" {
		manualURL = cfg.WebServer.PublicAddr + "/commands"
	}

	emb.Title = "Command List"
	emb.Description = fmt.Sprintf(
		"[**Here**](%s) you can find the full list of commands with all details in one document.",
		manualURL)

	for cat, catCmds := range cmds {
		commandHelpLines := ""
		for _, c := range catCmds {
			commandHelpLines += fmt.Sprintf("`%s` - *%s* `[%s]`\n", c.GetInvokes()[0], c.GetDescription(), c.GetDomainName())
		}
		emb.Fields = append(emb.Fields, &discordgo.MessageEmbedField{
			Name:  cat,
			Value: commandHelpLines,
		})
	}
}

func commandEmbed(cmd shireikan.Command, cfg *config.Config, emb *discordgo.MessageEmbed) {
	manualURL := static.CommandManualDocument
	if cfg.WebServer != nil && cfg.WebServer.Enabled && cfg.WebServer.PublicAddr != "" {
		manualURL = cfg.WebServer.PublicAddr + "/commands"
	}

	emb.Title = "Command Description"
	emb.Description = fmt.Sprintf(
		"[**Here**](%s#%s) you can find the command description online.",
		manualURL, cmd.GetInvokes()[0])

	emb.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "Invokes",
			Value:  strings.Join(cmd.GetInvokes(), "\n"),
			Inline: true,
		},
		{
			Name:   "Group",
			Value:  cmd.GetGroup(),
			Inline: true,
		},
		{
			Name:   "Domain Name",
			Value:  cmd.GetDomainName(),
			Inline: true,
		},
		{
			Name: "DM Capable",
			Value: util.BoolAsString(
				cmd.IsExecutableInDMChannels(), "Yes", "No"),
			Inline: true,
		},
		{
			Name:  "Description",
			Value: util.EnsureNotEmpty(cmd.GetDescription(), "`no description`"),
		},
		{
			Name:  "Usage",
			Value: util.EnsureNotEmpty(cmd.GetHelp(), "`no uage information`"),
		},
	}

	if spr := cmd.GetSubPermissionRules(); spr != nil {
		txt := "*`[E]` in front of permissions means `Explicit`, which means that this " +
			"permission must be explicitly allowed and can not be wild-carded.\n" +
			"`[D]` implies that wildecards will apply to this sub permission.*\n\n"

		for _, rule := range spr {
			expl := "D"
			if rule.Explicit {
				expl = "E"
			}

			txt = fmt.Sprintf("%s`[%s]` %s.%s - *%s*\n",
				txt, expl, cmd.GetDomainName(), rule.Term, rule.Description)
		}

		emb.Fields = append(emb.Fields, &discordgo.MessageEmbedField{
			Name:  "Sub Permission Rules",
			Value: txt,
		})
	}
}

func sendUserOrInChannel(ctx shireikan.Context, emb *discordgo.MessageEmbed) (err error) {
	userChan, err := ctx.GetSession().UserChannelCreate(ctx.GetUser().ID)
	if err != nil {
		return
	}

	_, err = ctx.GetSession().ChannelMessageSendEmbed(userChan.ID, emb)
	if err != nil {
		if strings.Contains(err.Error(), `"Cannot send messages to this user"`) {
			emb.Footer = &discordgo.MessageEmbedFooter{
				Text: "This message appears in this channel because you have " +
					"disabled DM's from guild members.",
			}
			_, err = ctx.GetSession().ChannelMessageSendEmbed(ctx.GetChannel().ID, emb)
			return
		}
	}

	return
}
