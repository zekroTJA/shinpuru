package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
	"github.com/zekroTJA/shireikan"
)

type CmdLogin struct {
}

func (c *CmdLogin) GetInvokes() []string {
	return []string{"login", "weblogin"}
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

	emb := &discordgo.MessageEmbed{}

	_, err = ctx.GetSession().ChannelMessageSendEmbed(ch.ID, emb)
	if discordutil.IsCanNotOpenDmToUserError(err) {
		err = util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
			"You need to enable DMs from users of this guild so that a secret authentication link "+
				"can be sent to you via DM.").Error()
	}

	return err
}
