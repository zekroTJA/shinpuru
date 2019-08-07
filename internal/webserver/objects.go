package webserver

import (
	"fmt"
	"time"

	"github.com/zekroTJA/shinpuru/internal/core"

	"github.com/zekroTJA/shinpuru/internal/util"

	"github.com/bwmarrin/discordgo"
)

type ListResponse struct {
	N    int         `json:"n"`
	Data interface{} `json:"data"`
}

type User struct {
	*discordgo.User

	AvatarURL string    `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
}

type Member struct {
	*discordgo.Member

	AvatarURL string    `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
}

type Guild struct {
	*discordgo.Guild

	SelfMember *Member   `json:"self_member"`
	IconURL    string    `json:"icon_url"`
	Members    []*Member `json:"members"`
}

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

type PermissionsResponse struct {
	Permissions core.PermissionArray `json:"permissions"`
}

type Report struct {
	*util.Report

	TypeName string    `json:"type_name"`
	Created  time.Time `json:"created"`
}

type GuildSettings struct {
	Prefix              string                          `json:"prefix"`
	Perms               map[string]core.PermissionArray `json:"perms"`
	AutoRole            string                          `json:"autorole"`
	ModLogChannel       string                          `json:"modlogchannel"`
	VoiceLogChannel     string                          `json:"voicelogchannel"`
	JoinMessageChannel  string                          `json:"joinmessagechannel"`
	JoinMessageText     string                          `json:"joinmessagetext"`
	LeaveMessageChannel string                          `json:"leavemessagechannel"`
	LeaveMessageText    string                          `json:"leavemessagetext"`
}

type PermissionsUpdate struct {
	Perm    string   `json:"perm"`
	RoleIDs []string `json:"role_ids"`
}

type ReasonRequest struct {
	Reason     string `json:"reason"`
	Attachment string `json:"attachment"`
}

type ReportRequest struct {
	*ReasonRequest

	Type int `json:"type"`
}

func GuildFromGuild(g *discordgo.Guild, m *discordgo.Member) *Guild {
	membs := make([]*Member, len(g.Members))
	for i, m := range g.Members {
		membs[i] = MemberFromMember(m)
	}

	return &Guild{
		Guild:      g,
		SelfMember: MemberFromMember(m),
		Members:    membs,
		IconURL:    getIconURL(g.ID, g.Icon),
	}
}

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

func MemberFromMember(m *discordgo.Member) *Member {
	created, _ := util.GetDiscordSnowflakeCreationTime(m.User.ID)
	return &Member{
		Member:    m,
		AvatarURL: m.User.AvatarURL(""),
		CreatedAt: created,
	}
}

func ReportFromReport(r *util.Report) *Report {
	rtype := util.ReportTypes[r.Type]
	return &Report{
		Report:   r,
		TypeName: rtype,
		Created:  r.GetTimestamp(),
	}
}

func getIconURL(guildID, iconHash string) string {
	if iconHash == "" {
		return ""
	}

	return fmt.Sprintf("https://cdn.discordapp.com/icons/%s/%s.png", guildID, iconHash)
}
