package redis

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/zekroTJA/shinpuru/internal/core/config"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/shared/models"

	"github.com/bwmarrin/snowflake"
	"github.com/go-redis/redis"
	"github.com/zekroTJA/shinpuru/internal/core/backup/backupmodels"
	"github.com/zekroTJA/shinpuru/internal/util/imgstore"
	"github.com/zekroTJA/shinpuru/internal/util/report"
	"github.com/zekroTJA/shinpuru/internal/util/tag"
	"github.com/zekroTJA/shinpuru/internal/util/vote"
	"github.com/zekroTJA/shinpuru/pkg/permissions"
	"github.com/zekroTJA/shinpuru/pkg/twitchnotify"
)

const (
	keySetting = "PROP"

	keyGuildPrefix        = "GUILD:PREFIX"
	keyGuildAutoRole      = "GUILD:AUTOROLE"
	keyGuildModLog        = "GUILD:MODLOG"
	keyGuildVoiceLog      = "GUILD:VOICELOG"
	keyGuildNotifyRole    = "GUILD:NOTROLE"
	keyGuildGhostPingMsg  = "GUILD:GPMSG"
	keyGuildJDoodleKey    = "GUILD:JDOODLE"
	keyGuildInviteBlock   = "GUILD:INVBLOCK"
	keyGuildBackupEnabled = "GUILD:BACKUP"
	keyGuildJoinMsg       = "GUILD:JOINMSG"
	keyGuildLeaveMsg      = "GUILD:LEAVEMSG"
	keyGuildMuteRole      = "GUILD:MUTEROLE"
	keyGuildColorReaction = "GUILD:COLORREACTION"

	keyKarmaState     = "KARMA:STATE"
	keyKarmaemotesInc = "KARMA:EMOTES:ENC"
	keyKarmaEmotesDec = "KARMA:EMOTES:DEC"
	keyKarmaTokens    = "KARMA:TOKENS"

	keyAntiraidState = "ANTIRAID:STATE"
	keyAntiraidLimit = "ANTIRAID:LIMIT"
	keyAntiraidBurst = "ANTIRAID:BURST"

	keyUserAPIToken = "USER:APITOKEN"

	keyAPISession = "API:SESSION"
)

// RedisMiddleware implements the Database interface for
// Redis.
//
// This driver can only be used as caching
// middleware and consumes another database driver.
// Incomming database requests are looked up in the cache
// and values are returned from cache instead of requesting
// the database if the value is existent. Otherwise, the
// value is requested from database and then stored to cache.
// On setting database values, values are set in database as
// same as in the cache.
type RedisMiddleware struct {
	client *redis.Client
	db     database.Database
}

func NewRedisMiddleware(config *config.DatabaseRedis, db database.Database) *RedisMiddleware {
	r := &RedisMiddleware{
		db: db,
	}

	r.client = redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.Type,
	})

	return r
}

// --- DATABASE INTERFACE IMPLEMENTATIONS -------------------------------------

func (r *RedisMiddleware) Connect(credentials ...interface{}) error {
	return r.db.Connect(credentials...)
}

func (r *RedisMiddleware) Close() {
	r.client.Close()
	r.db.Close()
}

func (r *RedisMiddleware) GetGuildPrefix(guildID string) (string, error) {
	var key = fmt.Sprintf("%s:%s", keyGuildPrefix, guildID)

	val, err := r.client.Get(key).Result()
	if err == redis.Nil {
		val, err = r.db.GetGuildPrefix(guildID)
		if err != nil {
			return "", err
		}

		err = r.client.Set(key, val, 0).Err()
		return val, err
	}
	if err != nil {
		return "", err
	}

	return val, nil
}

func (r *RedisMiddleware) SetGuildPrefix(guildID, newPrefix string) error {
	var key = fmt.Sprintf("%s:%s", keyGuildPrefix, guildID)

	if err := r.client.Set(key, newPrefix, 0).Err(); err != nil {
		return err
	}

	return r.db.SetGuildPrefix(guildID, newPrefix)
}

func (r *RedisMiddleware) GetGuildAutoRole(guildID string) (string, error) {
	var key = fmt.Sprintf("%s:%s", keyGuildAutoRole, guildID)

	val, err := r.client.Get(key).Result()
	if err == redis.Nil {
		val, err = r.db.GetGuildAutoRole(guildID)
		if err != nil {
			return "", err
		}

		err = r.client.Set(key, val, 0).Err()
		return val, err
	}
	if err != nil {
		return "", err
	}

	return val, nil
}

