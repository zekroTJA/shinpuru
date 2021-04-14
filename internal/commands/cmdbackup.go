package commands

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/core/backup"
	"github.com/zekroTJA/shinpuru/internal/core/backup/backupmodels"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/core/storage"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/acceptmsg"
	"github.com/zekroTJA/shireikan"
)

const (
	timeFormat = time.RFC1123
)

type CmdBackup struct {
}

func (c *CmdBackup) GetInvokes() []string {
	return []string{"backup", "backups", "bckp", "guildbackup"}
}

func (c *CmdBackup) GetDescription() string {
	return "Enable, disable and manage guild backups."
}

func (c *CmdBackup) GetHelp() string {
	return "`backup <enable|disable>` - enable or disable backups for your guild\n" +
		"`backup (list)` - list all saved backups\n" +
		"`backup restore <id>` - restore a backup\n" +
		"`backup purge` - delete all backups of the guild"
}

func (c *CmdBackup) GetGroup() string {
	return shireikan.GroupGuildAdmin
}

func (c *CmdBackup) GetDomainName() string {
	return "sp.guild.admin.backup"
}

func (c *CmdBackup) GetSubPermissionRules() []shireikan.SubPermission {
	return nil
}

func (c *CmdBackup) IsExecutableInDMChannels() bool {
	return false
}

func (c *CmdBackup) Exec(ctx shireikan.Context) error {
	if len(ctx.GetArgs()) > 0 {
		switch strings.ToLower(ctx.GetArgs().Get(0).AsString()) {
		case "e", "enable":
			return c.switchStatus(ctx, true)
		case "d", "disable":
			return c.switchStatus(ctx, false)
		case "r", "restore":
			return c.restore(ctx)
		case "purge", "clear":
			return c.purgeBackupsAccept(ctx)
		default:
			return c.list(ctx)
		}
	}
	return c.list(ctx)
}

func (c *CmdBackup) switchStatus(ctx shireikan.Context, enable bool) error {
	db, _ := ctx.GetObject(static.DiDatabase).(database.Database)

	err := db.SetGuildBackup(ctx.GetGuild().ID, enable)
	if err != nil {
		return err
	}

	if enable {
		return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID, "Enabled backup for this guild.\nA full guild backup *(incl. Members, Roles, Channels and Guild Settings)* "+
			"will be created every 12 hours. Only 10 backups per guild will be saved, so you will habe the backup files of the last 5 days.", "", static.ColorEmbedGreen).
			DeleteAfter(15 * time.Second).Error()
	}

	return util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID, "Backup creation disabled.\n"+
		"You will be still have access to created backups and be able to restore them.", "", static.ColorEmbedOrange).
		DeleteAfter(15 * time.Second).Error()
}

func (c *CmdBackup) getBackupsList(ctx shireikan.Context) ([]*backupmodels.Entry, string, error) {
	db, _ := ctx.GetObject(static.DiDatabase).(database.Database)

	backups, err := db.GetBackups(ctx.GetGuild().ID)
	if err != nil && database.IsErrDatabaseNotFound(err) {
		return nil, "", err
	}

	strBackupAll := "*no backups saved*"

	if len(backups) > 0 {
		sort.Slice(backups, func(i, j int) bool {
			return backups[i].Timestamp.Before(backups[j].Timestamp)
		})

		if len(backups) > 10 {
			backups = backups[0:10]
		}

		strBackups := make([]string, len(backups))

		for i, b := range backups {
			strBackups[i] = fmt.Sprintf("`%d` - %s - (ID: `%s`)", i, b.Timestamp.Format(timeFormat), b.FileID)
		}

		strBackupAll = strings.Join(strBackups, "\n")
	}

	return backups, strBackupAll, nil
}

func (c *CmdBackup) list(ctx shireikan.Context) error {
	db, _ := ctx.GetObject(static.DiDatabase).(database.Database)

	status, err := db.GetGuildBackup(ctx.GetGuild().ID)
	if err != nil && database.IsErrDatabaseNotFound(err) {
		return err
	}

	strStatus := ":x:  Backups **disabled**"
	if status {
		strStatus = ":white_check_mark:  Backups **enabled**"
	}

	_, strBackupAll, err := c.getBackupsList(ctx)
	if err != nil {
		return err
	}

	emb := &discordgo.MessageEmbed{
		Color:       static.ColorEmbedDefault,
		Title:       "Backups",
		Description: strStatus,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Saved Backups",
				Value: strBackupAll,
			},
		},
	}

	_, err = ctx.GetSession().ChannelMessageSendEmbed(ctx.GetChannel().ID, emb)
	return err
}

