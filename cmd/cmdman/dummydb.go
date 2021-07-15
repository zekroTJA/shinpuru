package main

import (
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/zekroTJA/shinpuru/internal/models"
	"github.com/zekroTJA/shinpuru/internal/services/backup/backupmodels"
	"github.com/zekroTJA/shinpuru/internal/util/tag"
	"github.com/zekroTJA/shinpuru/internal/util/vote"
	"github.com/zekroTJA/shinpuru/pkg/permissions"
	"github.com/zekroTJA/shinpuru/pkg/twitchnotify"
)

type dummyDB struct{}

func (db *dummyDB) Connect(credentials ...interface{}) (_ error) {
	return
}
func (db *dummyDB) Close() {
	return
}
func (db *dummyDB) GetSetting(setting string) (_ string, _ error) {
	return
}
func (db *dummyDB) SetSetting(setting, value string) (_ error) {
	return
}
func (db *dummyDB) GetGuildPrefix(guildID string) (_ string, _ error) {
	return
}
func (db *dummyDB) SetGuildPrefix(guildID, newPrefix string) (_ error) {
	return
}
func (db *dummyDB) GetGuildAutoRole(guildID string) (_ string, _ error) {
	return
}
func (db *dummyDB) SetGuildAutoRole(guildID, autoRoleID string) (_ error) {
	return
}
func (db *dummyDB) GetGuildModLog(guildID string) (_ string, _ error) {
	return
}
func (db *dummyDB) SetGuildModLog(guildID, chanID string) (_ error) {
	return
}
func (db *dummyDB) GetGuildVoiceLog(guildID string) (_ string, _ error) {
	return
}
func (db *dummyDB) SetGuildVoiceLog(guildID, chanID string) (_ error) {
	return
}
func (db *dummyDB) GetGuildVoiceLogIgnores(guildID string) (_ []string, _ error) {
	return
}
func (db *dummyDB) IsGuildVoiceLogIgnored(guildID, channelID string) (_ bool, _ error) {
	return
}
func (db *dummyDB) SetGuildVoiceLogIngore(guildID, channelID string) (_ error) {
	return
}
func (db *dummyDB) RemoveGuildVoiceLogIgnore(guildID, channelID string) (_ error) {
	return
}
func (db *dummyDB) GetGuildNotifyRole(guildID string) (_ string, _ error) {
	return
}
func (db *dummyDB) SetGuildNotifyRole(guildID, roleID string) (_ error) {
	return
}
func (db *dummyDB) GetGuildGhostpingMsg(guildID string) (_ string, _ error) {
	return
}
func (db *dummyDB) SetGuildGhostpingMsg(guildID, msg string) (_ error) {
	return
}
func (db *dummyDB) GetGuildPermissions(guildID string) (_ map[string]permissions.PermissionArray, _ error) {
	return
}
func (db *dummyDB) SetGuildRolePermission(guildID, roleID string, p permissions.PermissionArray) (_ error) {
	return
}
func (db *dummyDB) GetGuildJdoodleKey(guildID string) (_ string, _ error) {
	return
}
func (db *dummyDB) SetGuildJdoodleKey(guildID, key string) (_ error) {
	return
}
func (db *dummyDB) GetGuildBackup(guildID string) (_ bool, _ error) {
	return
}
func (db *dummyDB) SetGuildBackup(guildID string, enabled bool) (_ error) {
	return
}
func (db *dummyDB) GetGuildInviteBlock(guildID string) (_ string, _ error) {
	return
}
func (db *dummyDB) SetGuildInviteBlock(guildID string, data string) (_ error) {
	return
}
func (db *dummyDB) GetGuildJoinMsg(guildID string) (_ string, _ string, _ error) {
	return
}
func (db *dummyDB) SetGuildJoinMsg(guildID string, channelID string, msg string) (_ error) {
	return
}
func (db *dummyDB) GetGuildLeaveMsg(guildID string) (_ string, _ string, _ error) {
	return
}
func (db *dummyDB) SetGuildLeaveMsg(guildID string, channelID string, msg string) (_ error) {
	return
}
func (db *dummyDB) GetGuildColorReaction(guildID string) (_ bool, _ error) {
	return
}
func (db *dummyDB) SetGuildColorReaction(guildID string, enable bool) (_ error) {
	return
}
func (db *dummyDB) GetGuildMuteRole(guildID string) (_ string, _ error) {
	return
}
func (db *dummyDB) SetGuildMuteRole(guildID, roleID string) (_ error) {
	return
}
func (db *dummyDB) GetGuildLogDisable(guildID string) (_ bool, _ error) {
	return
}
func (db *dummyDB) SetGuildLogDisable(guildID string, enabled bool) (_ error) {
	return
}
func (db *dummyDB) GetUserOTAEnabled(userID string) (_ bool, _ error) {
	return
}
func (db *dummyDB) SetUserOTAEnabled(userID string, enabled bool) (_ error) {
	return
}
func (db *dummyDB) GetUserByRefreshToken(token string) (_ string, _ time.Time, _ error) {
	return
}
func (db *dummyDB) SetUserRefreshToken(userID, token string, expires time.Time) (_ error) {
	return
}
func (db *dummyDB) RevokeUserRefreshToken(userID string) (_ error) {
	return
}
func (db *dummyDB) CleanupExpiredRefreshTokens() (_ int64, _ error) {
	return
}
func (db *dummyDB) AddReport(rep *models.Report) (_ error) {
	return
}
func (db *dummyDB) DeleteReport(id snowflake.ID) (_ error) {
	return
}
func (db *dummyDB) GetReport(id snowflake.ID) (_ *models.Report, _ error) {
	return
}
func (db *dummyDB) GetReportsGuild(guildID string, offset, limit int) (_ []*models.Report, _ error) {
	return
}
func (db *dummyDB) GetReportsFiltered(guildID, memberID string, repType int) (_ []*models.Report, _ error) {
	return
}
func (db *dummyDB) GetReportsGuildCount(guildID string) (_ int, _ error) {
	return
}
func (db *dummyDB) GetReportsFilteredCount(guildID, memberID string, repType int) (_ int, _ error) {
	return
}
func (db *dummyDB) GetGuildUnbanRequests(guildID string) (_ []*models.UnbanRequest, _ error) {
	return
}
func (db *dummyDB) GetGuildUserUnbanRequests(userID, guildID string) (_ []*models.UnbanRequest, _ error) {
	return
}
func (db *dummyDB) GetUnbanRequest(id string) (_ *models.UnbanRequest, _ error) {
	return
}
func (db *dummyDB) AddUnbanRequest(request *models.UnbanRequest) (_ error) {
	return
}
func (db *dummyDB) UpdateUnbanRequest(request *models.UnbanRequest) (_ error) {
	return
}
func (db *dummyDB) GetVotes() (_ map[string]*vote.Vote, _ error) {
	return
}
func (db *dummyDB) AddUpdateVote(votes *vote.Vote) (_ error) {
	return
}
func (db *dummyDB) DeleteVote(voteID string) (_ error) {
	return
}
func (db *dummyDB) GetAllTwitchNotifies(twitchUserID string) (_ []*twitchnotify.DBEntry, _ error) {
	return
}
func (db *dummyDB) GetTwitchNotify(twitchUserID, guildID string) (_ *twitchnotify.DBEntry, _ error) {
	return
}
func (db *dummyDB) SetTwitchNotify(twitchNotify *twitchnotify.DBEntry) (_ error) {
	return
}
func (db *dummyDB) DeleteTwitchNotify(twitchUserID, guildID string) (_ error) {
	return
}
func (db *dummyDB) AddBackup(guildID, fileID string) (_ error) {
	return
}
func (db *dummyDB) DeleteBackup(guildID, fileID string) (_ error) {
	return
}
func (db *dummyDB) GetBackups(guildID string) (_ []*backupmodels.Entry, _ error) {
	return
}
func (db *dummyDB) GetGuilds() (_ []string, _ error) {
	return
}
func (db *dummyDB) AddTag(tag *tag.Tag) (_ error) {
	return
}
func (db *dummyDB) EditTag(tag *tag.Tag) (_ error) {
	return
}
func (db *dummyDB) GetTagByID(id snowflake.ID) (_ *tag.Tag, _ error) {
	return
}
func (db *dummyDB) GetTagByIdent(ident string, guildID string) (_ *tag.Tag, _ error) {
	return
}
func (db *dummyDB) GetGuildTags(guildID string) (_ []*tag.Tag, _ error) {
	return
}
func (db *dummyDB) DeleteTag(id snowflake.ID) (_ error) {
	return
}
func (db *dummyDB) SetAPIToken(token *models.APITokenEntry) (_ error) {
	return
}
func (db *dummyDB) GetAPIToken(userID string) (_ *models.APITokenEntry, _ error) {
	return
}
func (db *dummyDB) DeleteAPIToken(userID string) (_ error) {
	return
}
func (db *dummyDB) GetKarma(userID, guildID string) (_ int, _ error) {
	return
}
func (db *dummyDB) GetKarmaSum(userID string) (_ int, _ error) {
	return
}
func (db *dummyDB) GetKarmaGuild(guildID string, limit int) (_ []*models.GuildKarma, _ error) {
	return
}
func (db *dummyDB) SetKarma(userID, guildID string, val int) (_ error) {
	return
}
func (db *dummyDB) UpdateKarma(userID, guildID string, diff int) (_ error) {
	return
}
func (db *dummyDB) SetKarmaState(guildID string, state bool) (_ error) {
	return
}
func (db *dummyDB) GetKarmaState(guildID string) (_ bool, _ error) {
	return
}
func (db *dummyDB) SetKarmaEmotes(guildID, emotesInc, emotesDec string) (_ error) {
	return
}
func (db *dummyDB) GetKarmaEmotes(guildID string) (emotesInc, emotesDec string, err error) {
	return
}
func (db *dummyDB) SetKarmaTokens(guildID string, tokens int) (_ error) {
	return
}
func (db *dummyDB) GetKarmaTokens(guildID string) (_ int, _ error) {
	return
}
func (db *dummyDB) SetKarmaPenalty(guildID string, state bool) (_ error) {
	return
}
func (db *dummyDB) GetKarmaPenalty(guildID string) (_ bool, _ error) {
	return
}
func (db *dummyDB) GetKarmaBlockList(guildID string) (_ []string, _ error) {
	return
}
func (db *dummyDB) IsKarmaBlockListed(guildID, userID string) (_ bool, _ error) {
	return
}
func (db *dummyDB) AddKarmaBlockList(guildID, userID string) (_ error) {
	return
}
func (db *dummyDB) RemoveKarmaBlockList(guildID, userID string) (_ error) {
	return
}
func (db *dummyDB) GetKarmaRules(guildID string) (_ []*models.KarmaRule, _ error) {
	return
}
func (db *dummyDB) CheckKarmaRule(guildID, checksum string) (ok bool, err error) {
	return
}
func (db *dummyDB) AddOrUpdateKarmaRule(rule *models.KarmaRule) (_ error) {
	return
}
func (db *dummyDB) RemoveKarmaRule(guildID string, id snowflake.ID) (_ error) {
	return
}
func (db *dummyDB) SetLockChan(chanID, guildID, executorID, permissions string) (_ error) {
	return
}
func (db *dummyDB) GetLockChan(chanID string) (guildID, executorID, permissions string, err error) {
	return
}
func (db *dummyDB) GetLockChannels(guildID string) (chanIDs []string, err error) {
	return
}
func (db *dummyDB) DeleteLockChan(chanID string) (_ error) {
	return
}
func (db *dummyDB) SetAntiraidState(guildID string, state bool) (_ error) {
	return
}
func (db *dummyDB) GetAntiraidState(guildID string) (_ bool, _ error) {
	return
}
func (db *dummyDB) SetAntiraidRegeneration(guildID string, periodSecs int) (_ error) {
	return
}
func (db *dummyDB) GetAntiraidRegeneration(guildID string) (_ int, _ error) {
	return
}
func (db *dummyDB) SetAntiraidBurst(guildID string, burst int) (_ error) {
	return
}
func (db *dummyDB) GetAntiraidBurst(guildID string) (_ int, _ error) {
	return
}
func (db *dummyDB) AddToAntiraidJoinList(guildID, userID, userTag string) (_ error) {
	return
}
func (db *dummyDB) GetAntiraidJoinList(guildID string) (_ []*models.JoinLogEntry, _ error) {
	return
}
func (db *dummyDB) FlushAntiraidJoinList(guildID string) (_ error) {
	return
}
func (db *dummyDB) SetStarboardConfig(config *models.StarboardConfig) (_ error) {
	return
}
func (db *dummyDB) GetStarboardConfig(guildID string) (_ *models.StarboardConfig, _ error) {
	return
}
func (db *dummyDB) SetStarboardEntry(e *models.StarboardEntry) (err error) {
	return
}
func (db *dummyDB) RemoveStarboardEntry(msgID string) (_ error) {
	return
}
func (db *dummyDB) GetStarboardEntries(guildID string, sortBy models.StarboardSortBy, limit, offset int) (_ []*models.StarboardEntry, _ error) {
	return
}
func (db *dummyDB) GetStarboardEntry(messageID string) (_ *models.StarboardEntry, _ error) {
	return
}
func (db *dummyDB) GetGuildLogEntries(guildID string, offset, limit int, severity models.GuildLogSeverity) (_ []*models.GuildLogEntry, _ error) {
	return
}
func (db *dummyDB) GetGuildLogEntriesCount(guildID string, severity models.GuildLogSeverity) (_ int, _ error) {
	return
}
func (db *dummyDB) AddGuildLogEntry(entry *models.GuildLogEntry) (_ error) {
	return
}
func (db *dummyDB) DeleteLogEntry(guildID string, id snowflake.ID) (_ error) {
	return
}
func (db *dummyDB) DeleteLogEntries(guildID string) (_ error) {
	return
}
func (db *dummyDB) FlushGuildData(guildID string) (_ error) {
	return
}
