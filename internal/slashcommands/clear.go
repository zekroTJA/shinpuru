package slashcommands

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/acceptmsg"
	"github.com/zekroTJA/shinpuru/pkg/fetch"
	"github.com/zekrotja/ken"
)

type Clear struct{}

var (
	_ ken.Command             = (*Clear)(nil)
	_ permissions.PermCommand = (*Clear)(nil)
)

func (c *Clear) Name() string {
	return "clear"
}

func (c *Clear) Description() string {
	return "Clear messages in a channel."
}

func (c *Clear) Version() string {
	return "1.0.0"
}

func (c *Clear) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *Clear) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "last",
			Description: "Clears the last message",
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "amount",
			Description: "Clear a specified amount of messages",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "amount",
					Description: "Amount of messages to clear",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "user",
					Description: "Clear messages send by this User",
					Required:    false,
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "selected",
			Description: "Removes either messages selected with ❌ emote by you or all messages below the 🔻 emote by you",
		},
	}
}

func (c *Clear) Domain() string {
	return "sp.guild.mod.clear"
}

func (c *Clear) SubDomains() []permissions.SubPermission {
	return nil
}

func (c *Clear) Run(ctx *ken.Ctx) (err error) {
	if err = ctx.Defer(); err != nil {
		return
	}

	err = ctx.HandleSubCommands(
		ken.SubCommandHandler{"last", c.last},
		ken.SubCommandHandler{"amount", c.amount},
		ken.SubCommandHandler{"selected", c.selected},
	)

	return
}

func (c *Clear) last(ctx *ken.SubCommandCtx) (err error) {
	msglist, err := ctx.Session.ChannelMessages(ctx.Event.ChannelID, 2, "", "", "")
	if err != nil {
		return err
	}
	return c.delete(ctx, msglist)
}

func (c *Clear) amount(ctx *ken.SubCommandCtx) (err error) {

	amount := ctx.Options()[0].IntValue()
	var user *discordgo.User
	if len(ctx.Options()) > 1 {
		user = ctx.Options()[1].UserValue(nil)
	}

	if amount < 1 || amount > 99 {
		return util.SendEmbedError(ctx.Session, ctx.Event.ChannelID,
			"Number of messages is invald and must be between *(including)* 1 and 100.").
			DeleteAfter(8 * time.Second).Error()
	}

	var member *discordgo.Member
	if user != nil {
		member, err = fetch.FetchMember(ctx.Session, ctx.Event.GuildID, user.ID)
		if err != nil {
			return util.SendEmbedError(ctx.Session, ctx.Event.ChannelID,
				"Sorry, but the member can not be found on this guild. :cry:").
				DeleteAfter(8 * time.Second).Error()
		}
	}
	msglistUnfiltered, err := ctx.Session.ChannelMessages(ctx.Event.ChannelID, int(amount), "", "", "")
	if err != nil {
		return err
	}

	var msglist []*discordgo.Message
	if member != nil {
		for _, m := range msglistUnfiltered {
			if m.Author.ID == member.User.ID {
				msglist = append(msglist, m)
			}
		}
	} else {
		msglist = msglistUnfiltered
	}

	return c.delete(ctx, msglist)
}

