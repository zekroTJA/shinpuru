package backup

import (
	"bytes"
	"encoding/json"
	"errors"
	"time"

	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/services/backup/backupmodels"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/guildlog"
	"github.com/zekroTJA/shinpuru/internal/services/storage"
	"github.com/zekroTJA/shinpuru/internal/services/timeprovider"
	"github.com/zekroTJA/shinpuru/internal/util/snowflakenodes"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
	"github.com/zekroTJA/shinpuru/pkg/inline"
	"github.com/zekrotja/dgrs"
	"github.com/zekrotja/rogu"
	"github.com/zekrotja/rogu/log"
	"github.com/zekrotja/sop"

	"github.com/bwmarrin/discordgo"
)

const (
	// tickRate = 30 * time.Second
	tickRate = 12 * time.Hour
)

// GuildBackups provides functionalities to backup
// and restore a guild to and from a JSON file.
type GuildBackups struct {
	session *discordgo.Session
	db      database.Database
	gl      guildlog.Logger
	st      storage.Storage
	state   *dgrs.State
	tp      timeprovider.Provider
	log     rogu.Logger
}

// asyncWriteStatus writes the passed status to the
// passed channel in a new goroutine.
func asyncWriteStatus(c chan string, status string) {
	go func() {
		c <- status
	}()
}

// asyncWriteError writes the passed error to the
// passed channel in a new goroutine.
func asyncWriteError(c chan error, err error) {
	go func() {
		c <- err
	}()
}

func (bck *GuildBackups) guilds() (guilds []string, err error) {
	guilds, err = bck.db.GetGuilds()
	shardID, shardTotal := discordutil.GetShardOfSession(bck.session)
	if shardTotal > 1 {
		guilds = sop.Slice(guilds).Filter(func(v string, i int) bool {
			id, err := discordutil.GetShardOfGuild(v, shardTotal)
			return err == nil && id == shardID
		}).Unwrap()
	}
	return
}

// New initializes a new GuildBackups instance using
// the passed discordgo Session, database provider,
// and storage provider. Also, the ticker loop is
// initialized.
func New(container di.Container) *GuildBackups {
	bck := new(GuildBackups)
	bck.db = container.Get(static.DiDatabase).(database.Database)
	bck.gl = container.Get(static.DiGuildLog).(guildlog.Logger).Section("backup")
	bck.st = container.Get(static.DiObjectStorage).(storage.Storage)
	bck.session = container.Get(static.DiDiscordSession).(*discordgo.Session)
	bck.state = container.Get(static.DiState).(*dgrs.State)
	bck.tp = container.Get(static.DiTimeProvider).(timeprovider.Provider)
	bck.log = log.Tagged("GuildBackup")
	return bck
}

// BackupAllGuilds iterates through all guilds
// which have guild backups enabled and initiates
// the backup routines one after one.
// Guild backups are not created in new goroutines
// because of potential rate limit exceedance.
func (bck *GuildBackups) BackupAllGuilds() {
	guilds, err := bck.guilds()

	bck.log.Info().Fields("nGuilds", len(guilds)).Msg("Backing up guilds ...")

	if err != nil {
		bck.log.Error().Err(err).Msg("Failed getting guilds to back up")
		return
	}

	for _, g := range guilds {
		err = bck.BackupGuild(g)
		if err != nil {
			bck.log.Error().Err(err).Field("gid", g).Msg("Failed creating backup for guild")
			bck.gl.Errorf(g, "Failed creating guild backup: %s", err.Error())
		}
		time.Sleep(1 * time.Second)
	}
}

