package mysql

import "database/sql"

var migrationFuncs = []migrationFunc{
	migration_0,
	migration_1,
}

// VERSION 0
// - base state
func migration_0(m *sql.Tx) (err error) {
	return
}

// VERSION 1:
// - add property `deleted` to `starboardEntries`
func migration_1(m *sql.Tx) (err error) {
	_, err = m.Exec(
		"ALTER TABLE starboardEntries " +
			"ADD COLUMN deleted int(1) NOT NULL DEFAULT '0'")

	return
}
