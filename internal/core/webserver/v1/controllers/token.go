package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/core/webserver/auth"
	"github.com/zekroTJA/shinpuru/internal/core/webserver/v1/models"
	"github.com/zekroTJA/shinpuru/internal/util/static"
)

type TokenController struct {
	db    database.Database
	apith auth.APITokenHandler
}

func (c *TokenController) Setup(container di.Container, router fiber.Router) {
	c.db = container.Get(static.DiDatabase).(database.Database)
	c.apith = container.Get(static.DiAuthAPITokenHandler).(auth.APITokenHandler)

	router.Get("", c.getToken)
	router.Post("", c.postToken)
	router.Delete("", c.deleteToken)
}

func (c *TokenController) getToken(ctx *fiber.Ctx) error {
	uid := ctx.Locals("uid").(string)

	token, err := c.db.GetAPIToken(uid)
	if database.IsErrDatabaseNotFound(err) {
		return fiber.NewError(fiber.StatusNotFound, "no token found")
	} else if err != nil {
		return err
	}

	tokenResp := &models.APITokenResponse{
		Created:    token.Created,
		Expires:    token.Expires,
		Hits:       token.Hits,
		LastAccess: token.LastAccess,
	}

	return ctx.JSON(tokenResp)
}

func (c *TokenController) postToken(ctx *fiber.Ctx) error {
	uid := ctx.Locals("uid").(string)

	token, expires, err := c.apith.GetAPIToken(uid)
	if err != nil {
		return err
	}

	return ctx.JSON(&models.APITokenResponse{
		Created: time.Now(),
		Expires: expires,
		Token:   token,
	})
}

func (c *TokenController) deleteToken(ctx *fiber.Ctx) error {
	uid := ctx.Locals("uid").(string)

	err := c.apith.RevokeToken(uid)
	if err != nil {
		return err
	}

	return ctx.JSON(struct{}{})
}
