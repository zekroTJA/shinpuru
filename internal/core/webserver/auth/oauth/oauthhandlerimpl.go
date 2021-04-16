package oauth

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/core/webserver/auth"
	"github.com/zekroTJA/shinpuru/internal/core/webserver/v1/models"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
)

type OAuthHandlerImpl struct {
	session             *discordgo.Session
	accessTokenHandler  auth.AccessTokenHandler
	refreshTokenHandler auth.RefreshTokenHandler
}

func NewOAuthHandlerImpl(container di.Container) *OAuthHandlerImpl {
	return &OAuthHandlerImpl{
		session:             container.Get(static.DiDiscordSession).(*discordgo.Session),
		accessTokenHandler:  container.Get(static.DiAuthAccessTokenHandler).(auth.AccessTokenHandler),
		refreshTokenHandler: container.Get(static.DiAuthRefreshTokenHandler).(auth.RefreshTokenHandler),
	}
}

func (h *OAuthHandlerImpl) LoginFailedHandler(ctx *fiber.Ctx, status int, msg string) error {
	return fiber.ErrUnauthorized
}

func (h *OAuthHandlerImpl) LoginSuccessHandler(ctx *fiber.Ctx, uid string) error {
	user, _ := h.session.User(uid)
	if user == nil {
		return fiber.ErrUnauthorized
	}

	ctx.Locals("uid", uid)
	refreshToken, err := h.refreshTokenHandler.GetRefreshToken(uid)
	if err != nil {
		return err
	}

	expires := time.Now().Add(static.AuthSessionExpiration)
	ctx.Cookie(&fiber.Cookie{
		Name:     static.RefreshTokenCookieName,
		Value:    refreshToken,
		Path:     "/",
		Expires:  expires,
		HTTPOnly: true,
		Secure:   util.IsRelease(),
	})

	return ctx.Redirect("/", fiber.StatusTemporaryRedirect)
}

func (h *OAuthHandlerImpl) LogoutHandler(ctx *fiber.Ctx) error {
	if uid, ok := ctx.Locals("uid").(string); ok && uid != "" {
		if err := h.refreshTokenHandler.RevokeToken(uid); err != nil {
			return err
		}
	}

	ctx.ClearCookie(static.RefreshTokenCookieName)

	return ctx.JSON(models.Ok)
}
