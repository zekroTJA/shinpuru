package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/snowflake"

	"github.com/zekroTJA/shinpuru/internal/core"
	"github.com/zekroTJA/shinpuru/internal/util"
)

var reserved = []string{"create", "add", "edit", "delete", "remove", "rem", "raw"}

type CmdTag struct {
}

func (c *CmdTag) GetInvokes() []string {
	return []string{"tag", "t", "note"}
}

func (c *CmdTag) GetDescription() string {
	return "set texts as tags which can be fastly re-posted later"
}

func (c *CmdTag) GetHelp() string {
	return "`tag` - Display all created tags on the current guild\n" +
		"`tag create <identifier> <content>` - Create a tag\n" +
		"`tag edit <identifier|ID> <content>` - Edit a tag\n" +
		"`tag delete <identifier|ID>` - Delete a tag\n" +
		"`tag raw <identifier|ID>` - Display tags content as raw markdown\n" +
		"`tag <identifier|ID>` - Display tag"
}

func (c *CmdTag) GetGroup() string {
	return GroupChat
}

func (c *CmdTag) GetDomainName() string {
	return "sp.chat.tag"
}

func (c *CmdTag) Exec(args *CommandArgs) error {
	db := args.CmdHandler.db

	if len(args.Args) < 1 {
		tags, err := db.GetGuildTags(args.Guild.ID)
		if err != nil {
			return err
		}

		var resTxt string

		if len(tags) < 1 {
			resTxt = "*No tags defined.*"
		} else {
			tlist := make([]string, len(tags))
			for i, t := range tags {
				tlist[i] = t.AsEntry(args.Session)
			}
			resTxt = strings.Join(tlist, "\n")
		}

		_, err = util.SendEmbed(args.Session, args.Channel.ID,
			resTxt, "Tags", 0)
		return err
	}

	switch strings.ToLower(args.Args[0]) {
	case "create", "add":
		if err, ok := checkPermission(args, c.GetDomainName()+".create"); !ok || err != nil {
			return err
		}
		return c.addTag(args, db)
	case "edit":
		return c.editTag(args, db)
	case "delete", "remove", "rem":
		return c.deleteTag(args, db)
	case "raw":
		return c.getRawTag(args, db)
	default:
		return c.getTag(args, db)
	}
}

func (c *CmdTag) addTag(args *CommandArgs, db core.Database) error {
	if len(args.Args) < 3 {
		return printInvalidArguments(args)
	}

	ident := strings.ToLower(args.Args[1])

	for _, r := range reserved {
		if r == ident {
			msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
				"A tag sub command can not be used as tag identifier.")
			util.DeleteMessageLater(args.Session, msg, 6*time.Second)
			return err
		}
	}

	tag, err := db.GetTagByIdent(ident, args.Guild.ID)
	if tag != nil {
		msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
			fmt.Sprintf("The tag `%s` already exists. Use `tag edit %s` to edit this tag or use another tag name.",
				ident, ident))
		util.DeleteMessageLater(args.Session, msg, 6*time.Second)
		return err
	}
	if err != nil && err != core.ErrDatabaseNotFound {
		return err
	}

	now := time.Now()
	argsJoined := strings.Join(args.Args[:2], " ")
	contentOffset := strings.Index(args.Message.Content, argsJoined) + len(argsJoined) + 1
	content := args.Message.Content[contentOffset:]

	tag = &util.Tag{
		Content:   content,
		Created:   now,
		CreatorID: args.User.ID,
		GuildID:   args.Guild.ID,
		ID:        util.NodeTags.Generate(),
		Ident:     ident,
		LastEdit:  now,
	}

	if err = db.AddTag(tag); err != nil {
		return err
	}

	msg, err := util.SendEmbed(args.Session, args.Channel.ID,
		fmt.Sprintf("Tag `%s` was created with ID `%s`.", ident, tag.ID), "", util.ColorEmbedGreen)
	util.DeleteMessageLater(args.Session, msg, 8*time.Second)
	return err
}

