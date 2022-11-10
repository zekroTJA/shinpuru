package slashcommands

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/models"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
	"github.com/zekrotja/dgrs"
	"github.com/zekrotja/ken"
)

const nRoleOptions = 10

type Roleselect struct{}

var (
	_ ken.SlashCommand        = (*Say)(nil)
	_ permissions.PermCommand = (*Say)(nil)
)

func (c *Roleselect) Name() string {
	return "roleselect"
}

func (c *Roleselect) Description() string {
	return "Create a role selection."
}

func (c *Roleselect) Version() string {
	return "1.2.0"
}

func (c *Roleselect) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *Roleselect) Options() []*discordgo.ApplicationCommandOption {
	roleOptions := make([]*discordgo.ApplicationCommandOption, 0, nRoleOptions+1)

	for i := 0; i < nRoleOptions; i++ {
		roleOptions = append(roleOptions, &discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionRole,
			Name:        fmt.Sprintf("role%d", i+1),
			Description: fmt.Sprintf("Role %d", i+1),
			Required:    i == 0,
		})
	}

	options := []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "create",
			Description: "Create a message with attached role select buttons.",
			Options: append([]*discordgo.ApplicationCommandOption{{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "content",
				Description: "The content of the message.",
				Required:    true,
			}}, roleOptions...),
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "attach",
			Description: "Attach role select buttons to a shinpuru message (sent i.e. with /say)",
			Options: append([]*discordgo.ApplicationCommandOption{{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "id",
				Description: "The ID of the message.",
				Required:    true,
			}}, roleOptions...),
		},
	}

	return options
}

func (c *Roleselect) Domain() string {
	return "sp.guild.mod.roleselect"
}

func (c *Roleselect) SubDomains() []permissions.SubPermission {
	return nil
}

func (c *Roleselect) Run(ctx ken.Context) (err error) {
	err = ctx.HandleSubCommands(
		ken.SubCommandHandler{
			Name: "create",
			Run:  c.create,
		},
		ken.SubCommandHandler{
			Name: "attach",
			Run:  c.attach,
		},
	)

	return err
}

func (c *Roleselect) create(ctx ken.SubCommandContext) error {
	if err := ctx.Defer(); err != nil {
		return err
	}

	content := ctx.Options().GetByName("content").StringValue()

	fum := ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: content,
	})

	b := fum.AddComponents()

	roleIDs, err := c.attachRoleButtons(ctx, b)
	if err != nil {
		return err
	}

	roleSelects := mapRoleSelects(ctx.GetEvent().GuildID, fum.ChannelID, fum.ID, roleIDs)

	db := ctx.Get(static.DiDatabase).(database.Database)
	return db.AddRoleSelects(roleSelects)
}

func (c *Roleselect) attach(ctx ken.SubCommandContext) error {
	ctx.SetEphemeral(true)
	if err := ctx.Defer(); err != nil {
		return err
	}

	id := ctx.Options().GetByName("id").StringValue()

	st := ctx.Get(static.DiState).(*dgrs.State)

	msg, err := st.Message(ctx.GetEvent().ChannelID, id)
	if err != nil {
		if discordutil.IsErrCode(err, discordgo.ErrCodeUnknownMessage) {
			return ctx.FollowUpError("Message could not be found in this channel.", "").Error
		}
		return err
	}

	b := ctx.GetKen().Components().Add(msg.ID, msg.ChannelID)

	roleIDs, err := c.attachRoleButtons(ctx, b)
	if err != nil {
		return err
	}

	roleSelects := mapRoleSelects(ctx.GetEvent().GuildID, msg.ChannelID, msg.ID, roleIDs)

	db := ctx.Get(static.DiDatabase).(database.Database)
	err = db.AddRoleSelects(roleSelects)
	if err != nil {
		return err
	}

	return ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: "Role buttons have been attached.",
	}).DeleteAfter(6 * time.Second).Error
}

func (c *Roleselect) attachRoleButtons(
	ctx ken.SubCommandContext,
	b *ken.ComponentBuilder,
) ([]string, error) {
	roles := make([]*discordgo.Role, 0, nRoleOptions)
	for i := 0; i < nRoleOptions; i++ {
		r, ok := ctx.Options().GetByNameOptional(fmt.Sprintf("role%d", i+1))
		if ok {
			roles = append(roles, r.RoleValue(ctx))
		}
	}

	return util.AttachRoleSelectButtons(b, roles)
}

func mapRoleSelects(guildID, channelID, msgID string, roleIDs []string) []models.RoleSelect {
	roleSelects := make([]models.RoleSelect, 0, len(roleIDs))
	for _, rid := range roleIDs {
		roleSelects = append(roleSelects, models.RoleSelect{
			GuildID:   guildID,
			ChannelID: channelID,
			MessageID: msgID,
			RoleID:    rid,
		})
	}
	return roleSelects
}