func (r *RedisMiddleware) SetGuildAutoRole(guildID, autoRoleID string) error {
	var key = fmt.Sprintf("%s:%s", keyGuildAutoRole, guildID)

	if err := r.client.Set(key, autoRoleID, 0).Err(); err != nil {
		return err
	}

	return r.db.SetGuildAutoRole(guildID, autoRoleID)
}

func (r *RedisMiddleware) GetGuildModLog(guildID string) (string, error) {
	var key = fmt.Sprintf("%s:%s", keyGuildModLog, guildID)

	val, err := r.client.Get(key).Result()
	if err == redis.Nil {
		val, err = r.db.GetGuildModLog(guildID)
		if err != nil {
			return "", err
		}

		err = r.client.Set(key, val, 0).Err()
		return val, err
	}
	if err != nil {
		return "", err
	}

	return val, nil
}

func (r *RedisMiddleware) SetGuildModLog(guildID, chanID string) error {
	var key = fmt.Sprintf("%s:%s", keyGuildModLog, guildID)

	if err := r.client.Set(key, chanID, 0).Err(); err != nil {
		return err
	}

	return r.db.SetGuildModLog(guildID, chanID)
}

func (r *RedisMiddleware) GetGuildVoiceLog(guildID string) (string, error) {
	var key = fmt.Sprintf("%s:%s", keyGuildVoiceLog, guildID)

	val, err := r.client.Get(key).Result()
	if err == redis.Nil {
		val, err = r.db.GetGuildVoiceLog(guildID)
		if err != nil {
			return "", err
		}

		err = r.client.Set(key, val, 0).Err()
		return val, err
	}
	if err != nil {
		return "", err
	}

	return val, nil
}

func (r *RedisMiddleware) SetGuildVoiceLog(guildID, chanID string) error {
	var key = fmt.Sprintf("%s:%s", keyGuildVoiceLog, guildID)

	if err := r.client.Set(key, chanID, 0).Err(); err != nil {
		return err
	}

	return r.db.SetGuildVoiceLog(guildID, chanID)
}

func (r *RedisMiddleware) GetGuildNotifyRole(guildID string) (string, error) {
	var key = fmt.Sprintf("%s:%s", keyGuildNotifyRole, guildID)

	val, err := r.client.Get(key).Result()
	if err == redis.Nil {
		val, err = r.db.GetGuildNotifyRole(guildID)
		if err != nil {
			return "", err
		}

		err = r.client.Set(key, val, 0).Err()
		return val, err
	}
	if err != nil {
		return "", err
	}

	return val, nil
}

func (r *RedisMiddleware) SetGuildNotifyRole(guildID, roleID string) error {
	var key = fmt.Sprintf("%s:%s", keyGuildNotifyRole, guildID)

	if err := r.client.Set(key, roleID, 0).Err(); err != nil {
		return err
	}

	return r.db.SetGuildNotifyRole(guildID, roleID)
}

func (r *RedisMiddleware) GetGuildGhostpingMsg(guildID string) (string, error) {
	var key = fmt.Sprintf("%s:%s", keyGuildGhostPingMsg, guildID)

	val, err := r.client.Get(key).Result()
	if err == redis.Nil {
		val, err = r.db.GetGuildGhostpingMsg(guildID)
		if err != nil {
			return "", err
		}

		err = r.client.Set(key, val, 0).Err()
		return val, err
	}
	if err != nil {
		return "", err
	}

	return val, nil
}

func (r *RedisMiddleware) SetGuildGhostpingMsg(guildID, msg string) error {
	var key = fmt.Sprintf("%s:%s", keyGuildGhostPingMsg, guildID)

	if err := r.client.Set(key, msg, 0).Err(); err != nil {
		return err
	}

	return r.db.SetGuildGhostpingMsg(guildID, msg)
}

func (r *RedisMiddleware) GetGuildPermissions(guildID string) (map[string]permissions.PermissionArray, error) {
	return r.db.GetGuildPermissions(guildID)
}

func (r *RedisMiddleware) SetGuildRolePermission(guildID, roleID string, p permissions.PermissionArray) error {
	return r.db.SetGuildRolePermission(guildID, roleID, p)
}

