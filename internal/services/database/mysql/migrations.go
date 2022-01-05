package mysql

import (
	"database/sql"
)

var migrationFuncs = []migrationFunc{
	migration_0,
	migration_1,
	migration_2,
	migration_3,
	migration_4,
	migration_5,
	migration_6,
	migration_7,
	migration_8,
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

// VERSION 4:
// - add property `penalty` to `karmaSettings`
func migration_4(m *sql.Tx) (err error) {
	return createTableColumnIfNotExists(m,
		"karmaSettings", "`penalty` int(1) NOT NULL DEFAULT '0'")
}

// VERSION 5:
// - add property `timeout` to `reports`
func migration_5(m *sql.Tx) (err error) {
	return createTableColumnIfNotExists(m,
		"reports", "`timeout` timestamp NULL DEFAULT NULL")
}

// VERSION 6:
// - add property `timeout` to `reports`
func migration_6(m *sql.Tx) (err error) {
	_, err = m.Exec(`SELECT 1 FROM antiraidJoinlog WHERE iid > -1`)
	if err == nil {
		return
	}

	_, err = m.Exec(`ALTER TABLE antiraidJoinlog DROP PRIMARY KEY;`)
	if err != nil {
		return
	}
	_, err = m.Exec(`ALTER TABLE antiraidJoinlog ADD COLUMN iid int(11) NOT NULL AUTO_INCREMENT PRIMARY KEY;`)
	return
}

// VERSION 7:
// - add property `accountCreated` to `antiraidJoinlog`
func migration_7(m *sql.Tx) (err error) {
	return createTableColumnIfNotExists(m,
		"antiraidJoinlog", "`accountCreated` timestamp NOT NULL DEFAULT 0")
}

// VERSION 8:
// - add property `verified` to `users`
// - add property `requireUserVerification` to `guilds`
func migration_8(m *sql.Tx) (err error) {
	err = createTableColumnIfNotExists(m,
		"users", "`verified` text NOT NULL DEFAULT '0'")
	if err != nil {
		return
	}
	return createTableColumnIfNotExists(m,
		"guilds", "`requireUserVerification` text NOT NULL DEFAULT ''")
}
