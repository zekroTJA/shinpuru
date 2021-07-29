package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/zekroTJA/shinpuru/internal/models"
	"github.com/zekroTJA/shinpuru/internal/services/database"

	"github.com/go-redis/redis/v8"
	"github.com/zekroTJA/shinpuru/pkg/boolutil"
)

const (
	keySetting = "PROP"

	keyGuildPrefix          = "GUILD:PREFIX"
	keyGuildAutoRole        = "GUILD:AUTOROLE"
	keyGuildModLog          = "GUILD:MODLOG"
	keyGuildVoiceLog        = "GUILD:VOICELOG"
	keyGuildNotifyRole      = "GUILD:NOTROLE"
	keyGuildGhostPingMsg    = "GUILD:GPMSG"
	keyGuildJDoodleKey      = "GUILD:JDOODLE"
	keyGuildInviteBlock     = "GUILD:INVBLOCK"
	keyGuildBackupEnabled   = "GUILD:BACKUP"
	keyGuildJoinMsg         = "GUILD:JOINMSG"
	keyGuildLeaveMsg        = "GUILD:LEAVEMSG"
	keyGuildMuteRole        = "GUILD:MUTEROLE"
	keyGuildColorReaction   = "GUILD:COLORREACTION"
	keyGuildStarboardConfig = "GUILD:STARBOARDCONFIG"
	keyGuildLogEnable       = "GUILD:GUILDLOG"

	keyKarmaState       = "KARMA:STATE"
	keyKarmaemotesInc   = "KARMA:EMOTES:ENC"
	keyKarmaEmotesDec   = "KARMA:EMOTES:DEC"
	keyKarmaTokens      = "KARMA:TOKENS"
	keyKarmaPenalty     = "KARMA:PENALTY"
	keyKarmaBlockListed = "KARMA:BLOCKLISTED"

	keyAntiraidState = "ANTIRAID:STATE"
	keyAntiraidLimit = "ANTIRAID:LIMIT"
	keyAntiraidBurst = "ANTIRAID:BURST"

	keyUserAPIToken  = "USER:APITOKEN"
	keyUserEnableOTA = "USER:ENABLEOTA"

	keyAPISession = "API:SESSION"
)

type getterFunc func() (interface{}, error)

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
	database.Database

	client *redis.Client
}

var _ database.Database = (*RedisMiddleware)(nil)

func NewRedisMiddleware(db database.Database, rd *redis.Client) *RedisMiddleware {
	return &RedisMiddleware{
		Database: db,
		client:   rd,
	}
}

// --- DATABASE INTERFACE IMPLEMENTATIONS -------------------------------------

func (r *RedisMiddleware) Connect(credentials ...interface{}) error {
	return r.Database.Connect(credentials...)
}

func (r *RedisMiddleware) Close() {
	r.client.Close()
	r.Database.Close()
}

func (r *RedisMiddleware) GetGuildPrefix(guildID string) (string, error) {
	var key = fmt.Sprintf("%s:%s", keyGuildPrefix, guildID)

	val, err := r.client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		val, err = r.Database.GetGuildPrefix(guildID)
		if err != nil {
			return "", err
		}

		err = r.client.Set(context.Background(), key, val, 0).Err()
		return val, err
	}
	if err != nil {
		return "", err
	}

	return val, nil
}

func (r *RedisMiddleware) SetGuildPrefix(guildID, newPrefix string) error {
	var key = fmt.Sprintf("%s:%s", keyGuildPrefix, guildID)

	if err := r.client.Set(context.Background(), key, newPrefix, 0).Err(); err != nil {
		return err
	}

	return r.Database.SetGuildPrefix(guildID, newPrefix)
}

func (r *RedisMiddleware) GetGuildAutoRole(guildID string) ([]string, error) {
	var key = fmt.Sprintf("%s:%s", keyGuildAutoRole, guildID)

	valC, err := r.client.Get(context.Background(), key).Result()
	val := strings.Split(valC, ";")
	if err == redis.Nil {
		val, err = r.Database.GetGuildAutoRole(guildID)
		if err != nil {
			return nil, err
		}

		err = r.client.Set(context.Background(), key, strings.Join(val, ";"), 0).Err()
		return val, err
	}
	if err != nil {
		return nil, err
	}

	return val, nil
}

