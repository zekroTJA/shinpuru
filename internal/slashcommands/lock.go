package slashcommands

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekrotja/dgrs"
	"github.com/zekrotja/ken"
)

const allowMask = discordgo.PermissionAll - discordgo.PermissionSendMessages

type Lock struct{}

var (
	_ ken.SlashCommand        = (*Lock)(nil)
	_ permissions.PermCommand = (*Lock)(nil)
)

func (c *Lock) Name() string {
	return "lock"
}

func (c *Lock) Description() string {
	return "Lock or unlock a channel so that no messages can be sent anymore."
}

func (c *Lock) Version() string {
	return "1.0.0"
}

func (c *Lock) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *Lock) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:         discordgo.ApplicationCommandOptionChannel,
			Name:         "channel",
			Description:  "The channel to be locked or unlocked (selects current channel if not passed).",
			ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildText},
		},
	}
}

func (c *Lock) Domain() string {
	return "sp.guild.mod.lock"
}

func (c *Lock) SubDomains() []permissions.SubPermission {
	return nil
}

func (c *Lock) Run(ctx ken.Context) (err error) {
	if err = ctx.Defer(); err != nil {
		return
	}

	st := ctx.Get(static.DiState).(*dgrs.State)
	db := ctx.Get(static.DiDatabase).(database.Database)

	var ch *discordgo.Channel
	if chV, ok := ctx.Options().GetByNameOptional("channel"); ok {
		ch = chV.ChannelValue(ctx)
	} else {
		ch, err = st.Channel(ctx.GetEvent().ChannelID)
		if err != nil {
			return
		}
	}

	_, _, encodedPerms, err := db.GetLockChan(ch.ID)
	if database.IsErrDatabaseNotFound(err) {
		return c.lock(ch, ctx)
	} else if err == nil {
		return c.unlock(ch, ctx, encodedPerms)
	}

	return
}

func (c *Lock) lock(target *discordgo.Channel, ctx ken.Context) error {
	st := ctx.Get(static.DiState).(*dgrs.State)
	db := ctx.Get(static.DiDatabase).(database.Database)

	procMsg := ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: ":clock4: Locking channel...",
		Color:       static.ColorEmbedGray,
	}).Send()
	if procMsg.Error != nil {
		return procMsg.Error
	}

	encodedPerms, err := c.encodePermissionOverrides(target.PermissionOverwrites)
	if err != nil {
		return err
	}

	guildRoles, err := st.Roles(ctx.GetEvent().GuildID)
	if err != nil {
		return err
	}
	sort.Slice(guildRoles, func(i, j int) bool {
		return guildRoles[i].Position < guildRoles[j].Position
	})

	memberRoles := ctx.GetEvent().Member.Roles

	highest := 0
	rolesMap := make(map[string]*discordgo.Role)
	for _, r := range guildRoles {
		rolesMap[r.ID] = r
		for _, mr := range memberRoles {
			if r.ID != mr {
				continue
			}
			if r.Position > highest {
				highest = r.Position
			}
		}
	}

	// The info message needs to be sent before all permissions are set
	// to prevent occuring errors due to potential missing permissions.
	err = procMsg.EditEmbed(&discordgo.MessageEmbed{
		Description: fmt.Sprintf("This channel is chat-locked by %s.\nYou may not be able to chat "+
			"into this channel until the channel is unlocked again.", ctx.User().Mention()),
		Color: static.ColorEmbedOrange,
	})
	if err != nil {
		return err
	}

	self, err := st.SelfUser()
	if err != nil {
		return err
	}

	hasSetEveryone := false
	for _, po := range target.PermissionOverwrites {
		if po.Type == discordgo.PermissionOverwriteTypeRole {
			if r, ok := rolesMap[po.ID]; ok && r.Position < highest {
				if err = ctx.GetSession().ChannelPermissionSet(
					target.ID, po.ID, discordgo.PermissionOverwriteTypeRole, po.Allow&allowMask, po.Deny|discordgo.PermissionSendMessages); err != nil {
					return err
				}
			}
		}
		if po.Type == discordgo.PermissionOverwriteTypeMember && ctx.User().ID != po.ID && self.ID != po.ID {
			if err = ctx.GetSession().ChannelPermissionSet(
				target.ID, po.ID, discordgo.PermissionOverwriteTypeMember, po.Allow&allowMask, po.Deny|discordgo.PermissionSendMessages); err != nil {
				return err
			}
			if po.ID == target.GuildID {
				hasSetEveryone = true
			}
		}
	}

	if err = ctx.GetSession().ChannelPermissionSet(
		target.ID, self.ID, discordgo.PermissionOverwriteTypeMember, discordgo.PermissionSendMessages&discordgo.PermissionViewChannel, 0); err != nil {
		return err
	}

	if !hasSetEveryone {
		if err = ctx.GetSession().ChannelPermissionSet(
			target.ID, target.GuildID, discordgo.PermissionOverwriteTypeRole, 0, discordgo.PermissionSendMessages); err != nil {
			return err
		}
	}

	if err = db.SetLockChan(target.ID, target.GuildID, ctx.User().ID, encodedPerms); err != nil {
		return err
	}

	return nil
}

func (c *Lock) unlock(target *discordgo.Channel, ctx ken.Context, encodedPerms string) error {
	db := ctx.Get(static.DiDatabase).(database.Database)

	procMsg := ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: ":clock4: Locking channel...",
		Color:       static.ColorEmbedGray,
	}).Send()
	if procMsg.Error != nil {
		return procMsg.Error
	}

	permissionOverrides, err := c.decodePermissionOverrrides(encodedPerms)
	if err != nil {
		return err
	}

	failed := 0
	for _, po := range permissionOverrides {
		if err = ctx.GetSession().ChannelPermissionSet(target.ID, po.ID, po.Type, po.Allow, po.Deny); err != nil {
			failed++
		}
	}

	if err = db.DeleteLockChan(target.ID); err != nil {
		return err
	}

	if failed > 0 {
		return procMsg.EditEmbed(&discordgo.MessageEmbed{
			Description: fmt.Sprintf("This channel is now unlocked. You can now chat here again.\n*(Unlocked by %s)*\n\n"+
				"**Attention:** %d permission actions failed on reset!", ctx.User().Mention(), failed),
			Color: static.ColorEmbedOrange,
		})
	}

	return procMsg.EditEmbed(&discordgo.MessageEmbed{
		Description: fmt.Sprintf("This channel is now unlocked. You can now chat here again.\n*(Unlocked by %s)*", ctx.User().Mention()),
		Color:       static.ColorEmbedGreen,
	})
}

func (c *Lock) encodePermissionOverrides(po []*discordgo.PermissionOverwrite) (res string, err error) {
	buff := bytes.NewBuffer([]byte{})

	if err = json.NewEncoder(buff).Encode(po); err != nil {
		return
	}

	res = base64.StdEncoding.EncodeToString(buff.Bytes())

	return
}

func (c *Lock) decodePermissionOverrrides(data string) (po []*discordgo.PermissionOverwrite, err error) {
	po = make([]*discordgo.PermissionOverwrite, 0)

	dataBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return
	}

	err = json.NewDecoder(bytes.NewBuffer(dataBytes)).Decode(&po)

	return
}
