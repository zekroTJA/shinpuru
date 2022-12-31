package slashcommands

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/services/codeexec"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/jdoodle"
	"github.com/zekrotja/dgrs"
	"github.com/zekrotja/ken"
)

const (
	apiIDLen  = 32
	apiKeyLen = 64
)

type Exec struct {
	ken.EphemeralCommand
}

var (
	_ ken.SlashCommand        = (*Exec)(nil)
	_ permissions.PermCommand = (*Exec)(nil)
)

func (c *Exec) Name() string {
	return "exec"
}

func (c *Exec) Description() string {
	return "Setup code execution of code embeds."
}

func (c *Exec) Version() string {
	return "1.1.0"
}

func (c *Exec) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *Exec) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "setup",
			Description: "Setup code execution.",
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "reset",
			Description: "Disable code execution and remove stored credentials.",
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "enable",
			Description: "Set the enabled state true or false.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "enabled",
					Description: "The enabled state of the code execution.",
					Required:    true,
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "check",
			Description: "Show the status of the current code execution setup.",
		},
	}
}

func (c *Exec) Domain() string {
	return "sp.guild.config.exec"
}

func (c *Exec) SubDomains() []permissions.SubPermission {
	return nil
}

func (c *Exec) Run(ctx ken.Context) (err error) {
	if err = ctx.Defer(); err != nil {
		return
	}

	var ranEnable bool
	err = ctx.HandleSubCommands(ken.SubCommandHandler{"enable", func(ctx ken.SubCommandContext) error {
		ranEnable = true
		return c.enable(ctx)
	}})
	if err != nil {
		return
	}
	if ranEnable {
		return
	}

	execFact := ctx.Get(static.DiCodeExecFactory).(codeexec.Factory)
	if execFact.Name() == "ranna" {
		return ctx.FollowUpEmbed(&discordgo.MessageEmbed{
			Description: "Code execution is supplied by [ranna](https://github.com/ranna-go) in this instance, so " +
				"nothing is required to be set up. :wink:",
		}).Send().Error
	}

	err = ctx.HandleSubCommands(
		ken.SubCommandHandler{"setup", c.setup},
		ken.SubCommandHandler{"reset", c.reset},
		ken.SubCommandHandler{"check", c.check},
	)

	return
}

