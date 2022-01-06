package database

import (
	"errors"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/zekroTJA/shinpuru/internal/models"
	"github.com/zekroTJA/shinpuru/internal/services/backup/backupmodels"
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

	GetGuildAutoRole(guildID string) ([]string, error)
	SetGuildAutoRole(guildID string, autoRoleIDs []string) error

	GetGuildModLog(guildID string) (string, error)
	SetGuildModLog(guildID, chanID string) error

	GetGuildVoiceLog(guildID string) (string, error)
	SetGuildVoiceLog(guildID, chanID string) error

	GetGuildVoiceLogIgnores(guildID string) ([]string, error)
	IsGuildVoiceLogIgnored(guildID, channelID string) (bool, error)
	SetGuildVoiceLogIngore(guildID, channelID string) error
	RemoveGuildVoiceLogIgnore(guildID, channelID string) error

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

	GetGuildLogDisable(guildID string) (bool, error)
	SetGuildLogDisable(guildID string, enabled bool) error

	GetGuildAPI(guildID string) (*models.GuildAPISettings, error)
	SetGuildAPI(guildID string, settings *models.GuildAPISettings) error

	GetGuildVerificationRequired(guildID string) (bool, error)
	SetGuildVerificationRequired(guildID string, enable bool) error

	//////////////////////////////////////////////////////
	//// USER SETTINGS

	GetUserOTAEnabled(userID string) (bool, error)
	SetUserOTAEnabled(userID string, enabled bool) error

	GetUserVerified(userID string) (bool, error)
	SetUserVerified(userID string, enabled bool) error

	GetUserByRefreshToken(token string) (string, time.Time, error)
	SetUserRefreshToken(userID, token string, expires time.Time) error
	RevokeUserRefreshToken(userID string) error
	CleanupExpiredRefreshTokens() (int64, error)

	//////////////////////////////////////////////////////
	//// REPORTS

	AddReport(rep *models.Report) error
	DeleteReport(id snowflake.ID) error
	GetReport(id snowflake.ID) (*models.Report, error)
	GetReportsGuild(guildID string, offset, limit int) ([]*models.Report, error)
	GetReportsFiltered(guildID, memberID string, repType, offset, limit int) ([]*models.Report, error)
	GetReportsGuildCount(guildID string) (int, error)
	GetReportsFilteredCount(guildID, memberID string, repType int) (int, error)
	GetExpiredReports() ([]*models.Report, error)
	ExpireReports(id ...string) (err error)

	//////////////////////////////////////////////////////
	//// UNBAN REQUESTS

	GetGuildUnbanRequests(guildID string) ([]*models.UnbanRequest, error)
	GetGuildUserUnbanRequests(userID, guildID string) ([]*models.UnbanRequest, error)
	GetUnbanRequest(id string) (*models.UnbanRequest, error)
	AddUnbanRequest(request *models.UnbanRequest) error
	UpdateUnbanRequest(request *models.UnbanRequest) error

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

	SetKarmaPenalty(guildID string, state bool) error
	GetKarmaPenalty(guildID string) (bool, error)

	GetKarmaBlockList(guildID string) ([]string, error)
	IsKarmaBlockListed(guildID, userID string) (bool, error)
	AddKarmaBlockList(guildID, userID string) error
	RemoveKarmaBlockList(guildID, userID string) error

	GetKarmaRules(guildID string) ([]*models.KarmaRule, error)
	CheckKarmaRule(guildID, checksum string) (ok bool, err error)
	AddOrUpdateKarmaRule(rule *models.KarmaRule) error
	RemoveKarmaRule(guildID string, id snowflake.ID) error

	//////////////////////////////////////////////////////
	//// CHAN LOCK

	SetLockChan(chanID, guildID, executorID, permissions string) error
	GetLockChan(chanID string) (guildID, executorID, permissions string, err error)
	GetLockChannels(guildID string) (chanIDs []string, err error)
	DeleteLockChan(chanID string) error

	//////////////////////////////////////////////////////
	//// ANTI RAID

	SetAntiraidState(guildID string, state bool) error
	GetAntiraidState(guildID string) (bool, error)

	SetAntiraidRegeneration(guildID string, periodSecs int) error
	GetAntiraidRegeneration(guildID string) (int, error)

	SetAntiraidBurst(guildID string, burst int) error
	GetAntiraidBurst(guildID string) (int, error)

	AddToAntiraidJoinList(guildID, userID, userTag string, accountCreated time.Time) error
	GetAntiraidJoinList(guildID string) ([]*models.JoinLogEntry, error)
	FlushAntiraidJoinList(guildID string) error
	RemoveAntiraidJoinList(guildID, userID string) error

	//////////////////////////////////////////////////////
	//// STARBOARD

	SetStarboardConfig(config *models.StarboardConfig) error
	GetStarboardConfig(guildID string) (*models.StarboardConfig, error)
	SetStarboardEntry(e *models.StarboardEntry) (err error)
	RemoveStarboardEntry(msgID string) error
	GetStarboardEntries(guildID string, sortBy models.StarboardSortBy, limit, offset int) ([]*models.StarboardEntry, error)
	GetStarboardEntry(messageID string) (*models.StarboardEntry, error)

	//////////////////////////////////////////////////////
	//// GUILDLOG

	GetGuildLogEntries(guildID string, offset, limit int, severity models.GuildLogSeverity) ([]*models.GuildLogEntry, error)
	GetGuildLogEntriesCount(guildID string, severity models.GuildLogSeverity) (int, error)
	AddGuildLogEntry(entry *models.GuildLogEntry) error
	DeleteLogEntry(guildID string, id snowflake.ID) error
	DeleteLogEntries(guildID string) error

	//////////////////////////////////////////////////////
	//// FUNCTIONALITIES

	FlushGuildData(guildID string) error

	//////////////////////////////////////////////////////
	//// VERIFICATION QUEUE

	GetVerificationQueue(guildID string) ([]*models.VerificationQueueEntry, error)
	FlushVerificationQueue(guildID string) error
	AddVerificationQueue(e *models.VerificationQueueEntry) error
	RemoveVerificationQueue(guildID, userID string) error
}

// IsErrDatabaseNotFound returns true if the passed err
// is an ErrDatabaseNotFound.
func IsErrDatabaseNotFound(err error) bool {
	return err == ErrDatabaseNotFound
}
