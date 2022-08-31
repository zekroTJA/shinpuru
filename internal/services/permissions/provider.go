package permissions

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
	"github.com/zekroTJA/shinpuru/pkg/permissions"
	"github.com/zekrotja/ken"
)

type Provider interface {
	ken.MiddlewareBefore
	HandleWs(s discordutil.ISession, required string) fiber.Handler
	GetPermissions(s discordutil.ISession, guildID, userID string) (perm permissions.PermissionArray, overrideExplicits bool, err error)
	CheckPermissions(s discordutil.ISession, guildID, userID, dn string) (bool, bool, error)
	GetMemberPermission(s discordutil.ISession, guildID string, memberID string) (permissions.PermissionArray, error)
	CheckSubPerm(ctx ken.Context, subDN string, explicit bool, message ...string) (ok bool, err error)
}
