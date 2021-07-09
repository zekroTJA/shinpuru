package commands

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/fetch"
	"github.com/zekroTJA/shireikan"
)

type CmdClear struct {
}

func (c *CmdClear) GetInvokes() []string {
	return []string{"clear", "c", "purge"}
}

func (c *CmdClear) GetDescription() string {
	return "Clear messages in a channel."
}

func (c *CmdClear) GetHelp() string {
	return "`clear` - delete last message\n" +
		"`clear <n>` - clear an ammount of messages\n" +
		"`clear <n> <userResolvable>` - clear an ammount of messages by a specific user"
}

func (c *CmdClear) GetGroup() string {
	return shireikan.GroupModeration
}

func (c *CmdClear) GetDomainName() string {
	return "sp.guild.mod.clear"
}

func (c *CmdClear) GetSubPermissionRules() []shireikan.SubPermission {
	return nil
}

func (c *CmdClear) IsExecutableInDMChannels() bool {
	return true
}

func (c *CmdClear) Exec(ctx shireikan.Context) error {
	var msglist []*discordgo.Message
	var err error

	if len(ctx.GetArgs()) == 0 {
		msglist, err = ctx.GetSession().ChannelMessages(ctx.GetChannel().ID, 2, "", "", "")
	} else {
		var memb *discordgo.Member
		n, err := ctx.GetArgs().Get(0).AsInt()
		if err != nil {
			return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
				"Sorry, but the member can not be found on this guild. :cry:").
				DeleteAfter(8 * time.Second).Error()
		} else if n < 0 || n > 99 {
			return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
				"Number of messages is invald and must be between *(including)* 0 and 100.").
				DeleteAfter(8 * time.Second).Error()
		}

		// Account for command message itself
		n++

		if len(ctx.GetArgs()) >= 2 {
			memb, err = fetch.FetchMember(ctx.GetSession(), ctx.GetGuild().ID, ctx.GetArgs().Get(1).AsString())
			if err != nil {
				return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
					"Sorry, but the member can not be found on this guild. :cry:").
					DeleteAfter(8 * time.Second).Error()
			}
		}
		msglistUnfiltered, err := ctx.GetSession().ChannelMessages(ctx.GetChannel().ID, n, "", "", "")
		if err != nil {
			return err
		}

		if memb != nil {
			for _, m := range msglistUnfiltered {
				if m.Author.ID == memb.User.ID {
					msglist = append(msglist, m)
				}
			}
		} else {
			msglist = msglistUnfiltered
		}
	}

	if err != nil {
		return err
	}

	msgs := make([]string, len(msglist))
	for i, m := range msglist {
		msgs[i] = m.ID
	}

	err = ctx.GetSession().ChannelMessagesBulkDelete(ctx.GetChannel().ID, msgs)
	if err != nil {
		return err
	}

	multipleMsgs := ""
	if len(msgs) > 2 {
		multipleMsgs = "s"
	}

	return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
		fmt.Sprintf("Deleted %d message%s.", len(msgs)-1, multipleMsgs), "", static.ColorEmbedUpdated).
		DeleteAfter(6 * time.Second).Error()
}
