package slashcommands

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/xid"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/stringutil"
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
	return "Create a role selection message."
}

func (c *Roleselect) Version() string {
	return "1.1.0"
}

func (c *Roleselect) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *Roleselect) Options() []*discordgo.ApplicationCommandOption {
	roleOptions := make([]*discordgo.ApplicationCommandOption, 0, nRoleOptions+1)

	roleOptions = append(roleOptions, &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "content",
		Description: "The content of the message.",
		Required:    true,
	})

	for i := 0; i < nRoleOptions; i++ {
		roleOptions = append(roleOptions, &discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionRole,
			Name:        fmt.Sprintf("role%d", i+1),
			Description: fmt.Sprintf("Role %d", i+1),
			Required:    i == 0,
		})
	}

	return roleOptions
}

func (c *Roleselect) Domain() string {
	return "sp.guild.mod.roleselect"
}

func (c *Roleselect) SubDomains() []permissions.SubPermission {
	return nil
}

func (c *Roleselect) Run(ctx ken.Context) (err error) {
	if err = ctx.Defer(); err != nil {
		return
	}

	content := ctx.Options().GetByName("content").StringValue()

	type roleButton struct {
		Button *discordgo.Button
		RoleID string
	}

	var roleButtons []roleButton
	for i := 0; i < nRoleOptions; i++ {
		r, ok := ctx.Options().GetByNameOptional(fmt.Sprintf("role%d", i+1))
		if ok {
			role := r.RoleValue(ctx)
			roleButtons = append(roleButtons, roleButton{
				RoleID: role.ID,
				Button: &discordgo.Button{
					Label:    role.Name,
					Style:    discordgo.PrimaryButton,
					CustomID: xid.New().String(),
				},
			})
		}
	}

	nCols := len(roleButtons) / 5
	if len(roleButtons)%5 > 0 {
		nCols++
	}

	roleButtonsColumns := make([][]roleButton, nCols)
	for i, rb := range roleButtons {
		roleButtonsColumns[i/5] = append(roleButtonsColumns[i/5], rb)
	}

	fmt.Printf("%+v\n", roleButtonsColumns)
	b := ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: content,
	}).AddComponents()

	for _, rbs := range roleButtonsColumns {
		b.AddActionsRow(func(b ken.ComponentAssembler) {
			for _, rb := range rbs {
				b.Add(rb.Button, c.onRoleSelect(rb.RoleID))
			}
		})
	}

	_, err = b.Build()

	return err
}

func (c *Roleselect) onRoleSelect(roleID string) func(ctx ken.ComponentContext) bool {
	return func(ctx ken.ComponentContext) bool {
		ctx.Defer()

		if stringutil.ContainsAny(roleID, ctx.GetEvent().Member.Roles) {
			err := ctx.GetSession().GuildMemberRoleRemove(ctx.GetEvent().GuildID, ctx.User().ID, roleID)
			if err != nil {
				err = ctx.FollowUpError("Failed removing role.", "").DeleteAfter(4 * time.Second).Error
				return err == nil
			}
			err = ctx.FollowUpEmbed(&discordgo.MessageEmbed{
				Color:       static.ColorEmbedGreen,
				Description: fmt.Sprintf("Role <@&%s> has been removed.", roleID),
			}).DeleteAfter(4 * time.Second).Error
			return err == nil
		}

		err := ctx.GetSession().GuildMemberRoleAdd(ctx.GetEvent().GuildID, ctx.User().ID, roleID)
		if err != nil {
			err = ctx.FollowUpError("Failed adding role.", "").DeleteAfter(4 * time.Second).Error
			return err == nil
		}

		err = ctx.FollowUpEmbed(&discordgo.MessageEmbed{
			Color:       static.ColorEmbedGreen,
			Description: fmt.Sprintf("Role <@&%s> has been added.", roleID),
		}).DeleteAfter(4 * time.Second).Error

		return err == nil
	}
}