func (c *Clear) selected(ctx *ken.SubCommandCtx) (err error) {
	msgs, err := ctx.Session.ChannelMessages(ctx.Event.ChannelID, 100, "", "", "")
	if err != nil {
		return
	}

	var deleteAfterMsg *discordgo.Message
	var deleteAfterIdx int
	c.iterMsgsWithReactionFromUser(ctx.Session, msgs, "🔻", ctx.User().ID, func(m *discordgo.Message, i int) bool {
		deleteAfterMsg = m
		deleteAfterIdx = i
		return false
	})

	if deleteAfterMsg != nil {
		msgIds := make([]string, 0, deleteAfterIdx+1)
		for _, m := range msgs[0 : deleteAfterIdx+1] {
			msgIds = append(msgIds, m.ID)
		}

		amsg, err := acceptmsg.New().
			WithSession(ctx.Session).
			WithContent(
				fmt.Sprintf("Do you really want to delete all %d messages to message %s?", len(msgIds), deleteAfterMsg.ID)).
			LockOnUser(ctx.User().ID).
			DeleteAfterAnswer().
			DoOnAccept(func(m *discordgo.Message) (err error) {
				if err = ctx.Session.ChannelMessagesBulkDelete(ctx.Event.ChannelID, msgIds); err != nil {
					return
				}
				return util.SendEmbed(ctx.Session, ctx.Event.ChannelID,
					fmt.Sprintf("Deleted %d %s.", len(msgIds), util.Pluralize(len(msgIds), "message")), "", static.ColorEmbedUpdated).
					DeleteAfter(6 * time.Second).Error()
			}).
			AsFollowUp(ctx.Ctx)
		if err != nil {
			return err
		}
		return amsg.Error()
	}

	msgIds := make([]string, 0, len(msgs))
	c.iterMsgsWithReactionFromUser(ctx.Session, msgs, "❌", ctx.User().ID, func(m *discordgo.Message, i int) bool {
		msgIds = append(msgIds, m.ID)
		return true
	})

	if len(msgIds) > 0 {
		amsg, err := acceptmsg.New().
			WithSession(ctx.Session).
			WithContent(
				fmt.Sprintf("Do you really want to delete all %d selected messages?", len(msgIds))).
			LockOnUser(ctx.User().ID).
			DeleteAfterAnswer().
			DoOnAccept(func(m *discordgo.Message) (err error) {
				if err = ctx.Session.ChannelMessagesBulkDelete(ctx.Event.ChannelID, msgIds); err != nil {
					return
				}
				return util.SendEmbed(ctx.Session, ctx.Event.ChannelID,
					fmt.Sprintf("Deleted %d %s.", len(msgIds), util.Pluralize(len(msgIds), "message")), "", static.ColorEmbedUpdated).
					DeleteAfter(6 * time.Second).Error()
			}).
			AsFollowUp(ctx.Ctx)
		if err != nil {
			return err
		}
		return amsg.Error()
	}

	return util.SendEmbedError(ctx.Session, ctx.Event.ChannelID,
		"No message was either selected by you with the 🔻 emote nor was any with the ❌ emote.\n\n"+
			"**Explaination:**\n"+
			"You can either select single messages to be deleted with the ❌ emote or select a message with the 🔻 emote "+
			"and this message plus all messages sent after this message will be deleted after entering the `clear selected` command.").
		DeleteAfter(12 * time.Second).Error()
}

func (c *Clear) delete(ctx *ken.SubCommandCtx, msglist []*discordgo.Message) (err error) {
	if err != nil {
		return err
	}

	msgs := make([]string, len(msglist))
	for i, m := range msglist {
		msgs[i] = m.ID
	}

	err = ctx.Session.ChannelMessagesBulkDelete(ctx.Event.ChannelID, msgs)
	if err != nil {
		return err
	}

	return ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: fmt.Sprintf("Deleted %d %s.", len(msgs)-1, util.Pluralize(len(msgs)-1, "message")),
		Title:       "",
		Color:       static.ColorEmbedUpdated,
	}).Error
}

func (c *Clear) iterMsgsWithReactionFromUser(
	s *discordgo.Session,
	msgs []*discordgo.Message,
	name, userID string,
	action func(*discordgo.Message, int) bool,
) (err error) {
	for i, m := range msgs {
	reactionLoop:
		for _, r := range m.Reactions {
			if r.Emoji.Name == name {
				rUsers, err := s.MessageReactions(m.ChannelID, m.ID, name, 100, "", "")
				if err != nil {
					return err
				}
				for _, rUser := range rUsers {
					if rUser.ID == userID {
						if !action(m, i) {
							return nil
						}
						break reactionLoop
					}
				}
			}
		}
	}

	return
}
