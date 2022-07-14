package report

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/models"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
	"github.com/zekroTJA/shinpuru/pkg/multierror"
)

type Provider interface {
	PushReport(rep models.Report) (models.Report, error)
	PushKick(rep models.Report) (models.Report, error)
	PushBan(rep models.Report) (models.Report, error)
	PushMute(rep models.Report) (models.Report, error)
	RevokeMute(guildID, executorID, victimID, reason string) (emb *discordgo.MessageEmbed, err error)
	RevokeReport(rep models.Report, executorID, reason,
		wsPublicAddr string,
		db database.Database,
		s discordutil.ISession,
	) (emb *discordgo.MessageEmbed, err error)
	ExpireLastReport(guildID, victimID string, typ int) (err error)
	ExpireExpiredReports() (mErr *multierror.MultiError)
}
