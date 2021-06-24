package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/zekrotja/discordgo"

	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/acceptmsg"
	"github.com/zekroTJA/shinpuru/pkg/twitchnotify"
	"github.com/zekroTJA/shireikan"
)

type CmdTwitchNotify struct {
}

func (c *CmdTwitchNotify) GetInvokes() []string {
	return []string{"twitch", "tn", "twitchnotify"}
}

func (c *CmdTwitchNotify) GetDescription() string {
	return "Get notifications in channels when someone goes live on Twitch."
}

func (c *CmdTwitchNotify) GetHelp() string {
	return "`twitch` - list all currently monitored twitch channels\n" +
		"`twitch <twitchUsername>` - get notified in the current channel when the streamer goes online\n" +
		"`twitch remove <twitchUsername>` - remove monitor"
}

func (c *CmdTwitchNotify) GetGroup() string {
	return shireikan.GroupChat
}

func (c *CmdTwitchNotify) GetDomainName() string {
	return "sp.chat.twitch"
}

func (c *CmdTwitchNotify) GetSubPermissionRules() []shireikan.SubPermission {
	return nil
}

func (c *CmdTwitchNotify) IsExecutableInDMChannels() bool {
	return false
}

func (c *CmdTwitchNotify) Exec(ctx shireikan.Context) error {
	tnw, _ := ctx.GetObject(static.DiTwitchNotifyWorker).(*twitchnotify.NotifyWorker)
	db, _ := ctx.GetObject(static.DiDatabase).(database.Database)

	if tnw == nil {
		return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
			"This feature is disabled because no Twitch App ID was provided.").
			DeleteAfter(8 * time.Second).Error()
	}

	if len(ctx.GetArgs()) < 1 {
		nots, err := db.GetAllTwitchNotifies("")
		if err != nil {
			return err
		}

		notsStr := ""

		for _, not := range nots {
			if not.GuildID == ctx.GetGuild().ID {
				if tUser, err := tnw.GetUser(not.TwitchUserID, twitchnotify.IdentID); err == nil {
					notsStr += fmt.Sprintf(":white_small_square:  **%s** in <#%s>\n",
						tUser.DisplayName, not.ChannelID)
				}
			}
		}

		return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID, notsStr, "Currently monitored streamers", 0).
			Error()
	}

	tUser, err := tnw.GetUser(ctx.GetArgs().Get(len(ctx.GetArgs())-1).AsString(), twitchnotify.IdentLogin)
	if err != nil && err.Error() == "not found" {
		return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
			"Twitch user could not be found.").
			DeleteAfter(8 * time.Second).Error()
	} else if err != nil {
		return err
	}

	if len(ctx.GetArgs()) > 1 && strings.ToLower(ctx.GetArgs().Get(0).AsString()) == "remove" {
		nots, err := db.GetAllTwitchNotifies(tUser.ID)
		if err != nil {
			return err
		}

		var notify *twitchnotify.DBEntry
		for _, not := range nots {
			if not.GuildID == ctx.GetGuild().ID {
				notify = not
			}
		}

		if notify == nil {
			return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
				"Twitch user was nto set to be monitored on this guild.").
				DeleteAfter(8 * time.Second).Error()
		}

		err = db.DeleteTwitchNotify(notify.TwitchUserID, notify.GuildID)
		if err != nil {
			return err
		}

		return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
			"Twitch user removed from monitor.", "", 0).
			DeleteAfter(8 * time.Second).Error()
	}

	accMsg := acceptmsg.AcceptMessage{
		Session:        ctx.GetSession(),
		UserID:         ctx.GetUser().ID,
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
				util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
					"Maximum count of registered Twitch accounts has been reached.").
					DeleteAfter(8 * time.Second)
				return
			}

			err = db.SetTwitchNotify(&twitchnotify.DBEntry{
				ChannelID:    ctx.GetChannel().ID,
				GuildID:      ctx.GetGuild().ID,
				TwitchUserID: tUser.ID,
			})
			if err != nil {
				util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
					"Unexpected error while saving to database: ```\n"+err.Error()+"\n```").
					DeleteAfter(20 * time.Second)
				return
			}

			util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
				fmt.Sprintf("You will now get notifications in this channel when **%s** goes online on Twitch.", tUser.DisplayName), "", static.ColorEmbedUpdated).
				DeleteAfter(8 * time.Second)
		},
		DeclineFunc: func(m *discordgo.Message) {
			util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID, "Canceled.").
				DeleteAfter(8 * time.Second)
		},
	}
	accMsg.Send(ctx.GetChannel().ID)

	return nil
}
