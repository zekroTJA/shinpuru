package core

import (
	"database/sql"
	"errors"
	"fmt"

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

func (m *MySql) GetGuildPrefix(guildID string) (string, error) {
	var prefix string
	err := m.DB.QueryRow("SELECT prefix FROM guilds WHERE guildID = ?", guildID).Scan(&prefix)
	if err == sql.ErrNoRows {
		err = ErrDatabaseNotFound
	}
	return prefix, err
}

func (m *MySql) SetGuildPrefix(guildID, newPrefix string) error {
	res, err := m.DB.Exec("UPDATE guilds SET prefix = ? WHERE guildID = ?", newPrefix, guildID)
	if ar, err := res.RowsAffected(); ar == 0 {
		if err != nil {
			return err
		}
		_, err := m.DB.Exec("INSERT INTO guilds (guildID, prefix) VALUES (?, ?)", guildID, newPrefix)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return err
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
