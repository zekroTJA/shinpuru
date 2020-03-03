package database

import (
	"time"

	"github.com/zekroTJA/shinpuru/internal/core/config"

	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/snowflake"
	"github.com/go-redis/redis"
	"github.com/zekroTJA/shinpuru/internal/core/backup/backupmodels"
	"github.com/zekroTJA/shinpuru/internal/core/permissions"
	"github.com/zekroTJA/shinpuru/internal/core/twitchnotify"
	"github.com/zekroTJA/shinpuru/internal/util/imgstore"
	"github.com/zekroTJA/shinpuru/internal/util/report"
	"github.com/zekroTJA/shinpuru/internal/util/tag"
	"github.com/zekroTJA/shinpuru/internal/util/vote"
)

type RedisMiddleware struct {
	client *redis.Client
	db     Database
}

func NewRedisMiddleware(config *config.Config, db Database) *RedisMiddleware {
	r := &RedisMiddleware{
		db: db,
	}

	r.client = redis.NewClient(&redis.Options{
		Addr:     config.Database.Redis.Addr,
		Password: config.Database.Redis.Password,
		DB:       config.Database.Redis.Type,
	})

	return r
}

func (r *RedisMiddleware) Connect(credentials ...interface{}) error {
	return r.db.Connect(credentials...)
}

func (r *RedisMiddleware) Close() {
	r.client.Close()
	r.db.Close()
}

func (r *RedisMiddleware) GetGuildPrefix(guildID string) (string, error) {
	return r.db.GetGuildPrefix(guildID)
}

func (r *RedisMiddleware) SetGuildPrefix(guildID, newPrefix string) error {
	return r.db.SetGuildPrefix(guildID, newPrefix)
}

func (r *RedisMiddleware) GetGuildAutoRole(guildID string) (string, error) {
	return r.db.GetGuildAutoRole(guildID)
}

func (r *RedisMiddleware) SetGuildAutoRole(guildID, autoRoleID string) error {
	return r.db.SetGuildAutoRole(guildID, autoRoleID)
}

func (r *RedisMiddleware) GetGuildModLog(guildID string) (string, error) {
	return r.db.GetGuildModLog(guildID)
}

func (r *RedisMiddleware) SetGuildModLog(guildID, chanID string) error {
	return r.db.SetGuildModLog(guildID, chanID)
}

func (r *RedisMiddleware) GetGuildVoiceLog(guildID string) (string, error) {
	return r.db.GetGuildVoiceLog(guildID)
}

func (r *RedisMiddleware) SetGuildVoiceLog(guildID, chanID string) error {
	return r.db.SetGuildVoiceLog(guildID, chanID)
}

func (r *RedisMiddleware) GetGuildNotifyRole(guildID string) (string, error) {
	return r.db.GetGuildNotifyRole(guildID)
}

func (r *RedisMiddleware) SetGuildNotifyRole(guildID, roleID string) error {
	return r.db.SetGuildNotifyRole(guildID, roleID)
}

func (r *RedisMiddleware) GetGuildGhostpingMsg(guildID string) (string, error) {
	return r.db.GetGuildGhostpingMsg(guildID)
}

func (r *RedisMiddleware) SetGuildGhostpingMsg(guildID, msg string) error {
	return r.db.SetGuildGhostpingMsg(guildID, msg)
}

func (r *RedisMiddleware) GetGuildPermissions(guildID string) (map[string]permissions.PermissionArray, error) {
	return r.db.GetGuildPermissions(guildID)
}

func (r *RedisMiddleware) SetGuildRolePermission(guildID, roleID string, p permissions.PermissionArray) error {
	return r.db.SetGuildRolePermission(guildID, roleID, p)
}

func (r *RedisMiddleware) GetMemberPermission(s *discordgo.Session, guildID string, memberID string) (permissions.PermissionArray, error) {
	return r.db.GetMemberPermission(s, guildID, memberID)
}

func (r *RedisMiddleware) GetGuildJdoodleKey(guildID string) (string, error) {
	return r.db.GetGuildJdoodleKey(guildID)
}

func (r *RedisMiddleware) SetGuildJdoodleKey(guildID, key string) error {
	return r.db.SetGuildJdoodleKey(guildID, key)
}

func (r *RedisMiddleware) GetGuildBackup(guildID string) (bool, error) {
	return r.db.GetGuildBackup(guildID)
}

func (r *RedisMiddleware) SetGuildBackup(guildID string, enabled bool) error {
	return r.db.SetGuildBackup(guildID, enabled)
}

func (r *RedisMiddleware) GetGuildInviteBlock(guildID string) (string, error) {
	return r.db.GetGuildInviteBlock(guildID)
}

func (r *RedisMiddleware) SetGuildInviteBlock(guildID string, data string) error {
	return r.db.SetGuildInviteBlock(guildID, data)
}

func (r *RedisMiddleware) GetGuildJoinMsg(guildID string) (string, string, error) {
	return r.db.GetGuildJoinMsg(guildID)
}

func (r *RedisMiddleware) SetGuildJoinMsg(guildID string, channelID string, msg string) error {
	return r.db.SetGuildJoinMsg(guildID, channelID, msg)
}

func (r *RedisMiddleware) GetGuildLeaveMsg(guildID string) (string, string, error) {
	return r.db.GetGuildLeaveMsg(guildID)
}

