package webserver

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"

	"github.com/zekroTJA/shinpuru/internal/commands"
	"github.com/zekroTJA/shinpuru/internal/core/permissions"
	"github.com/zekroTJA/shinpuru/internal/util/imgstore"
	"github.com/zekroTJA/shinpuru/internal/util/report"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
)

// ListResponse wraps a list response object
// with the list as Data and N as len(Data).
type ListResponse struct {
	N    int         `json:"n"`
	Data interface{} `json:"data"`
}

// User extends a discordgo.User as reponse
// model.
type User struct {
	*discordgo.User

	AvatarURL string    `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
	BotOwner  bool      `json:"bot_owner"`
}

// Member extends a discordgo.Member as
// response model.
type Member struct {
	*discordgo.Member

	AvatarURL string    `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
	Dominance int       `json:"dominance"`
}

// Guild extends a discordgo.Guild as
// response model.
type Guild struct {
	*discordgo.Guild

	SelfMember *Member   `json:"self_member"`
	IconURL    string    `json:"icon_url"`
	Members    []*Member `json:"members"`
}

// GuildReduced is a Guild model with fewer
// details than Guild model.
type GuildReduced struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Icon        string              `json:"icon"`
	IconURL     string              `json:"icon_url"`
	Region      string              `json:"region"`
	OwnerID     string              `json:"owner_id"`
	JoinedAt    discordgo.Timestamp `json:"joined_at"`
	MemberCount int                 `json:"member_count"`
}

// PermissionsResponse wraps a
// permissions.PermissionsArra as response
// model.
type PermissionsResponse struct {
	Permissions permissions.PermissionArray `json:"permissions"`
}

// Report extends report.Report by TypeName
// and Created time.
type Report struct {
	*report.Report

	TypeName string    `json:"type_name"`
	Created  time.Time `json:"created"`
}

// GuildSettings is the response model for
// guil dsettings and preferences.
type GuildSettings struct {
	Prefix              string                                 `json:"prefix"`
	Perms               map[string]permissions.PermissionArray `json:"perms"`
	AutoRole            string                                 `json:"autorole"`
	ModLogChannel       string                                 `json:"modlogchannel"`
	VoiceLogChannel     string                                 `json:"voicelogchannel"`
	JoinMessageChannel  string                                 `json:"joinmessagechannel"`
	JoinMessageText     string                                 `json:"joinmessagetext"`
	LeaveMessageChannel string                                 `json:"leavemessagechannel"`
	LeaveMessageText    string                                 `json:"leavemessagetext"`
}

// PermissionsUpdate is the request model to
// update a permissions array.
type PermissionsUpdate struct {
	Perm    string   `json:"perm"`
	RoleIDs []string `json:"role_ids"`
}

// ReasonRequest is a request model wrapping a
// Reason and Attachment URL.
type ReasonRequest struct {
	Reason     string `json:"reason"`
	Attachment string `json:"attachment"`
}

// ReportRequest extends ReasonRequest by
// Type of report.
type ReportRequest struct {
	*ReasonRequest

	Type int `json:"type"`
}

// InviteSettingsRequest is the request model
// for setting the global invite setting.
type InviteSettingsRequest struct {
	GuildID    string `json:"guild_id"`
	Messsage   string `json:"message"`
	InviteCode string `json:"invite_code"`
}

// InviteSettingsResponse is the response model
// sent back when setting the global invite setting.
type InviteSettingsResponse struct {
	Guild     *Guild `json:"guild"`
	InviteURL string `json:"invite_url"`
	Message   string `json:"message"`
}

// Count is a simple response wrapper for a
// count number.
type Count struct {
	Count int `json:"count"`
}

