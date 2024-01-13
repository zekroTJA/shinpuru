package controllers

import (
	"context"
	"fmt"
	"runtime"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/models"
	"github.com/zekroTJA/shinpuru/internal/services/config"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/storage"
	"github.com/zekroTJA/shinpuru/internal/services/webserver/auth"
	apiModels "github.com/zekroTJA/shinpuru/internal/services/webserver/v1/models"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/embedded"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
	"github.com/zekrotja/dgrs"
	"github.com/zekrotja/ken"
)

type EtcController struct {
	session    *discordgo.Session
	cfg        config.Provider
	authMw     auth.Middleware
	st         *dgrs.State
	storage    storage.Storage
	db         database.Database
	cmdHandler *ken.Ken
	rd         *redis.Client
}

func (c *EtcController) Setup(container di.Container, router fiber.Router) {
	c.session = container.Get(static.DiDiscordSession).(*discordgo.Session)
	c.cfg = container.Get(static.DiConfig).(config.Provider)
	c.authMw = container.Get(static.DiAuthMiddleware).(auth.Middleware)
	c.st = container.Get(static.DiState).(*dgrs.State)
	c.storage = container.Get(static.DiObjectStorage).(storage.Storage)
	c.db = container.Get(static.DiDatabase).(database.Database)
	c.cmdHandler = container.Get(static.DiCommandHandler).(*ken.Ken)
	c.rd = container.Get(static.DiRedis).(*redis.Client)

	router.Get("/me", c.authMw.Handle, c.getMe)
	router.Get("/sysinfo", c.getSysinfo)
	router.Get("/privacyinfo", c.getPrivacyinfo)
	router.Get("/allpermissions", c.getAllPermissions)
	router.Get("/healthcheck", c.getHealthcheck)
}

// @Summary Me
// @Description Returns the user object of the currently authenticated user.
// @Tags Etc
// @Accept json
// @Produce json
// @Success 200 {object} apiModels.User
// @Router /me [get]
func (c *EtcController) getMe(ctx *fiber.Ctx) error {
	uid := ctx.Locals("uid").(string)

	user, err := c.st.User(uid)
	if err != nil {
		return err
	}

	caapchaVerified, err := c.db.GetUserVerified(uid)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	created, _ := discordutil.GetDiscordSnowflakeCreationTime(user.ID)

	res := &apiModels.User{
		User:            user,
		AvatarURL:       user.AvatarURL(""),
		CreatedAt:       created,
		BotOwner:        uid == c.cfg.Config().Discord.OwnerID,
		CaptchaVerified: caapchaVerified,
	}

	return ctx.JSON(res)
}

// @Summary System Information
// @Description Returns general global system information.
// @Tags Etc
// @Accept json
// @Produce json
// @Success 200 {object} apiModels.SystemInfo
// @Router /sysinfo [get]
func (c *EtcController) getSysinfo(ctx *fiber.Ctx) error {
	buildTS, _ := strconv.Atoi(embedded.AppDate)
	buildDate := time.Unix(int64(buildTS), 0)

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	uptime := int64(time.Since(util.StatsStartupTime).Seconds())

	self, err := c.st.SelfUser()
	if err != nil {
		return err
	}

	guilds, err := c.st.Guilds()
	if err != nil {
		return err
	}

	res := &apiModels.SystemInfo{
		Version:    embedded.AppVersion,
		CommitHash: embedded.AppCommit,
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

		BotUserID: self.ID,
		BotInvite: util.GetInviteLink(self.ID),

		Guilds: len(guilds),
	}

	return ctx.JSON(res)
}

// @Summary Privacy Information
// @Description Returns general global privacy information.
// @Tags Etc
// @Accept json
// @Produce json
// @Success 200 {object} models.Privacy
// @Router /privacyinfo [get]
func (c *EtcController) getPrivacyinfo(ctx *fiber.Ctx) error {
	return ctx.JSON(c.cfg.Config().Privacy)
}

// @Summary All Permissions
// @Description Return a list of all available permissions.
// @Tags Etc
// @Accept json
// @Produce json
// @Success 200 {array} string "Wrapped in models.ListResponse"
// @Router /allpermissions [get]
func (c *EtcController) getAllPermissions(ctx *fiber.Ctx) error {
	all := util.GetAllPermissions(c.cmdHandler)
	return ctx.JSON(apiModels.NewListResponse(all.Unwrap()))
}

// @Summary Healthcheck
// @Description General system healthcheck.
// @Tags Etc
// @Accept json
// @Produce json
// @Success 200 {array} string "Wrapped in models.ListResponse"
// @Router /healthcheck [get]
func (c *EtcController) getHealthcheck(ctx *fiber.Ctx) error {
	var hc models.HealthcheckResponse

	hc.Database = models.HealthcheckStatusFromError(c.db.Status())
	hc.Storage = models.HealthcheckStatusFromError(c.storage.Status())
	hc.Redis = models.HealthcheckStatusFromError(c.rd.Ping(context.Background()).Err())

	hc.Discord.Ok = atomic.LoadInt32(&util.ConnectedState) == 1
	if !hc.Discord.Ok {
		hc.Discord.Message = "gateway connection has been disconnected"
	}

	hc.AllOk = hc.Database.Ok && hc.Storage.Ok && hc.Redis.Ok && hc.Discord.Ok

	return ctx.JSON(hc)
}
