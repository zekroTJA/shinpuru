package mysql

import (
	"database/sql"
)

var migrationFuncs = []migrationFunc{
	migration_0,
	migration_1,
	migration_2,
	migration_3,
}

// VERSION 0:
// - base state
func migration_0(m *sql.Tx) (err error) {
	return
}

// VERSION 1:
// - add property `deleted` to `starboardEntries`
func migration_1(m *sql.Tx) (err error) {
	return createTableColumnIfNotExists(m,
		"starboardEntries", "`deleted` int(1) NOT NULL DEFAULT '0'")
}

// VERSION 2:
// - add property `karmaGain` to `starboardConfig`
func migration_2(m *sql.Tx) (err error) {
	return createTableColumnIfNotExists(m,
		"starboardConfig", "`karmaGain` int(16) NOT NULL DEFAULT '3'")
}

// VERSION 3:
// - add property `guildlog` to `guilds`
func migration_3(m *sql.Tx) (err error) {
	return createTableColumnIfNotExists(m,
		"guilds", "`guildlogDisable` text NOT NULL DEFAULT '0'")
}
