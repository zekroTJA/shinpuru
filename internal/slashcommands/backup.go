package slashcommands

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/snowflake"
	"github.com/zekroTJA/shinpuru/internal/services/backup"
	"github.com/zekroTJA/shinpuru/internal/services/backup/backupmodels"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/services/storage"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/acceptmsg"
	"github.com/zekrotja/ken"
)

const (
	timeFormat = time.RFC1123
)

type Backup struct{}

var (
	_ ken.SlashCommand        = (*Backup)(nil)
	_ permissions.PermCommand = (*Backup)(nil)
)

func (c *Backup) Name() string {
	return "backup"
}

func (c *Backup) Description() string {
	return "Manage guild backups."
}

func (c *Backup) Version() string {
	return "1.0.0"
}

func (c *Backup) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *Backup) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "state",
			Description: "Enable or disable the backup system.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "state",
					Description: "Dispaly or set the backup state to enabled or disabled",
					Required:    false,
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "list",
			Description: "List all stored backups.",
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "restore",
			Description: "Restore a backup.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "index",
					Description: "The index of the backup to be restored.",
					Required:    true,
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "purge",
			Description: "Delete all stored backups.",
		},
	}
}

func (c *Backup) Domain() string {
	return "sp.guild.admin.backup"
}

func (c *Backup) SubDomains() []permissions.SubPermission {
	return nil
}

func (c *Backup) Run(ctx *ken.Ctx) (err error) {
	if err = ctx.Defer(); err != nil {
		return
	}

	err = ctx.HandleSubCommands(
		ken.SubCommandHandler{"state", c.state},
		ken.SubCommandHandler{"list", c.list},
		ken.SubCommandHandler{"restore", c.restore},
		ken.SubCommandHandler{"purge", c.purge},
	)

	return
}

func (c *Backup) state(ctx *ken.SubCommandCtx) (err error) {
	db := ctx.Get(static.DiDatabase).(database.Database)

	var (
		state bool
		emb   *discordgo.MessageEmbed
	)
	subOpts := ctx.Event.ApplicationCommandData().Options[0].Options
	if len(subOpts) > 0 {
		state = subOpts[0].BoolValue()
		if err = db.SetGuildBackup(ctx.Event.GuildID, state); err != nil {
			return
		}
		emb = &discordgo.MessageEmbed{
			Color:       static.ColorEmbedOrange,
			Description: "The backup system is now **disabled**.",
		}
		if state {
			emb.Color = static.ColorEmbedGreen
			emb.Description = "The backup system is now **enabled**."
		}
	} else {
		state, err = db.GetGuildBackup(ctx.Event.GuildID)
		if err != nil {
			return
		}
		emb = &discordgo.MessageEmbed{
			Color:       static.ColorEmbedOrange,
			Description: "The backup system is currently **disabled**.",
		}
		if state {
			emb.Color = static.ColorEmbedGreen
			emb.Description = "The backup system is currently **enabled**."
		}
	}

	err = ctx.FollowUpEmbed(emb).Error
	return
}

func (c *Backup) list(ctx *ken.SubCommandCtx) (err error) {
	db, _ := ctx.Get(static.DiDatabase).(database.Database)

	status, err := db.GetGuildBackup(ctx.Event.GuildID)
	if err != nil && database.IsErrDatabaseNotFound(err) {
		return err
	}

	strStatus := ":x:  Backups **disabled**"
	if status {
		strStatus = ":white_check_mark:  Backups **enabled**"
	}

	_, strBackupAll, err := c.getBackupsList(ctx.Ctx)
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

	err = ctx.FollowUpEmbed(emb).Error
	return
}

func (c *Backup) restore(ctx *ken.SubCommandCtx) (err error) {
	db := ctx.Get(static.DiDatabase).(database.Database)

	i := ctx.Options().Get(0).IntValue()
	if err != nil {
		return err
	}

	if i < 0 {
		return ctx.FollowUpError("Index must be between 0 and 9 or a snowflake ID.", "").Error
	}

	var backup *backupmodels.Entry

	backups, err := db.GetBackups(ctx.Event.GuildID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return
	}

	if i < 10 {
		if len(backups)-1 < int(i) {
			return ctx.FollowUpError(
				fmt.Sprintf("There are only %d (index 0 to %d) backups you can chose from.",
					len(backups), len(backups)-1), "").Error
		}
		backup = backups[i]
	} else {
		sf := snowflake.ParseInt64(i).String()
		for _, b := range backups {
			if b.FileID == sf {
				backup = b
			}
		}
	}

	if backup == nil {
		return ctx.FollowUpError(
			fmt.Sprintf("Could not find any backup by this specifier: ```\n%d\n```", i), "").
			Error
	}

	accMsg := &acceptmsg.AcceptMessage{
		Session:        ctx.Session,
		DeleteMsgAfter: true,
		UserID:         ctx.User().ID,
		Embed: &discordgo.MessageEmbed{
			Color: static.ColorEmbedOrange,
			Description: fmt.Sprintf(":warning:  **WARNING**  :warning:\n\n"+
				"By pressing :white_check_mark:, the structure of this guild will be **reset** to the selected backup:\n\n"+
				"%s - (ID: `%s`)", backup.Timestamp.Format(timeFormat), backup.FileID),
		},
		DeclineFunc: func(m *discordgo.Message) (err error) {
			err = util.SendEmbedError(ctx.Session, ctx.Event.ChannelID, "Canceled.").
				DeleteAfter(6 * time.Second).Error()
			return
		},
		AcceptFunc: func(m *discordgo.Message) (err error) {
			c.proceedRestore(ctx.Ctx, backup.FileID)
			return
		},
	}

	if _, err = accMsg.AsFollowUp(ctx.Ctx); err != nil {
		return err
	}
	return accMsg.Error()
}

