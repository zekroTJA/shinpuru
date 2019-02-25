package commands

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/core"
	"github.com/zekroTJA/shinpuru/internal/util"
)

type CmdBackup struct {
	PermLvl int
}

func (c *CmdBackup) GetInvokes() []string {
	return []string{"backup", "bckp", "guildbackup"}
}

func (c *CmdBackup) GetDescription() string {
	return "enable, disable and manage guild backups"
}

func (c *CmdBackup) GetHelp() string {
	return "`backup <enable|disable>` - enable or disable backups for your guild\n" +
		"`backup (list)` - list all saved backups\n" +
		"`backup restore <id>` - restore a backup"
}

func (c *CmdBackup) GetGroup() string {
	return GroupGuildAdmin
}

func (c *CmdBackup) GetPermission() int {
	return c.PermLvl
}

func (c *CmdBackup) SetPermission(permLvl int) {
	c.PermLvl = permLvl
}

func (c *CmdBackup) Exec(args *CommandArgs) error {
	if len(args.Args) > 0 {
		switch strings.ToLower(args.Args[0]) {
		case "e", "enable":
			return c.switchStatus(args, true)
		case "d", "disable":
			return c.switchStatus(args, false)
		default:
			return c.list(args)
		}
	}
	return c.list(args)
}

func (c *CmdBackup) switchStatus(args *CommandArgs, enable bool) error {
	err := args.CmdHandler.db.SetGuildBackup(args.Guild.ID, enable)
	if err != nil {
		return err
	}

	if enable {
		msg, err := util.SendEmbed(args.Session, args.Channel.ID, "Enabled backup for this guild.\nA full guild backup *(incl. Members, Roles, Channels and Guild Settings)* "+
			"will be created every 12 hours. Only 10 backups per guild will be saved, so you will habe the backup files of the last 5 days.", "", util.ColorEmbedGreen)
		util.DeleteMessageLater(args.Session, msg, 15*time.Second)
		return err
	}

	msg, err := util.SendEmbed(args.Session, args.Channel.ID, "Backup creation disabled.\n"+
		"You will be still have access to created backups and be able to restore them.", "", util.ColorEmbedOrange)
	util.DeleteMessageLater(args.Session, msg, 15*time.Second)
	return err
}

func (c *CmdBackup) getBackupsList(args *CommandArgs) ([]*core.BackupEntry, string, error) {
	backups, err := args.CmdHandler.db.GetBackups(args.Guild.ID)
	if err != nil && core.IsErrDatabaseNotFound(err) {
		return nil, "", err
	}

	strBackupAll := "*no backups saved*"

	if len(backups) > 0 {
		sort.Slice(backups, func(i, j int) bool {
			return backups[i].Timestamp.Before(backups[j].Timestamp)
		})

		strBackups := make([]string, len(backups))

		for i, b := range backups {
			strBackups[i] = fmt.Sprintf("`%d` - %s - *`(ID: %s)`*", i, b.Timestamp.Format(time.RFC1123), b.FileID)
		}

		strBackupAll = strings.Join(strBackups, "\n")
	}

	return backups, strBackupAll, nil
}

func (c *CmdBackup) list(args *CommandArgs) error {
	status, err := args.CmdHandler.db.GetGuildBackup(args.Guild.ID)
	if err != nil && core.IsErrDatabaseNotFound(err) {
		return err
	}

	strStatus := ":x:  Backups **disabled**"
	if status {
		strStatus = ":white_check_mark:  Backups **enabled**"
	}

	_, strBackupAll, err := c.getBackupsList(args)
	if err != nil {
		return err
	}

	emb := &discordgo.MessageEmbed{
		Color:       util.ColorEmbedDefault,
		Title:       "Backups",
		Description: strStatus,
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:  "Saved Backups",
				Value: strBackupAll,
			},
		},
	}

	_, err = args.Session.ChannelMessageSendEmbed(args.Channel.ID, emb)
	return err
}
