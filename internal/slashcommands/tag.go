package slashcommands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/services/timeprovider"
	"github.com/zekroTJA/shinpuru/internal/util/snowflakenodes"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/internal/util/tag"
	"github.com/zekroTJA/shinpuru/pkg/acceptmsg/v2"
	"github.com/zekrotja/dgrs"
	"github.com/zekrotja/ken"
)

type Tag struct{}

var (
	_ ken.SlashCommand        = (*Tag)(nil)
	_ permissions.PermCommand = (*Tag)(nil)
)

func (c *Tag) Name() string {
	return "tag"
}

func (c *Tag) Description() string {
	return "Set texts as tags which can be fastly re-posted later."
}

func (c *Tag) Version() string {
	return "1.0.0"
}

func (c *Tag) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *Tag) Options() []*discordgo.ApplicationCommandOption {
	commonOpts := []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "name",
			Description: "The name of the Tag.",
			Required:    true,
		},
	}
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "show",
			Description: "Show the content of a tag.",
			Options:     commonOpts,
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "list",
			Description: "List created tags.",
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "set",
			Description: "Create or update a tag.",
			Options: append(commonOpts, []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "content",
					Description: "The content of the tag. You can use markdown as well as `\\n` for line breaks.",
					Required:    true,
				},
			}...),
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "delete",
			Description: "Delete a tag.",
			Options:     commonOpts,
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "raw",
			Description: "Show a raw tag.",
			Options:     commonOpts,
		},
	}
}

func (c *Tag) Domain() string {
	return "sp.chat.tag"
}

func (c *Tag) SubDomains() []permissions.SubPermission {
	return []permissions.SubPermission{
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

func (c *Tag) Run(ctx ken.Context) (err error) {
	if err = ctx.Defer(); err != nil {
		return
	}

	err = ctx.HandleSubCommands(
		ken.SubCommandHandler{"show", c.show},
		ken.SubCommandHandler{"list", c.list},
		ken.SubCommandHandler{"set", c.set},
		ken.SubCommandHandler{"delete", c.delete},
	)

	return
}

func (c *Tag) show(ctx ken.SubCommandContext) (err error) {
	db := ctx.Get(static.DiDatabase).(database.Database)
	st := ctx.Get(static.DiState).(*dgrs.State)

	ident := strings.ToLower(ctx.Options().GetByName("name").StringValue())

	tg, err := db.GetTagByIdent(ident, ctx.GetEvent().GuildID)
	if database.IsErrDatabaseNotFound(err) {
		return ctx.FollowUpError("Tag could not be found.", "").Send().Error
	}
	if err != nil {
		return
	}

	return ctx.FollowUpEmbed(tg.AsEmbed(st)).Send().Error
}

func (c *Tag) list(ctx ken.SubCommandContext) (err error) {
	db := ctx.Get(static.DiDatabase).(database.Database)
	st := ctx.Get(static.DiState).(*dgrs.State)

	tags, err := db.GetGuildTags(ctx.GetEvent().GuildID)
	if err != nil {
		return
	}

	tagsStr := make([]string, len(tags))
	for i, tag := range tags {
		tagsStr[i] = tag.AsEntry(st)
	}

	return ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Title:       "Registered Tags",
		Description: strings.Join(tagsStr, "\n"),
	}).Send().Error
}

func (c *Tag) set(ctx ken.SubCommandContext) (err error) {
	db := ctx.Get(static.DiDatabase).(database.Database)
	st := ctx.Get(static.DiState).(*dgrs.State)
	pmw := ctx.Get(static.DiPermissions).(*permissions.Permissions)

	ident := strings.ToLower(ctx.Options().GetByName("name").StringValue())
	content := ctx.Options().GetByName("content").StringValue()

	tg, err := db.GetTagByIdent(ident, ctx.GetEvent().GuildID)

	if database.IsErrDatabaseNotFound(err) {
		if tg.CreatorID != ctx.User().ID {
			ok, err := pmw.CheckSubPerm(ctx, "edit", true,
				"A tag with the same nam (created by another user) already exists and you do not have the permission to edit it.")
			if !ok {
				return err
			}
		}
		var creator *discordgo.User
		creator, err = st.User(tg.CreatorID)
		if err != nil {
			return err
		}
		emb := &discordgo.MessageEmbed{
			Color: static.ColorEmbedOrange,
			Description: fmt.Sprintf(
				"A tag with the name `%s` already assists - created by %s "+
					"- with the following content:\n%s\n"+
					"Do you really want to overwrite this tag?",
				tg.Ident, creator.Mention(), tg.RawContent(),
			),
		}
		_, err = acceptmsg.New().
			WithKen(ctx.GetKen()).
			WithEmbed(emb).
			LockOnUser(ctx.User().ID).
			DeleteAfterAnswer().
			DoOnAccept(func(cctx ken.ComponentContext) (err error) {
				if err = cctx.Defer(); err != nil {
					return
				}
				tg.Content = content
				if err = db.EditTag(tg); err != nil {
					return
				}
				return cctx.FollowUpEmbed(&discordgo.MessageEmbed{
					Description: fmt.Sprintf(
						"Tag has been updated.\nUse the command `/tag show %s` to use the tag.",
						tg.Ident),
				}).Send().Error
			}).
			AsFollowUp(ctx)
		return
	}

	if err != nil {
		return
	}

	ok, err := pmw.CheckSubPerm(ctx, "create", true,
		"You do not have the permission to create tags.")
	if !ok {
		return err
	}

	tp := ctx.Get(static.DiTimeProvider).(timeprovider.Provider)

	now := tp.Now()
	tg = tag.Tag{
		Content:   content,
		Created:   now,
		CreatorID: ctx.User().ID,
		GuildID:   ctx.GetEvent().GuildID,
		ID:        snowflakenodes.NodeTags.Generate(),
		Ident:     ident,
		LastEdit:  now,
	}
	if err = db.AddTag(tg); err != nil {
		return
	}

	return ctx.RespondEmbed(&discordgo.MessageEmbed{
		Description: fmt.Sprintf(
			"Tag has been created.\nUse the command `/tag show %s` to use the tag.",
			tg.Ident),
	})
}

func (c *Tag) delete(ctx ken.SubCommandContext) (err error) {
	db := ctx.Get(static.DiDatabase).(database.Database)
	pmw := ctx.Get(static.DiPermissions).(*permissions.Permissions)

	ident := strings.ToLower(ctx.Options().GetByName("name").StringValue())

	tg, err := db.GetTagByIdent(ident, ctx.GetEvent().GuildID)
	if database.IsErrDatabaseNotFound(err) {
		return ctx.FollowUpError("Tag could not be found.", "").Send().Error
	}
	if err != nil {
		return
	}

	if tg.CreatorID != ctx.User().ID {
		ok, err := pmw.CheckSubPerm(ctx, "delete", true,
			"A tag with the same nam (created by another user) already exists and you do not have the permission to edit it.")
		if !ok {
			return err
		}
	}

	if err = db.DeleteTag(tg.ID); err != nil {
		return
	}

	return ctx.RespondEmbed(&discordgo.MessageEmbed{
		Description: "Tag has been deleted.",
	})
}
