package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/zekroTJA/shinpuru/util"
)

type CmdVote struct {
}

func (c *CmdVote) GetInvokes() []string {
	return []string{"vote", "poll"}
}

func (c *CmdVote) GetDescription() string {
	return "create and manage polls"
}

func (c *CmdVote) GetHelp() string {
	return ""
}

func (c *CmdVote) GetGroup() string {
	return GroupChat
}

func (c *CmdVote) GetPermission() int {
	return 0
}

func (c *CmdVote) Exec(args *CommandArgs) error {

	if len(args.Args) > 0 && strings.ToLower(args.Args[0]) == "close" {
		var vote *util.Vote
		if len(args.Args) > 1 {
			vid := args.Args[1]
			for _, v := range util.VotesRunning {
				if v.GuildID == args.Guild.ID && v.ID == vid {
					vote = v
				}
			}
			if vote == nil {
				msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
					fmt.Sprintf("There is no open vote on this guild with the ID `%s`.", vid))
				util.DeleteMessageLater(args.Session, msg, 10*time.Second)
				return err
			}
		} else {
			vids := make([]string, 0)
			for _, v := range util.VotesRunning {
				if v.GuildID == args.Guild.ID && v.CreatorID == args.User.ID {
					vote = v
					vids = append(vids, v.ID)
				}
			}
			fmt.Println(vids)
			if len(vids) > 1 {
				msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
					fmt.Sprintf("You have open more votes than 1. Please select the ID of the vote to close it: ```\n%s\n```", strings.Join(vids, "\n")))
				util.DeleteMessageLater(args.Session, msg, 30*time.Second)
				return err
			} else if vote == nil {
				msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
					"You have no open votes on this guild. Please specify a specific vote ID to close another ones vote, if you have the permissions to do this.")
				util.DeleteMessageLater(args.Session, msg, 12*time.Second)
				return err
			}
		}
		permlvls, err := args.CmdHandler.db.GetGuildPermissions(args.Guild.ID)
		if err != nil {
			return err
		}
		permlvl, _ := permlvls[args.User.ID]
		fmt.Println(permlvls, permlvl)
		if vote.CreatorID != args.User.ID && permlvl < 6 && args.User.ID != args.Guild.OwnerID {
			msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
				"You do not have the permission to close another ones votes.")
			util.DeleteMessageLater(args.Session, msg, 6*time.Second)
			return err
		}
		return vote.Close(args.Session)
	}

	split := strings.Split(strings.Join(args.Args, " "), "|")
	if len(split) < 3 || len(split) > 11 {
		msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
			"Invalid arguments. Please use `help vote` go get help about how to use this command.")
		util.DeleteMessageLater(args.Session, msg, 10*time.Second)
		return err
	}
	for i, e := range split {
		if len(e) < 1 {
			msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
				"Description or possibilities can not be empty.")
			util.DeleteMessageLater(args.Session, msg, 10*time.Second)
			return err
		}
		split[i] = strings.Trim(e, " \t")
	}
	vote := &util.Vote{
		ID:            args.Message.ID,
		MsgID:         "",
		CreatorID:     args.User.ID,
		GuildID:       args.Guild.ID,
		ChannelID:     args.Channel.ID,
		Description:   split[0],
		Possibilities: split[1:],
		Ticks:         make([]*util.VoteTick, 0),
	}
	emb, err := vote.AsEmbed(args.Session, false)
	if err != nil {
		return err
	}
	msg, err := args.Session.ChannelMessageSendEmbed(args.Channel.ID, emb)
	if err != nil {
		return err
	}
	vote.MsgID = msg.ID
	err = vote.AddReactions(args.Session)
	if err != nil {
		return err
	}
	util.VotesRunning[vote.ID] = vote
	return nil
}
