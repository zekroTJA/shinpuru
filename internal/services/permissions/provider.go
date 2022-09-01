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

	// GetPermissions tries to fetch the permissions array of
	// the passed user of the specified guild. The merged
	// permissions array is returned as well as the override,
	// which is true when the specified user is the bot owner,
	// guild owner or an admin of the guild.
	GetPermissions(s discordutil.ISession, guildID, userID string) (perm permissions.PermissionArray, overrideExplicits bool, err error)

	// CheckPermissions tries to fetch the permissions of the specified user
	// on the specified guild and returns true, if the passed dn matches the
	// fetched permissions array. Also, the override status is returned as
	// well as errors occured during permissions fetching.
	CheckPermissions(s discordutil.ISession, guildID, userID, dn string) (bool, bool, error)

	// GetMemberPermissions returns a PermissionsArray based on the passed
	// members roles permissions rulesets for the given guild.
	GetMemberPermission(s discordutil.ISession, guildID string, memberID string) (permissions.PermissionArray, error)

	// CheckSubPerm takes a command context and checks is the given
	// subDN is permitted.
	CheckSubPerm(ctx ken.Context, subDN string, explicit bool, message ...string) (ok bool, err error)
}
