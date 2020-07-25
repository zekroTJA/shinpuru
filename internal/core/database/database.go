package database

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/snowflake"
	"github.com/zekroTJA/shinpuru/internal/core/backup/backupmodels"
	"github.com/zekroTJA/shinpuru/internal/core/permissions"
	"github.com/zekroTJA/shinpuru/internal/core/twitchnotify"
	"github.com/zekroTJA/shinpuru/internal/shared/models"
	"github.com/zekroTJA/shinpuru/internal/util/imgstore"
	"github.com/zekroTJA/shinpuru/internal/util/report"
	"github.com/zekroTJA/shinpuru/internal/util/tag"
	"github.com/zekroTJA/shinpuru/internal/util/vote"
)

// ErrDatabaseNotFound is returned when no value was
// found in the database for the specified request.
var ErrDatabaseNotFound = errors.New("value not found")

// Database describes functionalities of a database
// driver.
type Database interface {
	Connect(credentials ...interface{}) error
	Close()

	GetGuildPrefix(guildID string) (string, error)
	SetGuildPrefix(guildID, newPrefix string) error

	GetGuildAutoRole(guildID string) (string, error)
	SetGuildAutoRole(guildID, autoRoleID string) error

	GetGuildModLog(guildID string) (string, error)
	SetGuildModLog(guildID, chanID string) error

	GetGuildVoiceLog(guildID string) (string, error)
	SetGuildVoiceLog(guildID, chanID string) error

	GetGuildNotifyRole(guildID string) (string, error)
	SetGuildNotifyRole(guildID, roleID string) error

	GetGuildGhostpingMsg(guildID string) (string, error)
	SetGuildGhostpingMsg(guildID, msg string) error

	GetGuildPermissions(guildID string) (map[string]permissions.PermissionArray, error)
	SetGuildRolePermission(guildID, roleID string, p permissions.PermissionArray) error
	GetMemberPermission(s *discordgo.Session, guildID string, memberID string) (permissions.PermissionArray, error)

	GetGuildJdoodleKey(guildID string) (string, error)
	SetGuildJdoodleKey(guildID, key string) error

	GetGuildBackup(guildID string) (bool, error)
	SetGuildBackup(guildID string, enabled bool) error

	GetGuildInviteBlock(guildID string) (string, error)
	SetGuildInviteBlock(guildID string, data string) error

	GetGuildJoinMsg(guildID string) (string, string, error)
	SetGuildJoinMsg(guildID string, channelID string, msg string) error

	GetGuildLeaveMsg(guildID string) (string, string, error)
	SetGuildLeaveMsg(guildID string, channelID string, msg string) error

	AddReport(rep *report.Report) error
	DeleteReport(id snowflake.ID) error
	GetReport(id snowflake.ID) (*report.Report, error)
	GetReportsGuild(guildID string, offset, limit int) ([]*report.Report, error)
	GetReportsFiltered(guildID, memberID string, repType int) ([]*report.Report, error)
	GetReportsGuildCount(guildID string) (int, error)
	GetReportsFilteredCount(guildID, memberID string, repType int) (int, error)

	GetSetting(setting string) (string, error)
	SetSetting(setting, value string) error

	GetVotes() (map[string]*vote.Vote, error)

	AddUpdateVote(votes *vote.Vote) error
	DeleteVote(voteID string) error

	GetMuteRoles() (map[string]string, error)
	GetMuteRoleGuild(guildID string) (string, error)
	SetMuteRole(guildID, roleID string) error

	GetAllTwitchNotifies(twitchUserID string) ([]*twitchnotify.DBEntry, error)
	GetTwitchNotify(twitchUserID, guildID string) (*twitchnotify.DBEntry, error)
	SetTwitchNotify(twitchNotify *twitchnotify.DBEntry) error
	DeleteTwitchNotify(twitchUserID, guildID string) error

	AddBackup(guildID, fileID string) error
	DeleteBackup(guildID, fileID string) error
	GetBackups(guildID string) ([]*backupmodels.Entry, error)
	GetGuilds() ([]string, error)

	AddTag(tag *tag.Tag) error
	EditTag(tag *tag.Tag) error
	GetTagByID(id snowflake.ID) (*tag.Tag, error)
	GetTagByIdent(ident string, guildID string) (*tag.Tag, error)
	GetGuildTags(guildID string) ([]*tag.Tag, error)
	DeleteTag(id snowflake.ID) error

	SetAPIToken(token *models.APITokenEntry) error
	GetAPIToken(userID string) (*models.APITokenEntry, error)
	DeleteAPIToken(userID string) error

	GetKarma(userID, guildID string) (int, error)
	GetKarmaSum(userID string) (int, error)
	SetKarma(userID, guildID string, val int) error
	UpdateKarma(userID, guildID string, diff int) error

	// Deprecated
	GetImageData(id snowflake.ID) (*imgstore.Image, error)
	// Deprecated
	SaveImageData(image *imgstore.Image) error
	// Deprecated
	RemoveImageData(id snowflake.ID) error
}

// IsErrDatabaseNotFound returns true if the passed err
// is an ErrDatabaseNotFound.
func IsErrDatabaseNotFound(err error) bool {
	return err == ErrDatabaseNotFound
}
