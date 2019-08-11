package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/util"
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
	return "`vote <description> | <possibility1> | <possibility2> (| <possibility3> ...)` - create vote\n" +
		"`vote list` - display currentltly running votes\n" +
		"`vote expire <duration> (<voteID>)` - set expire to last created (or specified) vote\n" +
		"`vote close (<VoteID>|all)` - close your last vote, a vote by ID or all your open votes `[sp.chat.vote.close]`"
}

func (c *CmdVote) GetGroup() string {
	return GroupChat
}

func (c *CmdVote) GetDomainName() string {
	return "sp.chat.vote"
}
func (c *CmdVote) Exec(args *CommandArgs) error {

	if len(args.Args) > 0 {
		switch strings.ToLower(args.Args[0]) {

		case "close":
			var vote *util.Vote
			if len(args.Args) > 1 {
				if strings.ToLower(args.Args[1]) == "all" {
					var i int
					for _, v := range util.VotesRunning {
						if v.GuildID == args.Guild.ID && v.CreatorID == args.User.ID {
							go func(vC *util.Vote) {
								args.CmdHandler.db.DeleteVote(vC.ID)
								vC.Close(args.Session, util.VoteStateClosed)
							}(v)
							i++
						}
					}
					msg, err := util.SendEmbed(args.Session, args.Channel.ID, fmt.Sprintf("Closed %d votes.", i), "", 0)
					util.DeleteMessageLater(args.Session, msg, 5*time.Second)
					return err
				}
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
				if len(vids) > 1 {
					emb := &discordgo.MessageEmbed{
						Description: "You have open more votes than 1. Please select the ID of the vote to close it:",
						Color:       util.ColorEmbedError,
						Fields:      make([]*discordgo.MessageEmbedField, 0),
					}
					for _, v := range util.VotesRunning {
						if v.GuildID == args.Guild.ID && v.CreatorID == args.User.ID {
							emb.Fields = append(emb.Fields, v.AsField())
						}
					}
					msg, err := args.Session.ChannelMessageSendEmbed(args.Channel.ID, emb)
					util.DeleteMessageLater(args.Session, msg, 30*time.Second)
					return err
				} else if vote == nil {
					msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
						"You have no open votes on this guild. Please specify a specific vote ID to close another ones vote, if you have the permissions to do this.")
					util.DeleteMessageLater(args.Session, msg, 12*time.Second)
					return err
				}
			}

			ok, err := args.CmdHandler.CheckPermissions(args.Session, args.Guild.ID, args.User.ID, c.GetDomainName()+".close")
			if vote.CreatorID != args.User.ID && !ok && args.User.ID != args.Guild.OwnerID {
				msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
					"You do not have the permission to close another ones votes.")
				util.DeleteMessageLater(args.Session, msg, 6*time.Second)
				return err
			}

			err = args.CmdHandler.db.DeleteVote(vote.ID)
			if err != nil {
				return err
			}

			err = vote.Close(args.Session, util.VoteStateClosed)
			msg, err := util.SendEmbed(args.Session, args.Channel.ID,
				"Vote closed.", "", util.ColorEmbedGreen)
			util.DeleteMessageLater(args.Session, msg, 6*time.Second)
			return err

		case "list":
			return listVotes(args)

		case "expire", "expires":
			if len(args.Args) < 2 {
				msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
					"Please cpecify a expire duration!")
				util.DeleteMessageLater(args.Session, msg, 10*time.Second)
				return err
			}

			expireDuration, err := time.ParseDuration(args.Args[1])
			if err != nil {
				msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
					"Invalid duration format. Please take a look "+
						"[here](https://golang.org/pkg/time/#ParseDuration) how to format duration parameter.")
				util.DeleteMessageLater(args.Session, msg, 10*time.Second)
				return err
			}

			var vote *util.Vote
			if len(args.Args) > 2 {
				vid := args.Args[2]
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
				votes := make([]*util.Vote, 0)
				for _, v := range util.VotesRunning {
					if v.GuildID == args.Guild.ID && v.CreatorID == args.User.ID {
						votes = append(votes, v)
					}
				}
				if len(votes) == 0 {
					msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
						"There is no open vote on this guild created by you.")
					util.DeleteMessageLater(args.Session, msg, 6*time.Second)
					return err
				}

				vote = votes[len(votes)-1]
			}

			vote.SetExpire(args.Session, expireDuration)
			if err = args.CmdHandler.db.AddUpdateVote(vote); err != nil {
				return err
			}

			msg, err := util.SendEmbed(args.Session, args.Channel.ID,
				fmt.Sprintf("Vote will expire at %s.", vote.Expires.Format("01/02 15:04 MST")), "", util.ColorEmbedGreen)
			util.DeleteMessageLater(args.Session, msg, 10*time.Second)
			return err
		}

	} else {
		return listVotes(args)
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

	description, imgLink := util.ExtractImageURLFromMessage(split[0], args.Message.Attachments)

	vote := &util.Vote{
		ID:            args.Message.ID,
		MsgID:         "",
		CreatorID:     args.User.ID,
		GuildID:       args.Guild.ID,
		ChannelID:     args.Channel.ID,
		Description:   description,
		Possibilities: split[1:],
		ImageURL:      imgLink,
		Ticks:         make([]*util.VoteTick, 0),
	}
	emb, err := vote.AsEmbed(args.Session)
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
	err = args.CmdHandler.db.AddUpdateVote(vote)
	if err != nil {
		return err
	}
	util.VotesRunning[vote.ID] = vote
	return nil
}

func listVotes(args *CommandArgs) error {
	emb := &discordgo.MessageEmbed{
		Description: "Your open votes on this guild:",
		Color:       util.ColorEmbedDefault,
		Fields:      make([]*discordgo.MessageEmbedField, 0),
	}
	for _, v := range util.VotesRunning {
		if v.GuildID == args.Guild.ID && v.CreatorID == args.User.ID {
			emb.Fields = append(emb.Fields, v.AsField())
		}
	}
	if len(emb.Fields) == 0 {
		emb.Description = "You don't have any open votes on this guild."
	}
	_, err := args.Session.ChannelMessageSendEmbed(args.Channel.ID, emb)
	return err
}
