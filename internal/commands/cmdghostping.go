package commands

import (
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/core"
	"github.com/zekroTJA/shinpuru/internal/util"
)

const defMsgPattern = "{pinger} ghost pinged {pinged} with message:\n\n{msg}"

type CmdGhostping struct {
}

func (c *CmdGhostping) GetInvokes() []string {
	return []string{"ghost", "gp", "ghostping", "gping"}
}

func (c *CmdGhostping) GetDescription() string {
	return "Send a message when someone ghost pinged a member"
}

func (c *CmdGhostping) GetHelp() string {
	return "`ghost` - display current ghost ping settings\n" +
		"`ghost set (<msgPattern>)` - Set a ghost ping message pattern. If no 2nd argument is provided, the default pattern will be used.\n" +
		"`ghost reset` - reset message and disable ghost ping warnings\n\n" +
		"Usable variables in message pattern:\n" +
		"- `{@pinger}` - mention of the user sent the ghost ping\n" +
		"- `{pinger}` - username#discriminator of the user sent the ghost ping\n" +
		"- `{@pinged}` - mention of the user got ghost pinged\n" +
		"- `{pinged}` - username#discriminator of the user got ghost pinged\n" +
		"- `{msg}` - the content of the message which ghost pinged\n\n" +
		"Default message pattern:\n```\n" + defMsgPattern + "\n```"
}

func (c *CmdGhostping) GetGroup() string {
	return GroupModeration
}

func (c *CmdGhostping) GetDomainName() string {
	return "sp.guild.mod.ghostping"
}

func (c *CmdGhostping) Exec(args *CommandArgs) error {
	if len(args.Args) < 1 {
		return c.info(args)
	}

	switch strings.ToLower(args.Args[0]) {
	case "set", "setup", "create":
		return c.set(args)
	case "reset", "remove", "delete":
		return c.reset(args)
	default:
		return c.info(args)
	}
}

func (c *CmdGhostping) info(args *CommandArgs) error {
	emb := &discordgo.MessageEmbed{
		Color: util.ColorEmbedDefault,
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name: "Help",
				Value: "Set up or edit a ghost ping warn with `ghost set`. You can disable this with `ghost reset`." +
					"Enter `help ghost` for further information.",
			},
		},
	}

	gpMsg, err := args.CmdHandler.db.GetGuildGhostpingMsg(args.Guild.ID)
	if err != nil && !core.IsErrDatabaseNotFound(err) {
		return err
	}

	if gpMsg != "" {
		emb.Description = "Currently set message:\n\n" + gpMsg
	} else {
		emb.Description = "*Currently, no message was set so this function is disabled.*"
	}

	msg, err := args.Session.ChannelMessageSendEmbed(args.Channel.ID, emb)
	util.DeleteMessageLater(args.Session, msg, 15*time.Second)
	return err
}

func (c *CmdGhostping) set(args *CommandArgs) error {
	msgPattern := defMsgPattern

	if len(args.Args) > 2 {
		msgPattern = strings.Join(args.Args[1:], " ")
	}

	if err := args.CmdHandler.db.SetGuildGhostpingMsg(args.Guild.ID, msgPattern); err != nil {
		return err
	}

	msg, err := util.SendEmbed(args.Session, args.Channel.ID,
		"Set message pattern as ghost ping warn:\n"+msgPattern+"\n\n"+
			"*Use `ghost reset` to disable ghost ping warnings or use `help ghost` for further information.*", "", util.ColorEmbedUpdated)
	util.DeleteMessageLater(args.Session, msg, 15*time.Second)
	return err
}

func (c *CmdGhostping) reset(args *CommandArgs) error {
	if err := args.CmdHandler.db.SetGuildGhostpingMsg(args.Guild.ID, ""); err != nil {
		return err
	}

	msg, err := util.SendEmbed(args.Session, args.Channel.ID,
		"Warn message reset and ghost ping warnings disabled.", "", util.ColorEmbedUpdated)
	util.DeleteMessageLater(args.Session, msg, 8*time.Second)
	return err
}
