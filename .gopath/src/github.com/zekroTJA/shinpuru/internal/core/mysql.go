package core

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/zekroTJA/shinpuru/internal/util"

	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/snowflake"
	_ "github.com/go-sql-driver/mysql"
)

type MySQL struct {
	DB *sql.DB
}

func (m *MySQL) setup() {
	if MySqlDbSchemeB64 == "" {
		util.Log.Warning("MySql database scheme was not set on compiling. Database can not be checked for structure changes!")
		return
	}
	scheme, err := base64.StdEncoding.DecodeString(MySqlDbSchemeB64)
	if err != nil {
		util.Log.Fatal("failed decoding base64 database scheme: ", err)
		return
	}
	for _, query := range strings.Split(string(scheme), ";") {
		if ok, _ := regexp.MatchString(`\w`, query); ok {
			fmt.Println(query)
			_, err = m.DB.Exec(query)
			if err != nil {
				util.Log.Errorf("Failed executing setup database query '%s': %s", query, err.Error())
			}
		}
	}
}

func (m *MySQL) Connect(credentials ...interface{}) error {
	var err error
	creds := credentials[0].(*ConfigDatabaseCreds)
	if creds == nil {
		return errors.New("Database credentials from config were nil")
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", creds.User, creds.Password, creds.Host, creds.Database)
	m.DB, err = sql.Open("mysql", dsn)
	m.setup()
	return err
}

func (m *MySQL) Close() {
	if m.DB != nil {
		m.DB.Close()
	}
}

func (m *MySQL) getGuildSetting(guildID, key string) (string, error) {
	var value string
	err := m.DB.QueryRow("SELECT "+key+" FROM guilds WHERE guildID = ?", guildID).Scan(&value)
	if err == sql.ErrNoRows {
		err = ErrDatabaseNotFound
	}
	return value, err
}

func (m *MySQL) setGuildSetting(guildID, key string, value string) error {
	res, err := m.DB.Exec("UPDATE guilds SET "+key+" = ? WHERE guildID = ?", value, guildID)
	if err != nil {
		return err
	}
	if ar, err := res.RowsAffected(); ar == 0 {
		if err != nil {
			return err
		}
		_, err := m.DB.Exec("INSERT INTO guilds (guildID, "+key+") VALUES (?, ?)", guildID, value)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return err
}

func (m *MySQL) GetGuildPrefix(guildID string) (string, error) {
	val, err := m.getGuildSetting(guildID, "prefix")
	return val, err
}

func (m *MySQL) SetGuildPrefix(guildID, newPrefix string) error {
	return m.setGuildSetting(guildID, "prefix", newPrefix)
}

func (m *MySQL) GetGuildAutoRole(guildID string) (string, error) {
	val, err := m.getGuildSetting(guildID, "autorole")
	return val, err
}

func (m *MySQL) SetGuildAutoRole(guildID, autoRoleID string) error {
	return m.setGuildSetting(guildID, "autorole", autoRoleID)
}

func (m *MySQL) GetGuildModLog(guildID string) (string, error) {
	val, err := m.getGuildSetting(guildID, "modlogchanID")
	return val, err
}

func (m *MySQL) SetGuildModLog(guildID, chanID string) error {
	return m.setGuildSetting(guildID, "modlogchanID", chanID)
}

func (m *MySQL) GetGuildVoiceLog(guildID string) (string, error) {
	val, err := m.getGuildSetting(guildID, "voicelogchanID")
	return val, err
}

func (m *MySQL) SetGuildVoiceLog(guildID, chanID string) error {
	return m.setGuildSetting(guildID, "voicelogchanID", chanID)
}

func (m *MySQL) GetGuildNotifyRole(guildID string) (string, error) {
	val, err := m.getGuildSetting(guildID, "notifyRoleID")
	return val, err
}

func (m *MySQL) SetGuildNotifyRole(guildID, roleID string) error {
	return m.setGuildSetting(guildID, "notifyRoleID", roleID)
}

func (m *MySQL) GetGuildGhostpingMsg(guildID string) (string, error) {
	val, err := m.getGuildSetting(guildID, "ghostPingMsg")
	return val, err
}

func (m *MySQL) SetGuildGhostpingMsg(guildID, msg string) error {
	return m.setGuildSetting(guildID, "ghostPingMsg", msg)
}

func (m *MySQL) GetMemberPermissionLevel(s *discordgo.Session, guildID string, memberID string) (int, error) {
	guildPerms, err := m.GetGuildPermissions(guildID)
	if err != nil {
		return 0, err
	}
	member, err := s.GuildMember(guildID, memberID)
	if err != nil {
		return 0, err
	}
	maxPermLvl := 0
	if lvl, ok := guildPerms[guildID]; ok {
		maxPermLvl = lvl
	}
	for _, rID := range member.Roles {
		if lvl, ok := guildPerms[rID]; ok && lvl > maxPermLvl {
			maxPermLvl = lvl
		}
	}
	return maxPermLvl, err
}

func (m *MySQL) GetGuildPermissions(guildID string) (map[string]int, error) {
	results := make(map[string]int)
	rows, err := m.DB.Query("SELECT roleID, permission FROM permissions WHERE guildID = ?",
		guildID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var roleID string
		var permission int
		err := rows.Scan(&roleID, &permission)
		if err != nil {
			return nil, err
		}
		results[roleID] = permission
	}
	return results, nil
}

func (m *MySQL) SetGuildRolePermission(guildID, roleID string, permLvL int) error {
	res, err := m.DB.Exec("UPDATE permissions SET permission = ? WHERE roleID = ? AND guildID = ?",
		permLvL, roleID, guildID)
	if err != nil {
		return err
	}
	if ar, err := res.RowsAffected(); ar == 0 {
		if err != nil {
			return err
		}
		_, err := m.DB.Exec("INSERT INTO permissions (roleID, guildID, permission) VALUES (?, ?, ?)",
			roleID, guildID, permLvL)
		return err
	}
	return nil
}

func (m *MySQL) GetGuildJdoodleKey(guildID string) (string, error) {
	val, err := m.getGuildSetting(guildID, "jdoodleToken")
	return val, err
}

func (m *MySQL) SetGuildJdoodleKey(guildID, key string) error {
	return m.setGuildSetting(guildID, "jdoodleToken", key)
}

func (m *MySQL) GetGuildBackup(guildID string) (bool, error) {
	val, err := m.getGuildSetting(guildID, "backup")
	return val != "", err
}

func (m *MySQL) SetGuildBackup(guildID string, enabled bool) error {
	var val string
	if enabled {
		val = "1"
	}
	return m.setGuildSetting(guildID, "backup", val)
}

func (m *MySQL) GetSetting(setting string) (string, error) {
	var value string
	err := m.DB.QueryRow("SELECT value FROM settings WHERE setting = ?", setting).Scan(&value)
	if err == sql.ErrNoRows {
		err = ErrDatabaseNotFound
	}
	return value, err
}

func (m *MySQL) SetSetting(setting, value string) error {
	res, err := m.DB.Exec("UPDATE settings SET value = ? WHERE setting = ?", value, setting)
	if ar, err := res.RowsAffected(); ar == 0 {
		if err != nil {
			return err
		}
		_, err := m.DB.Exec("INSERT INTO settings (setting, value) VALUES (?, ?)", setting, value)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return err
}

func (m *MySQL) AddReport(rep *util.Report) error {
	_, err := m.DB.Exec("INSERT INTO reports (id, type, guildID, executorID, victimID, msg) VALUES (?, ?, ?, ?, ?, ?)",
		rep.ID, rep.Type, rep.GuildID, rep.ExecutorID, rep.VictimID, rep.Msg)
	return err
}

func (m *MySQL) DeleteReport(id snowflake.ID) error {
	_, err := m.DB.Exec("DELETE FROM reports WHERE id = ?", id)
	return err
}

func (m *MySQL) GetReport(id snowflake.ID) (*util.Report, error) {
	rep := new(util.Report)

	row := m.DB.QueryRow("SELECT id, type, guildID, executorID, victimID, msg FROM reports WHERE id = ?", id)
	err := row.Scan(&rep.ID, &rep.Type, &rep.GuildID, &rep.ExecutorID, &rep.VictimID, &rep.Msg)
	if err == sql.ErrNoRows {
		return nil, ErrDatabaseNotFound
	}

	return rep, err
}

func (m *MySQL) GetReportsGuild(guildID string) ([]*util.Report, error) {
	rows, err := m.DB.Query("SELECT id, type, guildID, executorID, victimID, msg FROM reports WHERE guildID = ?", guildID)
	var results []*util.Report
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		rep := new(util.Report)
		err := rows.Scan(&rep.ID, &rep.Type, &rep.GuildID, &rep.ExecutorID, &rep.VictimID, &rep.Msg)
		if err != nil {
			return nil, err
		}
		results = append(results, rep)
	}
	return results, nil
}

func (m *MySQL) GetReportsFiltered(guildID, memberID string, repType int) ([]*util.Report, error) {
	query := fmt.Sprintf(`SELECT id, type, guildID, executorID, victimID, msg FROM reports WHERE guildID = "%s"`, guildID)
	if memberID != "" {
		query += fmt.Sprintf(` AND victimID = "%s"`, memberID)
	}
	if repType != -1 {
		query += fmt.Sprintf(` AND type = %d`, repType)
	}
	rows, err := m.DB.Query(query)
	var results []*util.Report
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		rep := new(util.Report)
		err := rows.Scan(&rep.ID, &rep.Type, &rep.GuildID, &rep.ExecutorID, &rep.VictimID, &rep.Msg)
		if err != nil {
			return nil, err
		}
		results = append(results, rep)
	}
	return results, nil
}

func (m *MySQL) GetVotes() (map[string]*util.Vote, error) {
	rows, err := m.DB.Query("SELECT id, data FROM votes")
	results := make(map[string]*util.Vote)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var voteID, rawData string
		err := rows.Scan(&voteID, &rawData)
		if err != nil {
			util.Log.Error("An error occured reading vote from database: ", err)
			continue
		}
		vote, err := util.VoteUnmarshal(rawData)
		if err != nil {
			m.DeleteVote(rawData)
		} else {
			results[vote.ID] = vote
		}
	}
	return results, err
}

