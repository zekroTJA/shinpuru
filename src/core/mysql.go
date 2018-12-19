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

func (m *MySql) GetMemberPermissionLevel(guildID string, memberID string) (int, error) {
	var permLvl int
	err := m.DB.QueryRow("SELECT permlvl FROM guildmembers WHERE guilduserBlob = ?",
		guildID+"-"+memberID).Scan(&permLvl)
	if err == sql.ErrNoRows {
		err = ErrDatabaseNotFound
	}
	return permLvl, err
}
