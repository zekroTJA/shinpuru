package core

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/pkg/multierror"

	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/snowflake"
	_ "github.com/go-sql-driver/mysql"
)

type MySQL struct {
	DB *sql.DB
}

func (m *MySQL) setup() {
	mErr := multierror.New(nil)

	_, err := m.DB.Exec("CREATE TABLE IF NOT EXISTS `guilds` (" +
		"`iid` int(11) NOT NULL AUTO_INCREMENT," +
		"`guildID` text NOT NULL," +
		"`prefix` text NOT NULL," +
		"`autorole` text NOT NULL," +
		"`modlogchanID` text NOT NULL," +
		"`voicelogchanID` text NOT NULL," +
		"`muteRoleID` text NOT NULL," +
		"`ghostPingMsg` text NOT NULL," +
		"`jdoodleToken` text NOT NULL," +
		"`backup` text NOT NULL," +
		"`inviteBlock` text NOT NULL," +
		"`joinMsg` text NOT NULL," +
		"`leaveMsg` text NOT NULL," +
		"PRIMARY KEY (`iid`)" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;")
	mErr.Append(err)

	_, err = m.DB.Exec("CREATE TABLE IF NOT EXISTS `permissions` (" +
		"`iid` int(11) NOT NULL AUTO_INCREMENT," +
		"`roleID` text NOT NULL," +
		"`guildID` text NOT NULL," +
		"`permission` text NOT NULL," +
		"PRIMARY KEY (`iid`)" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;")
	mErr.Append(err)

	_, err = m.DB.Exec("CREATE TABLE IF NOT EXISTS `reports` (" +
		"`iid` int(11) NOT NULL AUTO_INCREMENT," +
		"`id` text NOT NULL," +
		"`type` int(11) NOT NULL," +
		"`guildID` text NOT NULL," +
		"`executorID` text NOT NULL," +
		"`victimID` text NOT NULL," +
		"`msg` text NOT NULL," +
		"`attachment` text NOT NULL," +
		"PRIMARY KEY (`iid`)" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;")
	mErr.Append(err)

	_, err = m.DB.Exec("CREATE TABLE IF NOT EXISTS `settings` (" +
		"`iid` int(11) NOT NULL AUTO_INCREMENT," +
		"`setting` text NOT NULL," +
		"`value` text NOT NULL," +
		"PRIMARY KEY (`iid`)" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;")
	mErr.Append(err)

	_, err = m.DB.Exec("CREATE TABLE IF NOT EXISTS `starboard` (" +
		"`iid` int(11) NOT NULL AUTO_INCREMENT," +
		"`guildID` text NOT NULL," +
		"`chanID` text NOT NULL," +
		"`enabled` tinyint(1) NOT NULL DEFAULT '1'," +
		"`minimum` int(11) NOT NULL DEFAULT '5'," +
		"PRIMARY KEY (`iid`)" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;")
	mErr.Append(err)

	_, err = m.DB.Exec("CREATE TABLE IF NOT EXISTS `votes` (" +
		"`iid` int(11) NOT NULL AUTO_INCREMENT," +
		"`id` text NOT NULL," +
		"`data` mediumtext NOT NULL," +
		"PRIMARY KEY (`iid`)" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;")
	mErr.Append(err)

	_, err = m.DB.Exec("CREATE TABLE IF NOT EXISTS `twitchnotify` (" +
		"`iid` int(11) NOT NULL AUTO_INCREMENT," +
		"`guildID` text NOT NULL," +
		"`channelID` text NOT NULL," +
		"`twitchUserID` text NOT NULL," +
		"PRIMARY KEY (`iid`)" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;")
	mErr.Append(err)

	_, err = m.DB.Exec("CREATE TABLE IF NOT EXISTS `backups` (" +
		"`iid` int(11) NOT NULL AUTO_INCREMENT," +
		"`guildID` text NOT NULL," +
		"`timestamp` bigint(20) NOT NULL," +
		"`fileID` text NOT NULL," +
		"PRIMARY KEY (`iid`)" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;")
	mErr.Append(err)

	_, err = m.DB.Exec("CREATE TABLE IF NOT EXISTS `tags` (" +
		"`iid` int(11) NOT NULL AUTO_INCREMENT," +
		"`id` text NOT NULL," +
		"`ident` text NOT NULL," +
		"`creatorID` text NOT NULL," +
		"`guildID` text NOT NULL," +
		"`content` text NOT NULL," +
		"`created` bigint(20) NOT NULL," +
		"`lastEdit` bigint(20) NOT NULL," +
		"PRIMARY KEY (`iid`)" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;")
	mErr.Append(err)

	_, err = m.DB.Exec("CREATE TABLE IF NOT EXISTS `sessions` (" +
		"`iid` int(11) NOT NULL AUTO_INCREMENT," +
		"`sessionkey` text NOT NULL," +
		"`userID` text NOT NULL," +
		"`expires` timestamp NOT NULL," +
		"PRIMARY KEY (`iid`)" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;")
	mErr.Append(err)

	if mErr.Len() > 0 {
		util.Log.Fatalf("Failed database setup: %s", mErr.Concat().Error())
	}
}

func (m *MySQL) Connect(credentials ...interface{}) error {
	var err error
	creds := credentials[0].(*ConfigDatabaseCreds)
	if creds == nil {
		return errors.New("Database credentials from config were nil")
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?collation=utf8mb4_unicode_ci&parseTime=true",
		creds.User, creds.Password, creds.Host, creds.Database)
	m.DB, err = sql.Open("mysql", dsn)
	m.setup()
	return err
}

func (m *MySQL) Close() {
	if m.DB != nil {
		m.DB.Close()
	}
}

func (m *MySQL) getGuildSetting(guildID, key string) (string, error) {
	var value string
	err := m.DB.QueryRow("SELECT "+key+" FROM guilds WHERE guildID = ?", guildID).Scan(&value)
	if err == sql.ErrNoRows {
		err = ErrDatabaseNotFound
	}
	return value, err
}

func (m *MySQL) setGuildSetting(guildID, key string, value string) error {
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

func (m *MySQL) GetGuildPrefix(guildID string) (string, error) {
	val, err := m.getGuildSetting(guildID, "prefix")
	return val, err
}

func (m *MySQL) SetGuildPrefix(guildID, newPrefix string) error {
	return m.setGuildSetting(guildID, "prefix", newPrefix)
}

func (m *MySQL) GetGuildAutoRole(guildID string) (string, error) {
	val, err := m.getGuildSetting(guildID, "autorole")
	return val, err
}

func (m *MySQL) SetGuildAutoRole(guildID, autoRoleID string) error {
	return m.setGuildSetting(guildID, "autorole", autoRoleID)
}

func (m *MySQL) GetGuildModLog(guildID string) (string, error) {
	val, err := m.getGuildSetting(guildID, "modlogchanID")
	return val, err
}

func (m *MySQL) SetGuildModLog(guildID, chanID string) error {
	return m.setGuildSetting(guildID, "modlogchanID", chanID)
}

func (m *MySQL) GetGuildVoiceLog(guildID string) (string, error) {
	val, err := m.getGuildSetting(guildID, "voicelogchanID")
	return val, err
}

func (m *MySQL) SetGuildVoiceLog(guildID, chanID string) error {
	return m.setGuildSetting(guildID, "voicelogchanID", chanID)
}

func (m *MySQL) GetGuildNotifyRole(guildID string) (string, error) {
	val, err := m.getGuildSetting(guildID, "notifyRoleID")
	return val, err
}

func (m *MySQL) SetGuildNotifyRole(guildID, roleID string) error {
	return m.setGuildSetting(guildID, "notifyRoleID", roleID)
}

func (m *MySQL) GetGuildGhostpingMsg(guildID string) (string, error) {
	val, err := m.getGuildSetting(guildID, "ghostPingMsg")
	return val, err
}

func (m *MySQL) SetGuildGhostpingMsg(guildID, msg string) error {
	return m.setGuildSetting(guildID, "ghostPingMsg", msg)
}

func (m *MySQL) GetMemberPermission(s *discordgo.Session, guildID string, memberID string) (PermissionArray, error) {
	guildPerms, err := m.GetGuildPermissions(guildID)
	if err != nil {
		return nil, err
	}
	member, err := s.GuildMember(guildID, memberID)
	if err != nil {
		return nil, err
	}

	var res PermissionArray
	for _, rID := range member.Roles {
		if p, ok := guildPerms[rID]; ok {
			if res == nil {
				res = p
			} else {
				res = res.Merge(p)
			}
		}
	}

	return res, nil
}

func (m *MySQL) GetGuildPermissions(guildID string) (map[string]PermissionArray, error) {
	results := make(map[string]PermissionArray)
	rows, err := m.DB.Query("SELECT roleID, permission FROM permissions WHERE guildID = ?",
		guildID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var roleID string
		var permission string
		err := rows.Scan(&roleID, &permission)
		if err != nil {
			return nil, err
		}
		results[roleID] = strings.Split(permission, ",")
	}
	return results, nil
}

func (m *MySQL) SetGuildRolePermission(guildID, roleID string, p PermissionArray) error {
	if len(p) == 0 {
		_, err := m.DB.Exec("DELETE FROM permissions WHERE roleID = ?", roleID)
		return err
	}

	pStr := strings.Join(p, ",")
	res, err := m.DB.Exec("UPDATE permissions SET permission = ? WHERE roleID = ? AND guildID = ?",
		pStr, roleID, guildID)
	if err != nil {
		return err
	}
	if ar, err := res.RowsAffected(); ar == 0 {
		if err != nil {
			return err
		}
		_, err := m.DB.Exec("INSERT INTO permissions (roleID, guildID, permission) VALUES (?, ?, ?)",
			roleID, guildID, pStr)
		return err
	}
	return nil
}

func (m *MySQL) GetGuildJdoodleKey(guildID string) (string, error) {
	val, err := m.getGuildSetting(guildID, "jdoodleToken")
	return val, err
}

func (m *MySQL) SetGuildJdoodleKey(guildID, key string) error {
	return m.setGuildSetting(guildID, "jdoodleToken", key)
}

func (m *MySQL) GetGuildBackup(guildID string) (bool, error) {
	val, err := m.getGuildSetting(guildID, "backup")
	return val != "", err
}

func (m *MySQL) SetGuildBackup(guildID string, enabled bool) error {
	var val string
	if enabled {
		val = "1"
	}
	return m.setGuildSetting(guildID, "backup", val)
}

func (m *MySQL) GetSetting(setting string) (string, error) {
	var value string
	err := m.DB.QueryRow("SELECT value FROM settings WHERE setting = ?", setting).Scan(&value)
	if err == sql.ErrNoRows {
		err = ErrDatabaseNotFound
	}
	return value, err
}

func (m *MySQL) SetSetting(setting, value string) error {
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

func (m *MySQL) AddReport(rep *util.Report) error {
	_, err := m.DB.Exec("INSERT INTO reports (id, type, guildID, executorID, victimID, msg, attachment) VALUES (?, ?, ?, ?, ?, ?, ?)",
		rep.ID, rep.Type, rep.GuildID, rep.ExecutorID, rep.VictimID, rep.Msg, rep.AttachmehtURL)
	return err
}

func (m *MySQL) DeleteReport(id snowflake.ID) error {
	_, err := m.DB.Exec("DELETE FROM reports WHERE id = ?", id)
	return err
}

func (m *MySQL) GetReport(id snowflake.ID) (*util.Report, error) {
	rep := new(util.Report)

	row := m.DB.QueryRow("SELECT id, type, guildID, executorID, victimID, msg, attachment FROM reports WHERE id = ?", id)
	err := row.Scan(&rep.ID, &rep.Type, &rep.GuildID, &rep.ExecutorID, &rep.VictimID, &rep.Msg, &rep.AttachmehtURL)
	if err == sql.ErrNoRows {
		return nil, ErrDatabaseNotFound
	}

	return rep, err
}

func (m *MySQL) GetReportsGuild(guildID string, offset, limit int) ([]*util.Report, error) {
	if limit == 0 {
		limit = 1000
	}

	rows, err := m.DB.Query(
		"SELECT id, type, guildID, executorID, victimID, msg, attachment "+
			"FROM reports WHERE guildID = ? "+
			"LIMIT ?, ?", guildID, offset, limit)
	var results []*util.Report
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		rep := new(util.Report)
		err := rows.Scan(&rep.ID, &rep.Type, &rep.GuildID, &rep.ExecutorID, &rep.VictimID, &rep.Msg, &rep.AttachmehtURL)
		if err != nil {
			return nil, err
		}
		results = append(results, rep)
	}
	return results, nil
}

func (m *MySQL) GetReportsFiltered(guildID, memberID string, repType int) ([]*util.Report, error) {
	if !util.IsNumber(guildID) || !util.IsNumber(memberID) {
		return nil, fmt.Errorf("invalid argument type")
	}

	query := fmt.Sprintf(`SELECT id, type, guildID, executorID, victimID, msg, attachment FROM reports WHERE guildID = "%s"`, guildID)
	if memberID != "" {
		query += fmt.Sprintf(` AND victimID = "%s"`, memberID)
	}
	if repType != -1 {
		query += fmt.Sprintf(` AND type = %d`, repType)
	}
	rows, err := m.DB.Query(query)
	var results []*util.Report
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		rep := new(util.Report)
		err := rows.Scan(&rep.ID, &rep.Type, &rep.GuildID, &rep.ExecutorID, &rep.VictimID, &rep.Msg, &rep.AttachmehtURL)
		if err != nil {
			return nil, err
		}
		results = append(results, rep)
	}
	return results, nil
}

func (m *MySQL) GetReportsGuildCount(guildID string) (count int, err error) {
	err = m.DB.QueryRow("SELECT COUNT(id) FROM reports WHERE guildID = ?", guildID).Scan(&count)
	return
}

func (m *MySQL) GetReportsFilteredCount(guildID, memberID string, repType int) (count int, err error) {
	if !util.IsNumber(guildID) || !util.IsNumber(memberID) {
		err = fmt.Errorf("invalid argument type")
		return
	}

	query := fmt.Sprintf(`SELECT COUNT(id) FROM reports WHERE guildID = "%s"`, guildID)
	if memberID != "" {
		query += fmt.Sprintf(` AND victimID = "%s"`, memberID)
	}
	if repType != -1 {
		query += fmt.Sprintf(` AND type = %d`, repType)
	}

	err = m.DB.QueryRow(query).Scan(&count)
	return
}

func (m *MySQL) GetVotes() (map[string]*util.Vote, error) {
	rows, err := m.DB.Query("SELECT id, data FROM votes")
	results := make(map[string]*util.Vote)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var voteID, rawData string
		err := rows.Scan(&voteID, &rawData)
		if err != nil {
			util.Log.Error("An error occured reading vote from database: ", err)
			continue
		}
		vote, err := util.VoteUnmarshal(rawData)
		if err != nil {
			m.DeleteVote(rawData)
		} else {
			results[vote.ID] = vote
		}
	}
	return results, err
}

func (m *MySQL) AddUpdateVote(vote *util.Vote) error {
	rawData, err := vote.Marshal()
	if err != nil {
		return err
	}
	res, err := m.DB.Exec("UPDATE votes SET data = ? WHERE id = ?", rawData, vote.ID)
	if ar, err := res.RowsAffected(); ar == 0 {
		if err != nil {
			return err
		}
		_, err := m.DB.Exec("INSERT INTO votes (id, data) VALUES (?, ?)", vote.ID, rawData)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return err
}

func (m *MySQL) DeleteVote(voteID string) error {
	_, err := m.DB.Exec("DELETE FROM votes WHERE id = ?", voteID)
	return err
}

func (m *MySQL) GetMuteRoles() (map[string]string, error) {
	rows, err := m.DB.Query("SELECT guildID, muteRoleID FROM guilds")
	results := make(map[string]string)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var guildID, roleID string
		err = rows.Scan(&guildID, &roleID)
		if err == nil {
			results[guildID] = roleID
		}
	}
	return results, nil
}

func (m *MySQL) GetMuteRoleGuild(guildID string) (string, error) {
	val, err := m.getGuildSetting(guildID, "muteRoleID")
	return val, err
}

func (m *MySQL) SetMuteRole(guildID, roleID string) error {
	return m.setGuildSetting(guildID, "muteRoleID", roleID)
}

func (m *MySQL) GetTwitchNotify(twitchUserID, guildID string) (*TwitchNotifyDBEntry, error) {
	t := &TwitchNotifyDBEntry{
		TwitchUserID: twitchUserID,
		GuildID:      guildID,
	}
	err := m.DB.QueryRow("SELECT channelID FROM twitchnotify WHERE twitchUserID = ? AND guildID = ?",
		twitchUserID, guildID).Scan(&t.ChannelID)
	if err == sql.ErrNoRows {
		err = ErrDatabaseNotFound
	}
	return t, err
}

func (m *MySQL) SetTwitchNotify(twitchNotify *TwitchNotifyDBEntry) error {
	res, err := m.DB.Exec("UPDATE twitchnotify SET channelID = ? WHERE twitchUserID = ? AND guildID = ?",
		twitchNotify.ChannelID, twitchNotify.TwitchUserID, twitchNotify.GuildID)
	if err != nil {
		return err
	}
	if ar, err := res.RowsAffected(); ar == 0 {
		if err != nil {
			return err
		}
		_, err := m.DB.Exec("INSERT INTO twitchnotify (twitchUserID, guildID, channelID) VALUES (?, ?, ?)",
			twitchNotify.TwitchUserID, twitchNotify.GuildID, twitchNotify.ChannelID)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return err
}

func (m *MySQL) DeleteTwitchNotify(twitchUserID, guildID string) error {
	_, err := m.DB.Exec("DELETE FROM twitchnotify WHERE twitchUserID = ? AND guildID = ?", twitchUserID, guildID)
	return err
}

func (m *MySQL) GetAllTwitchNotifies(twitchUserID string) ([]*TwitchNotifyDBEntry, error) {
	query := "SELECT twitchUserID, guildID, channelID FROM twitchnotify"
	if twitchUserID != "" {
		query += " WHERE twitchUserID = " + twitchUserID
	}
	rows, err := m.DB.Query(query)
	results := make([]*TwitchNotifyDBEntry, 0)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		t := new(TwitchNotifyDBEntry)
		err = rows.Scan(&t.TwitchUserID, &t.GuildID, &t.ChannelID)
		if err == nil {
			results = append(results, t)
		}
	}
	return results, nil
}

func (m *MySQL) AddBackup(guildID, fileID string) error {
	timestamp := time.Now().Unix()
	_, err := m.DB.Exec("INSERT INTO backups (guildID, timestamp, fileID) VALUES (?, ?, ?)", guildID, timestamp, fileID)
	return err
}

func (m *MySQL) DeleteBackup(guildID, fileID string) error {
	_, err := m.DB.Exec("DELETE FROM backups WHERE guildID = ? AND fileID = ?", guildID, fileID)
	return err
}

func (m *MySQL) GetGuildInviteBlock(guildID string) (string, error) {
	return m.getGuildSetting(guildID, "inviteBlock")
}

func (m *MySQL) SetGuildInviteBlock(guildID string, data string) error {
	return m.setGuildSetting(guildID, "inviteBlock", data)
}

func (m *MySQL) GetGuildJoinMsg(guildID string) (string, string, error) {
	data, err := m.getGuildSetting(guildID, "joinMsg")
	if err != nil {
		return "", "", err
	}
	if data == "" {
		return "", "", nil
	}

	i := strings.Index(data, "|")
	if i < 0 || len(data) < i+1 {
		return "", "", nil
	}

	return data[:i], data[i+1:], nil
}

func (m *MySQL) SetGuildJoinMsg(guildID string, channelID string, msg string) error {
	return m.setGuildSetting(guildID, "joinMsg", fmt.Sprintf("%s|%s", channelID, msg))
}

func (m *MySQL) GetGuildLeaveMsg(guildID string) (string, string, error) {
	data, err := m.getGuildSetting(guildID, "leaveMsg")
	if err != nil {
		return "", "", err
	}
	if data == "" {
		return "", "", nil
	}

	i := strings.Index(data, "|")
	if i < 0 || len(data) < i+1 {
		return "", "", nil
	}

	return data[:i], data[i+1:], nil
}

func (m *MySQL) SetGuildLeaveMsg(guildID string, channelID string, msg string) error {
	return m.setGuildSetting(guildID, "leaveMsg", fmt.Sprintf("%s|%s", channelID, msg))
}

func (m *MySQL) GetBackups(guildID string) ([]*BackupEntry, error) {
	rows, err := m.DB.Query("SELECT guildID, timestamp, fileID FROM backups WHERE guildID = ?", guildID)
	if err == sql.ErrNoRows {
		return nil, ErrDatabaseNotFound
	}
	if err != nil {
		return nil, err
	}

	backups := make([]*BackupEntry, 0)
	for rows.Next() {
		be := new(BackupEntry)
		var timeStampUnix int64
		err = rows.Scan(&be.GuildID, &timeStampUnix, &be.FileID)
		if err != nil {
			return nil, err
		}
		be.Timestamp = time.Unix(timeStampUnix, 0)
		backups = append(backups, be)
	}

	return backups, nil
}

func (m *MySQL) GetBackupGuilds() ([]string, error) {
	rows, err := m.DB.Query("SELECT guildID FROM guilds WHERE backup = '1'")
	if err == sql.ErrNoRows {
		return nil, ErrDatabaseNotFound
	}
	if err != nil {
		return nil, err
	}

	guilds := make([]string, 0)
	for rows.Next() {
		var s string
		err = rows.Scan(&s)
		if err != nil {
			return nil, err
		}
		guilds = append(guilds, s)
	}

	return guilds, err
}

func (m *MySQL) AddTag(tag *util.Tag) error {
	_, err := m.DB.Exec("INSERT INTO tags (id, ident, creatorID, guildID, content, created, lastEdit) VALUES "+
		"(?, ?, ?, ?, ?, ?, ?)", tag.ID, tag.Ident, tag.CreatorID, tag.GuildID, tag.Content, tag.Created.Unix(), tag.LastEdit.Unix())
	return err
}

func (m *MySQL) EditTag(tag *util.Tag) error {
	_, err := m.DB.Exec("UPDATE tags SET "+
		"ident = ?, creatorID = ?, guildID = ?, content = ?, created = ?, lastEdit = ? "+
		"WHERE id = ?", tag.Ident, tag.CreatorID, tag.GuildID, tag.Content, tag.Created.Unix(), tag.LastEdit.Unix(), tag.ID)
	if err == sql.ErrNoRows {
		return ErrDatabaseNotFound
	}
	return err
}

func (m *MySQL) GetTagByID(id snowflake.ID) (*util.Tag, error) {
	tag := new(util.Tag)
	var timestampCreated int64
	var timestampLastEdit int64

	row := m.DB.QueryRow("SELECT id, ident, creatorID, guildID, content, created, lastEdit FROM tags "+
		"WHERE id = ?", id)

	err := row.Scan(&tag.ID, &tag.Ident, &tag.CreatorID, &tag.GuildID,
		&tag.Content, &timestampCreated, &timestampLastEdit)
	if err == sql.ErrNoRows {
		return nil, ErrDatabaseNotFound
	}
	if err != nil {
		return nil, err
	}

	tag.Created = time.Unix(timestampCreated, 0)
	tag.LastEdit = time.Unix(timestampLastEdit, 0)

	return tag, nil
}

func (m *MySQL) GetTagByIdent(ident string, guildID string) (*util.Tag, error) {
	tag := new(util.Tag)
	var timestampCreated int64
	var timestampLastEdit int64

	row := m.DB.QueryRow("SELECT id, ident, creatorID, guildID, content, created, lastEdit FROM tags "+
		"WHERE ident = ? AND guildID = ?", ident, guildID)

	err := row.Scan(&tag.ID, &tag.Ident, &tag.CreatorID, &tag.GuildID,
		&tag.Content, &timestampCreated, &timestampLastEdit)
	if err == sql.ErrNoRows {
		return nil, ErrDatabaseNotFound
	}
	if err != nil {
		return nil, err
	}

	tag.Created = time.Unix(timestampCreated, 0)
	tag.LastEdit = time.Unix(timestampLastEdit, 0)

	return tag, nil
}

func (m *MySQL) GetGuildTags(guildID string) ([]*util.Tag, error) {
	rows, err := m.DB.Query("SELECT id, ident, creatorID, guildID, content, created, lastEdit FROM tags "+
		"WHERE guildID = ?", guildID)
	if err == sql.ErrNoRows {
		return nil, ErrDatabaseNotFound
	}
	if err != nil {
		return nil, err
	}

	tags := make([]*util.Tag, 0)
	var timestampCreated int64
	var timestampLastEdit int64
	for rows.Next() {
		tag := new(util.Tag)
		err = rows.Scan(&tag.ID, &tag.Ident, &tag.CreatorID, &tag.GuildID,
			&tag.Content, &timestampCreated, &timestampLastEdit)
		if err != nil {
			return nil, err
		}
		tag.Created = time.Unix(timestampCreated, 0)
		tag.LastEdit = time.Unix(timestampLastEdit, 0)
		tags = append(tags, tag)
	}

	return tags, nil
}

func (m *MySQL) DeleteTag(id snowflake.ID) error {
	_, err := m.DB.Exec("DELETE FROM tags WHERE id = ?", id)
	if err == sql.ErrNoRows {
		return ErrDatabaseNotFound
	}
	return err
}

func (m *MySQL) SetSession(key, userID string, expires time.Time) error {
	res, err := m.DB.Exec("UPDATE sessions SET sessionkey = ?, expires = ? WHERE userID = ?", key, expires, userID)
	if err != sql.ErrNoRows && err != nil {
		return err
	}

	if ar, err := res.RowsAffected(); ar == 0 {
		if err != nil {
			return err
		}
		_, err := m.DB.Exec("INSERT INTO sessions (sessionkey, userID, expires) VALUES (?, ?, ?)", key, userID, expires)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return err
}

func (m *MySQL) GetSession(key string) (string, error) {
	var userID string
	var expires time.Time
	err := m.DB.QueryRow("SELECT userID, expires FROM sessions WHERE sessionkey = ?", key).
		Scan(&userID, &expires)

	if err == sql.ErrNoRows {
		return "", ErrDatabaseNotFound
	}
	if err != nil {
		return "", err
	}

	if expires.Before(time.Now()) {
		return "", ErrDatabaseNotFound
	}

	return userID, nil
}

func (m *MySQL) DeleteSession(userID string) error {
	_, err := m.DB.Exec("DELETE FROM sessions WHERE userID = ?", userID)
	if err == sql.ErrNoRows {
		return ErrDatabaseNotFound
	}
	return err
}
