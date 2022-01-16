package slashcommands

import (
	"fmt"
	"regexp"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/snowflake"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
	"github.com/zekrotja/dgrs"
	"github.com/zekrotja/ken"
)

var linkRx = regexp.MustCompile(`^(?:https?:\/\/)?(?:www\.)?discord(?:app)?\.com\/channels\/\d{14,22}\/(\d{14,22})\/(\d{14,22}).*$`)

type Quote struct{}

var (
	_ ken.SlashCommand        = (*Quote)(nil)
	_ permissions.PermCommand = (*Quote)(nil)
)

func (c *Quote) Name() string {
	return "quote"
}

func (c *Quote) Description() string {
	return "Quote a message from any chat."
}

func (c *Quote) Version() string {
	return "1.0.0"
}

func (c *Quote) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *Quote) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "id",
			Description: "The message ID or URL to be quoted.",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "comment",
			Description: "Add a comment directly to the quote.",
		},
	}
}

func (c *Quote) Domain() string {
	return "sp.chat.quote"
}

func (c *Quote) SubDomains() []permissions.SubPermission {
	return nil
}

func (c *Quote) Run(ctx *ken.Ctx) (err error) {
	if err = ctx.Defer(); err != nil {
		return
	}

	st := ctx.Get(static.DiState).(*dgrs.State)

	var ident, comment string

	if resolved := ctx.Event.ApplicationCommandData().Resolved; resolved != nil {
		for ident = range resolved.Messages {
			break
		}
	}

	if ident == "" {
		ident = ctx.Options().GetByName("id").StringValue()
		if commentV, ok := ctx.Options().GetByNameOptional("comment"); ok {
			comment = commentV.StringValue()
		}
	}

	var quoteMsg *discordgo.Message
	var fum *ken.FollowUpMessage
	linkMatches := linkRx.FindAllStringSubmatch(ident, 2)
	if len(linkMatches) > 0 {
		messageID := linkMatches[0][2]
		channelID := linkMatches[0][1]
		quoteMsg, err = st.Message(channelID, messageID)
		if err != nil {
			return ctx.FollowUpError("Message could not be found.", "").Error
		}
	} else {
		messageID := ident
		if _, err = snowflake.ParseString(messageID); err != nil {
			return ctx.FollowUpError("Invlid message ID or link.", "").Error
		}

		msgSearchEmb := &discordgo.MessageEmbed{
			Color:       static.ColorEmbedGray,
			Description: fmt.Sprintf(":hourglass_flowing_sand:  Searching for message in channel <#%s>...", ctx.Event.ChannelID),
		}

		fum = ctx.FollowUpEmbed(msgSearchEmb)
		if fum.Error != nil {
			return fum.Error
		}

		chans, err := st.Channels(ctx.Event.GuildID, true)
		if err != nil {
			return err
		}

		var textChans []*discordgo.Channel
		for _, ch := range chans {
			if ch.Type == discordgo.ChannelTypeGuildText {
				textChans = append(textChans, ch)
			}
		}

		loopLen := len(textChans)
		results := make(chan *discordgo.Message, loopLen)
		timeout := make(chan bool, 1)
		timeOuted := false
		quoteMsg, _ = st.Message(ctx.Event.ChannelID, messageID)
		if quoteMsg == nil {
			msgSearchEmb.Description = ":hourglass_flowing_sand:  Searching for message in other channels..."
			fum.EditEmbed(msgSearchEmb)
			for _, ch := range textChans {
				go func(c *discordgo.Channel) {
					quoteMsg, _ := st.Message(c.ID, messageID)
					results <- quoteMsg
				}(ch)
			}
			timeoutTimer := time.AfterFunc(10*time.Second, func() {
				timeout <- true
			})
			func() {
				i := 0
				for {
					select {
					case fmsg := <-results:
						i++
						if i >= loopLen {
							return
						}
						if fmsg != nil {
							quoteMsg = fmsg
							timeoutTimer.Stop()
							return
						}
					case <-timeout:
						timeOuted = true
						return
					}
				}
			}()
		}

		isErr := true
		if timeOuted {
			msgSearchEmb.Description = "Searching worker timeout."
			msgSearchEmb.Color = static.ColorEmbedError
		} else if quoteMsg == nil {
			msgSearchEmb.Description = "Could not find any message with this ID. :disappointed:"
			msgSearchEmb.Color = static.ColorEmbedError
		} else if len(quoteMsg.Content) < 1 && len(quoteMsg.Attachments) < 1 {
			msgSearchEmb.Description = "Found messages content is empty. Maybe, it is an embed message which can not be quoted."
			msgSearchEmb.Color = static.ColorEmbedError
		} else {
			isErr = false
		}

		if isErr {
			return fum.EditEmbed(msgSearchEmb)
		}
	}

	ch, err := st.Channel(quoteMsg.ChannelID)
	if err != nil {
		return
	}

	// Sometimes, the Author can be nil on a message somehow
	// (see #342). Therefore, refrech message from API when
	// Author is nil. If Author is still nil, nah fuck it.
	if quoteMsg.Author == nil {
		quoteMsg, err = ctx.Session.ChannelMessage(quoteMsg.ChannelID, quoteMsg.ID)
		if err != nil {
			return err
		}
		if quoteMsg.Author == nil {
			quoteMsg.Author = &discordgo.User{
				ID:            "000000000000000000",
				Username:      "Discord doesn't want to give the author of this message :(",
				Avatar:        "",
				Discriminator: "0000",
			}
		}
		st.SetMessage(quoteMsg)
	}

	emb := &discordgo.MessageEmbed{
		Color: static.ColorEmbedDefault,
		Author: &discordgo.MessageEmbedAuthor{
			IconURL: quoteMsg.Author.AvatarURL(""),
			Name:    quoteMsg.Author.String(),
		},
		Description: quoteMsg.Content +
			fmt.Sprintf("\n\n*[jump to message](%s)*", discordutil.GetMessageLink(quoteMsg, ctx.Event.GuildID)),
		Footer: &discordgo.MessageEmbedFooter{
			IconURL: ctx.User().AvatarURL("16"),
			Text:    fmt.Sprintf("#%s - quoted by: %s", ch.Name, ctx.User().String()),
		},
		Timestamp: string(quoteMsg.Timestamp.Format(time.RFC3339)),
	}

	if len(quoteMsg.Attachments) > 0 {
		att := quoteMsg.Attachments[0]
		emb.Image = &discordgo.MessageEmbedImage{
			URL:      att.URL,
			ProxyURL: att.ProxyURL,
			Height:   att.Height,
			Width:    att.Width,
		}
	}

	if fum == nil {
		err = ctx.FollowUp(true, &discordgo.WebhookParams{
			Content: comment,
			Embeds:  []*discordgo.MessageEmbed{emb},
		}).Error
	} else {
		err = fum.Edit(&discordgo.WebhookEdit{
			Content: comment,
			Embeds:  []*discordgo.MessageEmbed{emb},
		})
	}

	return err
}
