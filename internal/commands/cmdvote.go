package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/core/middleware"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/imgstore"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/internal/util/vote"
	"github.com/zekroTJA/shireikan"
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
		"`vote close (<VoteID>|all) (nochart|nc)` - close your last vote, a vote by ID or all your open votes"
}

func (c *CmdVote) GetGroup() string {
	return shireikan.GroupChat
}

func (c *CmdVote) GetDomainName() string {
	return "sp.chat.vote"
}

func (c *CmdVote) GetSubPermissionRules() []shireikan.SubPermission {
	return []shireikan.SubPermission{
		{
			Term:        "close",
			Explicit:    true,
			Description: "Allows closing votes also from other users",
		},
	}
}

func (c *CmdVote) IsExecutableInDMChannels() bool {
	return false
}

func (c *CmdVote) Exec(ctx shireikan.Context) error {
	db, _ := ctx.GetObject("dbtnw").(database.Database)

	if len(ctx.GetArgs()) > 0 {
		switch strings.ToLower(ctx.GetArgs().Get(0).AsString()) {

		case "close":
			return c.close(ctx)

		case "list":
			return listVotes(ctx)

		case "expire", "expires":
			if len(ctx.GetArgs()) < 2 {
				return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
					"Please cpecify a expire duration!").
					DeleteAfter(8 * time.Second).Error()
			}

			expireDuration, err := time.ParseDuration(ctx.GetArgs().Get(1).AsString())
			if err != nil {
				return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
					"Invalid duration format. Please take a look "+
						"[here](https://golang.org/pkg/time/#ParseDuration) how to format duration parameter.").
					DeleteAfter(8 * time.Second).Error()
			}

			var ivote *vote.Vote
			if len(ctx.GetArgs()) > 2 {
				vid := ctx.GetArgs().Get(2).AsString()
				for _, v := range vote.VotesRunning {
					if v.GuildID == ctx.GetGuild().ID && v.ID == vid {
						ivote = v
					}
				}
				if ivote == nil {
					return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
						fmt.Sprintf("There is no open vote on this guild with the ID `%s`.", vid)).
						DeleteAfter(8 * time.Second).Error()
				}
			} else {
				votes := make([]*vote.Vote, 0)
				for _, v := range vote.VotesRunning {
					if v.GuildID == ctx.GetGuild().ID && v.CreatorID == ctx.GetUser().ID {
						votes = append(votes, v)
					}
				}
				if len(votes) == 0 {
					return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
						"There is no open vote on this guild created by you.").
						DeleteAfter(8 * time.Second).Error()
				}

				ivote = votes[len(votes)-1]
			}

			ivote.SetExpire(ctx.GetSession(), expireDuration)
			if err = db.AddUpdateVote(ivote); err != nil {
				return err
			}

			return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
				fmt.Sprintf("Vote will expire at %s.", ivote.Expires.Format("01/02 15:04 MST")), "", static.ColorEmbedGreen).
				DeleteAfter(8 * time.Second).Error()
		}

	} else {
		return listVotes(ctx)
	}

	split := strings.Split(strings.Join(ctx.GetArgs(), " "), "|")
	if len(split) < 3 || len(split) > 11 {
		return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
			"Invalid arguments. Please use `help vote` go get help about how to use this command.").
			DeleteAfter(8 * time.Second).Error()
	}
	for i, e := range split {
		if len(e) < 1 {
			return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
				"Description or possibilities can not be empty.").
				DeleteAfter(8 * time.Second).Error()
		}
		split[i] = strings.Trim(e, " \t")
	}

	description, imgLink := imgstore.ExtractFromMessage(split[0], ctx.GetMessage().Attachments)

	ivote := &vote.Vote{
		ID:            ctx.GetMessage().ID,
		MsgID:         "",
		CreatorID:     ctx.GetUser().ID,
		GuildID:       ctx.GetGuild().ID,
		ChannelID:     ctx.GetChannel().ID,
		Description:   description,
		Possibilities: split[1:],
		ImageURL:      imgLink,
		Ticks:         make(map[string]*vote.Tick),
	}

	emb, err := ivote.AsEmbed(ctx.GetSession())
	if err != nil {
		return err
	}

	msg, err := ctx.GetSession().ChannelMessageSendEmbed(ctx.GetChannel().ID, emb)
	if err != nil {
		return err
	}

	ivote.MsgID = msg.ID
	err = ivote.AddReactions(ctx.GetSession())
	if err != nil {
		return err
	}

	err = db.AddUpdateVote(ivote)
	if err != nil {
		return err
	}

	vote.VotesRunning[ivote.ID] = ivote
	return nil
}

