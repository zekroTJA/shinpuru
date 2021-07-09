package commands

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/dgrs"

	"github.com/zekroTJA/shinpuru/internal/services/codeexec"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/jdoodle"
	"github.com/zekroTJA/shireikan"
)

const (
	apiIDLen  = 32
	apiKeyLen = 64
)

type CmdExec struct {
}

func (c *CmdExec) GetInvokes() []string {
	return []string{"exec", "ex", "execute", "jdoodle"}
}

func (c *CmdExec) GetDescription() string {
	return "Setup code execution of code embeds."
}

func (c *CmdExec) GetHelp() string {
	return "`exec setup` - enter jdoodle setup\n" +
		"`exec reset` - disable and delete token from database\n" +
		"`exec check` - retrurns the number of tokens consumed this day\n"
}

func (c *CmdExec) GetGroup() string {
	return shireikan.GroupChat
}

func (c *CmdExec) GetDomainName() string {
	return "sp.chat.exec"
}

func (c *CmdExec) GetSubPermissionRules() []shireikan.SubPermission {
	return []shireikan.SubPermission{
		{
			Term:        "exec",
			Explicit:    false,
			Description: "Allows activating a code execution in chat via reaction",
		},
	}
}

func (c *CmdExec) IsExecutableInDMChannels() bool {
	return false
}

func (c *CmdExec) Exec(ctx shireikan.Context) error {
	execFact := ctx.GetObject(static.DiCodeExecFactory).(codeexec.Factory)
	if execFact.Name() == "ranna" {
		return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
			"Code execution is supplied by [ranna](https://github.com/ranna-go) in this instance, so "+
				"nothing is required to be set up. :wink:",
			"", 0).DeleteAfter(10 * time.Second).Error()
	}

	errHelpMsg := func(ctx shireikan.Context) error {
		return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
			"Invalid command arguments. Please use `help exec` to see how to use this command.").
			DeleteAfter(8 * time.Second).Error()
	}

	if len(ctx.GetArgs()) < 1 {
		return errHelpMsg(ctx)
	}

	switch strings.ToLower(ctx.GetArgs().Get(0).AsString()) {
	case "setup", "enable":
		return c.setup(ctx)
	case "reset", "remove":
		return c.reset(ctx)
	case "check", "stats":
		return c.check(ctx)
	default:
		return errHelpMsg(ctx)
	}
}

func (c *CmdExec) setup(ctx shireikan.Context) error {
	dmChan, err := ctx.GetSession().UserChannelCreate(ctx.GetUser().ID)
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
			err := util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
				"In order to setup [jsdoodle's](https://www.jdoodle.com) API, we need to get your jsdoodle API client ID and secret. "+
					"Because of security, we don't want that you send your credentials into a guilds chat, that would be done via DM.\n"+
					"So, please enable DM's for this guild to proceed.").
				DeleteAfter(15 * time.Second).Error()
			return err
		}
	}

	var removeHandler func()
	var state int
	var clientId, clientSecret string
	removeHandler = ctx.GetSession().AddHandler(func(s *discordgo.Session, e *discordgo.MessageCreate) {
		st := ctx.GetObject(static.DiState).(*dgrs.State)
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

			db, _ := ctx.GetObject(static.DiDatabase).(database.Database)
			err = db.SetGuildJdoodleKey(ctx.GetGuild().ID, clientId+"#"+clientSecret)
			if err != nil {
				util.SendEmbedError(ctx.GetSession(), dmChan.ID,
					"An unexpected error occured while saving the key. Please contact the host of this bot about this: ```\n"+err.Error()+"\n```")
			}

			util.SendEmbed(s, dmChan.ID, "API key set and system is enabled. :ok_hand:", "", static.ColorEmbedGreen)
		}

		if removeHandler != nil {
			removeHandler()
		}
	})

	return nil
}

func (c *CmdExec) reset(ctx shireikan.Context) error {
	db, _ := ctx.GetObject(static.DiDatabase).(database.Database)
	err := db.SetGuildJdoodleKey(ctx.GetGuild().ID, "")
	if err != nil {
		return err
	}

	return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
		"API key was deleted from database and system was disabled.", "", static.ColorEmbedYellow).
		DeleteAfter(8 * time.Second).Error()
}

func (c *CmdExec) check(ctx shireikan.Context) error {
	db, _ := ctx.GetObject(static.DiDatabase).(database.Database)
	key, err := db.GetGuildJdoodleKey(ctx.GetGuild().ID)
	if database.IsErrDatabaseNotFound(err) {
		return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
			"Code execution is not set up on this guild. Use `exec setup` to set up code execution.").
			DeleteAfter(6 * time.Second).
			Error()
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

	return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
		fmt.Sprintf("Today, you've spent **%d** tokens on this guild.", res.Used),
		"JDoodle API Token Statistics", 0).
		DeleteAfter(15 * time.Second).
		Error()
}
