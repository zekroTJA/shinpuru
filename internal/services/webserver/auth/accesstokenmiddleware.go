package auth

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/onetimeauth/v2"
)

var (
	errInvalidAccessToken = fiber.NewError(fiber.StatusUnauthorized, "invalid access token")
)

// AccessTokenMiddleware implements Middleware
// for Access and API tokens passed via an
// Authorization header
type AccessTokenMiddleware struct {
	ath   AccessTokenHandler
	apith APITokenHandler
	ota   onetimeauth.OneTimeAuth
}

// NewAccessTokenMiddleware initializes a new instance
// of AccessTokenMiddleware.
func NewAccessTokenMiddleware(container di.Container) *AccessTokenMiddleware {
	return &AccessTokenMiddleware{
		ath:   container.Get(static.DiAuthAccessTokenHandler).(AccessTokenHandler),
		apith: container.Get(static.DiAuthAPITokenHandler).(APITokenHandler),
		ota:   container.Get(static.DiOneTimeAuth).(onetimeauth.OneTimeAuth),
	}
}

func (m *AccessTokenMiddleware) Handle(ctx *fiber.Ctx) (err error) {
	ident, err := m.checkOta(ctx)
	if err != nil {
		return
	}
	if ident != "" {
		return next(ctx, ident)
	}

	authHeader := ctx.Get("authorization")
	if authHeader == "" {
		return errInvalidAccessToken
	}

	split := strings.Split(authHeader, " ")
	if len(split) < 2 {
		return errInvalidAccessToken
	}

	switch strings.ToLower(split[0]) {

	case "accesstoken":
		if ident, err = m.ath.ValidateAccessToken(split[1]); err != nil || ident == "" {
			return errInvalidAccessToken
		}

	case "bearer":
		if ident, err = m.apith.ValidateAPIToken(split[1]); err != nil || ident == "" {
			return fiber.ErrUnauthorized
		}

	default:
		return fiber.ErrUnauthorized
	}

	return next(ctx, ident)
}

func (m *AccessTokenMiddleware) checkOta(ctx *fiber.Ctx) (ident string, err error) {
	token := ctx.Query("ota_token")
	if token == "" {
		return
	}

	ident, err = m.ota.ValidateKey(token, ctx.Path())
	if err != nil {
		err = fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}
	return
}

func next(ctx *fiber.Ctx, ident string) error {
	ctx.Locals("uid", ident)
	return ctx.Next()
}
