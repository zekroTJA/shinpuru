package permissions

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zekroTJA/shinpuru/pkg/permissions"
	"github.com/zekrotja/ken"
)

type Provider interface {
	ken.MiddlewareBefore

	HandleWs(s Session, required string) fiber.Handler

	// GetPermissions tries to fetch the permissions array of
	// the passed user of the specified guild. The merged
	// permissions array is returned as well as the override,
	// which is true when the specified user is the bot owner,
	// guild owner or an admin of the guild.
	GetPermissions(
		s Session,
		guildID string,
		userID string,
	) (perm permissions.PermissionArray, overrideExplicits bool, err error)

	// CheckPermissions tries to fetch the permissions of the specified user
	// on the specified guild and returns true, if any of the passed dns match
	// the fetched permissions array. Also, the override status is returned as
	// well as errors occured during permissions fetching.
	//
	// If the userID matches the configured bot owner, all bot owner permissions
	// will be added to the fetched permissions array.
	//
	// If guildID is passed as non-mepty string, all configured guild owner
	// permissions will be added to the fetched permissions array as well.
	CheckPermissions(
		s Session,
		guildID string,
		userID string,
		dns ...string,
	) (bool, bool, error)

	// GetMemberPermissions returns a PermissionsArray based on the passed
	// members roles permissions rulesets for the given guild.
	GetMemberPermission(
		s Session,
		guildID string,
		memberID string,
	) (permissions.PermissionArray, error)

	// CheckSubPerm takes a command context and checks is the given
	// subDN is permitted.
	CheckSubPerm(
		ctx ken.Context,
		subDN string,
		explicit bool,
		message ...string,
	) (ok bool, err error)
}
