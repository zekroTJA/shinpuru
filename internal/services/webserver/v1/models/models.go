package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
	sharedmodels "github.com/zekroTJA/shinpuru/internal/models"
	"github.com/zekroTJA/shinpuru/internal/services/backup/backupmodels"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	permService "github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/util/imgstore"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
	"github.com/zekroTJA/shinpuru/pkg/permissions"
	"github.com/zekroTJA/shinpuru/pkg/versioncheck"
	"github.com/zekroTJA/shireikan"
	"github.com/zekrotja/ken"
)

var Ok = &Status{200}

type Status struct {
	Code int `json:"code"`
}

type State struct {
	State bool `json:"state"`
}

type AccessTokenResponse struct {
	Token   string    `json:"token"`
	Expires time.Time `json:"expires"`
}

type Error struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Context string `json:"context,omitempty"`
}

// ListResponse wraps a list response object
// with the list as Data and N as len(Data).
type ListResponse[T any] struct {
	N    int `json:"n"`
	Data []T `json:"data"`
}

func NewListResponse[T any](data []T) ListResponse[T] {
	return ListResponse[T]{len(data), data}
}

// User extends a discordgo.User as reponse
// model.
type User struct {
	*discordgo.User

	AvatarURL       string    `json:"avatar_url"`
	CreatedAt       time.Time `json:"created_at"`
	BotOwner        bool      `json:"bot_owner"`
	CaptchaVerified bool      `json:"captcha_verified"`
}

// FlatUser shrinks the user object to the only
// necessary parts for the web interface.
type FlatUser struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	AvatarURL     string `json:"avatar_url"`
	Bot           bool   `json:"bot"`
}

// Member extends a discordgo.Member as
// response model.
type Member struct {
	*discordgo.Member

	GuildName  string    `json:"guild_name,omitempty"`
	AvatarURL  string    `json:"avatar_url"`
	CreatedAt  time.Time `json:"created_at"`
	Dominance  int       `json:"dominance"`
	Karma      int       `json:"karma"`
	KarmaTotal int       `json:"karma_total"`
	ChatMuted  bool      `json:"chat_muted"`
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
	JoinedAt                 time.Time                   `json:"joined_at"`
	Splash                   string                      `json:"splash"`
	MemberCount              int                         `json:"member_count"`
	VerificationLevel        discordgo.VerificationLevel `json:"verification_level"`
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
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	Icon              string    `json:"icon"`
	IconURL           string    `json:"icon_url"`
	Region            string    `json:"region"`
	OwnerID           string    `json:"owner_id"`
	JoinedAt          time.Time `json:"joined_at"`
	MemberCount       int       `json:"member_count"`
	OnlineMemberCount int       `json:"online_member_count,omitempty"`
}

// PermissionsResponse wraps a
// permissions.PermissionsArra as response
// model.
type PermissionsResponse struct {
	Permissions permissions.PermissionArray `json:"permissions"`
}

// Report extends models.Report by TypeName
// and Created time.
type Report struct {
	*sharedmodels.Report

	TypeName string    `json:"type_name"`
	Created  time.Time `json:"created"`
	Executor *FlatUser `json:"executor,omitempty"`
	Victim   *FlatUser `json:"victim,omitempty"`
}

// GuildSettings is the response model for
// guild settings and preferences.
type GuildSettings struct {
	Prefix              string                                 `json:"prefix"`
	Perms               map[string]permissions.PermissionArray `json:"perms"`
	AutoRoles           []string                               `json:"autoroles"`
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
	Reason     string     `json:"reason"`
	Attachment string     `json:"attachment"`
	Timeout    *time.Time `json:"timeout"`
}

