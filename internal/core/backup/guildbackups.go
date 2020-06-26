package backup

import (
	"bytes"
	"encoding/json"
	"errors"
	"time"

	"github.com/zekroTJA/shinpuru/internal/core/backup/backupmodels"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/core/storage"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/snowflakenodes"
	"github.com/zekroTJA/shinpuru/internal/util/static"

	"github.com/bwmarrin/discordgo"
)

const (
	// tickRate = 5 * time.Minute
	tickRate = 12 * time.Hour
)

var (
	backupLocation = "./guildBackups"
)

type GuildBackups struct {
	ticker  *time.Ticker
	session *discordgo.Session
	db      database.Database
	st      storage.Storage
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

func New(s *discordgo.Session, db database.Database, st storage.Storage, loc string) *GuildBackups {
	if loc != "" {
		backupLocation = loc
	}

	bck := new(GuildBackups)
	bck.db = db
	bck.st = st
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
	guilds, err := bck.db.GetGuilds()

	util.Log.Infof("Initializing backups for %d guilds %+v\n", len(guilds), guilds)

	if err != nil {
		util.Log.Error("failed getting backup guilds: ", err)
		return
	}

	for _, g := range guilds {
		err = bck.Guild(g)
		if err != nil {
			util.Log.Error("failed creating backup for guild '%s': %s", g, err.Error())
		}
		time.Sleep(1 * time.Second)
	}
}

func (bck *GuildBackups) Guild(guildID string) error {
	if bck.session == nil {
		return errors.New("session is nil")
	}

	g, err := bck.session.Guild(guildID)
	if err != nil {
		return err
	}

	backup := new(backupmodels.Object)
	backup.Guild = &backupmodels.Guild{
		AfkChannelID:                g.AfkChannelID,
		AfkTimeout:                  g.AfkTimeout,
		DefaultMessageNotifications: g.DefaultMessageNotifications,
		Name:                        g.Name,
		VerificationLevel:           int(g.VerificationLevel),
	}

	for _, c := range g.Channels {
		backup.Channels = append(backup.Channels, &backupmodels.Channel{
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
		backup.Roles = append(backup.Roles, &backupmodels.Role{
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
		backup.Members = append(backup.Members, &backupmodels.Member{
			Deaf:  m.Deaf,
			ID:    m.User.ID,
			Mute:  m.Mute,
			Nick:  m.Nick,
			Roles: m.Roles,
		})
	}

	backupID := snowflakenodes.NodeBackup.Generate()

	buff := bytes.NewBuffer([]byte{})

	enc := json.NewEncoder(buff)
	enc.SetIndent("", "  ")
	err = enc.Encode(backup)
	if err != nil {
		return err
	}

	err = bck.st.PutObject(static.StorageBucketBackups, backupID.String(), buff, int64(buff.Len()), "application/json")
	if err != nil {
		return err
	}

	err = bck.db.AddBackup(g.ID, backupID.String())
	if err != nil {
		bck.st.DeleteObject(static.StorageBucketBackups, backupID.String())
		return err
	}

	cBackups, err := bck.db.GetBackups(g.ID)
	if err != nil {
		return err
	}

	if len(cBackups) > 10 {
		var lastEntry *backupmodels.Entry
		for _, b := range cBackups {
			if lastEntry == nil || b.Timestamp.Before(lastEntry.Timestamp) {
				lastEntry = b
			}
		}

		err = bck.st.DeleteObject(static.StorageBucketBackups, lastEntry.FileID)
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
	reader, _, err := bck.st.GetObject(static.StorageBucketBackups, fileID)
	if err != nil {
		return err
	}
	defer reader.Close()

	var backup backupmodels.Object

	dec := json.NewDecoder(reader)
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
