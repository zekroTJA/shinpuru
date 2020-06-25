package main

import (
	"bytes"
	"flag"
	"os"
	"path"

	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/snowflake"
	"github.com/zekroTJA/shinpuru/internal/core/backup/backupmodels"
	"github.com/zekroTJA/shinpuru/internal/core/config"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/core/storage"
	"github.com/zekroTJA/shinpuru/internal/inits"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/report"
	"github.com/zekroTJA/shinpuru/internal/util/static"
)

var (
	flagConfigLocation = flag.String("c", "config.yml", "The location of the main config file")
)

func main() {
	flag.Parse()

	cfg := inits.InitConfig(*flagConfigLocation, new(config.YAMLConfigParser))

	db := inits.InitDatabase(cfg.Database)
	defer db.Close()

	st := inits.InitStorage(cfg)

	session, err := discordgo.New(cfg.Discord.Token)
	if err != nil {
		util.Log.Fatal(err)
	}

	if err = session.Open(); err != nil {
		util.Log.Fatal(err)
	}
	defer session.Close()

	guilds := session.State.Guilds

	util.Log.Infof("Migrating report attachments for %d guilds...", len(guilds))
	for _, guild := range guilds {
		reports, err := db.GetReportsGuild(guild.ID, 0, 10000)
		if err != nil {
			util.Log.Fatal(err)
		}

		for _, rep := range reports {
			migrateReport(rep, db, st)
		}
	}

	util.Log.Infof("Migrating backup files for %d guilds...", len(guilds))
	for _, guild := range guilds {
		backups, err := db.GetBackups(guild.ID)
		if err != nil {
			util.Log.Fatal(err)
		}

		for _, backup := range backups {
			migrateBackup(backup, cfg.Discord.GuildBackupLoc, db, st)
		}
	}
}

func migrateReport(rep *report.Report, db database.Database, st storage.Storage) {
	if rep.AttachmehtURL == "" {
		return
	}

	id, err := snowflake.ParseString(rep.AttachmehtURL)
	if err != nil {
		util.Log.Error(err)
		return
	}

	img, err := db.GetImageData(id)
	if err != nil {
		if !database.IsErrDatabaseNotFound(err) {
			util.Log.Error(err)
		}
		return
	}

	util.Log.Infof("Migrating report attachment %s", img.Data)
	r := bytes.NewReader(img.Data)
	err = st.PutObject(static.StorageBucketImages, img.ID.String(), r, int64(img.Size), img.MimeType)
	if err != nil {
		util.Log.Error(err)
	}
}

func migrateBackup(backup *backupmodels.Entry, loc string, db database.Database, st storage.Storage) {
	fd := path.Join(loc, backup.FileID+".json")
	fh, err := os.Open(fd)
	if err != nil {
		util.Log.Error(err)
		return
	}
	defer fh.Close()

	stat, err := fh.Stat()
	if err != nil {
		util.Log.Error(err)
		return
	}

	util.Log.Infof("Migrating backup file %s", backup.FileID)
	err = st.PutObject(static.StorageBucketBackups, backup.FileID, fh, stat.Size(), "application/json")
	if err != nil {
		util.Log.Error(err)
		return
	}

	os.Remove(fd)
}
