package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v8"
	"github.com/sarulabs/di/v2"
	"github.com/sirupsen/logrus"
	"github.com/zekroTJA/shinpuru/internal/inits"
	"github.com/zekroTJA/shinpuru/internal/middleware"
	"github.com/zekroTJA/shinpuru/internal/models"
	"github.com/zekroTJA/shinpuru/internal/services/config"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/util/embedded"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekrotja/dgrs"
	"github.com/zekrotja/ken"
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

type DummyConfigProvider struct {
	cfg models.Config
}

var _ config.Provider = (*DummyConfigProvider)(nil)

func (cp *DummyConfigProvider) Parse() error {
	return nil
}

func (cp *DummyConfigProvider) Config() *models.Config {
	return &cp.cfg
}

func main() {
	flag.Parse()

	// Setting Release flag to true manually to prevent
	// registration of test command and exclude it in the
	// command manual.
	embedded.Release = "TRUE"

	diBuilder, _ := di.NewBuilder()

	config := &DummyConfigProvider{}
	diBuilder.Set(static.DiConfig, config)

	s, _ := discordgo.New()
	diBuilder.Set(static.DiDiscordSession, s)

	// Initialize command handler
	diBuilder.Add(di.Def{
		Name: static.DiCommandHandler,
		Build: func(ctn di.Container) (interface{}, error) {
			return inits.InitCommandHandler(ctn)
		},
	})

	diBuilder.Add(di.Def{
		Name: static.DiState,
		Build: func(ctn di.Container) (interface{}, error) {
			return (*dgrs.State)(nil), nil
		},
	})

	diBuilder.Add(di.Def{
		Name: static.DiRedis,
		Build: func(ctn di.Container) (interface{}, error) {
			return (*redis.Client)(nil), nil
		},
	})

	diBuilder.Set(static.DiPermissions, &permissions.Permissions{})
	diBuilder.Set(static.DiGhostpingIgnoreMiddleware, &middleware.GhostPingIgnoreMiddleware{})

	ctn := diBuilder.Build()

	cmdHandler := ctn.Get(static.DiCommandHandler).(*ken.Ken)
	if err := exportCommandManual(cmdHandler, *flagExportFile); err != nil {
		logrus.WithError(err).Fatal("Failed exporting command manual")
	}
	logrus.Info("Successfully exported command manual file to " + *flagExportFile)
}

// exportCommandManual generates a markdown text file
// from the registered command instances of the passed
// command handler. The file is then exported to the
// given file location. Occuring errors are returned.
func exportCommandManual(h *ken.Ken, fileName string) (err error) {
	var document strings.Builder

	fmt.Fprintf(&document, "> Auto generated command manual | %s\n\n", time.Now().Format(time.RFC1123))

	cmdInfo := h.GetCommandInfo()
	groups := groupCommandInfo(cmdInfo)

	writeCommandList(&document, groups)
	writeCommandDetailsList(&document, groups)

	filePath := path.Dir(fileName)
	if filePath != "." && filePath != "/" {
		if err := os.MkdirAll(filePath, 0744); err != nil {
			return err
		}
	}

	return ioutil.WriteFile(fileName, []byte(document.String()), 0644)
}

func groupCommandInfo(cmdInfo ken.CommandInfoList) (m map[string][]*ken.CommandInfo) {
	m = make(map[string][]*ken.CommandInfo)
	for _, info := range cmdInfo {
		groupName := groupNameFromDomain(info.Implementations["Domain"][0].(string))
		m[groupName] = append(m[groupName], info)
	}
	return
}

func groupNameFromDomain(domain string) (group string) {
	domainSplit := strings.Split(domain, ".")
	group = strings.Join(domainSplit[1:len(domainSplit)-1], " ")
	group = strings.ToUpper(group)
	return
}

