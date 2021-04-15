package controllers

import (
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/core/config"
	"github.com/zekroTJA/shinpuru/internal/core/webserver/auth"
	"github.com/zekroTJA/shinpuru/internal/core/webserver/v1/models"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
	"github.com/zekroTJA/shireikan"
)

type EtcController struct {
	session    *discordgo.Session
	cfg        *config.Config
	cmdHandler shireikan.Handler
	authMw     auth.Middleware
}

func (c *EtcController) Setup(container di.Container, router fiber.Router) {
	c.session = container.Get(static.DiDiscordSession).(*discordgo.Session)
	c.cfg = container.Get(static.DiConfig).(*config.Config)
	c.cmdHandler = container.Get(static.DiCommandHandler).(shireikan.Handler)
	c.authMw = container.Get(static.DiAuthMiddleware).(auth.Middleware)

	router.Get("/me", c.authMw.Handle, c.getMe)
	router.Get("/sysinfo", c.getSysinfo)
	router.Get("/commands", c.getCommands)
}

func (c *EtcController) getMe(ctx *fiber.Ctx) error {
	uid := ctx.Locals("uid").(string)

	user, err := c.session.User(uid)
	if err != nil {
		return err
	}

	created, _ := discordutil.GetDiscordSnowflakeCreationTime(user.ID)

	res := &models.User{
		User:      user,
		AvatarURL: user.AvatarURL(""),
		CreatedAt: created,
		BotOwner:  uid == c.cfg.Discord.OwnerID,
	}

	return ctx.JSON(res)
}

func (c *EtcController) getSysinfo(ctx *fiber.Ctx) error {
	buildTS, _ := strconv.Atoi(util.AppDate)
	buildDate := time.Unix(int64(buildTS), 0)

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	uptime := int64(time.Since(util.StatsStartupTime).Seconds())

	res := &models.SystemInfo{
		Version:    util.AppVersion,
		CommitHash: util.AppCommit,
		BuildDate:  buildDate,
		GoVersion:  runtime.Version(),

		Uptime:    uptime,
		UptimeStr: fmt.Sprintf("%d", uptime),

		OS:          runtime.GOOS,
		Arch:        runtime.GOARCH,
		CPUs:        runtime.NumCPU(),
		GoRoutines:  runtime.NumGoroutine(),
		StackUse:    memStats.StackInuse,
		StackUseStr: fmt.Sprintf("%d", memStats.StackInuse),
		HeapUse:     memStats.HeapInuse,
		HeapUseStr:  fmt.Sprintf("%d", memStats.HeapInuse),

		BotUserID: c.session.State.User.ID,
		BotInvite: util.GetInviteLink(c.session),

		Guilds: len(c.session.State.Guilds),
	}

	return ctx.JSON(res)
}

func (c *EtcController) getCommands(ctx *fiber.Ctx) error {
	cmdInstances := c.cmdHandler.GetCommandInstances()
	cmdInfos := make([]*models.CommandInfo, len(cmdInstances))

	for i, c := range cmdInstances {
		cmdInfo := models.GetCommandInfoFromCommand(c)
		cmdInfos[i] = cmdInfo
	}

	list := &models.ListResponse{N: len(cmdInfos), Data: cmdInfos}

	return ctx.JSON(list)
}