func (r *RedisMiddleware) GetGuildJdoodleKey(guildID string) (string, error) {
	var key = fmt.Sprintf("%s:%s", keyGuildJDoodleKey, guildID)

	val, err := r.client.Get(key).Result()
	if err == redis.Nil {
		val, err = r.db.GetGuildJdoodleKey(guildID)
		if err != nil {
			return "", err
		}

		err = r.client.Set(key, val, 0).Err()
		return val, err
	}
	if err != nil {
		return "", err
	}

	return val, nil
}

func (r *RedisMiddleware) SetGuildJdoodleKey(guildID, jdkey string) error {
	var key = fmt.Sprintf("%s:%s", keyGuildJDoodleKey, guildID)

	if err := r.client.Set(key, jdkey, 0).Err(); err != nil {
		return err
	}

	return r.db.SetGuildJdoodleKey(guildID, jdkey)
}

func (r *RedisMiddleware) GetGuildBackup(guildID string) (bool, error) {
	var key = fmt.Sprintf("%s:%s", keyGuildBackupEnabled, guildID)

	var val bool
	err := r.client.Get(key).Scan(&val)
	if err == redis.Nil {
		val, err = r.db.GetGuildBackup(guildID)
		if err != nil {
			return false, err
		}

		err = r.client.Set(key, val, 0).Err()
		return val, err
	}
	if err != nil {
		return false, err
	}

	return val, nil
}

func (r *RedisMiddleware) SetGuildBackup(guildID string, enabled bool) error {
	var key = fmt.Sprintf("%s:%s", keyGuildBackupEnabled, guildID)

	if err := r.client.Set(key, enabled, 0).Err(); err != nil {
		return err
	}

	return r.db.SetGuildBackup(guildID, enabled)
}

func (r *RedisMiddleware) GetGuildInviteBlock(guildID string) (string, error) {
	var key = fmt.Sprintf("%s:%s", keyGuildInviteBlock, guildID)

	val, err := r.client.Get(key).Result()
	if err == redis.Nil {
		val, err = r.db.GetGuildInviteBlock(guildID)
		if err != nil {
			return "", err
		}

		err = r.client.Set(key, val, 0).Err()
		return val, err
	}
	if err != nil {
		return "", err
	}

	return val, nil
}

func (r *RedisMiddleware) SetGuildInviteBlock(guildID string, data string) error {
	var key = fmt.Sprintf("%s:%s", keyGuildInviteBlock, guildID)

	if err := r.client.Set(key, data, 0).Err(); err != nil {
		return err
	}

	return r.db.SetGuildInviteBlock(guildID, data)
}

func (r *RedisMiddleware) GetGuildJoinMsg(guildID string) (string, string, error) {
	var key = fmt.Sprintf("%s:%s", keyGuildJoinMsg, guildID)

	var val1, val2 string

	raw, err := r.client.Get(key).Result()
	if err == redis.Nil {
		val1, val2, err = r.db.GetGuildJoinMsg(guildID)
		if err != nil {
			return "", "", err
		}

		err = r.client.Set(key, fmt.Sprintf("%s|%s", val1, val2), 0).Err()
		return val1, val2, err
	}
	if err != nil {
		return "", "", err
	}

	rawSplit := strings.Split(raw, "|")
	val1, val2 = rawSplit[0], rawSplit[1]

	return val1, val2, nil
}

func (r *RedisMiddleware) SetGuildJoinMsg(guildID string, channelID string, msg string) error {
	var key = fmt.Sprintf("%s:%s", keyGuildJoinMsg, guildID)

	if err := r.client.Set(key, fmt.Sprintf("%s|%s", channelID, msg), 0).Err(); err != nil {
		return err
	}

	return r.db.SetGuildJoinMsg(guildID, channelID, msg)
}

func (r *RedisMiddleware) GetGuildLeaveMsg(guildID string) (string, string, error) {
	var key = fmt.Sprintf("%s:%s", keyGuildLeaveMsg, guildID)

	var val1, val2 string

	raw, err := r.client.Get(key).Result()
	if err == redis.Nil {
		val1, val2, err = r.db.GetGuildLeaveMsg(guildID)
		if err != nil {
			return "", "", err
		}

		err = r.client.Set(key, fmt.Sprintf("%s|%s", val1, val2), 0).Err()
		return val1, val2, err
	}
	if err != nil {
		return "", "", err
	}

	rawSplit := strings.Split(raw, "|")
	val1, val2 = rawSplit[0], rawSplit[1]

	return val1, val2, nil
}

