package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/models"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/colors"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
	"github.com/zekroTJA/shinpuru/pkg/embedbuilder"
	"github.com/zekroTJA/shireikan"
	"github.com/zekrotja/dgrs"
)

type CmdGuild struct {
}

func (c *CmdGuild) GetInvokes() []string {
	return []string{"guild", "guildinfo", "g", "gi"}
}

func (c *CmdGuild) GetDescription() string {
	return "Outputs info about the current guild."
}

func (c *CmdGuild) GetHelp() string {
	return "`guild` - Prints guild info"
}

func (c *CmdGuild) GetGroup() string {
	return shireikan.GroupChat
}

func (c *CmdGuild) GetDomainName() string {
	return "sp.chat.guild"
}

func (c *CmdGuild) GetSubPermissionRules() []shireikan.SubPermission {
	return nil
}

func (c *CmdGuild) IsExecutableInDMChannels() bool {
	return false
}

func (c *CmdGuild) Exec(ctx shireikan.Context) (err error) {
	const maxGuildRoles = 30

	g := ctx.GetGuild()

	clr, err := colors.GetVibrantColorFromImageUrl(g.IconURL())
	if err != nil {
		clr = static.ColorEmbedDefault
	}

	st := ctx.GetObject(static.DiState).(*dgrs.State)
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

	db, _ := ctx.GetObject(static.DiDatabase).(database.Database)

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
		WithFooter(fmt.Sprintf("issued by %s", ctx.GetUser().String()), "", "").
		Build()

	_, err = ctx.GetSession().ChannelMessageSendEmbed(ctx.GetChannel().ID, emb)
	return
}

func (c *CmdGuild) wrapBool(b bool) string {
	if b {
		return ":white_check_mark:"
	}
	return ":x:"
}
