package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/core/config"
	"github.com/zekroTJA/shinpuru/internal/core/middleware"
	"github.com/zekroTJA/shinpuru/internal/inits"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shireikan"
)

var (
	flagExportFile = flag.String("o", "commandsManual.md", "output location of manual file")
)

//////////////////////////////////////////////////////////////////////
//
//   CMDMAN
//   ------
//   This tool initializes the command handler of shinpuru with all
//   commands and automatically generates a markdown manual from the
//   commands metadata which is then output in a file.
//   You can specify the output location with the parameter '-o'.
//
//////////////////////////////////////////////////////////////////////

func main() {
	flag.Parse()

	// Setting Release flag to true manually to prevent
	// registration of test command and exclude it in the
	// command manual.
	util.Release = "TRUE"

	config := &config.Config{
		Discord: &config.Discord{},
	}

	s, _ := discordgo.New()
	database := new(middleware.SqliteMiddleware)

	cmdHandler := inits.InitCommandHandler(s, config, database, nil, nil, nil, nil, nil)
	if err := exportCommandManual(cmdHandler, *flagExportFile); err != nil {
		util.Log.Fatal("Failed exporting command manual: ", err)
	}
	util.Log.Info("Successfully exported command manual file to " + *flagExportFile)
}

// exportCommandManual generates a markdown text file
// from the registered command instances of the passed
// command handler. The file is then exported to the
// given file location. Occuring errors are returned.
func exportCommandManual(cmdHandler shireikan.Handler, fileName string) error {
	document := "> Auto generated command manual | " + time.Now().Format(time.RFC1123) + "\n\n" +
		"# Explicit Sub Commands\n\n" +
		"The commands below have sub command permissions which must be set explicitly and can not be " +
		"applied by wildcards (`*`). So here you have them if you want to allow them for specific roles:\n\n"

	for _, cmd := range cmdHandler.GetCommandInstances() {
		if spr := cmd.GetSubPermissionRules(); spr != nil {
			document += fmt.Sprintf("**%s**\n\n", cmd.GetInvokes()[0])

			for _, perm := range spr {
				if perm.Explicit {
					document += fmt.Sprintf("- **`%s.%s`** - %s\n",
						cmd.GetDomainName(), perm.Term, perm.Description)
				}
			}

			document += "\n"
		}
	}

	document += "\n# Command List\n\n"

	cmdCats := make(map[string][]shireikan.Command)
	cmdDetails := "# Command Details\n\n"

	for _, cmd := range cmdHandler.GetCommandInstances() {
		if cat, ok := cmdCats[cmd.GetGroup()]; !ok {
			cmdCats[cmd.GetGroup()] = []shireikan.Command{cmd}
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
	cmdCatsSorted := make(map[string][]shireikan.Command)

	for _, cat := range catKeys {
		cmds := cmdCats[cat]

		sort.Slice(cmds, func(i, j int) bool {
			return cmds[i].GetInvokes()[0] < cmds[j].GetInvokes()[0]
		})

		cmdCatsSorted[cat] = cmds

		document += fmt.Sprintf("## %s\n", cat)
		cmdDetails += fmt.Sprintf("## %s\n\n", cat)

		for _, cmd := range cmds {
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
					"| Aliases | %s |\n"+
					"| DM Capable | %s |\n\n"+
					"**Usage**  \n"+
					"%s\n\n",
				cmd.GetInvokes()[0], cmd.GetDescription(), cmd.GetDomainName(), cmd.GetGroup(),
				aliases, util.BoolAsString(cmd.IsExecutableInDMChannels(), "Yes", "No"), help)

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
