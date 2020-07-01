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
	"github.com/zekroTJA/shinpuru/internal/core/permissions"
	"github.com/zekroTJA/shinpuru/internal/core/storage"
	"github.com/zekroTJA/shinpuru/internal/core/twitchnotify"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
	"github.com/zekroTJA/shinpuru/pkg/lctimer"
)

const (
	notifiedCmdsCleanupDelay = 5 * time.Minute
	notifiedCmdsExpireTime   = 6 * time.Hour
)

// CmdHandler provides functionalities to register Commands,
// manage registered command instances, manage command
// permissions, colelcting usage statistics and generating
// help pages.
type CmdHandler struct {
	registeredCmds         map[string]Command
	registeredCmdInstances []Command

	db     database.Database
	st     storage.Storage
	config *config.Config
	tnw    *twitchnotify.NotifyWorker
	bck    *backup.GuildBackups
	lct    *lctimer.LifeCycleTimer

	defAdminRules permissions.PermissionArray
	defUserRules  permissions.PermissionArray

	notifiedCmdMsgs *timedmap.TimedMap
}

// NewCmdHandler initializes a new instance of CmdHandler with the passed
// discord Session, database provider, storage provider, configuration,
// twitch notify worker and lifecycle timer.
func NewCmdHandler(s *discordgo.Session, db database.Database, st storage.Storage, config *config.Config, tnw *twitchnotify.NotifyWorker, lct *lctimer.LifeCycleTimer) *CmdHandler {
	cmd := &CmdHandler{
		registeredCmds:         make(map[string]Command),
		registeredCmdInstances: make([]Command, 0),
		db:                     db,
		st:                     st,
		config:                 config,
		tnw:                    tnw,
		lct:                    lct,
		bck:                    backup.New(s, db, st),
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

// RegisterCommand registers the passed command instance to
// the command handler.
func (c *CmdHandler) RegisterCommand(cmd Command) {
	c.registeredCmdInstances = append(c.registeredCmdInstances, cmd)
	for _, invoke := range cmd.GetInvokes() {
		if _, ok := c.registeredCmds[invoke]; ok {
			util.Log.Warningf("Command invoke '%s' was registered more than once!", invoke)
		}
		c.registeredCmds[invoke] = cmd
	}
}

// GetCommand returns the command instance by primary
// invoke or by alias.
func (c *CmdHandler) GetCommand(invoke string) (Command, bool) {
	cmd, ok := c.registeredCmds[invoke]
	return cmd, ok
}

// GetCommandListLen returns the ammount of registered
// command instances.
func (c *CmdHandler) GetCommandListLen() int {
	return len(c.registeredCmdInstances)
}

// IsBotOwner returns true if the passed userID is
// the userID specified as owner in the bots config.
func (c *CmdHandler) IsBotOwner(userID string) bool {
	return userID == c.config.Discord.OwnerID
}

// GetPermissions tries to fetch the permissions array of
// the passed user of the specified guild. The merged
// permissions array is returned as well as the override,
// which is true when the specified user is the bot owner,
// guild owner or an admin of the guild.
func (c *CmdHandler) GetPermissions(s *discordgo.Session, guildID, userID string) (perm permissions.PermissionArray, overrideExplicits bool, err error) {
	if guildID != "" {
		perm, err = c.db.GetMemberPermission(s, guildID, userID)
		if err != nil && !database.IsErrDatabaseNotFound(err) {
			return
		}
	} else {
		perm = make(permissions.PermissionArray, 0)
	}

	if c.IsBotOwner(userID) {
		perm = perm.Merge(permissions.PermissionArray{"+sp.*"}, false)
		overrideExplicits = true
	}

	if guildID != "" {
		guild, err := s.State.Guild(guildID)
		if err != nil {
			return permissions.PermissionArray{}, false, nil
		}

		member, _ := s.GuildMember(guildID, userID)

		if userID == guild.OwnerID || (member != nil && discordutil.IsAdmin(guild, member)) {
			perm = perm.Merge(c.defAdminRules, false)
			overrideExplicits = true
		}
	}

	perm = perm.Merge(c.defUserRules, false)

	return perm, overrideExplicits, nil
}

// CheckPermissions tries to fetch the permissions of the specified user
// on the specified guild and returns true, if the passed dn matches the
// fetched permissions array. Also, the override status is returned as
// well as errors occured during permissions fetching.
func (c *CmdHandler) CheckPermissions(s *discordgo.Session, guildID, userID, dn string) (bool, bool, error) {
	perms, overrideExplicits, err := c.GetPermissions(s, guildID, userID)
	if err != nil {
		return false, false, err
	}

	return perms.Check(dn), overrideExplicits, nil
}

// ExportCommandManual generates a markdown document with the
// description and help and details of all registered commands
// and saves it to the specified file directory.
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

// AddNotifiedCommandMsg marks the specified message ID
// to be excluded from the ghost ping detection if it
// includes a mention.
func (c *CmdHandler) AddNotifiedCommandMsg(msgID string) {
	c.notifiedCmdMsgs.Set(msgID, struct{}{}, notifiedCmdsExpireTime)
}

// GetNotifiedCommandMsgs returns the array of message
// IDs marked with AddNotifiedCommandMsg.
func (c *CmdHandler) GetNotifiedCommandMsgs() *timedmap.TimedMap {
	return c.notifiedCmdMsgs
}

// GetCmdInstances returns the list of all registered
// command instances.
func (c *CmdHandler) GetCmdInstances() []Command {
	return c.registeredCmdInstances
}
