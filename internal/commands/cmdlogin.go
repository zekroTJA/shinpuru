package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/core/config"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/onetimeauth"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
	"github.com/zekroTJA/shireikan"
)

type CmdLogin struct {
}

func (c *CmdLogin) GetInvokes() []string {
	return []string{"login", "weblogin", "token"}
}

func (c *CmdLogin) GetDescription() string {
	return "Get a link via DM to log into the shinpuru web interface."
}

func (c *CmdLogin) GetHelp() string {
	return "`login`"
}

func (c *CmdLogin) GetGroup() string {
	return shireikan.GroupEtc
}

func (c *CmdLogin) GetDomainName() string {
	return "sp.etc.login"
}

func (c *CmdLogin) GetSubPermissionRules() []shireikan.SubPermission {
	return nil
}

func (c *CmdLogin) IsExecutableInDMChannels() bool {
	return true
}

func (c *CmdLogin) Exec(ctx shireikan.Context) (err error) {
	var ch *discordgo.Channel

	if ctx.GetChannel().Type == discordgo.ChannelTypeGroupDM {
		ch = ctx.GetChannel()
	} else {
		if ch, err = ctx.GetSession().UserChannelCreate(ctx.GetUser().ID); err != nil {
			return
		}
	}

	cfg := ctx.GetObject("config").(*config.Config)
	ota := ctx.GetObject("onetimeauth").(*onetimeauth.OneTimeAuth)

	token, err := ota.GetKey(ctx.GetUser().ID)
	if err != nil {
		return
	}

	link := fmt.Sprintf("%s/ota?token=%s", cfg.WebServer.PublicAddr, token)
	emb := &discordgo.MessageEmbed{
		Color: static.ColorEmbedDefault,
		Description: "Click this [**this link**](" + link + ") and you will be automatically logged " +
			"in to the shinpuru web interface.\n\nThis link is only valid for **one minute** from now!",
	}

	_, err = ctx.GetSession().ChannelMessageSendEmbed(ch.ID, emb)
	if discordutil.IsCanNotOpenDmToUserError(err) {
		err = util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
			"You need to enable DMs from users of this guild so that a secret authentication link "+
				"can be sent to you via DM.").Error()
	}

	return err
}
