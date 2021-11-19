package controllers

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/services/kvcache"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/services/webserver/v1/models"
	_ "github.com/zekroTJA/shinpuru/internal/services/webserver/v1/models"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
	"github.com/zekroTJA/shinpuru/pkg/stringutil"
	"github.com/zekrotja/dgrs"
)

type ChannelController struct {
	session *discordgo.Session
	st      *dgrs.State
	pmw     *permissions.Permissions
	kv      kvcache.Provider
}

func (c *ChannelController) Setup(container di.Container, router fiber.Router) {
	c.session = container.Get(static.DiDiscordSession).(*discordgo.Session)
	c.st = container.Get(static.DiState).(*dgrs.State)
	c.pmw = container.Get(static.DiPermissions).(*permissions.Permissions)
	c.kv = container.Get(static.DiKVCache).(kvcache.Provider)

	router.Get("", c.getChannels)
	router.Post("/:id", c.pmw.HandleWs(c.session, "sp.chat.say"), c.postChannelMessage)
	router.Post("/:id/:msgid", c.pmw.HandleWs(c.session, "sp.chat.say"), c.postChannelMessage)
}

// @Summary Get Allowed Channels
// @Description Returns a list of channels the user has access to.
// @Tags Channels
// @Accept json
// @Produce json
// @Param guildid path string true "The ID of the guild."
// @Success 201 {object} discordgo.Message
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /channels/{guildid} [get]
func (c *ChannelController) getChannels(ctx *fiber.Ctx) (err error) {
	uid := ctx.Locals("uid").(string)
	gid := ctx.Params("guildid")

	guildChans, err := c.st.Channels(gid)
	if err != nil {
		return
	}

	chans := make([]*models.ChannelWithPermissions, 0)
	var perms int64
	for _, gc := range guildChans {
		if perms, err = c.getUserChannelPermissions(uid, gc.ID); err != nil {
			return
		}
		if perms&discordgo.PermissionViewChannel != 0 {
			chans = append(chans, &models.ChannelWithPermissions{
				Channel:  gc,
				CanRead:  true,
				CanWrite: perms&discordgo.PermissionSendMessages != 0,
			})
		}
	}

	return ctx.JSON(models.ListResponse{N: len(chans), Data: chans})
}

// @Summary Send Embed Message
// @Description Send an Embed Message into a specified Channel.
// @Tags Channels
// @Accept json
// @Produce json
// @Param guildid path string true "The ID of the guild."
// @Param id path string true "The ID of the channel."
// @Param payload body discordgo.MessageEmbed true "The message embed object."
// @Success 201 {object} discordgo.Message
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /channels/{guildid}/{id} [post]
func (c *ChannelController) postChannelMessage(ctx *fiber.Ctx) (err error) {
	uid := ctx.Locals("uid").(string)
	id := ctx.Params("id")
	msgid := ctx.Params("msgid")

	perms, err := c.getUserChannelPermissions(uid, id)
	if err != nil {
		return
	}

	if perms&discordgo.PermissionSendMessages != discordgo.PermissionSendMessages {
		return fiber.ErrForbidden
	}

	emb := new(discordgo.MessageEmbed)
	if err = ctx.BodyParser(emb); err != nil {
		return
	}

	ch, err := c.st.Channel(id)
	if err != nil {
		if discordutil.IsErrCode(err, discordgo.ErrCodeUnknownChannel) {
			err = fiber.ErrNotFound
		}
		return
	}

	gids, err := c.st.UserGuilds(uid)
	if err != nil {
		return
	}

	if !stringutil.ContainsAny(ch.GuildID, gids) {
		return fiber.ErrNotFound
	}

	memb, err := c.st.Member(ch.GuildID, uid)
	if err != nil {
		return
	}

	emb.Author = &discordgo.MessageEmbedAuthor{
		Name:    memb.User.String(),
		IconURL: memb.User.AvatarURL("16"),
	}

	var msg *discordgo.Message
	if msgid != "" {
		if _, err = c.st.Message(ch.ID, msgid); err != nil {
			if discordutil.IsErrCode(err, discordgo.ErrCodeUnknownMessage) {
				return fiber.ErrNotFound
			}
			return
		}
		msg, err = c.session.ChannelMessageEditEmbed(ch.ID, msgid, emb)
		ctx.Status(fiber.StatusOK)
	} else {
		msg, err = c.session.ChannelMessageSendEmbed(ch.ID, emb)
		ctx.Status(fiber.StatusCreated)
	}

	if err != nil {
		return
	}

	return ctx.JSON(msg)
}

// @Summary Update Embed Message
// @Description Update an Embed Message in a specified Channel with the given message ID.
// @Tags Channels
// @Accept json
// @Produce json
// @Param guildid path string true "The ID of the guild."
// @Param id path string true "The ID of the channel."
// @Param msgid path string true "The ID of the message."
// @Param payload body discordgo.MessageEmbed true "The message embed object."
// @Success 200 {object} discordgo.Message
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /channels/{guildid}/{id}/{msgid} [post]
//
// This is a dummy method for API doc generation.
func (*ChannelController) _(*fiber.Ctx) error {
	return nil
}

// --- HELPER ---

func (c *ChannelController) getUserChannelPermissions(uid, chid string) (perms int64, err error) {
	var ok bool
	cacheKey := fmt.Sprintf("userchanperms:%s:%s", uid, chid)
	if perms, ok = c.kv.Get(cacheKey).(int64); !ok {
		perms, err = c.session.UserChannelPermissions(uid, chid)
		if err != nil {
			return
		}
		c.kv.Set(cacheKey, perms, 10*time.Minute)
	}
	return
}
