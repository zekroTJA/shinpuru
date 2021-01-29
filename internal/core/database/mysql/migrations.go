package mysql

import "database/sql"

var migrationFuncs = []migrationFunc{
	migration_0,
}

func migration_0(m *sql.Tx) (err error) {
	return
}
