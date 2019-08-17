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

	"github.com/zekroTJA/shinpuru/internal/core"
	"github.com/zekroTJA/shinpuru/internal/util"
)

const (
	notifiedCmdsCleanupDelay = 5 * time.Minute
	notifiedCmdsExpireTime   = 6 * time.Hour
)

type CmdHandler struct {
	registeredCmds         map[string]Command
	registeredCmdInstances []Command

	db     core.Database
	config *core.Config
	tnw    *core.TwitchNotifyWorker
	bck    *core.GuildBackups
	lct    *core.LCTimer

	defAdminRules core.PermissionArray
	defUserRules  core.PermissionArray

	notifiedCmdMsgs *timedmap.TimedMap
}

func NewCmdHandler(s *discordgo.Session, db core.Database, config *core.Config, tnw *core.TwitchNotifyWorker, lct *core.LCTimer) *CmdHandler {
	cmd := &CmdHandler{
		registeredCmds:         make(map[string]Command),
		registeredCmdInstances: make([]Command, 0),
		db:                     db,
		config:                 config,
		tnw:                    tnw,
		lct:                    lct,
		bck:                    core.NewGuildBackups(s, db),
		notifiedCmdMsgs:        timedmap.New(notifiedCmdsCleanupDelay),
		defAdminRules:          util.DefaultAdminRules,
		defUserRules:           util.DefaultUserRules,
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

func (c *CmdHandler) GetPermissions(s *discordgo.Session, guildID, userID string) (core.PermissionArray, error) {
	if c.IsBotOwner(userID) {
		return core.PermissionArray{"+sp.*"}, nil
	}

	if guildID != "" {
		guild, err := s.Guild(guildID)
		if err != nil {
			return core.PermissionArray{}, nil
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

	if err != nil && !core.IsErrDatabaseNotFound(err) {
		return nil, err
	}

	perm = perm.Merge(c.defUserRules)

	return perm, nil
}

func (c *CmdHandler) CheckPermissions(s *discordgo.Session, guildID, userID, dn string) (bool, error) {
	perms, err := c.GetPermissions(s, guildID, userID)
	if err != nil {
		return false, err
	}

	return core.PermissionCheck(dn, perms), nil
}

func (c *CmdHandler) ExportCommandManual(fileName string) error {
	document := "> Auto generated command manual | " + time.Now().Format(time.RFC1123) + "\n\n" +
		"# Command List\n\n"

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
