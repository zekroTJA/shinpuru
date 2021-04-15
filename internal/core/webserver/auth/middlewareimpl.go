package auth

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/util/static"
)

var (
	errInvalidAccessToken = fiber.NewError(fiber.StatusUnauthorized, "invalid access token")
)

type MiddlewareImpl struct {
	ath AccessTokenHandler
}

func NewMiddlewareImpl(container di.Container) *MiddlewareImpl {
	return &MiddlewareImpl{
		ath: container.Get(static.DiAuthAccessTokenHandler).(AccessTokenHandler),
	}
}

func (m *MiddlewareImpl) Handle(ctx *fiber.Ctx) (err error) {
	authHeader := ctx.Get("authorization")
	if authHeader == "" {
		return fiber.ErrUnauthorized
	}

	split := strings.Split(authHeader, " ")
	if len(split) < 2 {
		return fiber.ErrUnauthorized
	}

	var ident string
	switch strings.ToLower(split[0]) {

	case "accesstoken":
		if ident, err = m.ath.ValidateAccessToken(split[1]); err != nil {
			return errInvalidAccessToken
		}

	default:
		return fiber.ErrUnauthorized

	}

	ctx.Locals("uid", ident)

	return
}
