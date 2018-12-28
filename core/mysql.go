package core

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/zekroTJA/shinpuru/util"

	"github.com/bwmarrin/discordgo"
	_ "github.com/go-sql-driver/mysql"
)

type MySql struct {
	DB *sql.DB
}

func (m *MySql) Connect(credentials ...interface{}) error {
	var err error
	creds := credentials[0].(*ConfigDatabase)
	if creds == nil {
		return errors.New("Database credentials from config were nil")
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", creds.User, creds.Password, creds.Host, creds.Database)
	m.DB, err = sql.Open("mysql", dsn)
	return err
}

func (m *MySql) Close() {
	if m.DB != nil {
		m.DB.Close()
	}
}

func (m *MySql) getGuildSetting(guildID, key string) (string, error) {
	var value string
	err := m.DB.QueryRow("SELECT "+key+" FROM guilds WHERE guildID = ?", guildID).Scan(&value)
	if err == sql.ErrNoRows {
		err = ErrDatabaseNotFound
	}
	return value, err
}

func (m *MySql) setGuildSetting(guildID, key string, value string) error {
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

func (m *MySql) GetGuildPrefix(guildID string) (string, error) {
	val, err := m.getGuildSetting(guildID, "prefix")
	return val, err
}

func (m *MySql) SetGuildPrefix(guildID, newPrefix string) error {
	return m.setGuildSetting(guildID, "prefix", newPrefix)
}

func (m *MySql) GetGuildAutoRole(guildID string) (string, error) {
	val, err := m.getGuildSetting(guildID, "autorole")
	return val, err
}

func (m *MySql) SetGuildAutoRole(guildID, autoRoleID string) error {
	return m.setGuildSetting(guildID, "autorole", autoRoleID)
}

func (m *MySql) GetGuildModLog(guildID string) (string, error) {
	val, err := m.getGuildSetting(guildID, "modlogchanID")
	return val, err
}

func (m *MySql) SetGuildModLog(guildID, chanID string) error {
	return m.setGuildSetting(guildID, "modlogchanID", chanID)
}

func (m *MySql) GetMemberPermissionLevel(s *discordgo.Session, guildID string, memberID string) (int, error) {
	guildPerms, err := m.GetGuildPermissions(guildID)
	if err != nil {
		return 0, err
	}
	member, err := s.GuildMember(guildID, memberID)
	if err != nil {
		return 0, err
	}
	maxPermLvl := 0
	for _, rID := range member.Roles {
		if lvl, ok := guildPerms[rID]; ok && lvl > maxPermLvl {
			maxPermLvl = lvl
		}
	}
	return maxPermLvl, err
}

func (m *MySql) GetGuildPermissions(guildID string) (map[string]int, error) {
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

func (m *MySql) SetGuildRolePermission(guildID, roleID string, permLvL int) error {
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

func (m *MySql) GetSetting(setting string) (string, error) {
	var value string
	err := m.DB.QueryRow("SELECT value FROM settings WHERE setting = ?", setting).Scan(&value)
	if err == sql.ErrNoRows {
		err = ErrDatabaseNotFound
	}
	return value, err
}

func (m *MySql) SetSetting(setting, value string) error {
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

func (m *MySql) AddReport(rep *util.Report) error {
	_, err := m.DB.Exec("INSERT INTO reports (id, type, guildID, executorID, victimID, msg) VALUES (?, ?, ?, ?, ?, ?)",
		rep.ID, rep.Type, rep.GuildID, rep.ExecutorID, rep.VictimID, rep.Msg)
	return err
}

// TODO: functionality
func (m *MySql) GetReportsGuild(guildID string) ([]*util.Report, error) {
	return nil, nil
}
