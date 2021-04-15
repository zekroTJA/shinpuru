package controllers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/core/middleware"
	"github.com/zekroTJA/shinpuru/internal/core/webserver/v1/models"
	"github.com/zekroTJA/shinpuru/internal/util/presence"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
)

type GlobalSettingsController struct {
	session *discordgo.Session
	db      database.Database
}

func (c *GlobalSettingsController) Setup(container di.Container, router fiber.Router) {
	c.session = container.Get(static.DiDiscordSession).(*discordgo.Session)
	c.db = container.Get(static.DiDatabase).(database.Database)

	pmw := container.Get(static.DiPermissionMiddleware).(*middleware.PermissionsMiddleware)

	router.Get("/presence", c.getPresence)
	router.Post("/presence", pmw.HandleWs(c.session, "sp.game"), c.postPresence)
	router.Get("/noguildinvite", c.getNoGuildInvites)
	router.Post("/noguildinvite", pmw.HandleWs(c.session, "sp.noguildinvite"), c.postPresence, c.postNoGuildInvites)
}

func (c *GlobalSettingsController) getPresence(ctx *fiber.Ctx) error {
	presenceRaw, err := c.db.GetSetting(static.SettingPresence)
	if err != nil {
		if database.IsErrDatabaseNotFound(err) {
			return ctx.JSON(&presence.Presence{
				Game:   static.StdMotd,
				Status: "online",
			})
		}
		return err
	}

	pre, err := presence.Unmarshal(presenceRaw)
	if err != nil {
		return err
	}

	return ctx.JSON(pre)
}

func (c *GlobalSettingsController) postPresence(ctx *fiber.Ctx) error {
	pre := new(presence.Presence)
	if err := ctx.BodyParser(pre); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := pre.Validate(); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := c.db.SetSetting(static.SettingPresence, pre.Marshal()); err != nil {
		return err
	}

	if err := c.session.UpdateStatusComplex(pre.ToUpdateStatusData()); err != nil {
		return err
	}

	return ctx.JSON(pre)
}

func (c *GlobalSettingsController) getNoGuildInvites(ctx *fiber.Ctx) error {
	var guildID, message, inviteCode string
	var err error

	if guildID, err = c.db.GetSetting(static.SettingWIInviteGuildID); err != nil {
		if err != nil && !database.IsErrDatabaseNotFound(err) {
			return err
		}
	}

	if guildID == "" {
		return ctx.JSON(&models.InviteSettingsResponse{
			Guild:     nil,
			InviteURL: "",
			Message:   "",
		})
	}

	if message, err = c.db.GetSetting(static.SettingWIInviteText); err != nil {
		if err != nil && !database.IsErrDatabaseNotFound(err) {
			return err
		}
	}

	if inviteCode, err = c.db.GetSetting(static.SettingWIInviteCode); err != nil {
		if err != nil && !database.IsErrDatabaseNotFound(err) {
			return err
		}
	}

	guild, err := discordutil.GetGuild(c.session, guildID)
	if err != nil {
		return err
	}

	invites, err := c.session.GuildInvites(guildID)
	if err != nil {
		return err
	}

	if inviteCode != "" {
		for _, inv := range invites {
			if inv.Inviter != nil && inv.Inviter.ID == c.session.State.User.ID && !inv.Revoked {
				inviteCode = inv.Code
				break
			}
		}
	}

	if inviteCode == "" {
		var channel *discordgo.Channel
		for _, c := range guild.Channels {
			if c.Type == discordgo.ChannelTypeGuildText {
				channel = c
				break
			}
		}
		if channel == nil {
			return fiber.NewError(fiber.StatusConflict, "could not find any channel to create invite for")
		}

		invite, err := c.session.ChannelInviteCreate(channel.ID, discordgo.Invite{
			Temporary: false,
		})
		if err != nil {
			return err
		}

		inviteCode = invite.Code
		if err = c.db.SetSetting(static.SettingWIInviteCode, inviteCode); err != nil {
			return err
		}
	}

	res := &models.InviteSettingsResponse{
		Guild:     models.GuildFromGuild(guild, nil, nil, ""),
		Message:   message,
		InviteURL: fmt.Sprintf("https://discord.gg/%s", inviteCode),
	}

	return ctx.JSON(res)
}

func (c *GlobalSettingsController) postNoGuildInvites(ctx *fiber.Ctx) error {
	req := new(models.InviteSettingsRequest)
	if err := ctx.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var err error

	if req.GuildID != "" {

		guild, err := discordutil.GetGuild(c.session, req.GuildID)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		if req.InviteCode != "" {
			invites, err := c.session.GuildInvites(req.GuildID)
			if err != nil {
				return err
			}

			var valid bool
			for _, inv := range invites {
				if inv.Code == req.InviteCode && !inv.Revoked {
					valid = true
					break
				}
			}

			if !valid {
				return fiber.NewError(fiber.StatusBadRequest, "invalid invite code")
			}
		} else {
			var channel *discordgo.Channel
			for _, c := range guild.Channels {
				if c.Type == discordgo.ChannelTypeGuildText {
					channel = c
					break
				}
			}
			if channel == nil {
				return fiber.NewError(fiber.StatusConflict, "could not find any channel to create invite for")
			}

			invite, err := c.session.ChannelInviteCreate(channel.ID, discordgo.Invite{
				Temporary: false,
			})
			if err != nil {
				return err
			}

			req.InviteCode = invite.Code
		}
	}

	if err = c.db.SetSetting(static.SettingWIInviteCode, req.InviteCode); err != nil {
		return err
	}

	if err = c.db.SetSetting(static.SettingWIInviteGuildID, req.GuildID); err != nil {
		return err
	}

	if err = c.db.SetSetting(static.SettingWIInviteText, req.Messsage); err != nil {
		return err
	}

	return ctx.JSON(struct{}{})
}