func (m *MySQL) AddUpdateVote(vote *util.Vote) error {
	rawData, err := vote.Marshal()
	if err != nil {
		return err
	}
	res, err := m.DB.Exec("UPDATE votes SET data = ? WHERE id = ?", rawData, vote.ID)
	if ar, err := res.RowsAffected(); ar == 0 {
		if err != nil {
			return err
		}
		_, err := m.DB.Exec("INSERT INTO votes (id, data) VALUES (?, ?)", vote.ID, rawData)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return err
}

func (m *MySQL) DeleteVote(voteID string) error {
	_, err := m.DB.Exec("DELETE FROM votes WHERE id = ?", voteID)
	return err
}

func (m *MySQL) GetMuteRoles() (map[string]string, error) {
	rows, err := m.DB.Query("SELECT guildID, muteRoleID FROM guilds")
	results := make(map[string]string)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var guildID, roleID string
		err = rows.Scan(&guildID, &roleID)
		if err == nil {
			results[guildID] = roleID
		}
	}
	return results, nil
}

func (m *MySQL) GetMuteRoleGuild(guildID string) (string, error) {
	val, err := m.getGuildSetting(guildID, "muteRoleID")
	return val, err
}

func (m *MySQL) SetMuteRole(guildID, roleID string) error {
	return m.setGuildSetting(guildID, "muteRoleID", roleID)
}