func (c *Exec) setup(ctx ken.SubCommandContext) (err error) {
	dmChan, err := ctx.GetSession().UserChannelCreate(ctx.User().ID)
	if err != nil {
		return err
	}

	err = util.SendEmbed(ctx.GetSession(), dmChan.ID,
		"We need a [jdoodle API](https://www.jdoodle.com/compiler-api) client ID and secret to enable code execution on this guild. These values will be \n"+
			"saved as clear text in our database to pass it to the API, so please, be careful which data you want to use, also, if we secure our \n"+
			"database as best as possible, we do not guarantee the safety of your data.\n\nPlease enter first your API **client ID** or enter `cancel` to return:", "", 0).
		Error()
	if err != nil {
		if strings.Contains(err.Error(), "Cannot send messages to this user") {
			err = ctx.FollowUpError("In order to setup [jsdoodle's](https://www.jdoodle.com) API, we need to get your jsdoodle API client ID and secret. "+
				"Because of security, we don't want that you send your credentials into a guilds chat, that would be done via DM.\n"+
				"So, please enable DM's for this guild to proceed.", "").Send().Error
		}
		return
	}

	ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: "Because you need to enter credentials, the setup is done in DM. " +
			"Please take a look into your DMs. 😉",
	}).Send()

	var removeHandler func()
	var state int
	var clientId, clientSecret string
	removeHandler = ctx.GetSession().AddHandler(func(s *discordgo.Session, e *discordgo.MessageCreate) {
		st := ctx.Get(static.DiState).(*dgrs.State)
		self, err := st.SelfUser()
		if err != nil {
			return
		}

		if e.ChannelID != dmChan.ID || e.Author.ID == self.ID {
			return
		}

		if strings.ToLower(e.Content) == "cancel" {
			util.SendEmbedError(s, dmChan.ID, "Setup canceled.")
		} else {
			switch state {
			case 0:
				clientId = e.Content
				if len(clientId) < apiIDLen {
					util.SendEmbedError(ctx.GetSession(), dmChan.ID,
						"Invalid API clientID, please enter again or enter `cancel` to exit.")
					return
				}
				state++
				util.SendEmbed(ctx.GetSession(), dmChan.ID, "Okay, now, please enter your API **secret** or enter `cancel` to exit:", "", 0)
				return
			case 1:
				clientSecret = e.Content
				if len(clientSecret) < apiKeyLen {
					util.SendEmbedError(ctx.GetSession(), dmChan.ID,
						"Invalid API secret, please enter again or enter `cancel` to exit.")
					return
				}
			}

			_, err := jdoodle.NewWrapper(clientId, clientSecret).CreditsSpent()
			if err != nil {
				util.SendEmbedError(ctx.GetSession(), dmChan.ID,
					"Sorry, but it seems like your entered credentials are not correct. Please try again entering your **clientID** or exit with `cancel`:")
				state = 0
				return
			}

			db, _ := ctx.Get(static.DiDatabase).(database.Database)
			err = db.SetGuildJdoodleKey(ctx.GetEvent().GuildID, clientId+"#"+clientSecret)
			if err != nil {
				util.SendEmbedError(ctx.GetSession(), dmChan.ID,
					"An unexpected error occured while saving the key. Please contact the host of this bot about this: ```\n"+err.Error()+"\n```")
				return
			}

			err = db.SetGuildCodeExecEnabled(ctx.GetEvent().GuildID, true)
			if err != nil {
				util.SendEmbedError(ctx.GetSession(), dmChan.ID,
					"An unexpected error occured while saving enabled state. Please contact the host of this bot about this: ```\n"+err.Error()+"\n```")
				return
			}

			util.SendEmbed(s, dmChan.ID, "API key set and system is enabled. :ok_hand:", "", static.ColorEmbedGreen)
		}

		if removeHandler != nil {
			removeHandler()
		}
	})

	return nil
}

func (c *Exec) reset(ctx ken.SubCommandContext) (err error) {
	db, _ := ctx.Get(static.DiDatabase).(database.Database)
	err = db.SetGuildJdoodleKey(ctx.GetEvent().GuildID, "")
	if err != nil {
		return err
	}

	err = db.SetGuildCodeExecEnabled(ctx.GetEvent().GuildID, false)
	if err != nil {
		return err
	}

	return ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: "API key was deleted from database and system was disabled.",
	}).Send().Error
}

func (c *Exec) enable(ctx ken.SubCommandContext) (err error) {
	enabled := ctx.Options().GetByName("enabled").BoolValue()

	db, _ := ctx.Get(static.DiDatabase).(database.Database)
	err = db.SetGuildCodeExecEnabled(ctx.GetEvent().GuildID, enabled)
	if err != nil {
		return err
	}

	stateStr := "disabled"
	if enabled {
		stateStr = "enabled"
	}

	return ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: fmt.Sprintf("Code execution has been **%s** on this guild.", stateStr),
	}).Send().Error
}

func (c *Exec) check(ctx ken.SubCommandContext) (err error) {
	db, _ := ctx.Get(static.DiDatabase).(database.Database)
	key, err := db.GetGuildJdoodleKey(ctx.GetEvent().GuildID)
	if database.IsErrDatabaseNotFound(err) {
		return ctx.FollowUpError(
			"Code execution is not set up on this guild. Use `exec setup` to set up code execution.", "").
			Send().
			Error
	}
	if err != nil {
		return err
	}

	split := strings.Split(key, "#")
	if len(split) < 2 {
		return errors.New("invalid jdoodle credentials")
	}
	clientId, clientSecret := split[0], split[1]

	res, err := jdoodle.NewWrapper(clientId, clientSecret).CreditsSpent()
	if err != nil {
		return err
	}

	return ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: fmt.Sprintf("Today, you've spent **%d** tokens on this guild.", res.Used),
		Title:       "JDoodle API Token Statistics",
	}).Send().Error
}
