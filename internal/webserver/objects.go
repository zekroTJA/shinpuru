package webserver

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type ListResponse struct {
	N    int         `json:"n"`
	Data interface{} `json:"data"`
}

type User struct {
	*discordgo.User

	AvatarURL string `json:"avatar_url"`
}

type Member struct {
	*discordgo.Member

	AvatarURL string `json:"avatar_url"`
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

type PermissionLvlResponse struct {
	Level int `json:"lvl"`
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
	return &Member{
		Member:    m,
		AvatarURL: m.User.AvatarURL(""),
	}
}

func getIconURL(guildID, iconHash string) string {
	if iconHash == "" {
		return ""
	}

	return fmt.Sprintf("https://cdn.discordapp.com/icons/%s/%s.png", guildID, iconHash)
}