func (m *MySQL) GetTwitchNotify(twitchUserID, guildID string) (*TwitchNotifyDBEntry, error) {
	t := &TwitchNotifyDBEntry{
		TwitchUserID: twitchUserID,
		GuildID:      guildID,
	}
	err := m.DB.QueryRow("SELECT channelID FROM twitchnotify WHERE twitchUserID = ? AND guildID = ?",
		twitchUserID, guildID).Scan(&t.ChannelID)
	if err == sql.ErrNoRows {
		err = ErrDatabaseNotFound
	}
	return t, err
}

func (m *MySQL) SetTwitchNotify(twitchNotify *TwitchNotifyDBEntry) error {
	res, err := m.DB.Exec("UPDATE twitchnotify SET channelID = ? WHERE twitchUserID = ? AND guildID = ?",
		twitchNotify.ChannelID, twitchNotify.TwitchUserID, twitchNotify.GuildID)
	if err != nil {
		return err
	}
	if ar, err := res.RowsAffected(); ar == 0 {
		if err != nil {
			return err
		}
		_, err := m.DB.Exec("INSERT INTO twitchnotify (twitchUserID, guildID, channelID) VALUES (?, ?, ?)",
			twitchNotify.TwitchUserID, twitchNotify.GuildID, twitchNotify.ChannelID)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return err
}

func (m *MySQL) DeleteTwitchNotify(twitchUserID, guildID string) error {
	_, err := m.DB.Exec("DELETE FROM twitchnotify WHERE twitchUserID = ? AND guildID = ?", twitchUserID, guildID)
	return err
}

func (m *MySQL) GetAllTwitchNotifies(twitchUserID string) ([]*TwitchNotifyDBEntry, error) {
	query := "SELECT twitchUserID, guildID, channelID FROM twitchnotify"
	if twitchUserID != "" {
		query += " WHERE twitchUserID = " + twitchUserID
	}
	rows, err := m.DB.Query(query)
	results := make([]*TwitchNotifyDBEntry, 0)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		t := new(TwitchNotifyDBEntry)
		err = rows.Scan(&t.TwitchUserID, &t.GuildID, &t.ChannelID)
		if err == nil {
			results = append(results, t)
		}
	}
	return results, nil
}

