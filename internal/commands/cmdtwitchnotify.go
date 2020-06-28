package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/zekroTJA/shinpuru/internal/core/twitchnotify"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/acceptmsg"
)

type CmdTwitchNotify struct {
}

func (c *CmdTwitchNotify) GetInvokes() []string {
	return []string{"twitch", "tn", "twitchnotify"}
}

func (c *CmdTwitchNotify) GetDescription() string {
	return "Get notifications in channels when someone goes live on twitch"
}

func (c *CmdTwitchNotify) GetHelp() string {
	return "`twitch` - list all currently monitored twitch channels\n" +
		"`twitch <twitchUsername>` - get notified in the current channel when the streamer goes online\n" +
		"`twitch remove <twitchUsername>` - remove monitor"
}

func (c *CmdTwitchNotify) GetGroup() string {
	return GroupChat
}

func (c *CmdTwitchNotify) GetDomainName() string {
	return "sp.chat.twitch"
}

func (c *CmdTwitchNotify) GetSubPermissionRules() []SubPermission {
	return nil
}

func (c *CmdTwitchNotify) Exec(args *CommandArgs) error {
	tnw := args.CmdHandler.tnw

	if tnw == nil {
		return util.SendEmbedError(args.Session, args.Channel.ID,
			"This feature is disabled because no Twitch App ID was provided.").
			DeleteAfter(8 * time.Second).Error()
	}

	if len(args.Args) < 1 {
		nots, err := args.CmdHandler.db.GetAllTwitchNotifies("")
		if err != nil {
			return err
		}

		notsStr := ""

		for _, not := range nots {
			if not.GuildID == args.Guild.ID {
				if tUser, err := tnw.GetUser(not.TwitchUserID, twitchnotify.IdentID); err == nil {
					notsStr += fmt.Sprintf(":white_small_square:  **%s** in <#%s>\n",
						tUser.DisplayName, not.ChannelID)
				}
			}
		}

		return util.SendEmbed(args.Session, args.Channel.ID, notsStr, "Currently monitored streamers", 0).
			Error()
	}

	tUser, err := tnw.GetUser(args.Args[len(args.Args)-1], twitchnotify.IdentLogin)
	if err != nil && err.Error() == "not found" {
		return util.SendEmbedError(args.Session, args.Channel.ID,
			"Twitch user could not be found.").
			DeleteAfter(8 * time.Second).Error()
	} else if err != nil {
		return err
	}

	if len(args.Args) > 1 && strings.ToLower(args.Args[0]) == "remove" {
		nots, err := args.CmdHandler.db.GetAllTwitchNotifies(tUser.ID)
		if err != nil {
			return err
		}

		var notify *twitchnotify.DBEntry
		for _, not := range nots {
			if not.GuildID == args.Guild.ID {
				notify = not
			}
		}

		if notify == nil {
			return util.SendEmbedError(args.Session, args.Channel.ID,
				"Twitch user was nto set to be monitored on this guild.").
				DeleteAfter(8 * time.Second).Error()
		}

		err = args.CmdHandler.db.DeleteTwitchNotify(notify.TwitchUserID, notify.GuildID)
		if err != nil {
			return err
		}

		return util.SendEmbed(args.Session, args.Channel.ID,
			"Twitch user removed from monitor.", "", 0).
			DeleteAfter(8 * time.Second).Error()
	}

	accMsg := acceptmsg.AcceptMessage{
		Session:        args.Session,
		UserID:         args.User.ID,
		DeleteMsgAfter: true,
		Embed: &discordgo.MessageEmbed{
			Color:       static.ColorEmbedDefault,
			Description: fmt.Sprintf("Do you want to get notifications in this channel when **%s** goes online on Twitch?", tUser.DisplayName),
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: tUser.AviURL,
			},
		},
		AcceptFunc: func(m *discordgo.Message) {
			err = tnw.AddUser(tUser)
			if err != nil {
				util.SendEmbedError(args.Session, args.Channel.ID,
					"Maximum count of registered Twitch accounts has been reached.").
					DeleteAfter(8 * time.Second)
				return
			}

			err = args.CmdHandler.db.SetTwitchNotify(&twitchnotify.DBEntry{
				ChannelID:    args.Channel.ID,
				GuildID:      args.Guild.ID,
				TwitchUserID: tUser.ID,
			})
			if err != nil {
				util.SendEmbedError(args.Session, args.Channel.ID,
					"Unexpected error while saving to database: ```\n"+err.Error()+"\n```").
					DeleteAfter(20 * time.Second)
				return
			}

			util.SendEmbed(args.Session, args.Channel.ID,
				fmt.Sprintf("You will now get notifications in this channel when **%s** goes online on Twitch.", tUser.DisplayName), "", static.ColorEmbedUpdated).
				DeleteAfter(8 * time.Second)
		},
		DeclineFunc: func(m *discordgo.Message) {
			util.SendEmbedError(args.Session, args.Channel.ID, "Canceled.").
				DeleteAfter(8 * time.Second)
		},
	}
	accMsg.Send(args.Channel.ID)

	return nil
}
