package inits

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v8"
	"github.com/sarulabs/di/v2"
	"github.com/sirupsen/logrus"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/slashcommands"
	"github.com/zekroTJA/shinpuru/internal/util/embedded"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/rediscmdstore"
	"github.com/zekrotja/dgrs"
	"github.com/zekrotja/ken"
	"github.com/zekrotja/ken/middlewares/cmdhelp"
	"github.com/zekrotja/ken/state"
	"github.com/zekrotja/ken/store"
)

func InitCommandHandler(container di.Container) (k *ken.Ken, err error) {
	session := container.Get(static.DiDiscordSession).(*discordgo.Session)
	st := container.Get(static.DiState).(*dgrs.State)
	perms := container.Get(static.DiPermissions).(*permissions.Permissions)
	rd, _ := container.Get(static.DiRedis).(*redis.Client)

	var cmdStore store.CommandStore
	if rd != nil {
		cmdStore = rediscmdstore.New(rd, fmt.Sprintf("snp:cmdstore:%s", embedded.AppCommit))
	}

	k, err = ken.New(session, ken.Options{
		State:              state.NewDgrs(st),
		CommandStore:       cmdStore,
		DependencyProvider: container,
		OnSystemError:      systemErrorHandler,
		OnCommandError:     commandErrorHandler,
		EmbedColors: ken.EmbedColors{
			Default: static.ColorEmbedDefault,
			Error:   static.ColorEmbedError,
		},
	})
	if err != nil {
		return
	}

	err = k.RegisterCommands(
		new(slashcommands.Autorole),
		new(slashcommands.Backup),
		new(slashcommands.Bug),
		new(slashcommands.Clear),
		new(slashcommands.Vote),
		new(slashcommands.Report),
		new(slashcommands.Mute),
		new(slashcommands.Perms),
		new(slashcommands.Chanstats),
		new(slashcommands.Exec),
		new(slashcommands.Say),
		new(slashcommands.Notify),
		new(slashcommands.Mvall),
		new(slashcommands.Lock),
		new(slashcommands.Inviteblock),
		new(slashcommands.Ghostping),
		new(slashcommands.Voicelog),
		new(slashcommands.Modlog),
		new(slashcommands.Announcements),
		new(slashcommands.Starboard),
		new(slashcommands.Colorreation),
		new(slashcommands.User),
		new(slashcommands.Twitchnotify),
		new(slashcommands.Tag),
		new(slashcommands.Presence),
		new(slashcommands.Login),
		new(slashcommands.Quote),
		new(slashcommands.Stats),
	)
	if err != nil {
		return
	}

	err = k.RegisterMiddlewares(
		perms,
		cmdhelp.New("help"),
	)

	return
}

func systemErrorHandler(context string, err error, args ...interface{}) {
	logrus.WithField("ctx", context).WithError(err).Error("ken error")
}

func commandErrorHandler(err error, ctx *ken.Ctx) {
	// Is ignored if interaction has already been responded
	ctx.Defer()

	if err == ken.ErrNotDMCapable {
		ctx.FollowUpError("This command can not be used in DMs.", "")
		return
	}

	ctx.FollowUpError(
		fmt.Sprintf("The command execution failed unexpectedly:\n```\n%s\n```", err.Error()),
		"Command execution failed")
}
