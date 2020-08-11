package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
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
	return "display list of command or get help for a specific command"
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

	emb := &discordgo.MessageEmbed{
		Color:  static.ColorEmbedDefault,
		Fields: make([]*discordgo.MessageEmbedField, 0),
	}

	if len(ctx.GetArgs()) == 0 {
		cmds := make(map[string][]shireikan.Command)
		for _, c := range cmdhandler.GetCommandInstances() {
			group := c.GetGroup()
			if _, ok := cmds[group]; !ok {
				cmds[group] = make([]shireikan.Command, 0)
			}
			cmds[group] = append(cmds[group], c)
		}

		emb.Title = "Command List"
		emb.Description = fmt.Sprintf(
			"[**Here**](%s) you can find the full list of commands with all details in one document.",
			static.CommandManualDocument)

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
	} else {
		cmd, ok := cmdhandler.GetCommand(ctx.GetArgs().Get(0).AsString())
		if !ok {
			return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
				fmt.Sprintf("Sorry, there is no command with the invoke `%s`", ctx.GetArgs().Get(0).AsString())).
				DeleteAfter(5 * time.Second).Error()
		}

		emb.Title = "Command Description"
		emb.Description = fmt.Sprintf(
			"[**Here**](%s#%s) you can find the command description online.",
			static.CommandManualDocument, cmd.GetInvokes()[0])

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

	userChan, err := ctx.GetSession().UserChannelCreate(ctx.GetUser().ID)
	if err != nil {
		return err
	}
	_, err = ctx.GetSession().ChannelMessageSendEmbed(userChan.ID, emb)
	if err != nil {
		if strings.Contains(err.Error(), `"Cannot send messages to this user"`) {
			emb.Footer = &discordgo.MessageEmbedFooter{
				Text: "Actually, this message appears in your DM, but you have deactivated receiving DMs from" +
					"server members, so I can not send you this message via DM and you see this here right now.",
			}
			_, err = ctx.GetSession().ChannelMessageSendEmbed(ctx.GetChannel().ID, emb)
			return err
		}
	}

	return err
}
