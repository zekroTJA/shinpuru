package controllers

import (
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/xid"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/webserver/auth"
	"github.com/zekroTJA/shinpuru/internal/services/webserver/v1/models"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/acceptmsg/v2"
	"github.com/zekroTJA/shinpuru/pkg/discordoauth/v2"
	"github.com/zekroTJA/timedmap"
	"github.com/zekrotja/dgrs"
	"github.com/zekrotja/ken"
	"github.com/zekrotja/rogu/log"
)

const pushcodeTimeout = 60 * time.Second

type AuthController struct {
	rth          auth.RefreshTokenHandler
	ath          auth.AccessTokenHandler
	authMw       auth.Middleware
	oauthHandler auth.RequestHandler

	st      State
	session Session

	cmdHandler   *ken.Ken
	discordOAuth *discordoauth.DiscordOAuth

	pushcodeSubs *timedmap.TimedMap
}

type pushCodeWaiter struct {
	mtx          sync.Mutex
	code         string
	am           *acceptmsg.AcceptMessage
	subscription func() error
	fin          chan *discordgo.Message
	closed       bool
}

func (pcw *pushCodeWaiter) close() bool {
	pcw.mtx.Lock()
	defer pcw.mtx.Unlock()

	if pcw.am != nil {
		pcw.am.Ken.Session().ChannelMessageEditComplex(&discordgo.MessageEdit{
			Channel: pcw.am.ChannelID,
			ID:      pcw.am.ID,
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Login",
					Description: "The code has been timed out.",
				},
			},
			Components: []discordgo.MessageComponent{},
		})
		pcw.am = nil
	}

	if !pcw.closed {
		close(pcw.fin)
		pcw.subscription()
		pcw.closed = true

		return true
	}

	return false
}

func (c *AuthController) Setup(container di.Container, router fiber.Router) {
	c.discordOAuth = container.Get(static.DiDiscordOAuthModule).(*discordoauth.DiscordOAuth)
	c.rth = container.Get(static.DiAuthRefreshTokenHandler).(auth.RefreshTokenHandler)
	c.ath = container.Get(static.DiAuthAccessTokenHandler).(auth.AccessTokenHandler)
	c.authMw = container.Get(static.DiAuthMiddleware).(auth.Middleware)
	c.st = container.Get(static.DiState).(*dgrs.State)
	c.session = container.Get(static.DiDiscordSession).(*discordgo.Session)
	c.oauthHandler = container.Get(static.DiOAuthHandler).(auth.RequestHandler)
	c.cmdHandler = container.Get(static.DiCommandHandler).(*ken.Ken)

	c.pushcodeSubs = timedmap.New(10 * time.Second)

	router.Get("/login", c.getLogin)
	router.Get("/oauthcallback", c.discordOAuth.HandlerCallback)
	router.Post("/accesstoken", c.postAccessToken)
	router.Post("/pushcode", c.pushCode)
	router.Get("/check", c.authMw.Handle, c.getCheck)
	router.Post("/logout", c.authMw.Handle, c.postLogout)
}

func (c *AuthController) getLogin(ctx *fiber.Ctx) error {
	state := make(map[string]string)

	if redirect := ctx.Query("redirect"); redirect != "" {
		state["redirect"] = redirect
	}

	return c.discordOAuth.HandlerInitWithState(ctx, state)
}

// @Summary Access Token Exchange
// @Description Exchanges the cookie-passed refresh token with a generated access token.
// @Tags Authorization
// @Accept json
// @Produce json
// @Success 200 {object} models.AccessTokenResponse
// @Failure 401 {object} models.Error
// @Router /auth/accesstoken [post]
func (c *AuthController) postAccessToken(ctx *fiber.Ctx) error {
	refreshToken := ctx.Cookies(static.RefreshTokenCookieName)
	if refreshToken == "" {
		return fiber.ErrUnauthorized
	}

	ident, err := c.rth.ValidateRefreshToken(refreshToken)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		ctlLog.Error().Err(err).Msg("Failed validating refresh token")
	}
	if ident == "" {
		return fiber.ErrUnauthorized
	}

	token, expires, err := c.ath.GetAccessToken(ident)
	if err != nil {
		return err
	}

	return ctx.JSON(&models.AccessTokenResponse{
		Token:   token,
		Expires: expires,
	})
}

