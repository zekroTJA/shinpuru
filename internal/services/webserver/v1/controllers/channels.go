package controllers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
	"github.com/zekroTJA/shinpuru/pkg/stringutil"
	"github.com/zekrotja/dgrs"
)

type ChannelController struct {
	session *discordgo.Session
	st      *dgrs.State
	pmw     *permissions.Permissions
}

func (c *ChannelController) Setup(container di.Container, router fiber.Router) {
	c.session = container.Get(static.DiDiscordSession).(*discordgo.Session)
	c.st = container.Get(static.DiState).(*dgrs.State)
	c.pmw = container.Get(static.DiPermissions).(*permissions.Permissions)

	router.Post("/:id", c.pmw.HandleWs(c.session, "sp.chat.say"), c.postChannelMessage)
	router.Post("/:id/:msgid", c.pmw.HandleWs(c.session, "sp.chat.say"), c.postChannelMessage)
}

// @Summary Send Embed Message
// @Description Send an Embed Message into a specified Channel.
// @Tags Channels
// @Accept json
// @Produce json
// @Param id path string true "The ID of the channel."
// @Param payload body discordgo.MessageEmbed true "The message embed object."
// @Success 201 {object} discordgo.Message
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /channels/{id} [post]
func (c *ChannelController) postChannelMessage(ctx *fiber.Ctx) (err error) {
	uid := ctx.Locals("uid").(string)
	id := ctx.Params("id")
	msgid := ctx.Params("msgid")

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
// @Param id path string true "The ID of the channel."
// @Param msgid path string true "The ID of the message."
// @Param payload body discordgo.MessageEmbed true "The message embed object."
// @Success 200 {object} discordgo.Message
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /channels/{id}/{msgid} [post]
//
// This is a dummy method for API doc generation.
func (*ChannelController) _(*fiber.Ctx) error {
	return nil
}
