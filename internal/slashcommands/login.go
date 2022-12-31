package slashcommands

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/services/config"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
	"github.com/zekroTJA/shinpuru/pkg/onetimeauth/v2"
	"github.com/zekroTJA/shinpuru/pkg/timerstack"
	"github.com/zekrotja/dgrs"
	"github.com/zekrotja/ken"
)

type Login struct {
	ken.EphemeralCommand
}

var (
	_ ken.SlashCommand        = (*Login)(nil)
	_ permissions.PermCommand = (*Login)(nil)
	_ ken.DmCapable           = (*Login)(nil)
)

func (c *Login) Name() string {
	return "login"
}

func (c *Login) Description() string {
	return "Receive a link via DM to log into the shinpuru web interface."
}

func (c *Login) Version() string {
	return "1.0.0"
}

func (c *Login) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *Login) IsDmCapable() bool {
	return true
}

func (c *Login) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{}
}

func (c *Login) Domain() string {
	return "sp.etc.login"
}

func (c *Login) SubDomains() []permissions.SubPermission {
	return nil
}

func (c *Login) Run(ctx ken.Context) (err error) {
	if err = ctx.Defer(); err != nil {
		return
	}

	st := ctx.Get(static.DiState).(*dgrs.State)
	cfg := ctx.Get(static.DiConfig).(config.Provider)
	ota := ctx.Get(static.DiOneTimeAuth).(onetimeauth.OneTimeAuth)
	db := ctx.Get(static.DiDatabase).(database.Database)

	ch, err := st.Channel(ctx.GetEvent().ChannelID)
	if err != nil {
		return
	}

	isDM := ch.Type == discordgo.ChannelTypeDM
	if !isDM {
		if ch, err = ctx.GetSession().UserChannelCreate(ctx.User().ID); err != nil {
			return
		}
	}

	enabled, err := db.GetUserOTAEnabled(ctx.User().ID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return
	}

	if !enabled {
		enableLink := fmt.Sprintf("%s/usersettings", cfg.Config().WebServer.PublicAddr)
		err = ctx.FollowUpError(
			"One Time Authorization is disabled by default. If you want to use it, you need "+
				"to enable it first in your [**user settings page**]("+enableLink+").", "").
			Send().Error
		return c.wrapDmError(ctx, err)
	}

	token, expires, err := ota.GetKey(ctx.User().ID, "login-via-dm")
	if err != nil {
		return
	}

	link := fmt.Sprintf("%s/api/ota?token=%s", cfg.Config().WebServer.PublicAddr, token)
	emb := &discordgo.MessageEmbed{
		Color: static.ColorEmbedDefault,
		Description: "Click this [**this link**](" + link + ") and you will be automatically logged " +
			"in to the shinpuru web interface.\n\nThis link is only valid for **a short time** from now!\n\n" +
			"Expires: `" + expires.Format(time.RFC1123) + "`",
	}

	var fEdit func(emb *discordgo.MessageEmbed) error
	if isDM {
		fum := ctx.FollowUpEmbed(emb).Send()
		err = fum.Error
		fEdit = func(emb *discordgo.MessageEmbed) error {
			return fum.EditEmbed(emb)
		}
	} else {
		var msg *discordgo.Message
		msg, err = ctx.GetSession().ChannelMessageSendEmbed(ch.ID, emb)
		fEdit = func(emb *discordgo.MessageEmbed) error {
			_, e := ctx.GetSession().ChannelMessageEditEmbed(ch.ID, msg.ID, emb)
			return e
		}
		if err == nil {
			err = ctx.FollowUpEmbed(&discordgo.MessageEmbed{
				Description: "The login token has been sent you via DM.",
			}).Send().Error
		}
	}
	if err != nil {
		return c.wrapDmError(ctx, err)
	}

	timerstack.New().After(1*time.Minute, func() bool {
		emb := &discordgo.MessageEmbed{
			Color:       static.ColorEmbedGray,
			Description: "The login link has expired.",
		}
		fEdit(emb)
		return true
	}).RunBlocking()

	return
}

func (c *Login) wrapDmError(ctx ken.Context, err error) error {
	if discordutil.IsCanNotOpenDmToUserError(err) {
		return ctx.FollowUpError(
			"You need to enable DMs from users of this guild so that a secret authentication link "+
				"can be sent to you via DM.", "").Send().Error
	}
	return err
}
