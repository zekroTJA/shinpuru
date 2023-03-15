package controllers

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	sharedmodels "github.com/zekroTJA/shinpuru/internal/models"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/services/webserver/v1/models"
	"github.com/zekroTJA/shinpuru/internal/util/snowflakenodes"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekrotja/dgrs"
	"github.com/zekrotja/sop"
)

type UnbanrequestsController struct {
	session *discordgo.Session
	db      database.Database
	pmw     *permissions.Permissions
	st      *dgrs.State
}

func (c *UnbanrequestsController) Setup(container di.Container, router fiber.Router) {
	c.session = container.Get(static.DiDiscordSession).(*discordgo.Session)
	c.db = container.Get(static.DiDatabase).(database.Database)
	c.pmw = container.Get(static.DiPermissions).(*permissions.Permissions)
	c.st = container.Get(static.DiState).(*dgrs.State)

	router.Get("", c.getUnbanrequests)
	router.Post("", c.postUnbanrequests)
	router.Get("/bannedguilds", c.getBannedGuilds)
}

// @Summary Get Unban Requests
// @Description Returns a list of unban requests created by the authenticated user.
// @Tags Unban Requests
// @Accept json
// @Produce json
// @Success 200 {array} models.RichUnbanRequest "Wrapped in models.ListResponse"
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /unbanrequests [get]
func (c *UnbanrequestsController) getUnbanrequests(ctx *fiber.Ctx) error {
	uid := ctx.Locals("uid").(string)

	requests, err := c.db.GetGuildUserUnbanRequests(uid, "")
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}
	if requests == nil {
		requests = make([]sharedmodels.UnbanRequest, 0)
	}

	self, err := c.st.User(uid)
	if err != nil {
		return err
	}

	res := sop.Map[sharedmodels.UnbanRequest](sop.Slice(requests),
		func(r sharedmodels.UnbanRequest, i int) models.RichUnbanRequest {
			r.Hydrate()
			rub := models.RichUnbanRequest{
				UnbanRequest: r,
				Creator:      models.FlatUserFromUser(self),
			}
			if proc, _ := c.st.User(rub.ProcessedBy); proc != nil {
				rub.Processor = models.FlatUserFromUser(proc)
			}
			return rub
		})

	return ctx.JSON(models.NewListResponse(res.Unwrap()))
}

// @Summary Create Unban Requests
// @Description Create an unban reuqest.
// @Tags Unban Requests
// @Accept json
// @Produce json
// @Param payload body sharedmodels.UnbanRequest true "The unban request payload."
// @Success 200 {object} models.RichUnbanRequest
// @Failure 400 {object} models.Error
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /unbanrequests [post]
func (c *UnbanrequestsController) postUnbanrequests(ctx *fiber.Ctx) error {
	uid := ctx.Locals("uid").(string)

	user, err := c.session.User(uid)
	if err != nil {
		return err
	}

	req := new(sharedmodels.UnbanRequest)
	if err := ctx.BodyParser(req); err != nil {
		return err
	}
	if err := req.Validate(); err != nil {
		return err
	}

	applicableReps, err := c.getUserApplicableReports(uid)
	if err != nil {
		return err
	}

	rep, i := sop.Slice(applicableReps).First(func(v sharedmodels.Report, i int) bool {
		return v.GuildID == req.GuildID
	})
	if i == -1 {
		return fiber.NewError(fiber.StatusBadRequest, "you are not able to create an unban request for this guild")
	}

	requests, err := c.db.GetGuildUserUnbanRequests(uid, req.GuildID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	if requests != nil {
		for _, r := range requests {
			if r.Status == sharedmodels.UnbanRequestStatePending {
				return fiber.NewError(fiber.StatusBadRequest, "there is still one open unban request to be proceed")
			}
		}
	}

	finalReq := sharedmodels.UnbanRequest{
		ID:       snowflakenodes.NodeUnbanRequests.Generate(),
		UserID:   uid,
		GuildID:  req.GuildID,
		UserTag:  user.String(),
		Message:  req.Message,
		Status:   sharedmodels.UnbanRequestStatePending,
		ReportID: rep.ID,
	}

	if err := c.db.AddUnbanRequest(finalReq); err != nil && !database.IsErrDatabaseNotFound(err) {
		return err
	}

	finalReq.Hydrate()

	return ctx.JSON(models.RichUnbanRequest{
		UnbanRequest: finalReq,
		Creator:      models.FlatUserFromUser(user),
	})
}

// @Summary Get Banned Guilds
// @Description Returns a list of guilds where the currently authenticated user is banned.
// @Tags Unban Requests
// @Accept json
// @Produce json
// @Success 200 {array} models.GuildReduced "Wrapped in models.ListResponse"
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /unbanrequests/bannedguilds [get]
func (c *UnbanrequestsController) getBannedGuilds(ctx *fiber.Ctx) error {
	uid := ctx.Locals("uid").(string)

	applicableReps, err := c.getUserApplicableReports(uid)
	if err != nil {
		return err
	}

	guilds := make(map[string]*models.GuildReduced)
	for _, rep := range applicableReps {
		if _, ok := guilds[rep.GuildID]; ok {
			continue
		}
		guild, err := c.st.Guild(rep.GuildID)
		if err != nil {
			return err
		}
		guilds[rep.GuildID] = models.GuildReducedFromGuild(guild)
	}

	guildsArr := make([]*models.GuildReduced, len(guilds))
	i := 0
	for _, g := range guilds {
		guildsArr[i] = g
		i++
	}

	return ctx.JSON(models.NewListResponse(guildsArr))
}

// --- HELPERS ------------

func (c *UnbanrequestsController) getUserApplicableReports(userID string) ([]sharedmodels.Report, error) {
	// Get all ban reports
	banReps, err := c.db.GetReportsFiltered(
		"", userID, sharedmodels.TypeBan, 0, 100000)
	if err != nil {
		if database.IsErrDatabaseNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	// Get all unban accepted reports
	unbanAcceptedReps, err := c.db.GetReportsFiltered(
		"", userID, sharedmodels.TypeUnban, 0, 100000)
	if err != nil {
		if database.IsErrDatabaseNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	// Get all unban rejected reports
	unbanRejectedReps, err := c.db.GetReportsFiltered(
		"", userID, sharedmodels.TypeUnbanRejected, 0, 100000)
	if err != nil {
		if database.IsErrDatabaseNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	applicableReps := make([]sharedmodels.Report, 0, len(banReps))
	for _, banRep := range banReps {
		_, i := sop.Slice(unbanAcceptedReps).First(func(v sharedmodels.Report, i int) bool {
			return v.GuildID == banRep.GuildID &&
				time.UnixMilli(banRep.ID.Time()).Before(time.UnixMilli(v.ID.Time()))
		})
		if i != -1 {
			continue
		}
		_, i = sop.Slice(unbanRejectedReps).First(func(v sharedmodels.Report, i int) bool {
			return v.GuildID == banRep.GuildID &&
				time.UnixMilli(banRep.ID.Time()).Before(time.UnixMilli(v.ID.Time())) &&
				time.UnixMilli(v.ID.Time()).
					Add(sharedmodels.UnbanRequestCooldown).
					After(time.Now())
		})
		if i != -1 {
			continue
		}
		applicableReps = append(applicableReps, banRep)
	}

	return applicableReps, nil
}
