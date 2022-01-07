package controllers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/gofiber/fiber/v2"
	"github.com/kataras/hcaptcha"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/services/config"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/webserver/auth"
	"github.com/zekroTJA/shinpuru/internal/services/webserver/v1/models"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/multierror"
	"github.com/zekrotja/dgrs"
)

type VerificationController struct {
	hc      *hcaptcha.Client
	cfg     config.Provider
	db      database.Database
	session *discordgo.Session
	authMw  auth.Middleware
	st      *dgrs.State
}

func (c *VerificationController) Setup(container di.Container, router fiber.Router) {
	c.session = container.Get(static.DiDiscordSession).(*discordgo.Session)
	c.cfg = container.Get(static.DiConfig).(config.Provider)
	c.authMw = container.Get(static.DiAuthMiddleware).(auth.Middleware)
	c.st = container.Get(static.DiState).(*dgrs.State)
	c.db = container.Get(static.DiDatabase).(database.Database)

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

	if err := c.db.SetUserVerified(uid, true); err != nil {
		return err
	}

	queue, err := c.db.GetVerificationQueue("", uid)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	mErr := multierror.New()
	for _, e := range queue {
		ok, err := c.db.RemoveVerificationQueue(e.GuildID, e.UserID)
		mErr.Append(err)
		if ok {
			mErr.Append(c.session.GuildMemberTimeout(e.GuildID, e.UserID, nil))
		}
	}

	if mErr.Len() != 0 {
		return mErr
	}

	return ctx.JSON(models.Ok)
}