// BackupGuild creates a backup of a single guild
// and writes the resulting JSON file to the specified
// storage. If the backup creation fails, the error is
// returned.
func (bck *GuildBackups) BackupGuild(guildID string) error {
	if bck.session == nil {
		return errors.New("session is nil")
	}

	g, err := bck.state.Guild(guildID, true)
	if err != nil {
		return err
	}

	backup := new(backupmodels.Object)
	backup.Timestamp = bck.tp.Now()

	backup.Guild = &backupmodels.Guild{
		AfkChannelID:                g.AfkChannelID,
		AfkTimeout:                  g.AfkTimeout,
		DefaultMessageNotifications: int(g.DefaultMessageNotifications),
		Name:                        g.Name,
		VerificationLevel:           int(g.VerificationLevel),
	}

	chans, err := bck.state.Channels(g.ID, true)
	if err != nil {
		return err
	}

	for _, c := range chans {
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

	members, err := bck.state.Members(g.ID, true)
	if err != nil {
		return err
	}
	for _, m := range members {
		backup.Members = append(backup.Members, &backupmodels.Member{
			Deaf:  m.Deaf,
			ID:    m.User.ID,
			Mute:  m.Mute,
			Nick:  m.Nick,
			Roles: m.Roles,
		})
	}

	backup.ID = snowflakenodes.NodeBackup.Generate().String()

	buff := bytes.NewBuffer([]byte{})

	enc := json.NewEncoder(buff)
	enc.SetIndent("", "  ")
	err = enc.Encode(backup)
	if err != nil {
		return err
	}

	err = bck.st.PutObject(static.StorageBucketBackups, backup.ID, buff, int64(buff.Len()), "application/json")
	if err != nil {
		return err
	}

	err = bck.db.AddBackup(g.ID, backup.ID)
	if err != nil {
		bck.st.DeleteObject(static.StorageBucketBackups, backup.ID)
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
				lastEntry = &b
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

// RestoreBackup tries to restore a guild structure by
// backup file specified via fileID. The current status
// is sent into the statusC channel and occured errors
// are pushed into the errorsC channel.
// Only if the initialization of this method has failed,
// an error is returned. The returned error does not
// represent the success result of the backup restore
// process.
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
	_, err = bck.session.GuildEdit(guildID, &discordgo.GuildParams{
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
	orderedRoles := make(map[int]*discordgo.Role)

	// CREATE ROLES
	asyncWriteStatus(statusC, "updating and creating roles")
	roles, err := bck.session.GuildRoles(guildID)
	if err != nil {
		return err
	}
	for _, r := range backup.Roles {
		var rObj *discordgo.Role
		for _, crObj := range roles {
			if crObj.ID == r.ID {
				rObj = crObj
			}
		}
		if rObj == nil {
			rObj, err = bck.session.GuildRoleCreate(guildID, &discordgo.RoleParams{
				Name:        r.Name,
				Color:       &r.Color,
				Hoist:       &r.Hoist,
				Permissions: &r.Permissions,
				Mentionable: &r.Mentionable,
			})
			if err != nil {
				asyncWriteError(errorsC, err)
				continue
			}
		}

		ids[r.ID] = rObj.ID
		orderedRoles[r.Position-1] = rObj
	}

	// RE-POSITION ROLES
	asyncWriteStatus(statusC, "re-position roles")
	or := sop.Map(sop.MapFlat(orderedRoles).Sort(func(p, q sop.Tuple[int, *discordgo.Role], i int) bool {
		return p.V1 < q.V1
	}), func(v sop.Tuple[int, *discordgo.Role], i int) *discordgo.Role {
		return v.V2
	}).Unwrap()
	_, err = bck.session.GuildRoleReorder(guildID, or)
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
					NSFW:                 &c.NSFW,
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
		ids[c.ID] = cObj.ID
	}
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
					NSFW:                 &c.NSFW,
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
			newRoles[i] = inline.NC(ids[r], r)
		}

		if mObj != nil {
			_, err = bck.session.GuildMemberEdit(guildID, m.ID, &discordgo.GuildMemberParams{
				Roles: &newRoles,
			})
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

// HardFlush removes all roles and channels
// of a guild.
func (bck *GuildBackups) HardFlush(guildID string) error {
	if bck.session == nil {
		return errors.New("session is nil")
	}

	g, err := bck.state.Guild(guildID, true)
	if err != nil {
		return err
	}

	for _, r := range g.Roles {
		bck.session.GuildRoleDelete(guildID, r.ID)
	}

	chans, err := bck.state.Channels(g.ID, true)
	if err != nil {
		return err
	}

	for _, c := range chans {
		bck.session.ChannelDelete(c.ID)
	}

	return nil
}
