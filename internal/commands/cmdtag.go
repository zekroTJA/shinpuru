package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/snowflake"

	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/snowflakenodes"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/internal/util/tag"
	"github.com/zekroTJA/shireikan"
)

var reserved = []string{"create", "add", "edit", "delete", "remove", "rem", "raw"}

type CmdTag struct {
}

func (c *CmdTag) GetInvokes() []string {
	return []string{"tag", "t", "note", "tags"}
}

func (c *CmdTag) GetDescription() string {
	return "Set texts as tags which can be fastly re-posted later."
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
	return shireikan.GroupChat
}

func (c *CmdTag) GetDomainName() string {
	return "sp.chat.tag"
}

func (c *CmdTag) GetSubPermissionRules() []shireikan.SubPermission {
	return []shireikan.SubPermission{
		{
			Term:        "create",
			Explicit:    true,
			Description: "Allows creating tags",
		},
		{
			Term:        "edit",
			Explicit:    true,
			Description: "Allows editing tags (of every user)",
		},
		{
			Term:        "delete",
			Explicit:    true,
			Description: "Allows deleting tags (of every user)",
		},
	}
}

func (c *CmdTag) IsExecutableInDMChannels() bool {
	return false
}

func (c *CmdTag) Exec(ctx shireikan.Context) error {
	db, _ := ctx.GetObject(static.DiDatabase).(database.Database)

	if len(ctx.GetArgs()) < 1 {
		tags, err := db.GetGuildTags(ctx.GetGuild().ID)
		if err != nil {
			return err
		}

		var resTxt string

		if len(tags) < 1 {
			resTxt = "*No tags defined.*"
		} else {
			tlist := make([]string, len(tags))
			for i, t := range tags {
				tlist[i] = t.AsEntry(ctx.GetSession())
			}
			resTxt = strings.Join(tlist, "\n")
		}

		return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
			resTxt, "Tags", 0).Error()
	}

	switch strings.ToLower(ctx.GetArgs().Get(0).AsString()) {
	case "create", "add":
		if err, ok := checkPermission(ctx, "!"+c.GetDomainName()+".create"); !ok || err != nil {
			return err
		}
		return c.addTag(ctx, db)
	case "edit":
		return c.editTag(ctx, db)
	case "delete", "remove", "rem":
		return c.deleteTag(ctx, db)
	case "raw":
		return c.getRawTag(ctx, db)
	default:
		return c.getTag(ctx, db)
	}
}

func (c *CmdTag) addTag(ctx shireikan.Context, db database.Database) error {
	if len(ctx.GetArgs()) < 3 {
		return printInvalidArguments(ctx)
	}

	ident := strings.ToLower(ctx.GetArgs().Get(1).AsString())

	for _, r := range reserved {
		if r == ident {
			return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
				"A tag sub command can not be used as tag identifier.").
				DeleteAfter(8 * time.Second).Error()
		}
	}

	itag, err := db.GetTagByIdent(ident, ctx.GetGuild().ID)
	if itag != nil {
		return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
			fmt.Sprintf("The tag `%s` already exists. Use `tag edit %s` to edit this tag or use another tag name.",
				ident, ident)).
			DeleteAfter(8 * time.Second).Error()
	}
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	now := time.Now()
	argsJoined := strings.Join(ctx.GetArgs()[:2], " ")
	contentOffset := strings.Index(ctx.GetMessage().Content, argsJoined) + len(argsJoined) + 1
	content := ctx.GetMessage().Content[contentOffset:]

	itag = &tag.Tag{
		Content:   content,
		Created:   now,
		CreatorID: ctx.GetUser().ID,
		GuildID:   ctx.GetGuild().ID,
		ID:        snowflakenodes.NodeTags.Generate(),
		Ident:     ident,
		LastEdit:  now,
	}

	if err = db.AddTag(itag); err != nil {
		return err
	}

	return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
		fmt.Sprintf("Tag `%s` was created with ID `%s`.", ident, itag.ID), "", static.ColorEmbedGreen).
		DeleteAfter(8 * time.Second).Error()
}

