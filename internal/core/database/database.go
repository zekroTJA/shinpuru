package database

import (
	"errors"

	"github.com/bwmarrin/snowflake"
	"github.com/zekroTJA/shinpuru/internal/core/backup/backupmodels"
	"github.com/zekroTJA/shinpuru/internal/shared/models"
	"github.com/zekroTJA/shinpuru/internal/util/imgstore"
	"github.com/zekroTJA/shinpuru/internal/util/report"
	"github.com/zekroTJA/shinpuru/internal/util/tag"
	"github.com/zekroTJA/shinpuru/internal/util/vote"
	"github.com/zekroTJA/shinpuru/pkg/permissions"
	"github.com/zekroTJA/shinpuru/pkg/twitchnotify"
)

// ErrDatabaseNotFound is returned when no value was
// found in the database for the specified request.
var ErrDatabaseNotFound = errors.New("value not found")

// Database describes functionalities of a database
// driver.
type Database interface {
	//////////////////////////////////////////////////////
	//// INITIALIZATION

	Connect(credentials ...interface{}) error
	Close()

	//////////////////////////////////////////////////////
	//// GUILD SETTINGS

	GetSetting(setting string) (string, error)
	SetSetting(setting, value string) error

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

	GetGuildColorReaction(guildID string) (bool, error)
	SetGuildColorReaction(guildID string, enable bool) error

	GetGuildMuteRole(guildID string) (string, error)
	SetGuildMuteRole(guildID, roleID string) error

	//////////////////////////////////////////////////////
	//// REPORTS

	AddReport(rep *report.Report) error
	DeleteReport(id snowflake.ID) error
	GetReport(id snowflake.ID) (*report.Report, error)
	GetReportsGuild(guildID string, offset, limit int) ([]*report.Report, error)
	GetReportsFiltered(guildID, memberID string, repType int) ([]*report.Report, error)
	GetReportsGuildCount(guildID string) (int, error)
	GetReportsFilteredCount(guildID, memberID string, repType int) (int, error)

	//////////////////////////////////////////////////////
	//// VOTES

	GetVotes() (map[string]*vote.Vote, error)
	AddUpdateVote(votes *vote.Vote) error
	DeleteVote(voteID string) error

	//////////////////////////////////////////////////////
	//// TWITCHNOTIFY

	GetAllTwitchNotifies(twitchUserID string) ([]*twitchnotify.DBEntry, error)
	GetTwitchNotify(twitchUserID, guildID string) (*twitchnotify.DBEntry, error)
	SetTwitchNotify(twitchNotify *twitchnotify.DBEntry) error
	DeleteTwitchNotify(twitchUserID, guildID string) error

	//////////////////////////////////////////////////////
	//// GUILD BACKUPS

	AddBackup(guildID, fileID string) error
	DeleteBackup(guildID, fileID string) error
	GetBackups(guildID string) ([]*backupmodels.Entry, error)
	GetGuilds() ([]string, error)

	//////////////////////////////////////////////////////
	//// TAGS

	AddTag(tag *tag.Tag) error
	EditTag(tag *tag.Tag) error
	GetTagByID(id snowflake.ID) (*tag.Tag, error)
	GetTagByIdent(ident string, guildID string) (*tag.Tag, error)
	GetGuildTags(guildID string) ([]*tag.Tag, error)
	DeleteTag(id snowflake.ID) error

	//////////////////////////////////////////////////////
	//// API TOKEN

	SetAPIToken(token *models.APITokenEntry) error
	GetAPIToken(userID string) (*models.APITokenEntry, error)
	DeleteAPIToken(userID string) error

	//////////////////////////////////////////////////////
	//// KARMA

	GetKarma(userID, guildID string) (int, error)
	GetKarmaSum(userID string) (int, error)
	GetKarmaGuild(guildID string, limit int) ([]*models.GuildKarma, error)
	SetKarma(userID, guildID string, val int) error
	UpdateKarma(userID, guildID string, diff int) error

	SetKarmaState(guildID string, state bool) error
	GetKarmaState(guildID string) (bool, error)

	SetKarmaEmotes(guildID, emotesInc, emotesDec string) error
	GetKarmaEmotes(guildID string) (emotesInc, emotesDec string, err error)

	SetKarmaTokens(guildID string, tokens int) error
	GetKarmaTokens(guildID string) (int, error)

	//////////////////////////////////////////////////////
	//// CHAN LOCK

	SetLockChan(chanID, guildID, executorID, permissions string) error
	GetLockChan(chanID string) (guildID, executorID, permissions string, err error)
	GetLockChannels(guildID string) (chanIDs []string, err error)
	DeleteLockChan(chanID string) error

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
