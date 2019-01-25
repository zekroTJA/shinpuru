package core

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/zekroTJA/shinpuru/internal/util"

	"github.com/bwmarrin/discordgo"
	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	DB *sql.DB
}

func (m *Sqlite) setup() {
	if SqliteDbSchemeB64 == "" {
		util.Log.Warning("sqlite database scheme was not set on compiling. Database can not be checked for structure changes!")
		return
	}
	scheme, err := base64.StdEncoding.DecodeString(SqliteDbSchemeB64)
	if err != nil {
		util.Log.Fatal("failed decoding base64 database scheme: ", err)
		return
	}
	for _, query := range strings.Split(string(scheme), ";") {
		if ok, _ := regexp.MatchString(`\w`, query); ok {
			_, err = m.DB.Exec(query)
			if err != nil {
				util.Log.Error("Failed executing setup database query: ", err)
			}
		}
	}
}

func (m *Sqlite) Connect(credentials ...interface{}) error {
	var err error
	creds := credentials[0].(*ConfigDatabaseFile)
	if creds == nil {
		return errors.New("Database credentials from config were nil")
	}
	dsn := fmt.Sprintf("file:" + creds.DBFile)
	m.DB, err = sql.Open("sqlite3", dsn)
	m.setup()
	return err
}

func (m *Sqlite) Close() {
	if m.DB != nil {
		m.DB.Close()
	}
}

func (m *Sqlite) getGuildSetting(guildID, key string) (string, error) {
	var value string
	err := m.DB.QueryRow("SELECT "+key+" FROM guilds WHERE guildID = ?", guildID).Scan(&value)
	if err == sql.ErrNoRows {
		err = ErrDatabaseNotFound
	}
	return value, err
}

func (m *Sqlite) setGuildSetting(guildID, key string, value string) error {
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

func (m *Sqlite) GetGuildPrefix(guildID string) (string, error) {
	val, err := m.getGuildSetting(guildID, "prefix")
	return val, err
}

func (m *Sqlite) SetGuildPrefix(guildID, newPrefix string) error {
	return m.setGuildSetting(guildID, "prefix", newPrefix)
}

func (m *Sqlite) GetGuildAutoRole(guildID string) (string, error) {
	val, err := m.getGuildSetting(guildID, "autorole")
	return val, err
}

func (m *Sqlite) SetGuildAutoRole(guildID, autoRoleID string) error {
	return m.setGuildSetting(guildID, "autorole", autoRoleID)
}

func (m *Sqlite) GetGuildModLog(guildID string) (string, error) {
	val, err := m.getGuildSetting(guildID, "modlogchanID")
	return val, err
}

func (m *Sqlite) SetGuildModLog(guildID, chanID string) error {
	return m.setGuildSetting(guildID, "modlogchanID", chanID)
}

func (m *Sqlite) GetGuildVoiceLog(guildID string) (string, error) {
	val, err := m.getGuildSetting(guildID, "voicelogchanID")
	return val, err
}

func (m *Sqlite) SetGuildVoiceLog(guildID, chanID string) error {
	return m.setGuildSetting(guildID, "voicelogchanID", chanID)
}

func (m *Sqlite) GetGuildNotifyRole(guildID string) (string, error) {
	val, err := m.getGuildSetting(guildID, "notifyRoleID")
	return val, err
}

func (m *Sqlite) SetGuildNotifyRole(guildID, roleID string) error {
	return m.setGuildSetting(guildID, "notifyRoleID", roleID)
}

func (m *Sqlite) GetMemberPermissionLevel(s *discordgo.Session, guildID string, memberID string) (int, error) {
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

func (m *Sqlite) GetGuildPermissions(guildID string) (map[string]int, error) {
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

func (m *Sqlite) SetGuildRolePermission(guildID, roleID string, permLvL int) error {
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

func (m *Sqlite) GetSetting(setting string) (string, error) {
	var value string
	err := m.DB.QueryRow("SELECT value FROM settings WHERE setting = ?", setting).Scan(&value)
	if err == sql.ErrNoRows {
		err = ErrDatabaseNotFound
	}
	return value, err
}

func (m *Sqlite) SetSetting(setting, value string) error {
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

func (m *Sqlite) AddReport(rep *util.Report) error {
	_, err := m.DB.Exec("INSERT INTO reports (id, type, guildID, executorID, victimID, msg) VALUES (?, ?, ?, ?, ?, ?)",
		rep.ID, rep.Type, rep.GuildID, rep.ExecutorID, rep.VictimID, rep.Msg)
	return err
}

func (m *Sqlite) GetReportsGuild(guildID string) ([]*util.Report, error) {
	rows, err := m.DB.Query("SELECT * FROM reports WHERE guildID = ?", guildID)
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

func (m *Sqlite) GetReportsFiltered(guildID, memberID string, repType int) ([]*util.Report, error) {
	query := fmt.Sprintf(`SELECT * FROM reports WHERE guildID = "%s"`, guildID)
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

func (m *Sqlite) GetVotes() (map[string]*util.Vote, error) {
	rows, err := m.DB.Query("SELECT * FROM votes")
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

func (m *Sqlite) AddUpdateVote(vote *util.Vote) error {
	rawData, err := vote.Marshal()
	if err != nil {
		return err
	}
	res, err := m.DB.Exec("UPDATE votes SET data = ? WHERE ID = ?", rawData, vote.ID)
	if ar, err := res.RowsAffected(); ar == 0 {
		if err != nil {
			return err
		}
		_, err := m.DB.Exec("INSERT INTO votes (ID, data) VALUES (?, ?)", vote.ID, rawData)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return err
}

func (m *Sqlite) DeleteVote(voteID string) error {
	_, err := m.DB.Exec("DELETE FROM votes WHERE ID = ?", voteID)
	return err
}

// func (m *Sqlite) SetVotes(updatedVotes []*util.Vote) error {
// 	dbVotes, err := m.GetVotes()
// 	if err != nil {
// 		return err
// 	}

// 	toDelete := make(map[string]*util.Vote)
// 	for _, dbV := range dbVotes {

// 	}

// 	return nil
// }

func (m *Sqlite) GetMuteRoles() (map[string]string, error) {
	rows, err := m.DB.Query("SELECT guildID, muteRoleID FROM guilds")
	results := make(map[string]string)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var guildID, roleID string
		err = rows.Scan(&guildID, &roleID)
		if err != nil {
			results[guildID] = roleID
		}
	}
	return results, nil
}

func (m *Sqlite) GetMuteRoleGuild(guildID string) (string, error) {
	val, err := m.getGuildSetting(guildID, "muteRoleID")
	return val, err
}

func (m *Sqlite) SetMuteRole(guildID, roleID string) error {
	return m.setGuildSetting(guildID, "muteRoleID", roleID)
}
