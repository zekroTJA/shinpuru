package controllers

import (
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/services/config"
	"github.com/zekroTJA/shinpuru/internal/services/webserver/v1/models"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/colors"
	"github.com/zekroTJA/shinpuru/pkg/etag"
	"github.com/zekroTJA/shireikan"
	"github.com/zekrotja/dgrs"
	"github.com/zekrotja/ken"
)

type UtilController struct {
	session          *discordgo.Session
	cfg              config.Provider
	legacyCmdHandler shireikan.Handler
	cmdHandler       *ken.Ken
	st               *dgrs.State
}

func (c *UtilController) Setup(container di.Container, router fiber.Router) {
	c.session = container.Get(static.DiDiscordSession).(*discordgo.Session)
	c.cfg = container.Get(static.DiConfig).(config.Provider)
	c.legacyCmdHandler = container.Get(static.DiLegacyCommandHandler).(shireikan.Handler)
	c.cmdHandler = container.Get(static.DiCommandHandler).(*ken.Ken)
	c.st = container.Get(static.DiState).(*dgrs.State)

	router.Get("/landingpageinfo", c.getLandingPageInfo)
	router.Get("/color/:hexcode", c.getColor)
	router.Get("/commands", c.getCommands)
	router.Get("/slashcommands", c.getSlashCommands)
	router.Get("/updateinfo", c.getUpdateInfo)
}

// @Summary Landing Page Info
// @Description Returns general information for the landing page like the local invite parameters.
// @Tags Utilities
// @Accept json
// @Produce json
// @Success 200 {object} models.LandingPageResponse
// @Router /util/landingpageinfo [get]
func (c *UtilController) getLandingPageInfo(ctx *fiber.Ctx) error {
	res := new(models.LandingPageResponse)

	publicInvites := c.cfg.Config().WebServer.LandingPage.ShowPublicInvites
	localInvite := c.cfg.Config().WebServer.LandingPage.ShowLocalInvite

	if publicInvites {
		res.PublicCanaryInvite = static.PublicCanaryInvite
		res.PublicMainInvite = static.PublicMainInvite
	}

	if localInvite {
		self, err := c.st.SelfUser()
		if err != nil {
			return err
		}
		res.LocalInvite = util.GetInviteLink(self.ID)
	}

	return ctx.JSON(res)
}

// @Summary Color Generator
// @Description Produces a square image of the given color and size.
// @Param hexcode path string true "Hex Code of the Color to produce"
// @Param size query int false "The dimension of the square image" default(24)
// @Tags Utilities
// @Accept json
// @Produce image/png
// @Success 200 {file} png image data
// @Router /util/color/{hexcode} [get]
func (c *UtilController) getColor(ctx *fiber.Ctx) error {
	hexcode := ctx.Params("hexcode")
	size := strings.ToLower(ctx.Query("size"))

	var xSize, ySize int
	var err error

	if size == "" {
		xSize, ySize = 24, 24
	} else if strings.Contains(size, "x") {
		split := strings.Split(size, "x")
		if len(split) != 2 {
			return fiber.NewError(fiber.StatusBadRequest, "invalid size parameter; must provide two size dimensions")
		}
		if xSize, err = strconv.Atoi(split[0]); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		if ySize, err = strconv.Atoi(split[1]); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	} else {
		if xSize, err = strconv.Atoi(size); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		ySize = xSize
	}

	if xSize < 1 || ySize < 1 || xSize > 5000 || ySize > 5000 {
		return fiber.NewError(fiber.StatusBadRequest, "invalid size parameter; value must be in range [1..5000]")
	}

	clr, err := colors.FromHex(hexcode)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	buff, err := colors.CreateImage(clr, xSize, ySize)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	data := buff.Bytes()

	etag := etag.Generate(data, false)

	ctx.Context().SetContentType("image/png")
	// 365 days browser caching
	ctx.Set("Cache-Control", "public, max-age=31536000, immutable")
	ctx.Set("ETag", etag)
	return ctx.Send(data)
}

// @Summary Command List
// @Description Returns a list of registered commands and their description.
// @Tags Utilities
// @Accept json
// @Produce json
// @Success 200 {array} models.CommandInfo "Wrapped in models.ListResponse"
// @Router /util/commands [get]
func (c *UtilController) getCommands(ctx *fiber.Ctx) error {
	cmdInstances := c.legacyCmdHandler.GetCommandInstances()
	cmdInfos := make([]*models.CommandInfo, len(cmdInstances))

	for i, c := range cmdInstances {
		cmdInfo := models.GetCommandInfoFromCommand(c)
		cmdInfos[i] = cmdInfo
	}

	list := &models.ListResponse{N: len(cmdInfos), Data: cmdInfos}

	return ctx.JSON(list)
}

// @Summary Slash Command List
// @Description Returns a list of registered slash commands and their description.
// @Tags Utilities
// @Accept json
// @Produce json
// @Success 200 {array} models.SlashCommandInfo "Wrapped in models.ListResponse"
// @Router /util/slashcommands [get]
func (c *UtilController) getSlashCommands(ctx *fiber.Ctx) error {
	cmdInfo := c.cmdHandler.GetCommandInfo()
	res := make([]*models.SlashCommandInfo, len(cmdInfo))

	for i, ci := range cmdInfo {
		res[i] = models.GetSlashCommandInfoFromCommand(ci)
	}

	return ctx.JSON(&models.ListResponse{N: len(res), Data: res})
}

// @Summary Update Information
// @Description Returns update information.
// @Tags Utilities
// @Accept json
// @Produce json
// @Success 200 {object} models.UpdateInfoResponse "Update info response"
// @Router /util/updateinfo [get]
func (c *UtilController) getUpdateInfo(ctx *fiber.Ctx) error {
	var res models.UpdateInfoResponse
	res.IsOld, res.Current, res.Latest = util.CheckForUpdate()
	res.CurrentStr = res.Current.String()
	res.LatestStr = res.Latest.String()

	return ctx.JSON(res)
}
