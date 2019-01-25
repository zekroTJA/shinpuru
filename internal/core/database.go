package core

import (
	"errors"

	"github.com/bwmarrin/discordgo"
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

	GetGuildPermissions(guildID string) (map[string]int, error)
	SetGuildRolePermission(guildID, roleID string, permLvL int) error

	AddReport(rep *util.Report) error
	GetReportsGuild(guildID string) ([]*util.Report, error)
	GetReportsFiltered(guildID, memberID string, repType int) ([]*util.Report, error)

	GetMemberPermissionLevel(s *discordgo.Session, guildID string, memberID string) (int, error)

	GetSetting(setting string) (string, error)
	SetSetting(setting, value string) error

	GetVotes() (map[string]*util.Vote, error)
	// SetVotes(votes []*util.Vote) error
	AddUpdateVote(votes *util.Vote) error
	DeleteVote(voteID string) error

	GetMuteRoles() (map[string]string, error)
	GetMuteRoleGuild(guildID string) (string, error)
	SetMuteRole(guildID, roleID string) error
}

func IsErrDatabaseNotFound(err error) bool {
	return err == ErrDatabaseNotFound
}