func (c *Backup) purge(ctx *ken.SubCommandCtx) (err error) {
	if err = ctx.Defer(); err != nil {
		return
	}

	am, err := acceptmsg.New().
		WithSession(ctx.Session).
		WithEmbed(&discordgo.MessageEmbed{
			Color: static.ColorEmbedOrange,
			Description: ":warning:  **WARNING**  :warning:\n\n" +
				"Do you really want to **purge __all__ backups** for this guild?",
		}).
		LockOnUser(ctx.User().ID).
		DeleteAfterAnswer().
		DoOnDecline(func(_ *discordgo.Message) (err error) {
			err = util.SendEmbedError(ctx.Session, ctx.Event.ChannelID, "Canceled.").
				DeleteAfter(6 * time.Second).Error()
			return
		}).
		DoOnAccept(func(_ *discordgo.Message) (err error) {
			c.purgeBackups(ctx.Ctx)
			return
		}).
		AsFollowUp(ctx.Ctx)

	if err != nil {
		return err
	}

	return am.Error()
}

// --- HELPERS ---

func (c *Backup) getBackupsList(ctx *ken.Ctx) ([]*backupmodels.Entry, string, error) {
	db, _ := ctx.Get(static.DiDatabase).(database.Database)

	backups, err := db.GetBackups(ctx.Event.GuildID)
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

func (c *Backup) proceedRestore(ctx *ken.Ctx, fileID string) {
	statusChan := make(chan string)
	errorsChan := make(chan error)

	statusMsg, _ := ctx.Session.ChannelMessageSendEmbed(ctx.Event.ChannelID,
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
					ctx.Session.ChannelMessageEditEmbed(statusMsg.ChannelID, statusMsg.ID, &discordgo.MessageEmbed{
						Color:       static.ColorEmbedGray,
						Description: status + "...",
					})
				case err, ok := <-errorsChan:
					if !ok || err == nil {
						continue
					}
					util.SendEmbedError(ctx.Session, ctx.Event.ChannelID,
						"An unexpected error occured while restoring backup (process will not be aborted): ```\n"+err.Error()+"\n```")
				}
			}
		}()
	}

	bck, _ := ctx.Get(static.DiBackupHandler).(*backup.GuildBackups)

	err := bck.RestoreBackup(ctx.Event.GuildID, fileID, statusChan, errorsChan)
	if err != nil {
		util.SendEmbedError(ctx.Session, ctx.Event.ChannelID,
			fmt.Sprintf("An unexpected error occured while restoring backup: ```\n%s\n```", err.Error()))
	}
}

func (c *Backup) purgeBackups(ctx *ken.Ctx) {
	db, _ := ctx.Get(static.DiDatabase).(database.Database)

	backups, err := db.GetBackups(ctx.Event.GuildID)
	if err != nil {
		util.SendEmbedError(ctx.Session, ctx.Event.GuildID,
			fmt.Sprintf("Failed getting backups: ```\n%s\n```", err.Error())).
			DeleteAfter(15 * time.Second).Error()
		return
	}

	var lnBackups = len(backups)
	if lnBackups < 1 {
		util.SendEmbedError(ctx.Session, ctx.Event.GuildID,
			"There are no backups saved to be purged.").
			DeleteAfter(8 * time.Second).Error()
		return
	}

	var success int
	for _, backup := range backups {
		if err = db.DeleteBackup(ctx.Event.GuildID, backup.FileID); err != nil {
			continue
		}

		st, _ := ctx.Get(static.DiObjectStorage).(storage.Storage)
		if err = st.DeleteObject(static.StorageBucketBackups, backup.FileID); err != nil {
			continue
		}
		success++
	}

	if success < lnBackups {
		util.SendEmbedError(ctx.Session, ctx.Event.ChannelID,
			fmt.Sprintf("Successfully purged `%d` of `%d` backups.\n`%d` backup purges failed.",
				success, lnBackups, lnBackups-success)).
			DeleteAfter(8 * time.Second).Error()
		return
	}

	util.SendEmbed(ctx.Session, ctx.Event.ChannelID,
		fmt.Sprintf("Successfully purged `%d` of `%d` backups.",
			success, lnBackups), "", static.ColorEmbedGreen).
		DeleteAfter(8 * time.Second).Error()
	return
}
