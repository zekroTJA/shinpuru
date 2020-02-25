package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/zekroTJA/timedmap"

	"github.com/zekroTJA/shinpuru/internal/core/backup"
	"github.com/zekroTJA/shinpuru/internal/core/config"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/core/lctimer"
	"github.com/zekroTJA/shinpuru/internal/core/permissions"
	"github.com/zekroTJA/shinpuru/internal/core/twitchnotify"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
)

const (
	notifiedCmdsCleanupDelay = 5 * time.Minute
	notifiedCmdsExpireTime   = 6 * time.Hour
)

type CmdHandler struct {
	registeredCmds         map[string]Command
	registeredCmdInstances []Command

	db     database.Database
	config *config.Config
	tnw    *twitchnotify.NotifyWorker
	bck    *backup.GuildBackups
	lct    *lctimer.LCTimer

	defAdminRules permissions.PermissionArray
	defUserRules  permissions.PermissionArray

	notifiedCmdMsgs *timedmap.TimedMap
}

func NewCmdHandler(s *discordgo.Session, db database.Database, config *config.Config, tnw *twitchnotify.NotifyWorker, lct *lctimer.LCTimer) *CmdHandler {
	cmd := &CmdHandler{
		registeredCmds:         make(map[string]Command),
		registeredCmdInstances: make([]Command, 0),
		db:                     db,
		config:                 config,
		tnw:                    tnw,
		lct:                    lct,
		bck:                    backup.New(s, db, config.Discord.GuildBackupLoc),
		notifiedCmdMsgs:        timedmap.New(notifiedCmdsCleanupDelay),
		defAdminRules:          static.DefaultAdminRules,
		defUserRules:           static.DefaultUserRules,
	}

	if config.Permissions != nil {
		if config.Permissions.DefaultAdminRules != nil {
			cmd.defAdminRules = config.Permissions.DefaultAdminRules
		}
		if config.Permissions.DefaultUserRules != nil {
			cmd.defUserRules = config.Permissions.DefaultUserRules
		}
	}

	return cmd
}

func (c *CmdHandler) RegisterCommand(cmd Command) {
	c.registeredCmdInstances = append(c.registeredCmdInstances, cmd)
	for _, invoke := range cmd.GetInvokes() {
		if _, ok := c.registeredCmds[invoke]; ok {
			util.Log.Warningf("Command invoke '%s' was registered more than once!", invoke)
		}
		c.registeredCmds[invoke] = cmd
	}
}

func (c *CmdHandler) GetCommand(invoke string) (Command, bool) {
	cmd, ok := c.registeredCmds[invoke]
	return cmd, ok
}

func (c *CmdHandler) GetCommandListLen() int {
	return len(c.registeredCmdInstances)
}

func (c *CmdHandler) IsBotOwner(userID string) bool {
	return userID == c.config.Discord.OwnerID
}

func (c *CmdHandler) GetPermissions(s *discordgo.Session, guildID, userID string) (permissions.PermissionArray, error) {
	if c.IsBotOwner(userID) {
		return permissions.PermissionArray{"+sp.*"}, nil
	}

	if guildID != "" {
		guild, err := s.Guild(guildID)
		if err != nil {
			return permissions.PermissionArray{}, nil
		}

		member, _ := s.GuildMember(guildID, userID)

		if member != nil {
			if util.IsAdmin(guild, member) {
				return c.defAdminRules, nil
			}
		}

		if userID == guild.OwnerID {
			return c.defAdminRules, nil
		}
	}

	perm, err := c.db.GetMemberPermission(s, guildID, userID)

	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return nil, err
	}

	perm = perm.Merge(c.defUserRules, false)

	return perm, nil
}

func (c *CmdHandler) CheckPermissions(s *discordgo.Session, guildID, userID, dn string) (bool, error) {
	perms, err := c.GetPermissions(s, guildID, userID)
	if err != nil {
		return false, err
	}

	return permissions.PermissionCheck(dn, perms), nil
}