func (c *CmdBackup) restore(ctx shireikan.Context) error {
	if len(ctx.GetArgs()) < 2 {
		return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID, "Please specify the index or the ID of the backup, you want to restore.").
			DeleteAfter(8 * time.Second).Error()
	}

	backups, _, err := c.getBackupsList(ctx)
	if err != nil {
		return err
	}

	i, err := ctx.GetArgs().Get(1).AsInt()
	if err != nil {
		return err
	}

	if i < 0 {
		return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID, "Argument must be an index between 0 and 9 or a snowflake ID.").
			DeleteAfter(8 * time.Second).Error()
	}

	var backup *backupmodels.Entry

	if i < 10 {
		if len(backups)-1 < i {
			return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
				fmt.Sprintf("There are only %d (index 0 to %d) backups you can chose from.", len(backups), len(backups)-1)).
				DeleteAfter(8 * time.Second).Error()
		}
		backup = backups[i]
	} else {
		for _, b := range backups {
			if b.FileID == ctx.GetArgs().Get(1).AsString() {
				backup = b
			}
		}
	}

	if backup == nil {
		return util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
			fmt.Sprintf("Could not find any backup by this specifier: ```\n%s\n```", ctx.GetArgs().Get(1).AsString())).
			DeleteAfter(8 * time.Second).Error()
	}

	accMsg := &acceptmsg.AcceptMessage{
		Session:        ctx.GetSession(),
		DeleteMsgAfter: true,
		UserID:         ctx.GetUser().ID,
		Embed: &discordgo.MessageEmbed{
			Color: static.ColorEmbedOrange,
			Description: fmt.Sprintf(":warning:  **WARNING**  :warning:\n\n"+
				"By pressing :white_check_mark:, the structure of this guild will be **reset** to the selected backup:\n\n"+
				"%s - (ID: `%s`)", backup.Timestamp.Format(timeFormat), backup.FileID),
		},
		DeclineFunc: func(m *discordgo.Message) {
			util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID, "Canceled.").
				DeleteAfter(6 * time.Second).Error()
			return
		},
		AcceptFunc: func(m *discordgo.Message) {
			c.proceedRestore(ctx, backup.FileID)
		},
	}

	_, err = accMsg.Send(ctx.GetChannel().ID)
	return err
}

func (c *CmdBackup) proceedRestore(ctx shireikan.Context, fileID string) {
	statusChan := make(chan string)
	errorsChan := make(chan error)

	statusMsg, _ := ctx.GetSession().ChannelMessageSendEmbed(ctx.GetChannel().ID,
		&discordgo.MessageEmbed{
			Color:       static.ColorEmbedGray,
			Description: "initializing backup restoring...",
		})

	if statusMsg != nil {
		go func() {
			for {
				select {
				case status, ok := <-statusChan:
					if !ok {
						continue
					}
					ctx.GetSession().ChannelMessageEditEmbed(statusMsg.ChannelID, statusMsg.ID, &discordgo.MessageEmbed{
						Color:       static.ColorEmbedGray,
						Description: status + "...",
					})
				case err, ok := <-errorsChan:
					if !ok || err == nil {
						continue
					}
					util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
						"An unexpected error occured while restoring backup (process will not be aborted): ```\n"+err.Error()+"\n```")
				}
			}
		}()
	}

	bck, _ := ctx.GetObject(static.DiBackupHandler).(*backup.GuildBackups)

	err := bck.RestoreBackup(ctx.GetGuild().ID, fileID, statusChan, errorsChan)
	if err != nil {
		util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
			fmt.Sprintf("An unexpected error occured while restoring backup: ```\n%s\n```", err.Error()))
	}
}

func (c *CmdBackup) purgeBackupsAccept(ctx shireikan.Context) error {
	_, err := acceptmsg.New().
		WithSession(ctx.GetSession()).
		WithEmbed(&discordgo.MessageEmbed{
			Color: static.ColorEmbedOrange,
			Description: ":warning:  **WARNING**  :warning:\n\n" +
				"Do you really want to **purge __all__ backups** for this guild?",
		}).
		LockOnUser(ctx.GetUser().ID).
		DeleteAfterAnswer().
		DoOnDecline(func(_ *discordgo.Message) {
			util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID, "Canceled.").
				DeleteAfter(6 * time.Second).Error()
			return
		}).
		DoOnAccept(func(_ *discordgo.Message) {
			c.purgeBackups(ctx)
		}).
		Send(ctx.GetChannel().ID)

	return err
}

func (c *CmdBackup) purgeBackups(ctx shireikan.Context) {
	db, _ := ctx.GetObject(static.DiDatabase).(database.Database)

	backups, err := db.GetBackups(ctx.GetGuild().ID)
	if err != nil {
		util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
			fmt.Sprintf("Failed getting backups: ```\n%s\n```", err.Error())).
			DeleteAfter(15 * time.Second).Error()
		return
	}

	var lnBackups = len(backups)
	if lnBackups < 1 {
		util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
			"There are no backups saved to be purged.").
			DeleteAfter(8 * time.Second).Error()
		return
	}

	var success int
	for _, backup := range backups {
		if err = db.DeleteBackup(ctx.GetGuild().ID, backup.FileID); err != nil {
			continue
		}

		st, _ := ctx.GetObject(static.DiObjectStorage).(storage.Storage)
		if err = st.DeleteObject(static.StorageBucketBackups, backup.FileID); err != nil {
			continue
		}
		success++
	}

	if success < lnBackups {
		util.SendEmbedError(ctx.GetSession(), ctx.GetChannel().ID,
			fmt.Sprintf("Successfully purged `%d` of `%d` backups.\n`%d` backup purges failed.",
				success, lnBackups, lnBackups-success)).
			DeleteAfter(8 * time.Second).Error()
		return
	}

	util.SendEmbed(ctx.GetSession(), ctx.GetChannel().ID,
		fmt.Sprintf("Successfully purged `%d` of `%d` backups.",
			success, lnBackups), "", static.ColorEmbedGreen).
		DeleteAfter(8 * time.Second).Error()
	return
}
