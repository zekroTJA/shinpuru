package core

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/snowflake"
	"github.com/zekroTJA/shinpuru/internal/util"
)

var ErrDatabaseNotFound = errors.New("value not found")

var (
	MySqlDbSchemeB64  = ""
	SqliteDbSchemeB64 = ""
)

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

	GetGuildPermissions(guildID string) (map[string]int, error)
	SetGuildRolePermission(guildID, roleID string, permLvL int) error

	GetGuildJdoodleKey(guildID string) (string, error)
	SetGuildJdoodleKey(guildID, key string) error

	GetGuildBackup(guildID string) (bool, error)
	SetGuildBackup(guildID string, enabled bool) error

	GetGuildInviteBlock(guildID string) (string, error)
	SetGuildInviteBlock(guildID string, data string) error

	GetGuildJoinMsg(guildID string) (string, string, error)
	SetGuildJoinMsg(guildID string, msg string, channelID string) error

	GetGuildLeaveMsg(guildID string) (string, string, error)
	SetGuildLeaveMsg(guildID string, msg string, channelID string) error

	AddReport(rep *util.Report) error
	DeleteReport(id snowflake.ID) error
	GetReport(id snowflake.ID) (*util.Report, error)
	GetReportsGuild(guildID string) ([]*util.Report, error)
	GetReportsFiltered(guildID, memberID string, repType int) ([]*util.Report, error)

	GetMemberPermissionLevel(s *discordgo.Session, guildID string, memberID string) (int, error)

	GetSetting(setting string) (string, error)
	SetSetting(setting, value string) error

	GetVotes() (map[string]*util.Vote, error)

	AddUpdateVote(votes *util.Vote) error
	DeleteVote(voteID string) error

	GetMuteRoles() (map[string]string, error)
	GetMuteRoleGuild(guildID string) (string, error)
	SetMuteRole(guildID, roleID string) error

	GetAllTwitchNotifies(twitchUserID string) ([]*TwitchNotifyDBEntry, error)
	GetTwitchNotify(twitchUserID, guildID string) (*TwitchNotifyDBEntry, error)
	SetTwitchNotify(twitchNotify *TwitchNotifyDBEntry) error
	DeleteTwitchNotify(twitchUserID, guildID string) error

	AddBackup(guildID, fileID string) error
	DeleteBackup(guildID, fileID string) error
	GetBackups(guildID string) ([]*BackupEntry, error)
	GetBackupGuilds() ([]string, error)

	AddTag(tag *util.Tag) error
	EditTag(tag *util.Tag) error
	GetTagByID(id snowflake.ID) (*util.Tag, error)
	GetTagByIdent(ident string, guildID string) (*util.Tag, error)
	GetGuildTags(guildID string) ([]*util.Tag, error)
	DeleteTag(id snowflake.ID) error
}

func IsErrDatabaseNotFound(err error) bool {
	return err == ErrDatabaseNotFound
}
