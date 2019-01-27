package commands

import (
	"fmt"
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
	return ""
}

func (c *CmdTwitchNotify) GetGroup() string {
	return GroupEtc
}

func (c *CmdTwitchNotify) GetPermission() int {
	return c.PermLvl
}

func (c *CmdTwitchNotify) SetPermission(permLvl int) {
	c.PermLvl = permLvl
}

func (c *CmdTwitchNotify) Exec(args *CommandArgs) error {
	tnw := args.CmdHandler.tnw

	if len(args.Args) > 0 {
		tUser, err := tnw.GetUser(args.Args[0], core.TwitchNotifyIdentLogin)
		if err != nil && err.Error() == "not found" {
			msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
				"Twitch user could not be found.")
			util.DeleteMessageLater(args.Session, msg, 8*time.Second)
			return err
		} else if err != nil {
			return err
		}

		// _, err := args.CmdHandler.db.GetTwitchNotify(tUser.ID, args.Guild.ID)

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
				err = args.CmdHandler.db.SetTwitchNotify(&core.TwitchNotifyDBEntry{})
			},
			DeclineFunc: func(m *discordgo.Message) {
				msg, _ := util.SendEmbedError(args.Session, args.Channel.ID, "Canceled.")
				util.DeleteMessageLater(args.Session, msg, 5*time.Second)
			},
		}
		accMsg.Send(args.Channel.ID)
	}
	return nil
}
