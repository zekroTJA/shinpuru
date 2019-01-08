package commands

import (
	"time"

	//"github.com/zekroTJA/shinpuru/core"
	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/util"
)

type CmdMute struct {
	PermLvl int
}

func (c *CmdMute) GetInvokes() []string {
	return []string{"mute", "m", "silence"}
}

func (c *CmdMute) GetDescription() string {
	return "Mute members in text channels"
}

func (c *CmdMute) GetHelp() string {
	return "`mute setup` - creates mute role and sets this role in every channel as muted\n" +
		"`mute <userResolvable>` - mute/unmute a user\n" +
		"`mute list` - display muted users on this guild"
}

func (c *CmdMute) GetGroup() string {
	return GroupModeration
}

func (c *CmdMute) GetPermission() int {
	return c.PermLvl
}

func (c *CmdMute) SetPermission(permLvl int) {
	c.PermLvl = permLvl
}

func (c *CmdMute) Exec(args *CommandArgs) error {
	if len(args.Args) < 1 {
		msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
			"Invalid arguments. Use `help mute` to get info how to use this command.")
		util.DeleteMessageLater(args.Session, msg, 8*time.Second)
		return err
	}

	switch args.Args[0] {
	case "setup":
		acmsg := &util.AcceptMessage{
			Session: args.Session,
			Embed: &discordgo.MessageEmbed{
				Color: util.ColorEmbedDefault,
				Title: "Warning",
				Description: "The follwoing will create a role with the name `shinpuru-muted` and will " +
					"set every channels *(which is visible to the bot)* permission for this role to " +
					"disallow write messages!",
			},
			UserID:         args.User.ID,
			DeleteMsgAfter: true,
			AcceptFunc: func(msg *discordgo.Message) {
				createdRole, err := args.Session.GuildRoleCreate(args.Guild.ID)
				if err != nil {
					msg, _ := util.SendEmbedError(args.Session, args.Channel.ID,
						"Failed creating mute role: ```\n"+err.Error()+"\n```")
					util.DeleteMessageLater(args.Session, msg, 30*time.Second)
					return
				}

				createdRole, err = args.Session.GuildRoleEdit(args.Guild.ID, createdRole.ID,
					"shinpuru-muted", 0, false, 0, false)
				if err != nil {
					msg, _ := util.SendEmbedError(args.Session, args.Channel.ID,
						"Failed editing mute role: ```\n"+err.Error()+"\n```")
					util.DeleteMessageLater(args.Session, msg, 30*time.Second)
					return
				}

				err = args.CmdHandler.db.SetMuteRole(args.Guild.ID, createdRole.ID)
				if err != nil {
					msg, _ := util.SendEmbedError(args.Session, args.Channel.ID,
						"Failed setting mute role in database: ```\n"+err.Error()+"\n```")
					util.DeleteMessageLater(args.Session, msg, 30*time.Second)
					return
				}

				err = util.MuteSetupChannels(args.Session, args.Guild.ID, createdRole.ID)
				if err != nil {
					msg, _ := util.SendEmbedError(args.Session, args.Channel.ID,
						"Failed updating channels: ```\n"+err.Error()+"\n```")
					util.DeleteMessageLater(args.Session, msg, 30*time.Second)
					return
				}

				msg, _ = util.SendEmbed(args.Session, args.Channel.ID,
					"Set up mute role and edited channel permissions.\nMaybe you need to increase the "+
						"position of the role to override other roles permission settings.",
					"", util.ColorEmbedUpdated)
				util.DeleteMessageLater(args.Session, msg, 12*time.Second)
			},
			DeclineFunc: func(msg *discordgo.Message) {
				msg, _ = util.SendEmbedError(args.Session, args.Channel.ID,
					"Setup canceled.")
				util.DeleteMessageLater(args.Session, msg, 5*time.Second)
			},
		}

		_, err := acmsg.Send(args.Channel.ID)
		return err
	}

	return nil
}