func (c *CmdHandler) ExportCommandManual(fileName string) error {
	document := "> Auto generated command manual | " + time.Now().Format(time.RFC1123) + "\n\n" +
		"# Explicit Sub Commands\n\n" +
		"The commands below have sub command permissions which must be set explicitly and can not be " +
		"applied by wildcards (`*`). So here you have them if you want to allow them for specific roles:\n\n"

	for _, cmd := range c.registeredCmdInstances {
		if spr := cmd.GetSubPermissionRules(); spr != nil {
			document += fmt.Sprintf("### %s\n\n", cmd.GetInvokes()[0])

			for _, perm := range spr {
				if perm.Explicit {
					document += fmt.Sprintf("- **`%s.%s`** - %s\n",
						cmd.GetDomainName(), perm.Term, perm.Description)
				}
			}
		}
	}

	document += "\n# Command List\n\n"

	cmdCats := make(map[string][]Command)
	cmdDetails := "# Command Details\n\n"

	for _, cmd := range c.registeredCmdInstances {
		if cat, ok := cmdCats[cmd.GetGroup()]; !ok {
			cmdCats[cmd.GetGroup()] = []Command{cmd}
		} else {
			cmdCats[cmd.GetGroup()] = append(cat, cmd)
		}
	}

	catKeys := make([]string, len(cmdCats))
	i := 0
	for cat := range cmdCats {
		catKeys[i] = cat
		i++
	}

	sort.Strings(catKeys)
	cmdCatsSorted := make(map[string][]Command)

	for _, cat := range catKeys {
		cmdCatsSorted[cat] = cmdCats[cat]
		document += fmt.Sprintf("## %s\n", cat)
		cmdDetails += fmt.Sprintf("## %s\n\n", cat)

		for _, cmd := range cmdCats[cat] {
			document += fmt.Sprintf("- [%s](#%s)\n", cmd.GetInvokes()[0], cmd.GetInvokes()[0])
			aliases := strings.Join(cmd.GetInvokes()[1:], ", ")
			help := strings.Replace(cmd.GetHelp(), "\n", "  \n", -1)
			cmdDetails += fmt.Sprintf(
				"### %s\n\n"+
					"> %s\n\n"+
					"| | |\n"+
					"|---|---|\n"+
					"| Domain Name | %s |\n"+
					"| Group | %s |\n"+
					"| Aliases | %s |\n\n"+
					"**Usage**  \n"+
					"%s\n\n", cmd.GetInvokes()[0], cmd.GetDescription(), cmd.GetDomainName(), cmd.GetGroup(), aliases, help)

			if spr := cmd.GetSubPermissionRules(); spr != nil {
				cmdDetails += "\n**Sub Permission Rules**\n"
				for _, perm := range spr {
					explicit := ""
					if perm.Explicit {
						explicit = "`[EXPLICIT]` "
					}
					cmdDetails += fmt.Sprintf("- **`%s.%s`** %s- %s\n",
						cmd.GetDomainName(), perm.Term, explicit, perm.Description)
				}
				cmdDetails += "\n\n"
			}
		}

		document += "\n"
	}

	document += cmdDetails

	filePath := path.Dir(fileName)
	if filePath != "." && filePath != "/" {
		if err := os.MkdirAll(filePath, 0744); err != nil {
			return err
		}
	}
	return ioutil.WriteFile(fileName, []byte(document), 0644)
}

func (c *CmdHandler) AddNotifiedCommandMsg(msgID string) {
	c.notifiedCmdMsgs.Set(msgID, struct{}{}, notifiedCmdsExpireTime)
}

func (c *CmdHandler) GetNotifiedCommandMsgs() *timedmap.TimedMap {
	return c.notifiedCmdMsgs
}

func (c *CmdHandler) GetCmdInstances() []Command {
	return c.registeredCmdInstances
}
