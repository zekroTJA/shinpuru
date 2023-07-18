package slashcommands

import (
	"fmt"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/xid"
	"github.com/zekroTJA/shinpuru/internal/services/backup"
	"github.com/zekroTJA/shinpuru/internal/services/backup/backupmodels"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/services/storage"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/acceptmsg/v2"
	"github.com/zekroTJA/shinpuru/pkg/logmsg"
	"github.com/zekrotja/ken"
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
	return "2.0.0"
}

func (c *Backup) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *Backup) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		// {
		// 	Type:        discordgo.ApplicationCommandOptionSubCommand,
		// 	Name:        "state",
		// 	Description: "Enable or disable the backup system.",
		// 	Options: []*discordgo.ApplicationCommandOption{
		// 		{
		// 			Type:        discordgo.ApplicationCommandOptionBoolean,
		// 			Name:        "state",
		// 			Description: "Dispaly or set the backup state to enabled or disabled",
		// 			Required:    false,
		// 		},
		// 	},
		// },
		// {
		// 	Type:        discordgo.ApplicationCommandOptionSubCommand,
		// 	Name:        "list",
		// 	Description: "List all stored backups.",
		// },
		// {
		// 	Type:        discordgo.ApplicationCommandOptionSubCommand,
		// 	Name:        "restore",
		// 	Description: "Restore a backup.",
		// 	Options: []*discordgo.ApplicationCommandOption{
		// 		{
		// 			Type:        discordgo.ApplicationCommandOptionInteger,
		// 			Name:        "index",
		// 			Description: "The index of the backup to be restored.",
		// 			Required:    true,
		// 		},
		// 	},
		// },
		// {
		// 	Type:        discordgo.ApplicationCommandOptionSubCommand,
		// 	Name:        "purge",
		// 	Description: "Delete all stored backups.",
		// },
	}
}

func (c *Backup) Domain() string {
	return "sp.guild.admin.backup"
}

func (c *Backup) SubDomains() []permissions.SubPermission {
	return nil
}

func (c *Backup) Run(ctx ken.Context) (err error) {
	if err = ctx.Defer(); err != nil {
		return
	}

	db, _ := ctx.Get(static.DiDatabase).(database.Database)
	st, _ := ctx.Get(static.DiObjectStorage).(storage.Storage)

	enabled, err := db.GetGuildBackup(ctx.GetEvent().GuildID)
	if err != nil && database.IsErrDatabaseNotFound(err) {
		return err
	}

	strStatus := ":x:  Backups **disabled**"
	if enabled {
		strStatus = ":white_check_mark:  Backups **enabled**"
	}

	entries, strBackupAll, err := c.getBackupsList(ctx)
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

	var unreg func() error
	fum := ctx.FollowUpEmbed(emb).Send()
	if fum.Error != nil {
		return err
	}

	cNext := make(chan string, 1)

	builder := fum.AddComponents().
		Condition(func(cctx ken.ComponentContext) bool {
			return cctx.User().ID == ctx.User().ID
		})

	if len(entries) != 0 {
		options := make([]discordgo.SelectMenuOption, 0, len(entries))
		for i, entry := range entries {
			options = append(options, discordgo.SelectMenuOption{
				Label: fmt.Sprintf("%d - %s", i, entry.TimestampFormatted()),
				Value: entry.FileID,
			})
		}

		builder.AddActionsRow(func(b ken.ComponentAssembler) {
			b.Add(discordgo.SelectMenu{
				CustomID:    xid.New().String(),
				Options:     options,
				Placeholder: "Select backup for Restore",
			}, func(ctx ken.ComponentContext) bool {
				vals := ctx.GetData().Values
				if len(vals) == 0 {
					return false
				}
				cNext <- vals[0]
				return true
			})
		})
	}

	builder.AddActionsRow(func(b ken.ComponentAssembler) {
		if enabled {
			b.Add(discordgo.Button{
				CustomID: xid.New().String(),
				Label:    "Disable Guild Backups",
				Style:    discordgo.DangerButton,
			}, func(ctx ken.ComponentContext) bool {
				if ctx.Defer() != nil {
					return false
				}

				err := db.SetGuildBackup(ctx.GetEvent().GuildID, false)
				if err != nil {
					return false
				}

				ctx.FollowUpEmbed(&discordgo.MessageEmbed{
					Description: "Guild backups are now disabled.",
					Color:       static.ColorEmbedOrange,
				}).Send()

				cNext <- ""
				return true
			})
		} else {
			b.Add(discordgo.Button{
				CustomID: xid.New().String(),
				Label:    "Enable Guild Backups",
				Style:    discordgo.SuccessButton,
			}, func(ctx ken.ComponentContext) bool {
				if ctx.Defer() != nil {
					return false
				}

				err := db.SetGuildBackup(ctx.GetEvent().GuildID, true)
				if err != nil {
					return false
				}

				ctx.FollowUpEmbed(&discordgo.MessageEmbed{
					Description: "Guild backups are now enabled.",
					Color:       static.ColorEmbedGreen,
				}).Send()

				cNext <- ""
				return true
			})
		}

		if len(entries) != 0 {
			b.Add(discordgo.Button{
				CustomID: xid.New().String(),
				Label:    "Purge all Backups",
				Style:    discordgo.DangerButton,
			}, func(ctx ken.ComponentContext) bool {
				c.purgeBackups(ctx, db, st)

				cNext <- ""
				return true
			})
		}

		b.Add(discordgo.Button{
			CustomID: xid.New().String(),
			Label:    "Cancel",
			Style:    discordgo.SecondaryButton,
		}, func(ctx ken.ComponentContext) bool {
			cNext <- ""
			return true
		})
	})

	unreg, err = builder.Build()
	if err != nil {
		return err
	}

	id := <-cNext

	unreg()
	fum.Delete()

	if id == "" {
		return nil
	}

	var entry backupmodels.Entry
	for _, entry = range entries {
		if entry.FileID == id {
			break
		}
	}

	if entry.FileID == "" {
		return ctx.FollowUpError(
			"Something went wrong. Please try again later.", "").Send().Error
	}

	bck := ctx.Get(static.DiBackupHandler).(*backup.GuildBackups)

	accMsg := &acceptmsg.AcceptMessage{
		Ken:            ctx.GetKen(),
		DeleteMsgAfter: true,
		UserID:         ctx.User().ID,
		Embed: &discordgo.MessageEmbed{
			Color: static.ColorEmbedOrange,
			Description: fmt.Sprintf(":warning:  **WARNING**  :warning:\n\n"+
				"By pressing \"Accept\", the structure of this guild will be **reset** to the selected backup:\n\n"+
				"%s - (ID: `%s`)", entry.TimestampFormatted(), entry.FileID),
		},
		DeclineFunc: func(cctx ken.ComponentContext) error {
			return cctx.RespondError("Canceled.", "")
		},
		AcceptFunc: func(cctx ken.ComponentContext) error {
			return c.proceedRestore(cctx, bck, entry.FileID)
		},
	}

	if _, err = accMsg.AsFollowUp(ctx); err != nil {
		return err
	}
	return accMsg.Error()
}

