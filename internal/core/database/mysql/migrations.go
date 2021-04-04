package mysql

import "database/sql"

var migrationFuncs = []migrationFunc{
	migration_0,
	migration_1,
}

func migration_0(m *sql.Tx) (err error) {
	return
}

func migration_1(m *sql.Tx) (err error) {
	_, err = m.Exec(
		"ALTER TABLE starboardEntries " +
			"ADD COLUMN deleted int(1) NOT NULL DEFAULT '0'")

	return
}
