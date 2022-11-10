package slashcommands

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/xid"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
	"github.com/zekroTJA/shinpuru/pkg/stringutil"
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

	b := ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: content,
	}).AddComponents()

	return c.attachRoleButtons(ctx, b)
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

	err = c.attachRoleButtons(ctx, b)
	if err != nil {
		return err
	}

	return ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: "Role buttons have been attached.",
	}).DeleteAfter(6 * time.Second).Error
}

func (c *Roleselect) attachRoleButtons(ctx ken.SubCommandContext, b *ken.ComponentBuilder) error {
	type roleButton struct {
		Button *discordgo.Button
		RoleID string
	}

	roleButtons := map[string]*discordgo.Button{}
	for i := 0; i < nRoleOptions; i++ {
		r, ok := ctx.Options().GetByNameOptional(fmt.Sprintf("role%d", i+1))
		if ok {
			role := r.RoleValue(ctx)
			roleButtons[role.ID] = &discordgo.Button{
				Label:    role.Name,
				Style:    discordgo.PrimaryButton,
				CustomID: xid.New().String(),
			}
		}
	}

	nCols := len(roleButtons) / 5
	if len(roleButtons)%5 > 0 {
		nCols++
	}

	roleButtonsColumns := make([][]roleButton, nCols)
	i := 0
	for id, b := range roleButtons {
		roleButtonsColumns[i/5] = append(roleButtonsColumns[i/5], roleButton{
			Button: b,
			RoleID: id,
		})
		i++
	}

	for _, rbs := range roleButtonsColumns {
		b.AddActionsRow(func(b ken.ComponentAssembler) {
			for _, rb := range rbs {
				b.Add(rb.Button, c.onRoleSelect(rb.RoleID))
			}
		})
	}

	_, err := b.Build()

	return err
}

func (c *Roleselect) onRoleSelect(roleID string) func(ctx ken.ComponentContext) bool {
	return func(ctx ken.ComponentContext) bool {
		ctx.SetEphemeral(true)
		ctx.Defer()

		if stringutil.ContainsAny(roleID, ctx.GetEvent().Member.Roles) {
			err := ctx.GetSession().GuildMemberRoleRemove(ctx.GetEvent().GuildID, ctx.User().ID, roleID)
			if err != nil {
				err = ctx.FollowUpError("Failed removing role.", "").DeleteAfter(10 * time.Second).Error
				return err == nil
			}
			err = ctx.FollowUpEmbed(&discordgo.MessageEmbed{
				Color:       static.ColorEmbedGreen,
				Description: fmt.Sprintf("Role <@&%s> has been removed.", roleID),
			}).DeleteAfter(10 * time.Second).Error
			return err == nil
		}

		err := ctx.GetSession().GuildMemberRoleAdd(ctx.GetEvent().GuildID, ctx.User().ID, roleID)
		if err != nil {
			err = ctx.FollowUpError("Failed adding role.", "").DeleteAfter(10 * time.Second).Error
			return err == nil
		}

		err = ctx.FollowUpEmbed(&discordgo.MessageEmbed{
			Color:       static.ColorEmbedGreen,
			Description: fmt.Sprintf("Role <@&%s> has been added.", roleID),
		}).DeleteAfter(10 * time.Second).Error

		return err == nil
	}
}
