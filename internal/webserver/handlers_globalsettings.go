package webserver

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/util/presence"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
)

// ---------------------------------------------------------------------------
// - GET /api/settings/presence

func (ws *WebServer) handlerGetPresence(ctx *routing.Context) error {
	presenceRaw, err := ws.db.GetSetting(static.SettingPresence)
	if err != nil {
		if database.IsErrDatabaseNotFound(err) {
			return jsonResponse(ctx, &presence.Presence{
				Game:   static.StdMotd,
				Status: "online",
			}, fasthttp.StatusOK)
		}
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	pre, err := presence.Unmarshal(presenceRaw)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	return jsonResponse(ctx, pre, fasthttp.StatusOK)
}

// ---------------------------------------------------------------------------
// - POST /api/settings/presence

func (ws *WebServer) handlerPostPresence(ctx *routing.Context) error {
	userID := ctx.Get("uid").(string)

	if ok, _, err := ws.pmw.CheckPermissions(ws.session, "", userID, "sp.game"); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	} else if !ok {
		return jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized)
	}

	pre := new(presence.Presence)
	if err := parseJSONBody(ctx, pre); err != nil {
		return jsonError(ctx, err, fasthttp.StatusBadRequest)
	}

	if err := pre.Validate(); err != nil {
		return jsonError(ctx, err, fasthttp.StatusBadRequest)
	}

	if err := ws.db.SetSetting(static.SettingPresence, pre.Marshal()); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	if err := ws.session.UpdateStatusComplex(pre.ToUpdateStatusData()); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	return jsonResponse(ctx, pre, fasthttp.StatusOK)
}

// ---------------------------------------------------------------------------
// - GET /api/settings/noguildinvite

func (ws *WebServer) handlerGetInviteSettings(ctx *routing.Context) error {
	var guildID, message, inviteCode string
	var err error

	if guildID, err = ws.db.GetSetting(static.SettingWIInviteGuildID); err != nil {
		if isErr, err := errInternalIgnoreNotFound(ctx, err); isErr {
			return err
		}
	}

	if guildID == "" {
		return jsonResponse(ctx, &InviteSettingsResponse{
			Guild:     nil,
			InviteURL: "",
			Message:   "",
		}, fasthttp.StatusOK)
	}

	if message, err = ws.db.GetSetting(static.SettingWIInviteText); err != nil {
		if isErr, err := errInternalIgnoreNotFound(ctx, err); isErr {
			return err
		}
	}

	if inviteCode, err = ws.db.GetSetting(static.SettingWIInviteCode); err != nil {
		if isErr, err := errInternalIgnoreNotFound(ctx, err); isErr {
			return err
		}
	}

	guild, err := discordutil.GetGuild(ws.session, guildID)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	invites, err := ws.session.GuildInvites(guildID)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	if inviteCode != "" {
		for _, inv := range invites {
			if inv.Inviter != nil && inv.Inviter.ID == ws.session.State.User.ID && !inv.Revoked {
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
			return jsonError(ctx, fmt.Errorf("could not find any channel to create invite for"), fasthttp.StatusConflict)
		}

		invite, err := ws.session.ChannelInviteCreate(channel.ID, discordgo.Invite{
			Temporary: false,
		})
		if err != nil {
			return jsonError(ctx, err, fasthttp.StatusInternalServerError)
		}

		inviteCode = invite.Code
		if err = ws.db.SetSetting(static.SettingWIInviteCode, inviteCode); err != nil {
			return jsonError(ctx, err, fasthttp.StatusInternalServerError)
		}
	}

	res := &InviteSettingsResponse{
		Guild:     GuildFromGuild(guild, nil, nil, ""),
		Message:   message,
		InviteURL: fmt.Sprintf("https://discord.gg/%s", inviteCode),
	}

	return jsonResponse(ctx, res, fasthttp.StatusOK)
}

// ---------------------------------------------------------------------------
// - POST /api/settings/noguildinvite

func (ws *WebServer) handlerPostInviteSettings(ctx *routing.Context) error {
	userID := ctx.Get("uid").(string)

	if ok, _, err := ws.pmw.CheckPermissions(ws.session, "", userID, "sp.noguildinvite"); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	} else if !ok {
		return jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized)
	}

	req := new(InviteSettingsRequest)
	if err := parseJSONBody(ctx, req); err != nil {
		return jsonError(ctx, err, fasthttp.StatusBadRequest)
	}

	var err error

	if req.GuildID != "" {

		guild, err := discordutil.GetGuild(ws.session, req.GuildID)
		if err != nil {
			return jsonError(ctx, err, fasthttp.StatusBadRequest)
		}

		if req.InviteCode != "" {
			invites, err := ws.session.GuildInvites(req.GuildID)
			if err != nil {
				return jsonError(ctx, err, fasthttp.StatusInternalServerError)
			}

			var valid bool
			for _, inv := range invites {
				if inv.Code == req.InviteCode && !inv.Revoked {
					valid = true
					break
				}
			}

			if !valid {
				return jsonError(ctx, fmt.Errorf("invalid invite code"), fasthttp.StatusBadRequest)
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
				return jsonError(ctx, fmt.Errorf("could not find any channel to create invite for"), fasthttp.StatusConflict)
			}

			invite, err := ws.session.ChannelInviteCreate(channel.ID, discordgo.Invite{
				Temporary: false,
			})
			if err != nil {
				return jsonError(ctx, err, fasthttp.StatusInternalServerError)
			}

			req.InviteCode = invite.Code
		}
	}

	if err = ws.db.SetSetting(static.SettingWIInviteCode, req.InviteCode); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	if err = ws.db.SetSetting(static.SettingWIInviteGuildID, req.GuildID); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	if err = ws.db.SetSetting(static.SettingWIInviteText, req.Messsage); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	return jsonResponse(ctx, nil, fasthttp.StatusOK)
}
