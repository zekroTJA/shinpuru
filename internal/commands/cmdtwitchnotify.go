package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/zekroTJA/shinpuru/internal/core"
	"github.com/zekroTJA/shinpuru/internal/util"
)

type CmdTwitchNotify struct {
	PermLvl int
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

func (c *CmdTwitchNotify) GetPermission() int {
	return c.PermLvl
}

func (c *CmdTwitchNotify) SetPermission(permLvl int) {
	c.PermLvl = permLvl
}

func (c *CmdTwitchNotify) Exec(args *CommandArgs) error {
	tnw := args.CmdHandler.tnw

	if tnw == nil {
		msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
			"This feature is disabled because no Twitch App ID was provided.")
		util.DeleteMessageLater(args.Session, msg, 8*time.Second)
		return err
	}

	if len(args.Args) < 1 {
		nots, err := args.CmdHandler.db.GetAllTwitchNotifies("")
		if err != nil {
			return err
		}

		notsStr := ""

		for _, not := range nots {
			if not.GuildID == args.Guild.ID {
				if tUser, err := tnw.GetUser(not.TwitchUserID, core.TwitchNotifyIdentID); err == nil {
					notsStr += fmt.Sprintf(":white_small_square:  **%s** in <#%s>\n",
						tUser.DisplayName, not.ChannelID)
				}
			}
		}

		_, err = util.SendEmbed(args.Session, args.Channel.ID, notsStr, "Currently monitored streamers", 0)
		return err
	}

	tUser, err := tnw.GetUser(args.Args[len(args.Args)-1], core.TwitchNotifyIdentLogin)
	if err != nil && err.Error() == "not found" {
		msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
			"Twitch user could not be found.")
		util.DeleteMessageLater(args.Session, msg, 8*time.Second)
		return err
	} else if err != nil {
		return err
	}

	if len(args.Args) > 1 && strings.ToLower(args.Args[0]) == "remove" {
		nots, err := args.CmdHandler.db.GetAllTwitchNotifies(tUser.ID)
		if err != nil {
			return err
		}

		var notify *core.TwitchNotifyDBEntry
		for _, not := range nots {
			if not.GuildID == args.Guild.ID {
				notify = not
			}
		}

		if notify == nil {
			msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
				"Twitch user was nto set to be monitored on this guild.")
			util.DeleteMessageLater(args.Session, msg, 10*time.Second)
			return err
		}

		err = args.CmdHandler.db.DeleteTwitchNotify(notify.TwitchUserID, notify.GuildID)
		if err != nil {
			return err
		}

		msg, err := util.SendEmbed(args.Session, args.Channel.ID,
			"Twitch user removed from monitor.", "", 0)
		util.DeleteMessageLater(args.Session, msg, 8*time.Second)
		return err
	}

	accMsg := util.AcceptMessage{
		Session:        args.Session,
		UserID:         args.User.ID,
		DeleteMsgAfter: true,
		Embed: &discordgo.MessageEmbed{
			Color:       util.ColorEmbedDefault,
			Description: fmt.Sprintf("Do you want to get notifications in this channel when **%s** goes online on Twitch?", tUser.DisplayName),
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: tUser.AviURL,
			},
		},
		AcceptFunc: func(m *discordgo.Message) {
			err = tnw.AddUser(tUser)
			if err != nil {
				msg, _ := util.SendEmbedError(args.Session, args.Channel.ID,
					"Maximum count of registered Twitch accounts has been reached.")
				util.DeleteMessageLater(args.Session, msg, 10*time.Second)
				return
			}

			err = args.CmdHandler.db.SetTwitchNotify(&core.TwitchNotifyDBEntry{
				ChannelID:    args.Channel.ID,
				GuildID:      args.Guild.ID,
				TwitchUserID: tUser.ID,
			})
			if err != nil {
				msg, _ := util.SendEmbedError(args.Session, args.Channel.ID,
					"Unexpected error while saving to database: ```\n"+err.Error()+"\n```")
				util.DeleteMessageLater(args.Session, msg, 30*time.Second)
				return
			}

			msg, _ := util.SendEmbed(args.Session, args.Channel.ID,
				fmt.Sprintf("You will now get notifications in this channel when **%s** goes online on Twitch.", tUser.DisplayName), "", util.ColorEmbedUpdated)
			util.DeleteMessageLater(args.Session, msg, 8*time.Second)
		},
		DeclineFunc: func(m *discordgo.Message) {
			msg, _ := util.SendEmbedError(args.Session, args.Channel.ID, "Canceled.")
			util.DeleteMessageLater(args.Session, msg, 5*time.Second)
		},
	}
	accMsg.Send(args.Channel.ID)

	return nil
}
