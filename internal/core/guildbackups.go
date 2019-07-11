package core

import (
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/zekroTJA/shinpuru/internal/util"

	"github.com/bwmarrin/discordgo"
)

const (
	tickRate       = 12 * time.Hour
	backupLocation = "./guildBackups"
)

type GuildBackups struct {
	ticker  *time.Ticker
	session *discordgo.Session
	db      Database
}

type BackupEntry struct {
	GuildID   string
	Timestamp time.Time
	FileID    string
}

type BackupObject struct {
	ID       string           `json:"id"`
	Guild    *BackupGuild     `json:"guild"`
	Channels []*BackupChannel `json:"channels"`
	Roles    []*BackupRole    `json:"roles"`
	Members  []*BackupMember  `json:"members"`
}

type BackupGuild struct {
	Name                        string `json:"name"`
	AfkChannelID                string `json:"afk_channel_id"`
	AfkTimeout                  int    `json:"afk_timeout"`
	VerificationLevel           int    `json:"verification_level"`
	DefaultMessageNotifications int    `json:"default_message_notifications"`
}

type BackupChannel struct {
	ID                   string                           `json:"id"`
	Name                 string                           `json:"name"`
	Topic                string                           `json:"topic"`
	Type                 int                              `json:"type"`
	NSFW                 bool                             `json:"nsfw"`
	Position             int                              `json:"position"`
	Bitrate              int                              `json:"bitrate"`
	UserLimit            int                              `json:"user_limit"`
	ParentID             string                           `json:"parent_id"`
	PermissionOverwrites []*discordgo.PermissionOverwrite `json:"permission_overwrites"`
}

type BackupRole struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Mentionable bool   `json:"mentionable"`
	Hoist       bool   `json:"hoist"`
	Color       int    `json:"color"`
	Position    int    `json:"position"`
	Permissions int    `json:"permissions"`
}

type BackupMember struct {
	ID    string   `json:"id"`
	Nick  string   `json:"nick"`
	Deaf  bool     `json:"deaf"`
	Mute  bool     `json:"mute"`
	Roles []string `json:"roles"`
}

func asyncWriteStatus(c chan string, status string) {
	go func() {
		c <- status
	}()
}

func asyncWriteError(c chan error, err error) {
	go func() {
		c <- err
	}()
}

func NewGuildBackups(s *discordgo.Session, db Database) *GuildBackups {
	bck := new(GuildBackups)
	bck.db = db
	bck.session = s
	bck.ticker = time.NewTicker(tickRate)
	go bck.initTickerLoop()
	return bck
}

func (bck *GuildBackups) initTickerLoop() {
	for {
		<-bck.ticker.C
		go bck.backupAllGuilds()
	}
}

func (bck *GuildBackups) backupAllGuilds() {
	guilds, err := bck.db.GetBackupGuilds()
	if err != nil {
		util.Log.Error("failed getting backup guilds: ", err)
		return
	}

	for _, g := range guilds {
		err = bck.BackupGuild(g)
		if err != nil {
			util.Log.Error("failed creating backup for guild '%s': %s", g, err.Error())
		}
		time.Sleep(1 * time.Second)
	}
}