// SystemInfo is the response model for a
// system info request.
type SystemInfo struct {
	Version    string    `json:"version"`
	CommitHash string    `json:"commit_hash"`
	BuildDate  time.Time `json:"build_date"`
	GoVersion  string    `json:"go_version"`

	Uptime    int64  `json:"uptime"`
	UptimeStr string `json:"uptime_str"`

	OS          string `json:"os"`
	Arch        string `json:"arch"`
	CPUs        int    `json:"cpus"`
	GoRoutines  int    `json:"go_routines"`
	StackUse    uint64 `json:"stack_use"`
	StackUseStr string `json:"stack_use_str"`
	HeapUse     uint64 `json:"heap_use"`
	HeapUseStr  string `json:"heap_use_str"`

	BotUserID string `json:"bot_user_id"`
	BotInvite string `json:"bot_invite"`

	Guilds int `json:"guilds"`
}

// Validate returns true, when the ReasonRequest is valid.
// Otherwise, false is returned and an error response is
// returned.
func (req *ReasonRequest) Validate(ctx *routing.Context) (bool, error) {
	if len(req.Reason) < 3 {
		return false, jsonError(ctx, errInvalidArguments, fasthttp.StatusBadRequest)
	}

	if req.Attachment != "" && !imgstore.ImgUrlSRx.MatchString(req.Attachment) {
		return false, jsonError(ctx,
			fmt.Errorf("attachment must be a valid url to a file with type of png, jpg, jpeg, gif, ico, tiff, img, bmp or mp4."),
			fasthttp.StatusBadRequest)
	}

	return true, nil
}

// GuildFromGuild returns a Guild model from the passed
// discordgo.Guild g, discordgo.Member m and cmdHandler.
func GuildFromGuild(g *discordgo.Guild, m *discordgo.Member, cmdHandler *commands.CmdHandler) *Guild {
	if g == nil {
		return nil
	}

	membs := make([]*Member, len(g.Members))
	for i, m := range g.Members {
		membs[i] = MemberFromMember(m)
	}

	selfmm := MemberFromMember(m)

	if m != nil {
		switch {
		case discordutil.IsAdmin(g, m):
			selfmm.Dominance = 1
		case g.OwnerID == m.User.ID:
			selfmm.Dominance = 2
		case cmdHandler.IsBotOwner(m.User.ID):
			selfmm.Dominance = 3
		}
	}

	return &Guild{
		Guild:      g,
		SelfMember: selfmm,
		Members:    membs,
		IconURL:    getIconURL(g.ID, g.Icon),
	}
}

// GuildReducedFromGuild returns a GuildReduced from the passed
// discordgo.Guild g.
func GuildReducedFromGuild(g *discordgo.Guild) *GuildReduced {
	return &GuildReduced{
		ID:          g.ID,
		Name:        g.Name,
		Icon:        g.Icon,
		IconURL:     getIconURL(g.ID, g.Icon),
		Region:      g.Region,
		OwnerID:     g.OwnerID,
		JoinedAt:    g.JoinedAt,
		MemberCount: g.MemberCount,
	}
}

// MemberFromMember returns a Member from the passed
// discordgo.Member m.
func MemberFromMember(m *discordgo.Member) *Member {
	if m == nil {
		return nil
	}

	created, _ := discordutil.GetDiscordSnowflakeCreationTime(m.User.ID)
	return &Member{
		Member:    m,
		AvatarURL: m.User.AvatarURL(""),
		CreatedAt: created,
	}
}

// ReportFromReport returns a Report from the passed
// report.Report r and publicAddr to generate an
// attachment URL.
func ReportFromReport(r *report.Report, publicAddr string) *Report {
	rtype := static.ReportTypes[r.Type]
	r.AttachmehtURL = imgstore.GetLink(r.AttachmehtURL, publicAddr)
	return &Report{
		Report:   r,
		TypeName: rtype,
		Created:  r.GetTimestamp(),
	}
}

// getIconURL returns the CDN URL of a Discord icon
// resource with the passed guildID and iconHash.
func getIconURL(guildID, iconHash string) string {
	if iconHash == "" {
		return ""
	}

	return fmt.Sprintf("https://cdn.discordapp.com/icons/%s/%s.png", guildID, iconHash)
}