func (r *RedisMiddleware) SetGuildAutoRole(guildID string, autoRoleIDs []string) error {
	var key = fmt.Sprintf("%s:%s", keyGuildAutoRole, guildID)

	if err := r.client.Set(context.Background(), key, strings.Join(autoRoleIDs, ";"), 0).Err(); err != nil {
		return err
	}

	return r.Database.SetGuildAutoRole(guildID, autoRoleIDs)
}

func (r *RedisMiddleware) GetGuildModLog(guildID string) (string, error) {
	var key = fmt.Sprintf("%s:%s", keyGuildModLog, guildID)

	val, err := r.client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		val, err = r.Database.GetGuildModLog(guildID)
		if err != nil {
			return "", err
		}

		err = r.client.Set(context.Background(), key, val, 0).Err()
		return val, err
	}
	if err != nil {
		return "", err
	}

	return val, nil
}

func (r *RedisMiddleware) SetGuildModLog(guildID, chanID string) error {
	var key = fmt.Sprintf("%s:%s", keyGuildModLog, guildID)

	if err := r.client.Set(context.Background(), key, chanID, 0).Err(); err != nil {
		return err
	}

	return r.Database.SetGuildModLog(guildID, chanID)
}

func (r *RedisMiddleware) GetGuildVoiceLog(guildID string) (string, error) {
	var key = fmt.Sprintf("%s:%s", keyGuildVoiceLog, guildID)

	val, err := r.client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		val, err = r.Database.GetGuildVoiceLog(guildID)
		if err != nil {
			return "", err
		}

		err = r.client.Set(context.Background(), key, val, 0).Err()
		return val, err
	}
	if err != nil {
		return "", err
	}

	return val, nil
}

func (r *RedisMiddleware) SetGuildVoiceLog(guildID, chanID string) error {
	var key = fmt.Sprintf("%s:%s", keyGuildVoiceLog, guildID)

	if err := r.client.Set(context.Background(), key, chanID, 0).Err(); err != nil {
		return err
	}

	return r.Database.SetGuildVoiceLog(guildID, chanID)
}

func (r *RedisMiddleware) GetGuildNotifyRole(guildID string) (string, error) {
	var key = fmt.Sprintf("%s:%s", keyGuildNotifyRole, guildID)

	val, err := r.client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		val, err = r.Database.GetGuildNotifyRole(guildID)
		if err != nil {
			return "", err
		}

		err = r.client.Set(context.Background(), key, val, 0).Err()
		return val, err
	}
	if err != nil {
		return "", err
	}

	return val, nil
}

func (r *RedisMiddleware) SetGuildNotifyRole(guildID, roleID string) error {
	var key = fmt.Sprintf("%s:%s", keyGuildNotifyRole, guildID)

	if err := r.client.Set(context.Background(), key, roleID, 0).Err(); err != nil {
		return err
	}

	return r.Database.SetGuildNotifyRole(guildID, roleID)
}

func (r *RedisMiddleware) GetGuildGhostpingMsg(guildID string) (string, error) {
	var key = fmt.Sprintf("%s:%s", keyGuildGhostPingMsg, guildID)

	val, err := r.client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		val, err = r.Database.GetGuildGhostpingMsg(guildID)
		if err != nil {
			return "", err
		}

		err = r.client.Set(context.Background(), key, val, 0).Err()
		return val, err
	}
	if err != nil {
		return "", err
	}

	return val, nil
}

func (r *RedisMiddleware) SetGuildGhostpingMsg(guildID, msg string) error {
	var key = fmt.Sprintf("%s:%s", keyGuildGhostPingMsg, guildID)

	if err := r.client.Set(context.Background(), key, msg, 0).Err(); err != nil {
		return err
	}

	return r.Database.SetGuildGhostpingMsg(guildID, msg)
}

func (r *RedisMiddleware) GetGuildJdoodleKey(guildID string) (string, error) {
	var key = fmt.Sprintf("%s:%s", keyGuildJDoodleKey, guildID)

	val, err := r.client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		val, err = r.Database.GetGuildJdoodleKey(guildID)
		if err != nil {
			return "", err
		}

		err = r.client.Set(context.Background(), key, val, 0).Err()
		return val, err
	}
	if err != nil {
		return "", err
	}

	return val, nil
}

func (r *RedisMiddleware) SetGuildJdoodleKey(guildID, jdkey string) error {
	var key = fmt.Sprintf("%s:%s", keyGuildJDoodleKey, guildID)

	if err := r.client.Set(context.Background(), key, jdkey, 0).Err(); err != nil {
		return err
	}

	return r.Database.SetGuildJdoodleKey(guildID, jdkey)
}