func (c *CmdTag) editTag(args *CommandArgs, db core.Database) error {
	if len(args.Args) < 3 {
		return printInvalidArguments(args)
	}

	tag, err, ok := getTag(args.Args[1], args, db)
	if !ok || err != nil {
		return err
	}

	if tag.CreatorID != args.User.ID {
		return printNotPermitted(args, "edit")
	}

	argsJoined := strings.Join(args.Args[:2], " ")
	contentOffset := strings.Index(args.Message.Content, argsJoined) + len(argsJoined) + 1
	tag.Content = args.Message.Content[contentOffset:]
	tag.LastEdit = time.Now()

	if err = db.EditTag(tag); err != nil {
		return err
	}

	msg, err := util.SendEmbed(args.Session, args.Channel.ID,
		fmt.Sprintf("Tag `%s` (ID `%s`) was updated.", tag.Ident, tag.ID), "", util.ColorEmbedGreen)
	util.DeleteMessageLater(args.Session, msg, 8*time.Second)

	return err
}

func (c *CmdTag) deleteTag(args *CommandArgs, db core.Database) error {
	if len(args.Args) < 2 {
		return printInvalidArguments(args)
	}

	tag, err, ok := getTag(args.Args[1], args, db)
	if !ok || err != nil {
		return err
	}

	ok, err = args.CmdHandler.CheckPermissions(args.Session, args.Guild.ID, args.User.ID, c.GetDomainName()+".delete")
	if err != nil {
		return err
	}

	if tag.CreatorID != args.User.ID && !ok {
		return printNotPermitted(args, "delete")
	}

	if err = db.DeleteTag(tag.ID); err != nil {
		return err
	}

	msg, err := util.SendEmbed(args.Session, args.Channel.ID,
		"Tag was deleted.", "", util.ColorEmbedGreen)
	util.DeleteMessageLater(args.Session, msg, 6*time.Second)

	return err
}

func (c *CmdTag) getRawTag(args *CommandArgs, db core.Database) error {
	if len(args.Args) < 2 {
		return printInvalidArguments(args)
	}

	tag, err, ok := getTag(args.Args[1], args, db)
	if !ok || err != nil {
		return err
	}

	_, err = args.Session.ChannelMessageSend(args.Channel.ID, tag.RawContent())
	return err
}

func (c CmdTag) getTag(args *CommandArgs, db core.Database) error {
	tag, err, ok := getTag(args.Args[0], args, db)
	if !ok || err != nil {
		return err
	}

	_, err = args.Session.ChannelMessageSendEmbed(args.Channel.ID, tag.AsEmbed(args.Session))
	return err
}

func getTag(ident string, args *CommandArgs, db core.Database) (*util.Tag, error, bool) {
	tag, err := db.GetTagByIdent(strings.ToLower(ident), args.Guild.ID)
	if err != nil && err != core.ErrDatabaseNotFound {
		return nil, err, false
	}
	if tag != nil {
		return tag, nil, true
	}

	id, err := snowflake.ParseString(ident)
	if err != nil {
		return nil, printTagNotFound(args), false
	}

	tag, err = db.GetTagByID(id)
	if err != nil && err != core.ErrDatabaseNotFound {
		return nil, err, false
	}

	if tag == nil || tag.GuildID != args.Guild.ID {
		return nil, printTagNotFound(args), false
	}

	return tag, nil, true
}

func printInvalidArguments(args *CommandArgs) error {
	msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
		"Invalid arguments. Use `help tag` to ge thelp about how to use this command.")
	util.DeleteMessageLater(args.Session, msg, 8*time.Second)
	return err
}

func printTagNotFound(args *CommandArgs) error {
	msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
		"Could not find any tag by the given identifier.")
	util.DeleteMessageLater(args.Session, msg, 6*time.Second)
	return err
}

func printNotPermitted(args *CommandArgs, t string) error {
	msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
		fmt.Sprintf("You are not permitted to %s this tag.", t))
	util.DeleteMessageLater(args.Session, msg, 6*time.Second)
	return err
}

func checkPermission(args *CommandArgs, dn string) (error, bool) {
	ok, err := args.CmdHandler.CheckPermissions(args.Session, args.Guild.ID, args.User.ID, dn)
	if err != nil {
		return err, false
	}

	if !ok {
		msg, err := util.SendEmbedError(args.Session, args.Channel.ID,
			"You are not permitted to use this command.")
		util.DeleteMessageLater(args.Session, msg, 6*time.Second)
		return err, false
	}

	return nil, true
}