// --- HELPERS ---

func (c *Backup) getBackupsList(ctx ken.Context) ([]backupmodels.Entry, string, error) {
	db, _ := ctx.Get(static.DiDatabase).(database.Database)

	backups, err := db.GetBackups(ctx.GetEvent().GuildID)
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
			strBackups[i] = b.StringIndexed(i)
		}

		strBackupAll = strings.Join(strBackups, "\n")
	}

	return backups, strBackupAll, nil
}

func (c *Backup) proceedRestore(ctx ken.ComponentContext, bck *backup.GuildBackups, fileID string) (err error) {
	if err = ctx.Defer(); err != nil {
		return err
	}

	statusChan := make(chan string)
	errorsChan := make(chan error)

	statusMsg, err := logmsg.NewWithSender(
		ctx.GetSession(),
		func(emb *discordgo.MessageEmbed) (*discordgo.Message, error) {
			fum := ctx.FollowUpEmbed(emb).Send()
			return fum.Message, fum.Error
		},
		&discordgo.MessageEmbed{
			Title: "Backup Restoration Status",
			Color: static.ColorEmbedGray,
		},
		statusChan,
		errorsChan,
		"initializing backup restoring...")
	if err != nil {
		return
	}
	defer statusMsg.Close("✔️ Backup restoration finished!")

	err = bck.RestoreBackup(ctx.GetEvent().GuildID, fileID, statusChan, errorsChan)

	return
}

func (c *Backup) purgeBackups(ctx ken.ComponentContext, db database.Database, st storage.Storage) {
	if err := ctx.Defer(); err != nil {
		return
	}

	backups, err := db.GetBackups(ctx.GetEvent().GuildID)
	if err != nil {
		ctx.FollowUpError(fmt.Sprintf("Failed getting backups: ```\n%s\n```", err.Error()), "").
			Send()
		return
	}

	var lnBackups = len(backups)
	if lnBackups < 1 {
		ctx.FollowUpError("There are no backups saved to be purged.", "").
			Send()
		return
	}

	var success int
	for _, backup := range backups {
		if err = db.DeleteBackup(ctx.GetEvent().GuildID, backup.FileID); err != nil {
			continue
		}

		if err = st.DeleteObject(static.StorageBucketBackups, backup.FileID); err != nil {
			continue
		}
		success++
	}

	if success < lnBackups {
		ctx.FollowUpError(fmt.Sprintf("Successfully purged `%d` of `%d` backups.\n`%d` backup purges failed.",
			success, lnBackups, lnBackups-success), "").Send()
		return
	}

	ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: fmt.Sprintf("Successfully purged `%d` of `%d` backups.",
			success, lnBackups),
		Color: static.ColorEmbedGreen,
	}).Send()
}