func (r *RedisMiddleware) SetGuildLeaveMsg(guildID string, channelID string, msg string) error {
	var key = fmt.Sprintf("%s:%s", keyGuildLeaveMsg, guildID)

	if err := r.client.Set(key, fmt.Sprintf("%s|%s", channelID, msg), 0).Err(); err != nil {
		return err
	}

	return r.db.SetGuildLeaveMsg(guildID, channelID, msg)
}

func (r *RedisMiddleware) GetGuildColorReaction(guildID string) (bool, error) {
	var key = fmt.Sprintf("%s:%s", keyGuildColorReaction, guildID)

	var val bool
	err := r.client.Get(key).Scan(&val)
	if err == redis.Nil {
		val, err = r.db.GetGuildColorReaction(guildID)
		if err != nil {
			return false, err
		}

		err = r.client.Set(key, val, 0).Err()
		return val, err
	}
	if err != nil {
		return false, err
	}

	return val, nil
}

func (r *RedisMiddleware) SetGuildColorReaction(guildID string, enabled bool) error {
	var key = fmt.Sprintf("%s:%s", keyGuildColorReaction, guildID)

	if err := r.client.Set(key, enabled, 0).Err(); err != nil {
		return err
	}

	return r.db.SetGuildColorReaction(guildID, enabled)
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
	var key = fmt.Sprintf("%s:%s", keySetting, setting)

	val, err := r.client.Get(key).Result()
	if err == redis.Nil {
		val, err = r.db.GetSetting(setting)
		if err != nil {
			return "", err
		}

		err = r.client.Set(key, val, 0).Err()
		return val, err
	}
	if err != nil {
		return "", err
	}

	return val, nil
}