func (c *CmdTag) editTag(ctx shireikan.Context, db database.Database) error {
	if len(ctx.GetArgs()) < 3 {
		return printInvalidArguments(ctx)
	}

	tag, err, ok := getTag(ctx.GetArgs().Get(1).AsString(), ctx, db)
	if !ok || err != nil {
		return err
	}

	pmw, _ := ctx.GetObject(static.DiPermissions).(*permissions.Permissions)
	ok, override, err := pmw.CheckPermissions(ctx.GetSession(), ctx.GetGuild().ID, ctx.GetUser().ID, "!"+c.GetDomainName()+".edit")
	if err != nil {
		return err
	}

	if tag.CreatorID != ctx.GetUser().ID && !ok && !override {
		return printNotPermitted(ctx, "edit")
	}

	argsJoined := strings.Join(ctx.GetArgs()[:2], " ")
	contentOffset := strings.Index(ctx.GetMessage().Content, argsJoined) + len(argsJoined) + 1
	tag.Content = ctx.GetMessage().Content[contentOffset:]
	tag.LastEdit = time.Now()

	if err = db.EditTag(tag); err != nil {
		return err
	}

	return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
		fmt.Sprintf("Tag `%s` (ID `%s`) was updated.", tag.Ident, tag.ID), "", static.ColorEmbedGreen).
		DeleteAfter(8 * time.Second).Error()
}

func (c *CmdTag) deleteTag(ctx shireikan.Context, db database.Database) error {
	if len(ctx.GetArgs()) < 2 {
		return printInvalidArguments(ctx)
	}

	itag, err, ok := getTag(ctx.GetArgs().Get(1).AsString(), ctx, db)
	if !ok || err != nil {
		return err
	}

	pmw, _ := ctx.GetObject(static.DiPermissions).(*permissions.Permissions)
	ok, override, err := pmw.CheckPermissions(ctx.GetSession(), ctx.GetGuild().ID, ctx.GetUser().ID, "!"+c.GetDomainName()+".delete")
	if err != nil {
		return err
	}

	if itag.CreatorID != ctx.GetUser().ID && !ok && !override {
		return printNotPermitted(ctx, "delete")
	}

	if err = db.DeleteTag(itag.ID); err != nil {
		return err
	}

	return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
		"Tag was deleted.", "", static.ColorEmbedGreen).
		DeleteAfter(8 * time.Second).Error()
}

func (c *CmdTag) getRawTag(ctx shireikan.Context, db database.Database) error {
	if len(ctx.GetArgs()) < 2 {
		return printInvalidArguments(ctx)
	}

	tag, err, ok := getTag(ctx.GetArgs().Get(1).AsString(), ctx, db)
	if !ok || err != nil {
		return err
	}

	_, err = ctx.GetSession().ChannelMessageSend(ctx.GetChannel().ID, tag.RawContent())
	return err
}

func (c CmdTag) getTag(ctx shireikan.Context, db database.Database) error {
	tag, err, ok := getTag(ctx.GetArgs().Get(0).AsString(), ctx, db)
	if !ok || err != nil {
		return err
	}

	_, err = ctx.GetSession().ChannelMessageSendEmbed(ctx.GetChannel().ID, tag.AsEmbed(ctx.GetSession()))
	return err
}

func getTag(ident string, ctx shireikan.Context, db database.Database) (*tag.Tag, error, bool) {
	itag, err := db.GetTagByIdent(strings.ToLower(ident), ctx.GetGuild().ID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return nil, err, false
	}
	if itag != nil {
		return itag, nil, true
	}

	id, err := snowflake.ParseString(ident)
	if err != nil {
		return nil, printTagNotFound(ctx), false
	}

	itag, err = db.GetTagByID(id)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return nil, err, false
	}

	if itag == nil || itag.GuildID != ctx.GetGuild().ID {
		return nil, printTagNotFound(ctx), false
	}

	return itag, nil, true
}

func printInvalidArguments(ctx shireikan.Context) error {
	return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
		"Invalid arguments. Use `help tag` to ge thelp about how to use this command.").
		DeleteAfter(8 * time.Second).Error()
}

func printTagNotFound(ctx shireikan.Context) error {
	return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
		"Could not find any tag by the given identifier.").
		DeleteAfter(8 * time.Second).Error()
}

func printNotPermitted(ctx shireikan.Context, t string) error {
	return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
		fmt.Sprintf("You are not permitted to %s this tag.", t)).
		DeleteAfter(8 * time.Second).Error()
}

func checkPermission(ctx shireikan.Context, dn string) (error, bool) {
	pmw, _ := ctx.GetObject(static.DiPermissions).(*permissions.Permissions)
	ok, override, err := pmw.CheckPermissions(ctx.GetSession(), ctx.GetGuild().ID, ctx.GetUser().ID, dn)
	if err != nil {
		return err, false
	}

	if !ok && !override {
		err := util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
			"You are not permitted to use this command.").
			DeleteAfter(8 * time.Second).Error()
		return err, false
	}

	return nil, true
}
