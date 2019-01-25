package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/core"
	"github.com/zekroTJA/shinpuru/internal/util"
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
		return c.setup(args)
	case "list":
		return c.list(args)
	default:
		return c.muteUnmute(args)
	}
}

func (c *CmdMute) setup(args *CommandArgs) error {
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
			var muteRole *discordgo.Role
			for _, r := range args.Guild.Roles {
				if r.Name == util.MutedRoleName {
					muteRole = r
				}
			}

			if muteRole == nil {
				muteRole, err := args.Session.GuildRoleCreate(args.Guild.ID)
				if err != nil {
					msg, _ := util.SendEmbedError(args.Session, args.Channel.ID,
						"Failed creating mute role: ```\n"+err.Error()+"\n```")
					util.DeleteMessageLater(args.Session, msg, 30*time.Second)
					return
				}

				muteRole, err = args.Session.GuildRoleEdit(args.Guild.ID, muteRole.ID,
					util.MutedRoleName, 0, false, 0, false)
				if err != nil {
					msg, _ := util.SendEmbedError(args.Session, args.Channel.ID,
						"Failed editing mute role: ```\n"+err.Error()+"\n```")
					util.DeleteMessageLater(args.Session, msg, 30*time.Second)
					return
				}
			}

			err := args.CmdHandler.db.SetMuteRole(args.Guild.ID, muteRole.ID)
			if err != nil {
				msg, _ := util.SendEmbedError(args.Session, args.Channel.ID,
					"Failed setting mute role in database: ```\n"+err.Error()+"\n```")
				util.DeleteMessageLater(args.Session, msg, 30*time.Second)
				return
			}

			err = util.MuteSetupChannels(args.Session, args.Guild.ID, muteRole.ID)
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

func (c *CmdMute) muteUnmute(args *CommandArgs) error {
	victim, err := util.FetchMember(args.Session, args.Guild.ID, args.Args[0])
	if err != nil {
		msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
			"Could not fetch any user by the passed resolvable.")
		util.DeleteMessageLater(args.Session, msg, 6*time.Second)
		return err
	}

	muteRoleID, err := args.CmdHandler.db.GetMuteRoleGuild(args.Guild.ID)
	if core.IsErrDatabaseNotFound(err) {
		msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
			"Mute command is not set up. Please enter `mute setup`.")
		util.DeleteMessageLater(args.Session, msg, 6*time.Second)
		return err
	} else if err != nil {
		return err
	}

	repType := util.IndexOfStrArray("MUTE", util.ReportTypes)
	repID := util.ReportNodes[repType].Generate()

	var roleExists bool
	for _, r := range args.Guild.Roles {
		if r.ID == muteRoleID && !roleExists {
			roleExists = true
		}
	}
	if !roleExists {
		msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
			"Mute role does not exist on this guild. Please enter `mute setup`.")
		util.DeleteMessageLater(args.Session, msg, 6*time.Second)
		return err
	}

	var victimIsMuted bool
	for _, rID := range victim.Roles {
		if rID == muteRoleID && !victimIsMuted {
			victimIsMuted = true
		}
	}
	if victimIsMuted {
		err := args.Session.GuildMemberRoleRemove(args.Guild.ID, victim.User.ID, muteRoleID)
		if err != nil {
			return err
		}
		emb := &discordgo.MessageEmbed{
			Title: "Case " + repID.String(),
			Color: util.ReportColors[repType],
			Fields: []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{
					Inline: true,
					Name:   "Executor",
					Value:  fmt.Sprintf("<@%s>", args.User.ID),
				},
				&discordgo.MessageEmbedField{
					Inline: true,
					Name:   "Victim",
					Value:  fmt.Sprintf("<@%s>", victim.User.ID),
				},
				&discordgo.MessageEmbedField{
					Name:  "Type",
					Value: "UNMUTE",
				},
				&discordgo.MessageEmbedField{
					Name:  "Description",
					Value: "MANUAL UNMUTE",
				},
			},
			Timestamp: time.Unix(repID.Time()/1000, 0).Format("2006-01-02T15:04:05.000Z"),
		}
		args.Session.ChannelMessageSendEmbed(args.Channel.ID, emb)
		if modlogChan, err := args.CmdHandler.db.GetGuildModLog(args.Guild.ID); err == nil {
			args.Session.ChannelMessageSendEmbed(modlogChan, emb)
		}
		dmChan, err := args.Session.UserChannelCreate(victim.User.ID)
		if err == nil {
			args.Session.ChannelMessageSendEmbed(dmChan.ID, emb)
		}
		return err
	}

	err = args.Session.GuildMemberRoleAdd(args.Guild.ID, victim.User.ID, muteRoleID)
	if err != nil {
		return err
	}

	reason := "no reason set"
	if len(args.Args) > 1 {
		reason = strings.Join(args.Args[1:], " ")
	}

	rep := &util.Report{
		ID:         repID,
		Type:       repType,
		GuildID:    args.Guild.ID,
		ExecutorID: args.User.ID,
		VictimID:   victim.User.ID,
		Msg:        reason,
	}
	err = args.CmdHandler.db.AddReport(rep)
	if err != nil {
		util.SendEmbedError(args.Session, args.Channel.ID,
			"Failed creating report: ```\n"+err.Error()+"\n```")
	} else {
		args.Session.ChannelMessageSendEmbed(args.Channel.ID, rep.AsEmbed())
		if modlogChan, err := args.CmdHandler.db.GetGuildModLog(args.Guild.ID); err == nil {
			args.Session.ChannelMessageSendEmbed(modlogChan, rep.AsEmbed())
		}
		dmChan, err := args.Session.UserChannelCreate(victim.User.ID)
		if err == nil {
			args.Session.ChannelMessageSendEmbed(dmChan.ID, rep.AsEmbed())
		}
	}

	return err
}

func (c *CmdMute) list(args *CommandArgs) error {
	muteRoleID, err := args.CmdHandler.db.GetMuteRoleGuild(args.Guild.ID)
	if err != nil {
		return err
	}

	emb := &discordgo.MessageEmbed{
		Color:       util.ColorEmbedGray,
		Description: "Fetching muted members...",
		Fields:      make([]*discordgo.MessageEmbedField, 0),
	}

	msg, err := args.Session.ChannelMessageSendEmbed(args.Channel.ID, emb)
	if err != nil {
		return err
	}

	muteReports, err := args.CmdHandler.db.GetReportsFiltered(args.Guild.ID, "",
		util.IndexOfStrArray("MUTE", util.ReportTypes))

	muteReportsMap := make(map[string]*util.Report)
	for _, r := range muteReports {
		muteReportsMap[r.VictimID] = r
	}

	for _, m := range args.Guild.Members {
		if util.IndexOfStrArray(muteRoleID, m.Roles) > -1 {
			if r, ok := muteReportsMap[m.User.ID]; ok {
				emb.Fields = append(emb.Fields, &discordgo.MessageEmbedField{
					Name: fmt.Sprintf("CaseID: %d", r.ID),
					Value: fmt.Sprintf("<@%s> since `%s` with reason:\n%s",
						m.User.ID, r.GetTimestamp().Format(time.RFC1123), r.Msg),
				})
			}
		}
	}

	emb.Color = util.ColorEmbedDefault
	emb.Description = ""

	_, err = args.Session.ChannelMessageEditEmbed(args.Channel.ID, msg.ID, emb)
	return err
}
