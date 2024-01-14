package report

import (
	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/snowflake"
	"github.com/zekroTJA/shinpuru/internal/models"
	"time"
)

type TimeProvider interface {
	Now() time.Time
}

type Database interface {
	AddReport(rep models.Report) error
	DeleteReport(id snowflake.ID) error
	GetReportsFiltered(guildID, memberID string, repType models.ReportType, offset, limit int) ([]models.Report, error)
	GetExpiredReports() ([]models.Report, error)
	ExpireReports(id ...string) (err error)
	GetGuildModLog(guildID string) (string, error)
}

type Session interface {
	ChannelMessageSendEmbed(channelID string, embed *discordgo.MessageEmbed, options ...discordgo.RequestOption) (*discordgo.Message, error)
	UserChannelCreate(recipientID string, options ...discordgo.RequestOption) (st *discordgo.Channel, err error)
	GuildMemberDeleteWithReason(guildID, userID, reason string, options ...discordgo.RequestOption) (err error)
	GuildBanCreateWithReason(guildID, userID, reason string, days int, options ...discordgo.RequestOption) (err error)
	GuildMemberTimeout(guildID string, userID string, until *time.Time, options ...discordgo.RequestOption) (err error)
	GuildBanDelete(guildID, userID string, options ...discordgo.RequestOption) (err error)
}

type State interface {
	Guild(id string, hydrate ...bool) (v *discordgo.Guild, err error)
	Member(guildID, memberID string, forceNoFetch ...bool) (v *discordgo.Member, err error)
}