func (bck *GuildBackups) BackupGuild(guildID string) error {
	if bck.session == nil {
		return errors.New("session is nil")
	}

	g, err := bck.session.Guild(guildID)
	if err != nil {
		return err
	}

	backup := new(BackupObject)
	backup.Guild = &BackupGuild{
		AfkChannelID:                g.AfkChannelID,
		AfkTimeout:                  g.AfkTimeout,
		DefaultMessageNotifications: g.DefaultMessageNotifications,
		Name:                        g.Name,
		VerificationLevel:           int(g.VerificationLevel),
	}

	for _, c := range g.Channels {
		backup.Channels = append(backup.Channels, &BackupChannel{
			Bitrate:   c.Bitrate,
			ID:        c.ID,
			NSFW:      c.NSFW,
			Name:      c.Name,
			ParentID:  c.ParentID,
			Position:  c.Position,
			Topic:     c.Topic,
			Type:      int(c.Type),
			UserLimit: c.UserLimit,
		})
	}

	for _, r := range g.Roles {
		if r.ID == guildID {
			continue
		}
		backup.Roles = append(backup.Roles, &BackupRole{
			Color:       r.Color,
			Hoist:       r.Hoist,
			ID:          r.ID,
			Mentionable: r.Mentionable,
			Name:        r.Name,
			Permissions: r.Permissions,
			Position:    r.Position,
		})
	}

	for _, m := range g.Members {
		backup.Members = append(backup.Members, &BackupMember{
			Deaf:  m.Deaf,
			ID:    m.User.ID,
			Mute:  m.Mute,
			Nick:  m.Nick,
			Roles: m.Roles,
		})
	}

	if _, err := os.Stat(backupLocation); os.IsNotExist(err) {
		err = os.MkdirAll(backupLocation, os.ModePerm)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	backupID := util.NodeBackup.Generate()
	backupFileName := backupLocation + "/" + backupID.String() + ".json"

	f, err := os.Create(backupFileName)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	err = enc.Encode(backup)
	if err != nil {
		f.Close()
		os.Remove(backupFileName)
		return err
	}

	err = bck.db.AddBackup(g.ID, backupID.String())
	if err != nil {
		return err
	}

	cBackups, err := bck.db.GetBackups(g.ID)
	if err != nil {
		return err
	}

	if len(cBackups) > 10 {
		var lastEntry *BackupEntry
		for _, b := range cBackups {
			if lastEntry == nil || b.Timestamp.Before(lastEntry.Timestamp) {
				lastEntry = b
			}
		}

		err = os.Remove(backupLocation + "/" + lastEntry.FileID + ".json")
		if err != nil {
			return err
		}

		err = bck.db.DeleteBackup(g.ID, lastEntry.FileID)
		if err != nil {
			return err
		}
	}

	return err
}

func (bck *GuildBackups) RestoreBackup(guildID, fileID string, statusC chan string, errorsC chan error) error {
	defer func() {
		close(statusC)
		close(errorsC)
	}()

	if bck.session == nil {
		return errors.New("session is nil")
	}

	asyncWriteStatus(statusC, "reading backup file")
	f, err := os.Open(backupLocation + "/" + fileID + ".json")
	if err != nil {
		return err
	}
	defer f.Close()

	var backup BackupObject

	dec := json.NewDecoder(f)
	err = dec.Decode(&backup)
	if err != nil {
		return err
	}

	if bck.session == nil {
		return errors.New("session is nil")
	}

	// EDIT GUILD
	asyncWriteStatus(statusC, "editing guild")
	_verificationLevel := discordgo.VerificationLevel(backup.Guild.VerificationLevel)
	_, err = bck.session.GuildEdit(guildID, discordgo.GuildParams{
		Name:                        backup.Guild.Name,
		AfkChannelID:                backup.Guild.AfkChannelID,
		AfkTimeout:                  backup.Guild.AfkTimeout,
		VerificationLevel:           &_verificationLevel,
		DefaultMessageNotifications: backup.Guild.DefaultMessageNotifications,
	})
	if err != nil {
		return err
	}

	// backup ID - created ID
	ids := make(map[string]string)
	channelsPos := make(map[string]int)
	orderedRoles := make([]*discordgo.Role, len(backup.Roles))

	// CREATE ROLES
	asyncWriteStatus(statusC, "updating and creating roles")
	for _, r := range backup.Roles {
		roles, err := bck.session.GuildRoles(guildID)
		if err != nil {
			return err
		}

		var rObj *discordgo.Role
		for _, crObj := range roles {
			if crObj.ID == r.ID {
				rObj = crObj
			}
		}

		if rObj == nil {
			rObj, err = bck.session.GuildRoleCreate(guildID)
		}
		_, err = bck.session.GuildRoleEdit(guildID, rObj.ID, r.Name, r.Color,
			r.Hoist, r.Permissions, r.Mentionable)
		if err != nil {
			asyncWriteError(errorsC, err)
			continue
		}

		ids[r.ID] = rObj.ID
		orderedRoles[r.Position-1] = rObj
	}

	// RE-POSITION ROLES
	asyncWriteStatus(statusC, "re-position roles")
	_, err = bck.session.GuildRoleReorder(guildID, orderedRoles)
	if err != nil {
		asyncWriteError(errorsC, err)
	}

	// CREATE CATEGORIES
	asyncWriteStatus(statusC, "updating and creating categories")
	for _, c := range backup.Channels {
		if c.Type != int(discordgo.ChannelTypeGuildCategory) {
			continue
		}
		cObj, _ := bck.session.Channel(c.ID)
		if cObj != nil && cObj.GuildID != guildID {
			cObj = nil
		}

		for _, po := range c.PermissionOverwrites {
			po.ID = ids[po.ID]
		}

		if cObj == nil {
			cObj, err = bck.session.GuildChannelCreateComplex(guildID,
				discordgo.GuildChannelCreateData{
					Name:                 c.Name,
					PermissionOverwrites: c.PermissionOverwrites,
					Type:                 discordgo.ChannelTypeGuildCategory,
				})
		} else {
			_, err = bck.session.ChannelEditComplex(cObj.ID,
				&discordgo.ChannelEdit{
					Name:                 c.Name,
					PermissionOverwrites: c.PermissionOverwrites,
				})
		}
		if err != nil {
			asyncWriteError(errorsC, err)
			continue
		}
		ids[c.ID] = cObj.ID
		channelsPos[cObj.ID] = c.Position
	}

	// CREATE CHANNELS AND ADD TO CATEGORIES
	asyncWriteStatus(statusC, "updating and creating channels")
	for _, c := range backup.Channels {
		if c.Type == int(discordgo.ChannelTypeGuildCategory) {
			continue
		}

		cObj, _ := bck.session.Channel(c.ID)
		if cObj != nil && cObj.GuildID != guildID {
			cObj = nil
		}

		for _, po := range c.PermissionOverwrites {
			po.ID = ids[po.ID]
		}

		if cObj == nil {
			cObj, err = bck.session.GuildChannelCreateComplex(guildID,
				discordgo.GuildChannelCreateData{
					Bitrate:              c.Bitrate,
					NSFW:                 c.NSFW,
					Name:                 c.Name,
					ParentID:             ids[c.ParentID],
					PermissionOverwrites: c.PermissionOverwrites,
					Topic:                c.Topic,
					Type:                 discordgo.ChannelType(c.Type),
					UserLimit:            c.UserLimit,
				})
		} else {
			_, err = bck.session.ChannelEditComplex(cObj.ID,
				&discordgo.ChannelEdit{
					Bitrate:              c.Bitrate,
					NSFW:                 c.NSFW,
					Name:                 c.Name,
					ParentID:             ids[c.ParentID],
					PermissionOverwrites: c.PermissionOverwrites,
					Topic:                c.Topic,
					UserLimit:            c.UserLimit,
				})
		}
		if err != nil {
			asyncWriteError(errorsC, err)
			continue
		}

		channelsPos[cObj.ID] = c.Position
	}

	// RE-POSITION CHANNELS
	asyncWriteStatus(statusC, "re-positioning channels")
	for cID, pos := range channelsPos {
		_, err = bck.session.ChannelEditComplex(cID, &discordgo.ChannelEdit{
			Position: pos,
		})
		if err != nil {
			return err
		}
	}

	// UPDATE MEMBERS
	asyncWriteStatus(statusC, "updating members")
	for _, m := range backup.Members {
		mObj, _ := bck.session.GuildMember(guildID, m.ID)

		newRoles := make([]string, len(m.Roles))
		for i, r := range m.Roles {
			newRoles[i] = ids[r]
		}

		if mObj != nil {
			err = bck.session.GuildMemberEdit(guildID, m.ID, newRoles)
			if err != nil {
				asyncWriteError(errorsC, err)
				continue
			}

			err = bck.session.GuildMemberNickname(guildID, m.ID, m.Nick)
			if err != nil {
				asyncWriteError(errorsC, err)
				continue
			}
		}
	}

	return nil
}

func (bck *GuildBackups) HardFlush(guildID string) error {
	if bck.session == nil {
		return errors.New("session is nil")
	}

	g, err := bck.session.Guild(guildID)
	if err != nil {
		return err
	}

	for _, r := range g.Roles {
		bck.session.GuildRoleDelete(guildID, r.ID)
	}

	for _, c := range g.Channels {
		bck.session.ChannelDelete(c.ID)
	}

	return nil
}