func (c *CmdVote) close(ctx shireikan.Context) error {
	db, _ := ctx.GetObject("db").(database.Database)

	args := ctx.GetArgs()

	state := vote.VoteStateClosed
	if len(args) > 1 {
		i := args.IndexOf("nc")
		if i == -1 {
			i = args.IndexOf("nochart")
		}
		if i > -1 {
			state = vote.VoteStateClosedNC
			args = args.Splice(i, 1)
		}
	}

	var ivote *vote.Vote
	if len(ctx.GetArgs()) > 1 {
		if strings.ToLower(ctx.GetArgs().Get(1).AsString()) == "all" {
			var i int
			for _, v := range vote.VotesRunning {
				if v.GuildID == ctx.GetGuild().ID && v.CreatorID == ctx.GetUser().ID {
					go func(vC *vote.Vote) {
						db.DeleteVote(vC.ID)
						vC.Close(ctx.GetSession(), state)
					}(v)
					i++
				}
			}
			return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID, fmt.Sprintf("Closed %d votes.", i), "", 0).
				DeleteAfter(8 * time.Second).Error()
		}
		vid := ctx.GetArgs().Get(1).AsString()
		for _, v := range vote.VotesRunning {
			if v.GuildID == ctx.GetGuild().ID && v.ID == vid {
				ivote = v
			}
		}
		if ivote == nil {
			return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
				fmt.Sprintf("There is no open vote on this guild with the ID `%s`.", vid)).
				DeleteAfter(8 * time.Second).Error()
		}
	} else {
		vids := make([]string, 0)
		for _, v := range vote.VotesRunning {
			if v.GuildID == ctx.GetGuild().ID && v.CreatorID == ctx.GetUser().ID {
				ivote = v
				vids = append(vids, v.ID)
			}
		}
		if len(vids) > 1 {
			emb := &discordgo.MessageEmbed{
				Description: "You have open more votes than 1. Please select the ID of the vote to close it:",
				Color:       static.ColorEmbedError,
				Fields:      make([]*discordgo.MessageEmbedField, 0),
			}
			for _, v := range vote.VotesRunning {
				if v.GuildID == ctx.GetGuild().ID && v.CreatorID == ctx.GetUser().ID {
					emb.Fields = append(emb.Fields, v.AsField())
				}
			}
			return util.SendEmbedRaw(ctx.GetSession(), ctx.GetChannel().ID, emb).
				DeleteAfter(30 * time.Second).
				Error()
		} else if ivote == nil {
			return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
				"You have no open votes on this guild. Please specify a specific vote ID to close another ones vote, if you have the permissions to do this.").
				DeleteAfter(12 * time.Second).Error()
		}
	}

	pmw, _ := ctx.GetObject("pmw").(*middleware.PermissionsMiddleware)
	ok, override, err := pmw.CheckPermissions(ctx.GetSession(), ctx.GetGuild().ID, ctx.GetUser().ID, "!"+c.GetDomainName()+".close")
	if ivote.CreatorID != ctx.GetUser().ID && !ok && !override {
		return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
			"You do not have the permission to close another ones votes.").
			DeleteAfter(8 * time.Second).Error()
	}

	err = db.DeleteVote(ivote.ID)
	if err != nil {
		return err
	}

	err = ivote.Close(ctx.GetSession(), state)
	return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
		"Vote closed.", "", static.ColorEmbedGreen).
		DeleteAfter(8 * time.Second).Error()
}

func listVotes(ctx shireikan.Context) error {
	emb := &discordgo.MessageEmbed{
		Description: "Your open votes on this guild:",
		Color:       static.ColorEmbedDefault,
		Fields:      make([]*discordgo.MessageEmbedField, 0),
	}
	for _, v := range vote.VotesRunning {
		if v.GuildID == ctx.GetGuild().ID && v.CreatorID == ctx.GetUser().ID {
			emb.Fields = append(emb.Fields, v.AsField())
		}
	}
	if len(emb.Fields) == 0 {
		emb.Description = "You don't have any open votes on this guild."
	}
	_, err := ctx.GetSession().ChannelMessageSendEmbed(ctx.GetChannel().ID, emb)
	return err
}
