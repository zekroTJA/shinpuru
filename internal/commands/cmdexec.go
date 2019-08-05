package commands

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/zekroTJA/shinpuru/internal/util"
)

const (
	apiIDLen  = 32
	apiKeyLen = 64

	apiUrl = "https://api.jdoodle.com/v1/credit-spent"
)

type CmdExec struct {
}

func (c *CmdExec) GetInvokes() []string {
	return []string{"exec", "ex", "execute", "jdoodle"}
}

func (c *CmdExec) GetDescription() string {
	return "setup code execution of code embeds"
}

func (c *CmdExec) GetHelp() string {
	return "`exec setup` - enter jdoodle setup\n" +
		"`exec reset` - disable and delete token from database\n"
}

func (c *CmdExec) GetGroup() string {
	return GroupChat
}

func (c *CmdExec) GetDomainName() string {
	return "sp.chat.exec"
}

func (c *CmdExec) Exec(args *CommandArgs) error {
	errHelpMsg := func(args *CommandArgs) error {
		msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
			"Invalid command arguments. Please use `help exec` to see how to use this command.")
		util.DeleteMessageLater(args.Session, msg, 8*time.Second)
		return err
	}

	if len(args.Args) < 1 {
		return errHelpMsg(args)
	}

	switch strings.ToLower(args.Args[0]) {
	case "setup":
		return c.setup(args)
	case "reset":
		return c.reset(args)
	default:
		return errHelpMsg(args)
	}
}

func (c *CmdExec) setup(args *CommandArgs) error {
	dmChan, err := args.Session.UserChannelCreate(args.User.ID)
	if err != nil {
		return err
	}

	_, err = util.SendEmbed(args.Session, dmChan.ID,
		"We need an [jsdoodle API](https://www.jdoodle.com/compiler-api) client ID and secret to enable code execution on this guild. These values will be \n"+
			"saved as clear text in our database to pass it to the API, so please, be careful which data you want to use, also, if we secure our \n"+
			"database as best as possible, we do not guarantee the safety of your data.\n\nPlease enter first your API **client ID** or enter `cancel` to return:", "", 0)
	if err != nil {
		if strings.Contains(err.Error(), "Cannot send messages to this user") {
			msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
				"In order to setup [jsdoodle's](https://www.jdoodle.com) API, we need to get your jsdoodle API client ID and secret. "+
					"Because of security, we don't want that you send your credentials into a guilds chat, that would be done via DM.\n"+
					"So, please enable DM's for this guild to proceed.")
			util.DeleteMessageLater(args.Session, msg, 15*time.Second)
			return err
		}
	}

	var removeHandler func()
	var state int
	var clientID, token string
	removeHandler = args.Session.AddHandler(func(s *discordgo.Session, e *discordgo.MessageCreate) {
		if e.ChannelID != dmChan.ID || e.Author.ID == s.State.User.ID {
			return
		}

		if strings.ToLower(e.Content) == "cancel" {
			util.SendEmbedError(s, dmChan.ID, "Setup canceled.")
		} else {
			switch state {
			case 0:
				clientID = e.Content
				if len(clientID) < apiIDLen {
					util.SendEmbedError(args.Session, dmChan.ID,
						"Invalid API clientID, please enter again or enter `cancel` to exit.")
					return
				}
				state++
				util.SendEmbed(args.Session, dmChan.ID, "Okay, now, please enter your API **secret** or enter `cancel` to exit:", "", 0)
				return
			case 1:
				token = e.Content
				if len(token) < apiKeyLen {
					util.SendEmbedError(args.Session, dmChan.ID,
						"Invalid API secret, please enter again or enter `cancel` to exit.")
					return
				}
			}

			bodyBuffer, _ := json.Marshal(map[string]string{
				"clientId":     clientID,
				"clientSecret": token,
			})
			res, err := http.Post(apiUrl, "application/json", bytes.NewReader(bodyBuffer))
			if err != nil {
				util.SendEmbedError(args.Session, dmChan.ID,
					"An unexpected error occured while saving the key. Please contact the host of this bot about this: ```\n"+err.Error()+"\n```")
				return
			}

			if res.StatusCode != 200 {
				util.SendEmbedError(args.Session, dmChan.ID,
					"Sorry, but it seems like your entered credentials are not correct. Please try again entering your **clientID** or exit with `cancel`:")
				state = 0
				return
			}

			err = args.CmdHandler.db.SetGuildJdoodleKey(args.Guild.ID, clientID+"#"+token)
			if err != nil {
				util.SendEmbedError(args.Session, dmChan.ID,
					"An unexpected error occured while saving the key. Please contact the host of this bot about this: ```\n"+err.Error()+"\n```")
			}

			util.SendEmbed(s, dmChan.ID, "API key set and system is enabled. :ok_hand:", "", util.ColorEmbedGreen)
		}

		if removeHandler != nil {
			removeHandler()
		}
	})

	return nil
}

func (c *CmdExec) reset(args *CommandArgs) error {
	err := args.CmdHandler.db.SetGuildJdoodleKey(args.Guild.ID, "")
	if err != nil {
		return err
	}

	msg, err := util.SendEmbed(args.Session, args.Channel.ID,
		"API key was deleted from database and system was disabled.", "", util.ColorEmbedYellow)
	util.DeleteMessageLater(args.Session, msg, 8*time.Second)
	return err
}