func writeCommandList(document *strings.Builder, groups map[string][]*ken.CommandInfo) {
	document.WriteString("# Command List\n\n")

	for group, cmds := range groups {
		fmt.Fprintf(document, "## %s\n\n", group)
		for _, cmd := range cmds {
			fmt.Fprintf(document, "- [%s](#%s)\n", cmd.ApplicationCommand.Name, cmd.ApplicationCommand.Name)
		}
		document.WriteString("\n")
	}
}

func writeCommandDetailsList(document *strings.Builder, groups map[string][]*ken.CommandInfo) {
	document.WriteString("# Command Details\n\n")

	for group, cmds := range groups {
		fmt.Fprintf(document, "## %s\n\n", group)
		for _, cmd := range cmds {
			writeCommandDetails(document, cmd)
		}
		document.WriteString("\n---\n")
	}
}

func writeCommandDetails(document *strings.Builder, cmd *ken.CommandInfo) {
	var dmCapable bool
	if v, ok := cmd.Implementations["IsDmCapable"]; ok {
		dmCapable, _ = v[0].(bool)
	}

	domain := cmd.Implementations["Domain"][0].(string)

	fmt.Fprintf(document, "### %s\n\n", cmd.ApplicationCommand.Name)
	fmt.Fprintf(document, "%s\n\n", cmd.ApplicationCommand.Description)
	fmt.Fprintf(document,
		"| | |\n"+
			"|--|--|\n"+
			"| Domain Name | %s |\n"+
			"| Version | %s |\n"+
			"| DM Capable | %t |\n"+
			"\n\n",
		domain,
		cmd.ApplicationCommand.Version,
		dmCapable,
	)

	subDNSI := cmd.Implementations["SubDomains"]
	if len(subDNSI) != 0 {
		subDNS, ok := subDNSI[0].([]permissions.SubPermission)
		if ok && len(subDNS) != 0 {
			document.WriteString("#### Sub Permission Rules\n\n")
			for _, perm := range subDNS {
				explicit := ""
				if perm.Explicit {
					explicit = "`[EXPLICIT]` "
				}
				fmt.Fprintf(document, "- **`%s`** %s- %s\n",
					getTermAssembly(domain, perm.Term), explicit, perm.Description)
			}
			document.WriteString("\n\n")
		}
	}

	if len(cmd.ApplicationCommand.Options) == 0 {
		return
	}

	if cmd.ApplicationCommand.Options[0].Type == discordgo.ApplicationCommandOptionSubCommand {
		document.WriteString("#### Sub Commands\n\n")
		for _, sub := range cmd.ApplicationCommand.Options {
			writeSubCommand(document, sub)
		}
	} else {
		document.WriteString("#### Arguments\n\n")
		writeArguments(document, cmd.ApplicationCommand.Options)
	}
}

func writeSubCommand(document *strings.Builder, sub *discordgo.ApplicationCommandOption) {
	fmt.Fprintf(document, "##### `%s`\n", sub.Name)
	fmt.Fprintf(document, "%s\n", sub.Description)
	if len(sub.Options) > 0 {
		document.WriteString("**Arguments**\n")
		writeArguments(document, sub.Options)
	}
	document.WriteString("\n")
}

func writeArguments(document *strings.Builder, options []*discordgo.ApplicationCommandOption) {
	fmt.Fprintf(document,
		"| Name | Type | Required | Description | Choises |\n"+
			"|------|------|----------|-------------|---------|\n")
	for _, opt := range options {
		fmt.Fprintf(document,
			"| %s | `%s` | `%t` | %s | %s |",
			opt.Name,
			opt.Type.String(),
			opt.Required,
			opt.Description,
			formatChoises(opt.Choices),
		)
	}
}

func formatChoises(choises []*discordgo.ApplicationCommandOptionChoice) string {
	var sb strings.Builder
	for _, ch := range choises {
		fmt.Fprintf(&sb, "- `%s` (`%v`)</br>", ch.Name, ch.Value)
	}
	return sb.String()
}

func getTermAssembly(domain, term string) string {
	if strings.HasPrefix(term, "/") {
		return term[1:]
	}
	return domain + "." + term
}
