package commands

import (
	"fmt"
	"time"

	"github.com/zekroTJA/shinpuru/internal/config"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
	"github.com/zekroTJA/shinpuru/pkg/onetimeauth/v2"
	"github.com/zekroTJA/shinpuru/pkg/timerstack"
	"github.com/zekroTJA/shireikan"
	"github.com/zekrotja/discordgo"
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

	if ctx.GetChannel().Type == discordgo.ChannelTypeDM {
		ch = ctx.GetChannel()
	} else {
		if ch, err = ctx.GetSession().UserChannelCreate(ctx.GetUser().ID); err != nil {
			return
		}
	}

	cfg := ctx.GetObject(static.DiConfig).(*config.Config)
	ota := ctx.GetObject(static.DiOneTimeAuth).(onetimeauth.OneTimeAuth)
	db := ctx.GetObject(static.DiDatabase).(database.Database)

	enabled, err := db.GetUserOTAEnabled(ctx.GetUser().ID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return
	}

	if !enabled {
		enableLink := fmt.Sprintf("%s/usersettings", cfg.WebServer.PublicAddr)
		err = util.SendEmbedError(ctx.GetSession(), ch.ID,
			"One Time Authorization is disabled by default. If you want to use it, you need "+
				"to enable it first in your [**user settings page**]("+enableLink+").").Error()
		return c.wrapDmError(ctx, err)
	}

	token, expires, err := ota.GetKey(ctx.GetUser().ID)
	if err != nil {
		return
	}

	link := fmt.Sprintf("%s/api/ota?token=%s", cfg.WebServer.PublicAddr, token)
	emb := &discordgo.MessageEmbed{
		Color: static.ColorEmbedDefault,
		Description: "Click this [**this link**](" + link + ") and you will be automatically logged " +
			"in to the shinpuru web interface.\n\nThis link is only valid for **a short time** from now!\n\n" +
			"Expires: `" + expires.Format(time.RFC1123) + "`",
	}

	msg, err := ctx.GetSession().ChannelMessageSendEmbed(ch.ID, emb)
	if err != nil {
		return c.wrapDmError(ctx, err)
	}

	timerstack.New().After(1*time.Minute, func() bool {
		emb := &discordgo.MessageEmbed{
			Color:       static.ColorEmbedGray,
			Description: "The login link has expired.",
		}
		ctx.GetSession().ChannelMessageEditEmbed(ch.ID, msg.ID, emb)
		return true
	}).RunBlocking()

	return err
}

func (c *CmdLogin) wrapDmError(ctx shireikan.Context, err error) error {
	if discordutil.IsCanNotOpenDmToUserError(err) {
		return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
			"You need to enable DMs from users of this guild so that a secret authentication link "+
				"can be sent to you via DM.").Error()
	}
	return err
}
