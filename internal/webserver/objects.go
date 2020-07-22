package webserver

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dgrijalva/jwt-go"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"

	"github.com/zekroTJA/shinpuru/internal/commands"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/core/permissions"
	"github.com/zekroTJA/shinpuru/internal/util"
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
	ID                       string                      `json:"id"`
	Name                     string                      `json:"name"`
	Icon                     string                      `json:"icon"`
	Region                   string                      `json:"region"`
	AfkChannelID             string                      `json:"afk_channel_id"`
	OwnerID                  string                      `json:"owner_id"`
	JoinedAt                 discordgo.Timestamp         `json:"joined_at"`
	Splash                   string                      `json:"splash"`
	MemberCount              int                         `json:"member_count"`
	VerificationLevel        discordgo.VerificationLevel `json:"verification_level"`
	EmbedEnabled             bool                        `json:"embed_enabled"`
	Large                    bool                        `json:"large"`
	Unavailable              bool                        `json:"unavailable"`
	MfaLevel                 discordgo.MfaLevel          `json:"mfa_level"`
	Description              string                      `json:"description"`
	Banner                   string                      `json:"banner"`
	PremiumTier              discordgo.PremiumTier       `json:"premium_tier"`
	PremiumSubscriptionCount int                         `json:"premium_subscription_count"`

	Roles    []*discordgo.Role    `json:"roles"`
	Channels []*discordgo.Channel `json:"channels"`

	SelfMember         *Member   `json:"self_member"`
	IconURL            string    `json:"icon_url"`
	BackupsEnabled     bool      `json:"backups_enabled"`
	LatestBackupEntry  time.Time `json:"latest_backup_entry"`
	InviteBlockEnabled bool      `json:"invite_block_enabled"`
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
// guild settings and preferences.
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

// APITokenResponse wraps the reponse model of
// an apit token request.
type APITokenResponse struct {
	Created    time.Time `json:"created"`
	Expires    time.Time `json:"expires"`
	LastAccess time.Time `json:"last_access"`
	Hits       int       `json:"hits"`
	Token      string    `json:"token,omitempty"`
}

// APITokenClaims extends the standard JWT claims
// by private claims used for api tokens.
type APITokenClaims struct {
	jwt.StandardClaims

	Salt string `json:"sp_salt,omitempty"`
}

// SessionTokenClaims extends the standard JWT
// claims by information used for session tokens.
//
// Currently, no additional information is
// extended but this wrapper is used tho to
// be able to add session information later.
type SessionTokenClaims struct {
	jwt.StandardClaims
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
func GuildFromGuild(g *discordgo.Guild, m *discordgo.Member, cmdHandler *commands.CmdHandler, db database.Database) *Guild {
	if g == nil {
		return nil
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

	ng := &Guild{
		AfkChannelID:             g.AfkChannelID,
		Banner:                   g.Banner,
		Channels:                 g.Channels,
		Description:              g.Description,
		EmbedEnabled:             g.EmbedEnabled,
		ID:                       g.ID,
		Icon:                     g.Icon,
		JoinedAt:                 g.JoinedAt,
		Large:                    g.Large,
		MemberCount:              g.MemberCount,
		MfaLevel:                 g.MfaLevel,
		Name:                     g.Name,
		OwnerID:                  g.OwnerID,
		PremiumSubscriptionCount: g.PremiumSubscriptionCount,
		PremiumTier:              g.PremiumTier,
		Region:                   g.Region,
		Roles:                    g.Roles,
		Splash:                   g.Splash,
		Unavailable:              g.Unavailable,
		VerificationLevel:        g.VerificationLevel,

		SelfMember: selfmm,
		IconURL:    getIconURL(g.ID, g.Icon),
	}

	if db != nil {
		var err error

		ng.BackupsEnabled, err = db.GetGuildBackup(g.ID)
		if err != nil && !database.IsErrDatabaseNotFound(err) {
			util.Log.Errorf("failed getting backup status of guild %s: %s", g.ID, err.Error())
		}

		backupEntries, err := db.GetBackups(g.ID)
		if err != nil && !database.IsErrDatabaseNotFound(err) {
			util.Log.Errorf("failed getting backup entries of guild %s: %s", g.ID, err.Error())
		} else {
			for _, e := range backupEntries {
				if e.Timestamp.After(ng.LatestBackupEntry) {
					ng.LatestBackupEntry = e.Timestamp
				}
			}
		}

		status, err := db.GetGuildInviteBlock(g.ID)
		if err != nil && !database.IsErrDatabaseNotFound(err) {
			util.Log.Errorf("failed getting inviteblock status of guild %s: %s", g.ID, err.Error())
		} else {
			ng.InviteBlockEnabled = status != ""
		}
	}

	return ng
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

// APITokenClaimsFromMap creates an APITokenClaims
// model from given jwt.MapClaims.
func APITokenClaimsFromMap(m jwt.MapClaims) APITokenClaims {
	c := APITokenClaims{
		StandardClaims: standardClaimsFromMap(m),
	}

	c.Salt, _ = m["sp_salt"].(string)

	return c
}

// SessionTokenClaimsFromMap creates an SessionTokenClaims
// model from given jwt.MapClaims.
func SessionTokenClaimsFromMap(m jwt.MapClaims) SessionTokenClaims {
	c := SessionTokenClaims{
		StandardClaims: standardClaimsFromMap(m),
	}

	return c
}

// standardClaimsFromMap creates a jwt.StandardClaims
// model from the given jwt.MapClaims.
func standardClaimsFromMap(m jwt.MapClaims) jwt.StandardClaims {
	c := jwt.StandardClaims{}

	c.Issuer, _ = m["iss"].(string)
	c.Subject, _ = m["sub"].(string)
	c.ExpiresAt, _ = m["exp"].(int64)
	c.NotBefore, _ = m["nbf"].(int64)
	c.IssuedAt, _ = m["iat"].(int64)

	return c
}

// getIconURL returns the CDN URL of a Discord icon
// resource with the passed guildID and iconHash.
func getIconURL(guildID, iconHash string) string {
	if iconHash == "" {
		return ""
	}

	return fmt.Sprintf("https://cdn.discordapp.com/icons/%s/%s.png", guildID, iconHash)
}