// @Summary Authorization Check
// @Description Returns OK if the request is authorized.
// @Tags Authorization
// @Accept json
// @Produce json
// @Success 200 {object} models.Status
// @Failure 401 {object} models.Error
// @Router /auth/check [get]
func (c *AuthController) getCheck(ctx *fiber.Ctx) error {
	return ctx.JSON(models.Ok)
}

// @Summary Logout
// @Description Reovkes the currently used access token and clears the refresh token.
// @Tags Authorization
// @Accept json
// @Produce json
// @Success 200 {object} models.Status
// @Router /auth/logout [post]
func (c *AuthController) postLogout(ctx *fiber.Ctx) error {
	uid := ctx.Locals("uid").(string)

	err := c.rth.RevokeToken(uid)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	ctx.ClearCookie(static.RefreshTokenCookieName)

	return ctx.JSON(models.Ok)
}

// @Summary Pushcode
// @Description Send a login push code resulting in a long-fetch request waiting for the code to be sent to shinpurus DMs.
// @Tags Authorization
// @Accept json
// @Produce json
// @Param payload body models.PushCodeRequest true "The push code."
// @Success 200 {object} models.Status
// @Success 400 {object} models.Status
// @Success 410 {object} models.Status
// @Router /auth/pushcode [post]
func (c *AuthController) pushCode(ctx *fiber.Ctx) (err error) {
	var req models.PushCodeRequest
	if err = ctx.BodyParser(&req); err != nil {
		return
	}

	if req.Code == "" {
		return fiber.NewError(fiber.StatusBadRequest, "empty code")
	}

	ipaddr := ctx.IP()
	if ipaddr == "" {
		// When the IP address is empty, which might happen, just
		// generate a new pcw for each request to avoid conflicts.
		ipaddr = xid.New().String()
	}

	pcw, ok := c.pushcodeSubs.GetValue(ipaddr).(*pushCodeWaiter)
	if !ok {
		pcw = new(pushCodeWaiter)
		c.pushcodeSubs.Set(ipaddr, pcw, pushcodeTimeout, func(_ any) {
			pcw.close()
		})

		pcw.code = req.Code
		pcw.fin = make(chan *discordgo.Message)
		pcw.subscription = c.st.Subscribe("dms", func(scan func(v any) error) {
			var msg discordgo.Message
			if err = scan(&msg); err != nil {
				ctlLog.Error().Err(err).Msg("failed scanning message from 'dms' event bus")
				return
			}
			if msg.Content == pcw.code && msg.Author != nil {
				am, err := acceptmsg.New().
					WithKen(c.cmdHandler).
					DeleteAfterAnswer().WithEmbed(&discordgo.MessageEmbed{
					Title: "Login",
					Description: "Do you really want to log in to the web interface using this " +
						"login code?\n\n⚠️ **Never __ever__ enter a login code here you got from someone else!**\n" +
						"If you got this login code from someone else, press `Cancel` or do nothing!",
					Color: static.ColorEmbedOrange,
				}).WithAcceptButton(discordgo.Button{
					Label: "Accept",
					Style: discordgo.SuccessButton,
				}).WithDeclineButton(discordgo.Button{
					Label: "Cancel",
					Style: discordgo.DangerButton,
				}).DoOnAccept(func(ctx ken.ComponentContext) error {
					pcw.am = nil
					pcw.fin <- &msg
					return nil
				}).Send(msg.ChannelID)
				if err == nil {
					pcw.am = am
				}
			}
		})
	} else {
		log.Debug().Field("ipaddr", ipaddr).Msg("Reusing pushcode handler for this client")
		pcw.code = req.Code
	}

	res := <-pcw.fin
	if res == nil {
		err = fiber.NewError(fiber.StatusGone, "timeout")
		return
	}

	c.pushcodeSubs.Remove(ipaddr)
	if pcw.close() {
		util.SendEmbed(c.session, res.ChannelID,
			"You are now being logged in!", "", static.ColorEmbedGreen)
	}

	err = c.oauthHandler.BindRefreshToken(ctx, res.Author.ID)
	if err != nil {
		return
	}

	return ctx.JSON(models.Ok)
}