func (r *RedisMiddleware) GetGuildBackup(guildID string) (bool, error) {
	var key = fmt.Sprintf("%s:%s", keyGuildBackupEnabled, guildID)

	var val bool
	err := r.client.Get(context.Background(), key).Scan(&val)
	if err == redis.Nil {
		val, err = r.Database.GetGuildBackup(guildID)
		if err != nil {
			return false, err
		}

		err = r.client.Set(context.Background(), key, val, 0).Err()
		return val, err
	}
	if err != nil {
		return false, err
	}

	return val, nil
}

func (r *RedisMiddleware) SetGuildBackup(guildID string, enabled bool) error {
	var key = fmt.Sprintf("%s:%s", keyGuildBackupEnabled, guildID)

	if err := r.client.Set(context.Background(), key, enabled, 0).Err(); err != nil {
		return err
	}

	return r.Database.SetGuildBackup(guildID, enabled)
}

func (r *RedisMiddleware) GetGuildInviteBlock(guildID string) (string, error) {
	var key = fmt.Sprintf("%s:%s", keyGuildInviteBlock, guildID)

	val, err := r.client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		val, err = r.Database.GetGuildInviteBlock(guildID)
		if err != nil {
			return "", err
		}

		err = r.client.Set(context.Background(), key, val, 0).Err()
		return val, err
	}
	if err != nil {
		return "", err
	}

	return val, nil
}

func (r *RedisMiddleware) SetGuildInviteBlock(guildID string, data string) error {
	var key = fmt.Sprintf("%s:%s", keyGuildInviteBlock, guildID)

	if err := r.client.Set(context.Background(), key, data, 0).Err(); err != nil {
		return err
	}

	return r.Database.SetGuildInviteBlock(guildID, data)
}

func (r *RedisMiddleware) GetGuildJoinMsg(guildID string) (string, string, error) {
	var key = fmt.Sprintf("%s:%s", keyGuildJoinMsg, guildID)

	var val1, val2 string

	raw, err := r.client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		val1, val2, err = r.Database.GetGuildJoinMsg(guildID)
		if err != nil {
			return "", "", err
		}

		err = r.client.Set(context.Background(), key, fmt.Sprintf("%s|%s", val1, val2), 0).Err()
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

	if err := r.client.Set(context.Background(), key, fmt.Sprintf("%s|%s", channelID, msg), 0).Err(); err != nil {
		return err
	}

	return r.Database.SetGuildJoinMsg(guildID, channelID, msg)
}

func (r *RedisMiddleware) GetGuildLeaveMsg(guildID string) (string, string, error) {
	var key = fmt.Sprintf("%s:%s", keyGuildLeaveMsg, guildID)

	var val1, val2 string

	raw, err := r.client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		val1, val2, err = r.Database.GetGuildLeaveMsg(guildID)
		if err != nil {
			return "", "", err
		}

		err = r.client.Set(context.Background(), key, fmt.Sprintf("%s|%s", val1, val2), 0).Err()
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

	if err := r.client.Set(context.Background(), key, fmt.Sprintf("%s|%s", channelID, msg), 0).Err(); err != nil {
		return err
	}

	return r.Database.SetGuildLeaveMsg(guildID, channelID, msg)
}

func (r *RedisMiddleware) GetGuildColorReaction(guildID string) (bool, error) {
	var key = fmt.Sprintf("%s:%s", keyGuildColorReaction, guildID)

	var val bool
	err := r.client.Get(context.Background(), key).Scan(&val)
	if err == redis.Nil {
		val, err = r.Database.GetGuildColorReaction(guildID)
		if err != nil {
			return false, err
		}

		err = r.client.Set(context.Background(), key, val, 0).Err()
		return val, err
	}
	if err != nil {
		return false, err
	}

	return val, nil
}

func (r *RedisMiddleware) SetGuildColorReaction(guildID string, enabled bool) error {
	var key = fmt.Sprintf("%s:%s", keyGuildColorReaction, guildID)

	if err := r.client.Set(context.Background(), key, enabled, 0).Err(); err != nil {
		return err
	}

	return r.Database.SetGuildColorReaction(guildID, enabled)
}

func (r *RedisMiddleware) GetSetting(setting string) (string, error) {
	var key = fmt.Sprintf("%s:%s", keySetting, setting)

	val, err := r.client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		val, err = r.Database.GetSetting(setting)
		if err != nil {
			return "", err
		}

		err = r.client.Set(context.Background(), key, val, 0).Err()
		return val, err
	}
	if err != nil {
		return "", err
	}

	return val, nil
}

