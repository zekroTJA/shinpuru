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
	"github.com/sarulabs/di/v2"
	"github.com/sirupsen/logrus"
	"github.com/zekroTJA/shinpuru/internal/config"
	"github.com/zekroTJA/shinpuru/internal/inits"
	"github.com/zekroTJA/shinpuru/internal/middleware"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/database/sqlite"
	"github.com/zekroTJA/shinpuru/internal/util/embedded"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/stringutil"
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
	embedded.Release = "TRUE"

	diBuilder, _ := di.NewBuilder()

	config := &config.Config{
		Discord: &config.Discord{},
		Logging: &config.Logging{},
	}
	diBuilder.Set(static.DiConfig, config)

	s, _ := discordgo.New()
	diBuilder.Set(static.DiDiscordSession, s)

	// Initialize database middleware and shutdown routine
	diBuilder.Add(di.Def{
		Name: static.DiDatabase,
		Build: func(ctn di.Container) (interface{}, error) {
			return new(sqlite.SqliteMiddleware), nil
		},
		Close: func(obj interface{}) error {
			database := obj.(database.Database)
			logrus.Info("Shutting down database connection...")
			database.Close()
			return nil
		},
	})

	// Initialize command handler
	diBuilder.Add(di.Def{
		Name: static.DiCommandHandler,
		Build: func(ctn di.Container) (interface{}, error) {
			return inits.InitCommandHandler(ctn), nil
		},
	})

	diBuilder.Set(static.DiPermissionMiddleware, &middleware.PermissionsMiddleware{})
	diBuilder.Set(static.DiGhostpingIgnoreMiddleware, &middleware.GhostPingIgnoreMiddleware{})

	ctn := diBuilder.Build()

	cmdHandler := ctn.Get(static.DiCommandHandler).(shireikan.Handler)
	if err := exportCommandManual(cmdHandler, *flagExportFile); err != nil {
		logrus.WithError(err).Fatal("Failed exporting command manual")
	}
	logrus.Info("Successfully exported command manual file to " + *flagExportFile)
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
					document += fmt.Sprintf("- **`%s`** - %s\n",
						getTermAssembly(cmd, perm.Term), perm.Description)
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
				aliases, stringutil.FromBool(cmd.IsExecutableInDMChannels(), "Yes", "No"), help)

			if spr := cmd.GetSubPermissionRules(); spr != nil {
				cmdDetails += "\n**Sub Permission Rules**\n"
				for _, perm := range spr {
					explicit := ""
					if perm.Explicit {
						explicit = "`[EXPLICIT]` "
					}
					cmdDetails += fmt.Sprintf("- **`%s`** %s- %s\n",
						getTermAssembly(cmd, perm.Term), explicit, perm.Description)
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

func getTermAssembly(cmd shireikan.Command, term string) string {
	if strings.HasPrefix(term, "/") {
		return term[1:]
	}
	return cmd.GetDomainName() + "." + term
}
