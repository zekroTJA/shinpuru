package controllers

import (
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/config"
	"github.com/zekroTJA/shinpuru/internal/services/webserver/auth"
	"github.com/zekroTJA/shinpuru/internal/services/webserver/v1/models"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
)

type EtcController struct {
	session *discordgo.Session
	cfg     *config.Config
	authMw  auth.Middleware
}

func (c *EtcController) Setup(container di.Container, router fiber.Router) {
	c.session = container.Get(static.DiDiscordSession).(*discordgo.Session)
	c.cfg = container.Get(static.DiConfig).(*config.Config)
	c.authMw = container.Get(static.DiAuthMiddleware).(auth.Middleware)

	router.Get("/me", c.authMw.Handle, c.getMe)
	router.Get("/sysinfo", c.getSysinfo)
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
