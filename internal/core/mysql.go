package core

import (
	"database/sql"
	"errors"
	"fmt"

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

func (m *MySql) GetMemberPermissionLevel(guildID string, memberID string) (int, error) {
	var permLvl int
	err := m.DB.QueryRow("SELECT permlvl FROM guildmembers WHERE guilduserBlob = ?",
		guildID+"-"+memberID).Scan(&permLvl)
	if err == sql.ErrNoRows {
		err = ErrDatabaseNotFound
	}
	return permLvl, err
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
