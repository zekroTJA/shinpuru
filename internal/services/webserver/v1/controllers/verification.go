package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kataras/hcaptcha"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/services/config"
	"github.com/zekroTJA/shinpuru/internal/services/verification"
	"github.com/zekroTJA/shinpuru/internal/services/webserver/v1/models"
	"github.com/zekroTJA/shinpuru/internal/util/static"
)

type VerificationController struct {
	hc  *hcaptcha.Client
	cfg config.Provider
	vs  verification.Provider
}

func (c *VerificationController) Setup(container di.Container, router fiber.Router) {
	c.cfg = container.Get(static.DiConfig).(config.Provider)
	c.vs = container.Get(static.DiVerification).(verification.Provider)

	c.hc = hcaptcha.New(c.cfg.Config().WebServer.Captcha.SecretKey)

	router.Get("/sitekey", c.getSitekey)
	router.Post("/verify", c.postVerify)
}

// @Summary Sitekey
// @Description Returns the sitekey for the captcha
// @Tags Verification
// @Accept json
// @Produce json
// @Success 200 {object} models.CaptchaSiteKey
// @Router /verification/sitekey [get]
func (c *VerificationController) getSitekey(ctx *fiber.Ctx) error {
	res := models.CaptchaSiteKey{
		SiteKey: c.cfg.Config().WebServer.Captcha.SiteKey,
	}

	return ctx.JSON(res)
}

// @Summary Verify
// @Description Verify a returned verification token.
// @Tags Verification
// @Accept json
// @Produce json
// @Success 200 {object} models.User
// @Router /verification/verify [post]
func (c *VerificationController) postVerify(ctx *fiber.Ctx) error {
	uid := ctx.Locals("uid").(string)

	var req models.CaptchaVerificationRequest
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	res := c.hc.VerifyToken(req.Token)
	if !res.Success {
		return fiber.ErrForbidden
	}

	err := c.vs.Verify(uid)
	if err != nil {
		return err
	}

	return ctx.JSON(models.Ok)
}