func (r *RedisMiddleware) SetSetting(setting, value string) error {
	var key = fmt.Sprintf("%s:%s", keySetting, setting)

	if err := r.client.Set(context.Background(), key, value, 0).Err(); err != nil {
		return err
	}

	return r.Database.SetSetting(setting, value)
}

func (r *RedisMiddleware) GetGuildMuteRole(guildID string) (string, error) {
	var key = fmt.Sprintf("%s:%s", keyGuildMuteRole, guildID)

	val, err := r.client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		val, err = r.Database.GetGuildMuteRole(guildID)
		if err != nil {
			return "", err
		}

		err = r.client.Set(context.Background(), key, val, 0).Err()
		return val, err
	}
	if err != nil {
		return "", err
	}

	return val, nil
}

func (r *RedisMiddleware) SetGuildMuteRole(guildID, roleID string) error {
	var key = fmt.Sprintf("%s:%s", keyGuildMuteRole, guildID)

	if err := r.client.Set(context.Background(), key, roleID, 0).Err(); err != nil {
		return err
	}

	return r.Database.SetGuildMuteRole(guildID, roleID)
}

func (m *RedisMiddleware) SetAPIToken(token *models.APITokenEntry) (err error) {
	var key = fmt.Sprintf("%s:%s", keyUserAPIToken, token.UserID)

	data, err := json.Marshal(token)
	if err != nil {
		return
	}

	if err = m.client.Set(context.Background(), key, data, 0).Err(); err != nil {
		return
	}

	return m.Database.SetAPIToken(token)
}