// ReportRequest extends ReasonRequest by
// Type of report.
type ReportRequest struct {
	*ReasonRequest

	Type sharedmodels.ReportType `json:"type"`
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

type LandingPageResponse struct {
	LocalInvite        string `json:"localinvite"`
	PublicMainInvite   string `json:"publicmaininvite"`
	PublicCanaryInvite string `json:"publiccaranyinvite"`
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

// GuildKarmaEntry wraps a Member model and karma
// value for an entry of the karma scoreboard
// of a guild.
type GuildKarmaEntry struct {
	Member *Member `json:"member"`
	Value  int     `json:"value"`
}

// CommandInfo wraps a shireikan.Command object
// containing all information of a command
// instance.
type CommandInfo struct {
	Invokes            []string                  `json:"invokes"`
	Description        string                    `json:"description"`
	Help               string                    `json:"help"`
	Group              string                    `json:"group"`
	DomainName         string                    `json:"domain_name"`
	SubPermissionRules []shireikan.SubPermission `json:"sub_permission_rules"`
	IsExecutableInDM   bool                      `json:"is_executable_in_dm"`
}

// SlashCommandInfo wraps a slash command object
// containing all information of a slash command
// instance.
type SlashCommandInfo struct {
	Name        string                                `json:"name"`
	Description string                                `json:"description"`
	Version     string                                `json:"version"`
	Options     []*discordgo.ApplicationCommandOption `json:"options"`
	Domain      string                                `json:"domain"`
	SubDomains  []permService.SubPermission           `json:"subdomains"`
	DmCapable   bool                                  `json:"dm_capable"`
	Group       string                                `json:"group"`
}

// KarmaSettings wraps settings properties for
// guild karma settings.
type KarmaSettings struct {
	State          bool     `json:"state"`
	EmotesIncrease []string `json:"emotes_increase"`
	EmotesDecrease []string `json:"emotes_decrease"`
	Tokens         int      `json:"tokens"`
	Penalty        bool     `json:"penalty"`
}

// AntiraidSettings wraps settings properties for
// guild antiraid settings.
type AntiraidSettings struct {
	State              bool `json:"state"`
	RegenerationPeriod int  `json:"regeneration_period"`
	Burst              int  `json:"burst"`
	Verification       bool `json:"verification"`
}

type UsersettingsOTA struct {
	Enabled bool `json:"enabled"`
}

type UsersettingsPrivacy struct {
	StarboardOptout bool `json:"starboard_optout"`
}

// StarboardEntryResponse wraps a starboard entry
// as response model containing hydrated information
// of the author.
type StarboardEntryResponse struct {
	*sharedmodels.StarboardEntry

	MessageURL     string `json:"message_url"`
	AuthorUsername string `json:"author_username"`
	AvatarURL      string `json:"author_avatar_url"`
}

type PermissionsMap map[string]permissions.PermissionArray

type EnableStatus struct {
	Enabled bool `json:"enabled"`
}

type FlushGuildRequest struct {
	Validation string `json:"validation"`
	LeaveAfter bool   `json:"leave_after"`
}

type SearchResult struct {
	Guilds  []*GuildReduced `json:"guilds"`
	Members []*Member       `json:"members"`
}

type GuildAPISettingsRequest struct {
	sharedmodels.GuildAPISettings
	NewToken   string `json:"token"`
	ResetToken bool   `json:"reset_token"`
}

type AntiraidActionType int

const (
	AntiraidActionTypeKick = iota
	AntiraidActionTypeBan
)

type AntiraidAction struct {
	Type AntiraidActionType `json:"type"`
	IDs  []string           `json:"ids"`
}

type ChannelWithPermissions struct {
	*discordgo.Channel

	CanRead  bool `json:"can_read"`
	CanWrite bool `json:"can_write"`
}

type CaptchaSiteKey struct {
	SiteKey string `json:"sitekey"`
}

type CaptchaVerificationRequest struct {
	Token string `json:"token"`
}

type CodeExecSettings struct {
	EnableStatus

	Type                string   `json:"type"`
	TypesOptions        []string `json:"types_options,omitempty"`
	JdoodleClientId     string   `json:"jdoodle_clientid,omitempty"`
	JdoodleClientSecret string   `json:"jdoodle_clientsecret,omitempty"`
}

type PushCodeRequest struct {
	Code string `json:"code"`
}

type UpdateInfoResponse struct {
	Current    versioncheck.Semver `json:"current"`
	CurrentStr string              `json:"current_str"`
	Latest     versioncheck.Semver `json:"latest"`
	LatestStr  string              `json:"latest_str"`
	IsOld      bool                `json:"isold"`
}

// Validate returns true, when the ReasonRequest is valid.
// Otherwise, false is returned and an error response is
// returned.
func (req *ReasonRequest) Validate(acceptEmptyReason bool) (bool, error) {
	if !acceptEmptyReason && len(req.Reason) < 3 {
		return false, errors.New("invalid argument")
	}

	if req.Attachment != "" && !imgstore.ImgUrlSRx.MatchString(req.Attachment) {
		return false, fmt.Errorf("attachment must be a valid url to a file with type of png, jpg, jpeg, gif, ico, tiff, img, bmp or mp4.")
	}

	return true, nil
}

// GuildFromGuild returns a Guild model from the passed
// discordgo.Guild g, discordgo.Member m and cmdHandler.
func GuildFromGuild(g *discordgo.Guild, m *discordgo.Member, db database.Database, botOwnerID string) (ng *Guild, err error) {
	if g == nil {
		return
	}

	selfmm := MemberFromMember(m)

	if m != nil {
		switch {
		case discordutil.IsAdmin(g, m):
			selfmm.Dominance = 1
		case g.OwnerID == m.User.ID:
			selfmm.Dominance = 2
		case botOwnerID == m.User.ID:
			selfmm.Dominance = 3
		}
	}

	ng = &Guild{
		AfkChannelID:             g.AfkChannelID,
		Banner:                   g.Banner,
		Channels:                 g.Channels,
		Description:              g.Description,
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
		IconURL:    g.IconURL(),
	}

	if db != nil {
		selfmm.Karma, err = db.GetKarma(m.User.ID, g.ID)
		if !database.IsErrDatabaseNotFound(err) && err != nil {
			return
		}

		selfmm.KarmaTotal, err = db.GetKarmaSum(m.User.ID)
		if !database.IsErrDatabaseNotFound(err) && err != nil {
			return
		}

		ng.BackupsEnabled, err = db.GetGuildBackup(g.ID)
		if err != nil && !database.IsErrDatabaseNotFound(err) {
			return
		}

		var backupEntries []*backupmodels.Entry
		backupEntries, err = db.GetBackups(g.ID)
		if err != nil && !database.IsErrDatabaseNotFound(err) {
			return
		} else {
			for _, e := range backupEntries {
				if e.Timestamp.After(ng.LatestBackupEntry) {
					ng.LatestBackupEntry = e.Timestamp
				}
			}
		}

		status, err := db.GetGuildInviteBlock(g.ID)
		if err != nil && !database.IsErrDatabaseNotFound(err) {
			logrus.WithError(err).WithField("gid", g.ID).Error("Failed getting inviteblock status")
		} else {
			ng.InviteBlockEnabled = status != ""
		}
	}

	return
}

// GuildReducedFromGuild returns a GuildReduced from the passed
// discordgo.Guild g.
func GuildReducedFromGuild(g *discordgo.Guild) *GuildReduced {
	return &GuildReduced{
		ID:          g.ID,
		Name:        g.Name,
		Icon:        g.Icon,
		IconURL:     g.IconURL(),
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
// models.Report r and publicAddr to generate an
// attachment URL.
func ReportFromReport(r *sharedmodels.Report, publicAddr string) *Report {
	rtype := sharedmodels.ReportTypes[r.Type]
	r.AttachmentURL = imgstore.GetLink(r.AttachmentURL, publicAddr)
	return &Report{
		Report:   r,
		TypeName: rtype,
		Created:  r.GetTimestamp(),
	}
}

func GetCommandInfoFromCommand(cmd shireikan.Command) *CommandInfo {
	return &CommandInfo{
		Invokes:            cmd.GetInvokes(),
		Description:        cmd.GetDescription(),
		DomainName:         cmd.GetDomainName(),
		Group:              cmd.GetGroup(),
		Help:               cmd.GetHelp(),
		IsExecutableInDM:   cmd.IsExecutableInDMChannels(),
		SubPermissionRules: cmd.GetSubPermissionRules(),
	}
}

func GetSlashCommandInfoFromCommand(cmd *ken.CommandInfo) (ci *SlashCommandInfo) {
	ci = new(SlashCommandInfo)

	ci.Name = cmd.ApplicationCommand.Name
	ci.Description = cmd.ApplicationCommand.Description
	ci.Options = cmd.ApplicationCommand.Options
	ci.Version = cmd.ApplicationCommand.Version
	ci.Domain = cmd.Implementations["Domain"][0].(string)
	ci.SubDomains = cmd.Implementations["SubDomains"][0].([]permService.SubPermission)

	if v, ok := cmd.Implementations["IsDmCapable"]; ok && len(v) != 0 {
		ci.DmCapable = v[0].(bool)
	}

	domainSplit := strings.Split(ci.Domain, ".")
	ci.Group = strings.Join(domainSplit[1:len(domainSplit)-1], " ")
	ci.Group = strings.ToUpper(ci.Group)

	return
}

// FlatUserFromUser returns the reduced FlatUser object
// from the given user object.
func FlatUserFromUser(u *discordgo.User) (fu *FlatUser) {
	return &FlatUser{
		ID:            u.ID,
		Username:      u.Username,
		Discriminator: u.Discriminator,
		AvatarURL:     u.AvatarURL(""),
		Bot:           u.Bot,
	}
}