func (m *MySQL) AddBackup(guildID, fileID string) error {
	timestamp := time.Now().Unix()
	_, err := m.DB.Exec("INSERT INTO backups (guildID, timestamp, fileID) VALUES (?, ?, ?)", guildID, timestamp, fileID)
	return err
}

func (m *MySQL) DeleteBackup(guildID, fileID string) error {
	_, err := m.DB.Exec("DELETE FROM backups WHERE guildID = ? AND fileID = ?", guildID, fileID)
	return err
}

func (m *MySQL) GetGuildInviteBlock(guildID string) (string, error) {
	return m.getGuildSetting(guildID, "inviteBlock")
}

func (m *MySQL) SetGuildInviteBlock(guildID string, data string) error {
	return m.setGuildSetting(guildID, "inviteBlock", data)
}

func (m *MySQL) GetBackups(guildID string) ([]*BackupEntry, error) {
	rows, err := m.DB.Query("SELECT guildID, timestamp, fileID FROM backups WHERE guildID = ?", guildID)
	if err == sql.ErrNoRows {
		return nil, ErrDatabaseNotFound
	}
	if err != nil {
		return nil, err
	}

	backups := make([]*BackupEntry, 0)
	for rows.Next() {
		be := new(BackupEntry)
		var timeStampUnix int64
		err = rows.Scan(&be.GuildID, &timeStampUnix, &be.FileID)
		if err != nil {
			return nil, err
		}
		be.Timestamp = time.Unix(timeStampUnix, 0)
		backups = append(backups, be)
	}

	return backups, nil
}

func (m *MySQL) GetBackupGuilds() ([]string, error) {
	rows, err := m.DB.Query("SELECT guildID FROM guilds WHERE backup = '1'")
	if err == sql.ErrNoRows {
		return nil, ErrDatabaseNotFound
	}
	if err != nil {
		return nil, err
	}

	guilds := make([]string, 0)
	for rows.Next() {
		var s string
		err = rows.Scan(&s)
		if err != nil {
			return nil, err
		}
		guilds = append(guilds, s)
	}

	return guilds, err
}
