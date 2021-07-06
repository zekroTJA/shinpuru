package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/dgrs"

	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
	"github.com/zekroTJA/shireikan"
)

type CmdQuote struct {
}

func (c *CmdQuote) GetInvokes() []string {
	return []string{"quote", "q"}
}

func (c *CmdQuote) GetDescription() string {
	return "Quote a message from any chat."
}

func (c *CmdQuote) GetHelp() string {
	return "`quote <msgID/msgURL> (<comment>)`"
}

func (c *CmdQuote) GetGroup() string {
	return shireikan.GroupChat
}

func (c *CmdQuote) GetDomainName() string {
	return "sp.chat.quote"
}

func (c *CmdQuote) GetSubPermissionRules() []shireikan.SubPermission {
	return nil
}

func (c *CmdQuote) IsExecutableInDMChannels() bool {
	return false
}

func (c *CmdQuote) Exec(ctx shireikan.Context) error {
	args := ctx.GetArgs()

	if len(args) < 1 {
		return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
			"Please enter a message ID or URL which should be quoted.").
			DeleteAfter(8 * time.Second).Error()
	}

	if strings.HasPrefix(ctx.GetArgs().Get(0).AsString(), "https://discordapp.com/channels/") ||
		strings.HasPrefix(ctx.GetArgs().Get(0).AsString(), "https://discord.com/channels/") ||
		strings.HasPrefix(ctx.GetArgs().Get(0).AsString(), "https://canary.discordapp.com/channels/") ||
		strings.HasPrefix(ctx.GetArgs().Get(0).AsString(), "https://canary.discord.com/channels/") {

		urlSplit := strings.Split(ctx.GetArgs().Get(0).AsString(), "/")
		args[0] = urlSplit[len(urlSplit)-1]
	}

	comment := strings.Join(ctx.GetArgs()[1:], " ")

	msgSearchEmb := &discordgo.MessageEmbed{
		Color:       static.ColorEmbedGray,
		Description: fmt.Sprintf(":hourglass_flowing_sand:  Searching for message in channel <#%s>...", ctx.GetChannel().ID),
	}

	msgSearch, err := ctx.GetSession().ChannelMessageSendEmbed(ctx.GetChannel().ID, msgSearchEmb)
	if err != nil {
		return err
	}

	st := ctx.GetObject(static.DiState).(*dgrs.State)

	chans, err := st.Channels(ctx.GetGuild().ID, true)
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
	quoteMsg, _ := st.Message(ctx.GetChannel().ID, ctx.GetArgs().Get(0).AsString())
	if quoteMsg == nil {
		msgSearchEmb.Description = ":hourglass_flowing_sand:  Searching for message in other channels..."
		ctx.GetSession().ChannelMessageEditEmbed(ctx.GetChannel().ID, msgSearch.ID, msgSearchEmb)
		for _, ch := range textChans {
			go func(c *discordgo.Channel) {
				quoteMsg, _ := st.Message(c.ID, ctx.GetArgs().Get(0).AsString())
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

	if timeOuted {
		msgSearchEmb.Description = "Searching worker timeout."
		msgSearchEmb.Color = static.ColorEmbedError
		_, err := ctx.GetSession().ChannelMessageEditEmbed(ctx.GetChannel().ID, msgSearch.ID, msgSearchEmb)
		discordutil.DeleteMessageLater(ctx.GetSession(), msgSearch, 5*time.Second)
		return err
	}

	if quoteMsg == nil {
		msgSearchEmb.Description = "Could not find any message with this ID. :disappointed:"
		msgSearchEmb.Color = static.ColorEmbedError
		_, err := ctx.GetSession().ChannelMessageEditEmbed(ctx.GetChannel().ID, msgSearch.ID, msgSearchEmb)
		discordutil.DeleteMessageLater(ctx.GetSession(), msgSearch, 5*time.Second)
		return err
	}

	if len(quoteMsg.Content) < 1 && len(quoteMsg.Attachments) < 1 {
		msgSearchEmb.Description = "Found messages content is empty. Maybe, it is an embed message itself, which can not be quoted."
		msgSearchEmb.Color = static.ColorEmbedError
		_, err := ctx.GetSession().ChannelMessageEditEmbed(ctx.GetChannel().ID, msgSearch.ID, msgSearchEmb)
		discordutil.DeleteMessageLater(ctx.GetSession(), msgSearch, 8*time.Second)
		return err
	}

	quoteMsgChannel, err := st.Channel(quoteMsg.ChannelID)
	if err != nil {
		return err
	}

	msgSearchEmb = &discordgo.MessageEmbed{
		Color: static.ColorEmbedDefault,
		Author: &discordgo.MessageEmbedAuthor{
			IconURL: quoteMsg.Author.AvatarURL(""),
			Name:    quoteMsg.Author.Username + "#" + quoteMsg.Author.Discriminator,
		},
		Description: quoteMsg.Content +
			fmt.Sprintf("\n\n*[jump to message](%s)*", discordutil.GetMessageLink(quoteMsg, ctx.GetGuild().ID)),
		Footer: &discordgo.MessageEmbedFooter{
			IconURL: ctx.GetUser().AvatarURL("16"),
			Text:    fmt.Sprintf("#%s - quoted by: %s#%s", quoteMsgChannel.Name, ctx.GetUser().Username, ctx.GetUser().Discriminator),
		},
		Timestamp: string(quoteMsg.Timestamp),
	}

	if len(quoteMsg.Attachments) > 0 {
		att := quoteMsg.Attachments[0]
		msgSearchEmb.Image = &discordgo.MessageEmbedImage{
			URL:      att.URL,
			ProxyURL: att.ProxyURL,
			Height:   att.Height,
			Width:    att.Width,
		}
	}

	if comment != "" {
		ctx.GetSession().ChannelMessageEdit(ctx.GetChannel().ID, msgSearch.ID,
			fmt.Sprintf("**%s:**\n%s\n", ctx.GetUser().String(), comment))
	}

	ctx.GetSession().ChannelMessageEditEmbed(ctx.GetChannel().ID, msgSearch.ID, msgSearchEmb)
	return nil
}