func (r *RedisMiddleware) SetGuildLeaveMsg(guildID string, channelID string, msg string) error {
	return r.db.SetGuildLeaveMsg(guildID, channelID, msg)
}

func (r *RedisMiddleware) AddReport(rep *report.Report) error {
	return r.db.AddReport(rep)
}

func (r *RedisMiddleware) DeleteReport(id snowflake.ID) error {
	return r.db.DeleteReport(id)
}

func (r *RedisMiddleware) GetReport(id snowflake.ID) (*report.Report, error) {
	return r.db.GetReport(id)
}

func (r *RedisMiddleware) GetReportsGuild(guildID string, offset, limit int) ([]*report.Report, error) {
	return r.db.GetReportsGuild(guildID, offset, limit)
}

func (r *RedisMiddleware) GetReportsFiltered(guildID, memberID string, repType int) ([]*report.Report, error) {
	return r.db.GetReportsFiltered(guildID, memberID, repType)
}

func (r *RedisMiddleware) GetReportsGuildCount(guildID string) (int, error) {
	return r.db.GetReportsGuildCount(guildID)
}

func (r *RedisMiddleware) GetReportsFilteredCount(guildID, memberID string, repType int) (int, error) {
	return r.db.GetReportsFilteredCount(guildID, memberID, repType)
}

func (r *RedisMiddleware) GetSetting(setting string) (string, error) {
	return r.db.GetSetting(setting)
}

func (r *RedisMiddleware) SetSetting(setting, value string) error {
	return r.db.SetSetting(setting, value)
}

func (r *RedisMiddleware) GetVotes() (map[string]*vote.Vote, error) {
	return r.db.GetVotes()
}

func (r *RedisMiddleware) AddUpdateVote(votes *vote.Vote) error {
	return r.db.AddUpdateVote(votes)
}

func (r *RedisMiddleware) DeleteVote(voteID string) error {
	return r.db.DeleteVote(voteID)
}

func (r *RedisMiddleware) GetMuteRoles() (map[string]string, error) {
	return r.db.GetMuteRoles()
}

func (r *RedisMiddleware) GetMuteRoleGuild(guildID string) (string, error) {
	return r.db.GetMuteRoleGuild(guildID)
}

func (r *RedisMiddleware) SetMuteRole(guildID, roleID string) error {
	return r.db.SetMuteRole(guildID, roleID)
}

func (r *RedisMiddleware) GetAllTwitchNotifies(twitchUserID string) ([]*twitchnotify.TwitchNotifyDBEntry, error) {
	return r.db.GetAllTwitchNotifies(twitchUserID)
}

func (r *RedisMiddleware) GetTwitchNotify(twitchUserID, guildID string) (*twitchnotify.TwitchNotifyDBEntry, error) {
	return r.db.GetTwitchNotify(twitchUserID, guildID)
}

func (r *RedisMiddleware) SetTwitchNotify(twitchNotify *twitchnotify.TwitchNotifyDBEntry) error {
	return r.db.SetTwitchNotify(twitchNotify)
}

func (r *RedisMiddleware) DeleteTwitchNotify(twitchUserID, guildID string) error {
	return r.db.DeleteTwitchNotify(twitchUserID, guildID)
}

func (r *RedisMiddleware) AddBackup(guildID, fileID string) error {
	return r.db.AddBackup(guildID, fileID)
}

func (r *RedisMiddleware) DeleteBackup(guildID, fileID string) error {
	return r.db.DeleteBackup(guildID, fileID)
}

func (r *RedisMiddleware) GetBackups(guildID string) ([]*backupmodels.Entry, error) {
	return r.db.GetBackups(guildID)
}

func (r *RedisMiddleware) GetGuilds() ([]string, error) {
	return r.db.GetGuilds()
}

func (r *RedisMiddleware) AddTag(tag *tag.Tag) error {
	return r.db.AddTag(tag)
}

func (r *RedisMiddleware) EditTag(tag *tag.Tag) error {
	return r.db.EditTag(tag)
}

func (r *RedisMiddleware) GetTagByID(id snowflake.ID) (*tag.Tag, error) {
	return r.db.GetTagByID(id)
}

func (r *RedisMiddleware) GetTagByIdent(ident string, guildID string) (*tag.Tag, error) {
	return r.db.GetTagByIdent(ident, guildID)
}

func (r *RedisMiddleware) GetGuildTags(guildID string) ([]*tag.Tag, error) {
	return r.db.GetGuildTags(guildID)
}

func (r *RedisMiddleware) DeleteTag(id snowflake.ID) error {
	return r.db.DeleteTag(id)
}

func (r *RedisMiddleware) SetSession(key, userID string, expires time.Time) error {
	return r.db.SetSession(key, userID, expires)
}

func (r *RedisMiddleware) GetSession(key string) (string, error) {
	return r.db.GetSession(key)
}

func (r *RedisMiddleware) DeleteSession(userID string) error {
	return r.db.DeleteSession(userID)
}

func (r *RedisMiddleware) GetImageData(id snowflake.ID) (*imgstore.Image, error) {
	return r.db.GetImageData(id)
}

func (r *RedisMiddleware) SaveImageData(image *imgstore.Image) error {
	return r.db.SaveImageData(image)
}

func (r *RedisMiddleware) RemoveImageData(id snowflake.ID) error {
	return r.db.RemoveImageData(id)
}
