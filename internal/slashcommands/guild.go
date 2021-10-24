package slashcommands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/models"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/colors"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
	"github.com/zekroTJA/shinpuru/pkg/embedbuilder"
	"github.com/zekrotja/dgrs"
	"github.com/zekrotja/ken"
)

type Guild struct{}

var (
	_ ken.Command             = (*Guild)(nil)
	_ permissions.PermCommand = (*Guild)(nil)
)

func (c *Guild) Name() string {
	return "guild"
}

func (c *Guild) Description() string {
	return "Displays information about the current guild."
}

func (c *Guild) Version() string {
	return "1.0.0"
}

func (c *Guild) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *Guild) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{}
}

func (c *Guild) Domain() string {
	return "sp.chat.guild"
}

func (c *Guild) SubDomains() []permissions.SubPermission {
	return nil
}

func (c *Guild) Run(ctx *ken.Ctx) (err error) {
	if err = ctx.Defer(); err != nil {
		return
	}

	const maxGuildRoles = 30

	st := ctx.Get(static.DiState).(*dgrs.State)
	db := ctx.Get(static.DiDatabase).(database.Database)

	g, err := st.Guild(ctx.Event.GuildID)
	if err != nil {
		return
	}

	clr, err := colors.GetVibrantColorFromImageUrl(g.IconURL())
	if err != nil {
		clr = static.ColorEmbedDefault
	}

	gChans, err := st.Channels(g.ID, true)

	nTextChans, nVoiceChans, nCategoryChans := 0, 0, 0
	for _, c := range gChans {
		switch c.Type {
		case discordgo.ChannelTypeGuildCategory:
			nCategoryChans++
		case discordgo.ChannelTypeGuildVoice:
			nVoiceChans++
		default:
			nTextChans++
		}
	}
	chans := fmt.Sprintf("Category Channels: `%d`\nText Channels: `%d`\nVoice Channels: `%d`",
		nCategoryChans, nTextChans, nVoiceChans)

	lenRoles := len(g.Roles) - 1
	if lenRoles > maxGuildRoles {
		lenRoles = maxGuildRoles + 1
	}
	roles := make([]string, lenRoles)
	i := 0
	for _, r := range g.Roles {
		if r.ID == g.ID {
			continue
		}
		if i == maxGuildRoles {
			roles[i] = "..."
			break
		}
		roles[i] = r.Mention()
		i++
	}

	createdTime, err := discordutil.GetDiscordSnowflakeCreationTime(g.ID)
	if err != nil {
		return
	}
	if err != nil {
		return
	}

	totalReportCount := 0
	reportCounts := make([]string, len(models.ReportTypes))
	for i, typ := range models.ReportTypes {
		c, err := db.GetReportsFilteredCount(g.ID, "", i)
		if err != nil {
			return err
		}
		reportCounts[i] = fmt.Sprintf("%s: `%d`", typ, c)
		totalReportCount += c
	}

	prefix, err := db.GetGuildPrefix(g.ID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return
	}
	if prefix == "" {
		prefix = "unset"
	} else {
		prefix = fmt.Sprintf("`%s`", prefix)
	}

	backupsEnabled, err := db.GetGuildBackup(g.ID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return
	}
	inviteBlockEnabled, err := db.GetGuildInviteBlock(g.ID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return
	}
	antiraidEnabled, err := db.GetAntiraidState(g.ID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return
	}
	guildSecurity := fmt.Sprintf("%s Backups Enabled\n%s Inviteblock Enabled\n%s Antiraid Enabled\n",
		c.wrapBool(backupsEnabled), c.wrapBool(inviteBlockEnabled != ""), c.wrapBool(antiraidEnabled))

	emb := embedbuilder.New().
		WithThumbnail(g.IconURL(), "", 100, 100).
		WithColor(clr).
		AddField("Name", g.Name).
		AddField("ID", fmt.Sprintf("```\n%s\n```", g.ID)).
		AddField("Created", createdTime.Format(time.RFC1123)).
		AddField("Guild Prefix", prefix).
		AddField("Owner", fmt.Sprintf("<@%s>", g.OwnerID)).
		AddField(fmt.Sprintf("Channels (%d)", len(gChans)), chans).
		AddField("Server Region", g.Region).
		AddField("Member Count", fmt.Sprintf("State: %d / Approx.: %d", g.MemberCount, g.ApproximateMemberCount)).
		AddField(fmt.Sprintf("Reports (%d)", totalReportCount), strings.Join(reportCounts, "\n")).
		AddField("Guild Security", guildSecurity).
		AddField(fmt.Sprintf("Roles (%d)", len(g.Roles)-1), strings.Join(roles, ", ")).
		WithFooter(fmt.Sprintf("issued by %s", ctx.User().String()), "", "").
		Build()

	return ctx.FollowUpEmbed(emb).Error
}

func (c *Guild) wrapBool(b bool) string {
	if b {
		return ":white_check_mark:"
	}
	return ":x:"
}
