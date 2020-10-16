package middleware

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/zekroTJA/shinpuru/pkg/multierror"

	"github.com/zekroTJA/shinpuru/internal/core/config"
	"github.com/zekroTJA/shinpuru/internal/util"

	_ "github.com/mattn/go-sqlite3"
)

// SqliteMiddleware implements the Database interface
// for SQLite.
type SqliteMiddleware struct {
	*MysqlMiddleware
}

func (m *SqliteMiddleware) setup() {
	mErr := multierror.New(nil)

	_, err := m.db.Exec("CREATE TABLE IF NOT EXISTS `guilds` (" +
		"`guildID` varchar(25) NOT NULL PRIMARY KEY," +
		"`prefix` text NOT NULL DEFAULT ''," +
		"`autorole` text NOT NULL DEFAULT ''," +
		"`modlogchanID` text NOT NULL DEFAULT ''," +
		"`voicelogchanID` text NOT NULL DEFAULT ''," +
		"`muteRoleID` text NOT NULL DEFAULT ''," +
		"`notifyRoleID` text NOT NULL DEFAULT ''," +
		"`ghostPingMsg` text NOT NULL DEFAULT ''," +
		"`jdoodleToken` text NOT NULL DEFAULT ''," +
		"`backup` text NOT NULL DEFAULT ''," +
		"`inviteBlock` text NOT NULL DEFAULT ''," +
		"`joinMsg` text NOT NULL DEFAULT ''," +
		"`leaveMsg` text NOT NULL DEFAULT ''," +
		"`colorReaction` text NOT NULL DEFAULT ''" +
		");")
	mErr.Append(err)

	_, err = m.db.Exec("CREATE TABLE IF NOT EXISTS `permissions` (" +
		"`roleID` varchar(25) NOT NULL PRIMARY KEY," +
		"`guildID` text NOT NULL DEFAULT ''," +
		"`permission` text NOT NULL DEFAULT ''" +
		");")
	mErr.Append(err)

	_, err = m.db.Exec("CREATE TABLE IF NOT EXISTS `reports` (" +
		"`id` varchar(25) NOT NULL PRIMARY KEY," +
		"`type` int(11) NOT NULL DEFAULT '3'," +
		"`guildID` text NOT NULL DEFAULT ''," +
		"`executorID` text NOT NULL DEFAULT ''," +
		"`victimID` text NOT NULL DEFAULT ''," +
		"`msg` text NOT NULL DEFAULT ''," +
		"`attachment` text NOT NULL DEFAULT ''" +
		");")
	mErr.Append(err)

	_, err = m.db.Exec("CREATE TABLE IF NOT EXISTS `settings` (" +
		"`iid` INTEGER PRIMARY KEY AUTOINCREMENT," +
		"`setting` text NOT NULL DEFAULT ''," +
		"`value` text NOT NULL DEFAULT ''" +
		");")
	mErr.Append(err)

	_, err = m.db.Exec("CREATE TABLE IF NOT EXISTS `votes` (" +
		"`id` varchar(25) NOT NULL PRIMARY KEY," +
		"`data` mediumtext NOT NULL DEFAULT ''" +
		");")
	mErr.Append(err)

	_, err = m.db.Exec("CREATE TABLE IF NOT EXISTS `twitchnotify` (" +
		"`iid` INTEGER PRIMARY KEY AUTOINCREMENT," +
		"`guildID` text NOT NULL DEFAULT ''," +
		"`channelID` text NOT NULL DEFAULT ''," +
		"`twitchUserID` text NOT NULL DEFAULT ''" +
		");")
	mErr.Append(err)

	_, err = m.db.Exec("CREATE TABLE IF NOT EXISTS `backups` (" +
		"`iid` INTEGER PRIMARY KEY AUTOINCREMENT," +
		"`guildID` text NOT NULL DEFAULT ''," +
		"`timestamp` bigint(20) NOT NULL DEFAULT 0," +
		"`fileID` text NOT NULL DEFAULT ''" +
		");")
	mErr.Append(err)

	_, err = m.db.Exec("CREATE TABLE IF NOT EXISTS `tags` (" +
		"`id` varchar(25) NOT NULL PRIMARY KEY," +
		"`ident` text NOT NULL DEFAULT ''," +
		"`creatorID` text NOT NULL DEFAULT ''," +
		"`guildID` text NOT NULL DEFAULT ''," +
		"`content` text NOT NULL DEFAULT ''," +
		"`created` bigint(20) NOT NULL DEFAULT 0," +
		"`lastEdit` bigint(20) NOT NULL DEFAULT 0" +
		");")
	mErr.Append(err)

	_, err = m.db.Exec("CREATE TABLE IF NOT EXISTS `apitokens` (" +
		"`userID` varchar(25) NOT NULL PRIMARY KEY," +
		"`salt` text NOT NULL," +
		"`created` timestamp NOT NULL," +
		"`expires` timestamp NOT NULL," +
		"`lastAccess` timestamp NOT NULL," +
		"`hits` bigint(20) NOT NULL" +
		");")
	mErr.Append(err)

	_, err = m.db.Exec("CREATE TABLE IF NOT EXISTS `karma` (" +
		"`iid` INTEGER PRIMARY KEY AUTOINCREMENT," +
		"`guildID` text NOT NULL DEFAULT ''," +
		"`userID` text NOT NULL DEFAULT ''," +
		"`value` bigint(20) NOT NULL DEFAULT '0'" +
		");")
	mErr.Append(err)

	_, err = m.db.Exec("CREATE TABLE IF NOT EXISTS `karmaSettings` (" +
		"`guildID` varchar(25) NOT NULL PRIMARY KEY," +
		"`state` int(1) NOT NULL DEFAULT '1'," +
		"`emotes` text NOT NULL DEFAULT ''," +
		"`tokens` bigint(20) NOT NULL DEFAULT '1'" +
		");")
	mErr.Append(err)

	_, err = m.db.Exec("CREATE TABLE IF NOT EXISTS `chanlock` (" +
		"`chanID` varchar(25) NOT NULL PRIMARY KEY," +
		"`guildID` text NOT NULL DEFAULT ''," +
		"`executorID` text NOT NULL DEFAULT ''," +
		"`permissions` text NOT NULL DEFAULT ''" +
		");")
	mErr.Append(err)

	if mErr.Len() > 0 {
		util.Log.Fatalf("Failed database setup: %s", mErr.Concat().Error())
	}
}

func (m *SqliteMiddleware) Connect(credentials ...interface{}) error {
	m.MysqlMiddleware = new(MysqlMiddleware)

	var err error

	creds := credentials[0].(*config.DatabaseFile)
	if creds == nil {
		return errors.New("Database credentials from config were nil")
	}

	dsn := fmt.Sprintf("file:" + creds.DBFile)
	m.db, err = sql.Open("sqlite3", dsn)
	m.setup()

	return err
}

func (m *SqliteMiddleware) Close() {
	if m.db != nil {
		m.db.Close()
	}
}
