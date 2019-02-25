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
	tickRate       = 1 * time.Minute // 12 * time.Hour
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
	ID        string `json:"id"`
	Name      string `json:"name"`
	Topic     string `json:"topic"`
	Type      int    `json:"type"`
	NSFW      bool   `json:"nsfw"`
	Position  int    `json:"position"`
	Bitrate   int    `json:"bitrate"`
	UserLimit int    `json:"user_limit"`
	ParentID  string `json:"parent_id"`
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

	backupID := util.BackupNode.Generate()
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