func (m *RedisMiddleware) GetAPIToken(userID string) (t *models.APITokenEntry, err error) {
	var key = fmt.Sprintf("%s:%s", keyUserAPIToken, userID)

	resStr, err := m.client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		if t, err = m.Database.GetAPIToken(userID); err != nil {
			return
		}
		var resB []byte
		resB, err = json.Marshal(t)
		if err != nil {
			return
		}
		if err = m.client.Set(context.Background(), key, resB, 0).Err(); err != nil {
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

	if err = m.client.Del(context.Background(), key).Err(); err != nil {
		return
	}

	return m.Database.DeleteAPIToken(userID)
}

func (m *RedisMiddleware) SetKarmaState(guildID string, state bool) error {
	var key = fmt.Sprintf("%s:%s", keyKarmaState, guildID)

	if err := m.client.Set(context.Background(), key, state, 0).Err(); err != nil {
		return err
	}

	return m.Database.SetKarmaState(guildID, state)
}

func (m *RedisMiddleware) GetKarmaState(guildID string) (bool, error) {
	var key = fmt.Sprintf("%s:%s", keyKarmaState, guildID)

	var val bool
	err := m.client.Get(context.Background(), key).Scan(&val)
	if err == redis.Nil {
		val, err = m.Database.GetKarmaState(guildID)
		if err != nil {
			return false, err
		}

		err = m.client.Set(context.Background(), key, val, 0).Err()
		return val, err
	}
	if err != nil {
		return false, err
	}

	return val, nil
}

func (m *RedisMiddleware) SetKarmaEmotes(guildID, emotesInc, emotesDec string) error {
	var key = fmt.Sprintf("%s:%s", keyKarmaemotesInc, guildID)
	if err := m.client.Set(context.Background(), key, emotesInc, 0).Err(); err != nil {
		return err
	}

	key = fmt.Sprintf("%s:%s", keyKarmaEmotesDec, guildID)
	if err := m.client.Set(context.Background(), key, emotesDec, 0).Err(); err != nil {
		return err
	}

	return m.Database.SetKarmaEmotes(guildID, emotesInc, emotesDec)
}

func (m *RedisMiddleware) GetKarmaEmotes(guildID string) (emotesInc, emotesDec string, err error) {
	var keyEnc = fmt.Sprintf("%s:%s", keyKarmaemotesInc, guildID)
	emotesInc, err1 := m.client.Get(context.Background(), keyEnc).Result()

	var keyDec = fmt.Sprintf("%s:%s", keyKarmaEmotesDec, guildID)
	emotesDec, err2 := m.client.Get(context.Background(), keyDec).Result()

	if err1 == redis.Nil || err2 == redis.Nil {
		emotesInc, emotesDec, err = m.Database.GetKarmaEmotes(guildID)
		if err != nil {
			return
		}

		if err = m.client.Set(context.Background(), keyEnc, emotesInc, 0).Err(); err != nil {
			return
		}
		if err = m.client.Set(context.Background(), keyDec, emotesDec, 0).Err(); err != nil {
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

	if err := m.client.Set(context.Background(), key, tokens, 0).Err(); err != nil {
		return err
	}

	return m.Database.SetKarmaTokens(guildID, tokens)
}

func (m *RedisMiddleware) GetKarmaTokens(guildID string) (int, error) {
	var key = fmt.Sprintf("%s:%s", keyKarmaTokens, guildID)

	var val int
	err := m.client.Get(context.Background(), key).Scan(&val)
	if err == redis.Nil {
		val, err = m.Database.GetKarmaTokens(guildID)
		if err != nil {
			return 0, err
		}

		err = m.client.Set(context.Background(), key, val, 0).Err()
		return val, err
	}
	if err != nil {
		return 0, err
	}

	return val, nil
}

func (m *RedisMiddleware) SetKarmaPenalty(guildID string, state bool) error {
	var key = fmt.Sprintf("%s:%s", keyKarmaPenalty, guildID)

	if err := m.client.Set(context.Background(), key, state, 0).Err(); err != nil {
		return err
	}

	return m.Database.SetKarmaPenalty(guildID, state)
}

func (m *RedisMiddleware) GetKarmaPenalty(guildID string) (bool, error) {
	var key = fmt.Sprintf("%s:%s", keyKarmaPenalty, guildID)

	var val bool
	err := m.client.Get(context.Background(), key).Scan(&val)
	if err == redis.Nil {
		val, err = m.Database.GetKarmaPenalty(guildID)
		if err != nil {
			return false, err
		}

		err = m.client.Set(context.Background(), key, val, 0).Err()
		return val, err
	}
	if err != nil {
		return false, err
	}

	return val, nil
}

func (m *RedisMiddleware) IsKarmaBlockListed(guildID, userID string) (ok bool, err error) {
	var key = fmt.Sprintf("%s:%s:%s", keyKarmaBlockListed, guildID, userID)

	err = m.client.Get(context.Background(), key).Scan(&ok)
	if err == redis.Nil {
		ok, err = m.Database.IsKarmaBlockListed(guildID, userID)
		if err != nil {
			return
		}

		err = m.client.SetBit(context.Background(), key, 0, boolutil.AsInt(ok)).Err()
		return
	}
	if err != nil {
		return
	}

	return
}

func (m *RedisMiddleware) AddKarmaBlockList(guildID, userID string) (err error) {
	var key = fmt.Sprintf("%s:%s:%s", keyKarmaBlockListed, guildID, userID)

	if err = m.client.Set(context.Background(), key, true, 0).Err(); err != nil {
		return
	}

	return m.Database.AddKarmaBlockList(guildID, userID)
}

func (m *RedisMiddleware) RemoveKarmaBlockList(guildID, userID string) (err error) {
	var key = fmt.Sprintf("%s:%s:%s", keyKarmaBlockListed, guildID, userID)

	if err = m.client.Set(context.Background(), key, false, 0).Err(); err != nil {
		return
	}

	return m.Database.RemoveKarmaBlockList(guildID, userID)
}

func (m *RedisMiddleware) SetAntiraidState(guildID string, state bool) error {
	var key = fmt.Sprintf("%s:%s", keyAntiraidState, guildID)

	if err := m.client.Set(context.Background(), key, state, 0).Err(); err != nil {
		return err
	}

	return m.Database.SetAntiraidState(guildID, state)
}

func (m *RedisMiddleware) GetAntiraidState(guildID string) (bool, error) {
	var key = fmt.Sprintf("%s:%s", keyAntiraidState, guildID)

	var val bool
	err := m.client.Get(context.Background(), key).Scan(&val)
	if err == redis.Nil {
		val, err = m.Database.GetAntiraidState(guildID)
		if err != nil {
			return false, err
		}

		err = m.client.Set(context.Background(), key, val, 0).Err()
		return val, err
	}
	if err != nil {
		return false, err
	}

	return val, nil
}

func (m *RedisMiddleware) SetAntiraidRegeneration(guildID string, limit int) error {
	var key = fmt.Sprintf("%s:%s", keyAntiraidLimit, guildID)

	if err := m.client.Set(context.Background(), key, limit, 0).Err(); err != nil {
		return err
	}

	return m.Database.SetKarmaTokens(guildID, limit)
}

func (m *RedisMiddleware) GetAntiraidRegeneration(guildID string) (int, error) {
	var key = fmt.Sprintf("%s:%s", keyAntiraidLimit, guildID)

	var val int
	err := m.client.Get(context.Background(), key).Scan(&val)
	if err == redis.Nil {
		val, err = m.Database.GetAntiraidRegeneration(guildID)
		if err != nil {
			return 0, err
		}

		err = m.client.Set(context.Background(), key, val, 0).Err()
		return val, err
	}
	if err != nil {
		return 0, err
	}

	return val, nil
}

func (m *RedisMiddleware) SetAntiraidBurst(guildID string, burst int) error {
	var key = fmt.Sprintf("%s:%s", keyAntiraidBurst, guildID)

	if err := m.client.Set(context.Background(), key, burst, 0).Err(); err != nil {
		return err
	}

	return m.Database.SetAntiraidBurst(guildID, burst)
}

func (m *RedisMiddleware) GetAntiraidBurst(guildID string) (int, error) {
	var key = fmt.Sprintf("%s:%s", keyAntiraidBurst, guildID)

	var val int
	err := m.client.Get(context.Background(), key).Scan(&val)
	if err == redis.Nil {
		val, err = m.Database.GetAntiraidBurst(guildID)
		if err != nil {
			return 0, err
		}

		err = m.client.Set(context.Background(), key, val, 0).Err()
		return val, err
	}
	if err != nil {
		return 0, err
	}

	return val, nil
}

func (m *RedisMiddleware) GetUserOTAEnabled(userID string) (bool, error) {
	var key = fmt.Sprintf("%s:%s", keyUserEnableOTA, userID)

	var val bool
	err := m.client.Get(context.Background(), key).Scan(&val)
	if err == redis.Nil {
		val, err = m.Database.GetUserOTAEnabled(userID)
		if err != nil {
			return false, err
		}

		err = m.client.Set(context.Background(), key, val, 0).Err()
		return val, err
	}
	if err != nil {
		return false, err
	}

	return val, nil
}

func (m *RedisMiddleware) SetUserOTAEnabled(userID string, enabled bool) error {
	var key = fmt.Sprintf("%s:%s", keyUserEnableOTA, userID)

	if err := m.client.Set(context.Background(), key, enabled, 0).Err(); err != nil {
		return err
	}

	return m.Database.SetUserOTAEnabled(userID, enabled)
}

func (m *RedisMiddleware) GetStarboardConfig(guildID string) (config *models.StarboardConfig, err error) {
	var key = fmt.Sprintf("%s:%s", keyGuildStarboardConfig, guildID)

	config = new(models.StarboardConfig)
	var configB []byte
	err = m.client.Get(context.Background(), key).Scan(&configB)
	if err == redis.Nil {
		config, err = m.Database.GetStarboardConfig(guildID)
		if err != nil {
			return
		}
		if configB, err = json.Marshal(config); err != nil {
			return
		}
		err = m.client.Set(context.Background(), key, configB, 0).Err()
		return
	}
	if err != nil {
		return
	}

	err = json.Unmarshal(configB, config)
	return
}

func (m *RedisMiddleware) SetStarboardConfig(config *models.StarboardConfig) (err error) {
	var key = fmt.Sprintf("%s:%s", keyGuildStarboardConfig, config.GuildID)
	configB, err := json.Marshal(config)
	if err != nil {
		return
	}
	err = m.client.Set(context.Background(), key, configB, 0).Err()
	return
}

func (r *RedisMiddleware) GetGuildLogDisable(guildID string) (bool, error) {
	var key = fmt.Sprintf("%s:%s", keyGuildLogEnable, guildID)

	var val bool
	err := r.client.Get(context.Background(), key).Scan(&val)
	if err == redis.Nil {
		val, err = r.Database.GetGuildLogDisable(guildID)
		if err != nil {
			return false, err
		}

		err = r.client.Set(context.Background(), key, val, 0).Err()
		return val, err
	}
	if err != nil {
		return false, err
	}

	return val, nil
}

func (r *RedisMiddleware) SetGuildLogDisable(guildID string, enabled bool) error {
	var key = fmt.Sprintf("%s:%s", keyGuildLogEnable, guildID)

	if err := r.client.Set(context.Background(), key, enabled, 0).Err(); err != nil {
		return err
	}

	return r.Database.SetGuildLogDisable(guildID, enabled)
}
