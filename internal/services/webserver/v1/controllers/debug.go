package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/services/webserver/v1/models"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekrotja/dgrs"
)

type DebugController struct {
	st dgrs.IState
}

func (c *DebugController) Setup(container di.Container, router fiber.Router) {
	c.st = container.Get(static.DiState).(*dgrs.State)

	router.Get("", c.get)
}

func (c *DebugController) get(ctx *fiber.Ctx) error {
	g, err := c.st.Guild("526196711962705925")
	if err != nil {
		return err
	}

	g.MemberCount = 1
	err = c.st.SetGuild(g)
	if err != nil {
		return err
	}

	return ctx.JSON(models.Ok)
}
