package permissions

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"

	"github.com/zekroTJA/shinpuru/internal/services/config"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
	"github.com/zekroTJA/shinpuru/pkg/permissions"
	"github.com/zekroTJA/shinpuru/pkg/roleutil"
	"github.com/zekroTJA/shinpuru/pkg/stringutil"
	"github.com/zekrotja/dgrs"
	"github.com/zekrotja/ken"
)

// Permissions is a command handler middleware
// processing permissions for command execution.
//
// Implements the shireikan.Middleware interface and
// exposes functions to check permissions.
type Permissions struct {
	db  database.Database
	cfg config.Provider
	st  *dgrs.State
}

var _ ken.MiddlewareBefore = (*Permissions)(nil)

// NewPermissions returns a new PermissionsMiddleware
// instance with the passed database and config instances.
func NewPermissions(container di.Container) *Permissions {
	return &Permissions{
		db:  container.Get(static.DiDatabase).(database.Database),
		cfg: container.Get(static.DiConfig).(config.Provider),
		st:  container.Get(static.DiState).(*dgrs.State),
	}
}

func (m *Permissions) Before(ctx *ken.Ctx) (next bool, err error) {
	cmd, ok := ctx.Command.(PermCommand)
	if !ok {
		next = true
		return
	}

	if m.db == nil {
		m.db, _ = ctx.Get(static.DiDatabase).(database.Database)
	}

	if m.cfg == nil {
		m.cfg, _ = ctx.Get(static.DiConfig).(config.Provider)
	}

	if ctx.User() == nil {
		return
	}

	ok, _, err = m.CheckPermissions(ctx.Session, ctx.Event.GuildID, ctx.User().ID, cmd.Domain())

	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return false, err
	}

	if !ok {
		err = ctx.FollowUpError("You are not permitted to use this command!", "Missing Permission").
			DeleteAfter(8 * time.Second).Error
		return
	}

	next = true
	return
}

func (pmw *Permissions) HandleWs(s *discordgo.Session, required string) fiber.Handler {
	if !stringutil.ContainsAny(required, static.AdditionalPermissions) {
		static.AdditionalPermissions = append(static.AdditionalPermissions, required)
	}

	return func(ctx *fiber.Ctx) error {
		uid, _ := ctx.Locals("uid").(string)
		guildID := ctx.Params("guildid")

		if uid == "" {
			return fiber.ErrForbidden
		}

		ok, _, err := pmw.CheckPermissions(s, guildID, uid, required)
		if err != nil {
			return err
		}
		if !ok {
			return fiber.ErrForbidden
		}

		return ctx.Next()
	}
}

// GetPermissions tries to fetch the permissions array of
// the passed user of the specified guild. The merged
// permissions array is returned as well as the override,
// which is true when the specified user is the bot owner,
// guild owner or an admin of the guild.
func (m *Permissions) GetPermissions(s *discordgo.Session, guildID, userID string) (perm permissions.PermissionArray, overrideExplicits bool, err error) {
	if guildID != "" {
		perm, err = m.GetMemberPermission(s, guildID, userID)
		if err != nil && !database.IsErrDatabaseNotFound(err) {
			return
		}
	} else {
		perm = make(permissions.PermissionArray, 0)
	}

	if m.cfg.Config().Discord.OwnerID == userID {
		perm = perm.Merge(permissions.PermissionArray{"+sp.*"}, false)
		overrideExplicits = true
	}

	if guildID != "" {
		guild, err := m.st.Guild(guildID, true)
		if err != nil {
			return permissions.PermissionArray{}, false, nil
		}

		member, _ := s.GuildMember(guildID, userID)

		if userID == guild.OwnerID || (member != nil && discordutil.IsAdmin(guild, member)) {
			var defAdminRoles []string
			defAdminRoles = m.cfg.Config().Permissions.DefaultAdminRules
			if defAdminRoles == nil {
				defAdminRoles = static.DefaultAdminRules
			}

			perm = perm.Merge(defAdminRoles, false)
			overrideExplicits = true
		}
	}

	var defUserRoles []string
	defUserRoles = m.cfg.Config().Permissions.DefaultUserRules
	if defUserRoles == nil {
		defUserRoles = static.DefaultUserRules
	}

	perm = perm.Merge(defUserRoles, false)

	return perm, overrideExplicits, nil
}

// CheckPermissions tries to fetch the permissions of the specified user
// on the specified guild and returns true, if the passed dn matches the
// fetched permissions array. Also, the override status is returned as
// well as errors occured during permissions fetching.
func (m *Permissions) CheckPermissions(s *discordgo.Session, guildID, userID, dn string) (bool, bool, error) {
	perms, overrideExplicits, err := m.GetPermissions(s, guildID, userID)
	if err != nil {
		return false, false, err
	}

	return perms.Check(dn), overrideExplicits, nil
}

// GetMemberPermissions returns a PermissionsArray based on the passed
// members roles permissions rulesets for the given guild.
func (m *Permissions) GetMemberPermission(s *discordgo.Session, guildID string, memberID string) (permissions.PermissionArray, error) {
	guildPerms, err := m.db.GetGuildPermissions(guildID)
	if err != nil {
		return nil, err
	}

	membRoles, err := roleutil.GetSortedMemberRoles(s, guildID, memberID, false, true)
	if err != nil {
		return nil, err
	}

	var res permissions.PermissionArray
	for _, r := range membRoles {
		if p, ok := guildPerms[r.ID]; ok {
			if res == nil {
				res = p
			} else {
				res = res.Merge(p, true)
			}
		}
	}

	return res, nil
}
