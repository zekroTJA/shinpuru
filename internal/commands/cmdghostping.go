package commands

import (
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
	"github.com/zekroTJA/shireikan"
)

const defMsgPattern = "{pinger} ghost pinged {pinged} with message:\n\n{msg}"

type CmdGhostping struct {
}

func (c *CmdGhostping) GetInvokes() []string {
	return []string{"ghost", "gp", "ghostping", "gping"}
}

func (c *CmdGhostping) GetDescription() string {
	return "Send a message when someone ghost pinged a member."
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
	return shireikan.GroupModeration
}

func (c *CmdGhostping) GetDomainName() string {
	return "sp.guild.mod.ghostping"
}

func (c *CmdGhostping) GetSubPermissionRules() []shireikan.SubPermission {
	return nil
}

func (c *CmdGhostping) IsExecutableInDMChannels() bool {
	return false
}

func (c *CmdGhostping) Exec(ctx shireikan.Context) error {
	if len(ctx.GetArgs()) < 1 {
		return c.info(ctx)
	}

	switch strings.ToLower(ctx.GetArgs().Get(0).AsString()) {
	case "set", "setup", "create":
		return c.set(ctx)
	case "reset", "remove", "delete":
		return c.reset(ctx)
	default:
		return c.info(ctx)
	}
}

func (c *CmdGhostping) info(ctx shireikan.Context) error {
	emb := &discordgo.MessageEmbed{
		Color: static.ColorEmbedDefault,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name: "Help",
				Value: "Set up or edit a ghost ping warn with `ghost set`. You can disable this with `ghost reset`." +
					"Enter `help ghost` for further information.",
			},
		},
	}

	db, _ := ctx.GetObject(static.DiDatabase).(database.Database)
	gpMsg, err := db.GetGuildGhostpingMsg(ctx.GetGuild().ID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	if gpMsg != "" {
		emb.Description = "Currently set message:\n\n" + gpMsg
	} else {
		emb.Description = "*Currently, no message was set so this function is disabled.*"
	}

	msg, err := ctx.GetSession().ChannelMessageSendEmbed(ctx.GetChannel().ID, emb)
	discordutil.DeleteMessageLater(ctx.GetSession(), msg, 15*time.Second)
	return err
}

func (c *CmdGhostping) set(ctx shireikan.Context) error {
	msgPattern := defMsgPattern

	if len(ctx.GetArgs()) > 2 {
		msgPattern = strings.Join(ctx.GetArgs()[1:], " ")
	}

	db, _ := ctx.GetObject(static.DiDatabase).(database.Database)
	if err := db.SetGuildGhostpingMsg(ctx.GetGuild().ID, msgPattern); err != nil {
		return err
	}

	return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
		"Set message pattern as ghost ping warn:\n"+msgPattern+"\n\n"+
			"*Use `ghost reset` to disable ghost ping warnings or use `help ghost` for further information.*", "", static.ColorEmbedUpdated).
		DeleteAfter(15 * time.Second).Error()
}

func (c *CmdGhostping) reset(ctx shireikan.Context) error {
	db, _ := ctx.GetObject(static.DiDatabase).(database.Database)
	if err := db.SetGuildGhostpingMsg(ctx.GetGuild().ID, ""); err != nil {
		return err
	}

	return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
		"Warn message reset and ghost ping warnings disabled.", "", static.ColorEmbedUpdated).
		DeleteAfter(8 * time.Second).Error()
}