func (r *RedisMiddleware) SetSetting(setting, value string) error {
	var key = fmt.Sprintf("%s:%s", keySetting, setting)

	if err := r.client.Set(key, value, 0).Err(); err != nil {
		return err
	}

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

func (r *RedisMiddleware) GetGuildMuteRole(guildID string) (string, error) {
	var key = fmt.Sprintf("%s:%s", keyGuildMuteRole, guildID)

	val, err := r.client.Get(key).Result()
	if err == redis.Nil {
		val, err = r.db.GetGuildMuteRole(guildID)
		if err != nil {
			return "", err
		}

		err = r.client.Set(key, val, 0).Err()
		return val, err
	}
	if err != nil {
		return "", err
	}

	return val, nil
}

func (r *RedisMiddleware) SetGuildMuteRole(guildID, roleID string) error {
	var key = fmt.Sprintf("%s:%s", keyGuildMuteRole, guildID)

	if err := r.client.Set(key, roleID, 0).Err(); err != nil {
		return err
	}

	return r.db.SetGuildMuteRole(guildID, roleID)
}

func (r *RedisMiddleware) GetAllTwitchNotifies(twitchUserID string) ([]*twitchnotify.DBEntry, error) {
	return r.db.GetAllTwitchNotifies(twitchUserID)
}

func (r *RedisMiddleware) GetTwitchNotify(twitchUserID, guildID string) (*twitchnotify.DBEntry, error) {
	return r.db.GetTwitchNotify(twitchUserID, guildID)
}

func (r *RedisMiddleware) SetTwitchNotify(twitchNotify *twitchnotify.DBEntry) error {
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

func (r *RedisMiddleware) GetImageData(id snowflake.ID) (*imgstore.Image, error) {
	return r.db.GetImageData(id)
}

func (r *RedisMiddleware) SaveImageData(image *imgstore.Image) error {
	return r.db.SaveImageData(image)
}

func (r *RedisMiddleware) RemoveImageData(id snowflake.ID) error {
	return r.db.RemoveImageData(id)
}

func (m *RedisMiddleware) SetAPIToken(token *models.APITokenEntry) (err error) {
	var key = fmt.Sprintf("%s:%s", keyUserAPIToken, token.UserID)

	data, err := json.Marshal(token)
	if err != nil {
		return
	}

	if err = m.client.Set(key, data, 0).Err(); err != nil {
		return
	}

	return m.db.SetAPIToken(token)
}

func (m *RedisMiddleware) GetAPIToken(userID string) (t *models.APITokenEntry, err error) {
	var key = fmt.Sprintf("%s:%s", keyUserAPIToken, userID)

	resStr, err := m.client.Get(key).Result()
	if err == redis.Nil {
		if t, err = m.db.GetAPIToken(userID); err != nil {
			return
		}
		var resB []byte
		resB, err = json.Marshal(t)
		if err != nil {
			return
		}
		if err = m.client.Set(key, resB, 0).Err(); err != nil {
			return
		}
		return
	}

	t = new(models.APITokenEntry)
	err = json.Unmarshal([]byte(resStr), t)

	return
}

func (m *RedisMiddleware) DeleteAPIToken(userID string) (err error) {
	var key = fmt.Sprintf("%s:%s", keyUserAPIToken, userID)

	if err = m.client.Del(key).Err(); err != nil {
		return
	}

	return m.db.DeleteAPIToken(userID)
}

func (m *RedisMiddleware) GetKarma(userID, guildID string) (int, error) {
	return m.db.GetKarma(userID, guildID)
}

func (m *RedisMiddleware) GetKarmaSum(userID string) (int, error) {
	return m.db.GetKarmaSum(userID)
}

func (m *RedisMiddleware) GetKarmaGuild(guildID string, limit int) ([]*models.GuildKarma, error) {
	return m.db.GetKarmaGuild(guildID, limit)
}

func (m *RedisMiddleware) SetKarma(userID, guildID string, val int) error {
	return m.db.SetKarma(userID, guildID, val)
}

func (m *RedisMiddleware) UpdateKarma(userID, guildID string, diff int) error {
	return m.db.UpdateKarma(userID, guildID, diff)
}

func (m *RedisMiddleware) SetKarmaState(guildID string, state bool) error {
	var key = fmt.Sprintf("%s:%s", keyKarmaState, guildID)

	if err := m.client.Set(key, state, 0).Err(); err != nil {
		return err
	}

	return m.db.SetKarmaState(guildID, state)
}

func (m *RedisMiddleware) GetKarmaState(guildID string) (bool, error) {
	var key = fmt.Sprintf("%s:%s", keyKarmaState, guildID)

	var val bool
	err := m.client.Get(key).Scan(&val)
	if err == redis.Nil {
		val, err = m.db.GetKarmaState(guildID)
		if err != nil {
			return false, err
		}

		err = m.client.Set(key, val, 0).Err()
		return val, err
	}
	if err != nil {
		return false, err
	}

	return val, nil
}

func (m *RedisMiddleware) SetKarmaEmotes(guildID, emotesInc, emotesDec string) error {
	var key = fmt.Sprintf("%s:%s", keyKarmaemotesInc, guildID)
	if err := m.client.Set(key, emotesInc, 0).Err(); err != nil {
		return err
	}

	key = fmt.Sprintf("%s:%s", keyKarmaEmotesDec, guildID)
	if err := m.client.Set(key, emotesDec, 0).Err(); err != nil {
		return err
	}

	return m.db.SetKarmaEmotes(guildID, emotesInc, emotesDec)
}

func (m *RedisMiddleware) GetKarmaEmotes(guildID string) (emotesInc, emotesDec string, err error) {
	var keyEnc = fmt.Sprintf("%s:%s", keyKarmaemotesInc, guildID)
	emotesInc, err1 := m.client.Get(keyEnc).Result()

	var keyDec = fmt.Sprintf("%s:%s", keyKarmaEmotesDec, guildID)
	emotesDec, err2 := m.client.Get(keyDec).Result()

	if err1 == redis.Nil || err2 == redis.Nil {
		emotesInc, emotesDec, err = m.db.GetKarmaEmotes(guildID)
		if err != nil {
			return
		}

		if err = m.client.Set(keyEnc, emotesInc, 0).Err(); err != nil {
			return
		}
		if err = m.client.Set(keyDec, emotesDec, 0).Err(); err != nil {
			return
		}
	}
	if err != nil {
		return
	}

	return
}

func (m *RedisMiddleware) SetKarmaTokens(guildID string, tokens int) error {
	var key = fmt.Sprintf("%s:%s", keyKarmaTokens, guildID)

	if err := m.client.Set(key, tokens, 0).Err(); err != nil {
		return err
	}

	return m.db.SetKarmaTokens(guildID, tokens)
}

func (m *RedisMiddleware) GetKarmaTokens(guildID string) (int, error) {
	var key = fmt.Sprintf("%s:%s", keyKarmaTokens, guildID)

	var val int
	err := m.client.Get(key).Scan(&val)
	if err == redis.Nil {
		val, err = m.db.GetKarmaTokens(guildID)
		if err != nil {
			return 0, err
		}

		err = m.client.Set(key, val, 0).Err()
		return val, err
	}
	if err != nil {
		return 0, err
	}

	return val, nil
}

func (m *RedisMiddleware) SetLockChan(chanID, guildID, executorID, permissions string) error {
	return m.db.SetLockChan(chanID, guildID, executorID, permissions)
}

func (m *RedisMiddleware) GetLockChan(chanID string) (guildID, executorID, permissions string, err error) {
	return m.db.GetLockChan(chanID)
}

func (m *RedisMiddleware) GetLockChannels(guildID string) (chanIDs []string, err error) {
	return m.db.GetLockChannels(guildID)
}

func (m *RedisMiddleware) DeleteLockChan(chanID string) error {
	return m.db.DeleteLockChan(chanID)
}

func (m *RedisMiddleware) SetAntiraidState(guildID string, state bool) error {
	var key = fmt.Sprintf("%s:%s", keyAntiraidState, guildID)

	if err := m.client.Set(key, state, 0).Err(); err != nil {
		return err
	}

	return m.db.SetAntiraidState(guildID, state)
}

func (m *RedisMiddleware) GetAntiraidState(guildID string) (bool, error) {
	var key = fmt.Sprintf("%s:%s", keyAntiraidState, guildID)

	var val bool
	err := m.client.Get(key).Scan(&val)
	if err == redis.Nil {
		val, err = m.db.GetAntiraidState(guildID)
		if err != nil {
			return false, err
		}

		err = m.client.Set(key, val, 0).Err()
		return val, err
	}
	if err != nil {
		return false, err
	}

	return val, nil
}

func (m *RedisMiddleware) SetAntiraidRegeneration(guildID string, limit int) error {
	var key = fmt.Sprintf("%s:%s", keyAntiraidLimit, guildID)

	if err := m.client.Set(key, limit, 0).Err(); err != nil {
		return err
	}

	return m.db.SetKarmaTokens(guildID, limit)
}

func (m *RedisMiddleware) GetAntiraidRegeneration(guildID string) (int, error) {
	var key = fmt.Sprintf("%s:%s", keyAntiraidLimit, guildID)

	var val int
	err := m.client.Get(key).Scan(&val)
	if err == redis.Nil {
		val, err = m.db.GetAntiraidRegeneration(guildID)
		if err != nil {
			return 0, err
		}

		err = m.client.Set(key, val, 0).Err()
		return val, err
	}
	if err != nil {
		return 0, err
	}

	return val, nil
}

func (m *RedisMiddleware) SetAntiraidBurst(guildID string, burst int) error {
	var key = fmt.Sprintf("%s:%s", keyAntiraidBurst, guildID)

	if err := m.client.Set(key, burst, 0).Err(); err != nil {
		return err
	}

	return m.db.SetAntiraidBurst(guildID, burst)
}

func (m *RedisMiddleware) GetAntiraidBurst(guildID string) (int, error) {
	var key = fmt.Sprintf("%s:%s", keyAntiraidBurst, guildID)

	var val int
	err := m.client.Get(key).Scan(&val)
	if err == redis.Nil {
		val, err = m.db.GetAntiraidBurst(guildID)
		if err != nil {
			return 0, err
		}

		err = m.client.Set(key, val, 0).Err()
		return val, err
	}
	if err != nil {
		return 0, err
	}

	return val, nil
}

func (m *RedisMiddleware) AddToAntiraidJoinList(guildID, userID, userTag string) error {
	return m.db.AddToAntiraidJoinList(guildID, userID, userTag)
}

func (m *RedisMiddleware) GetAntiraidJoinList(guildID string) ([]*models.JoinLogEntry, error) {
	return m.db.GetAntiraidJoinList(guildID)
}

func (m *RedisMiddleware) FlushAntiraidJoinList(guildID string) error {
	return m.db.FlushAntiraidJoinList(guildID)
}

func (m *RedisMiddleware) GetGuildUnbanRequests(guildID string) ([]*report.UnbanRequest, error) {
	return m.db.GetGuildUnbanRequests(guildID)
}

func (m *RedisMiddleware) GetGuildUserUnbanRequests(userID, guildID string) ([]*report.UnbanRequest, error) {
	return m.db.GetGuildUserUnbanRequests(userID, guildID)
}

func (m *RedisMiddleware) GetUnbanRequest(id string) (*report.UnbanRequest, error) {
	return m.db.GetUnbanRequest(id)
}

func (m *RedisMiddleware) AddUnbanRequest(request *report.UnbanRequest) error {
	return m.db.AddUnbanRequest(request)
}

func (m *RedisMiddleware) UpdateUnbanRequest(request *report.UnbanRequest) error {
	return m.db.UpdateUnbanRequest(request)
}
