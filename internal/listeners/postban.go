package listeners

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/xid"
	"github.com/sarulabs/di/v2"
	"github.com/sirupsen/logrus"
	"github.com/zekroTJA/shinpuru/internal/models"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/guildlog"
	"github.com/zekroTJA/shinpuru/internal/services/permissions"
	"github.com/zekroTJA/shinpuru/internal/services/report"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
	"github.com/zekrotja/ken"
)

type ListenerPostBan struct {
	ken *ken.Ken
	db  database.Database
	gl  guildlog.Logger
	rep report.Provider
	pmw permissions.Provider
}

func NewListenerPostBan(ctn di.Container) ListenerPostBan {
	return ListenerPostBan{
		ken: ctn.Get(static.DiCommandHandler).(*ken.Ken),
		db:  ctn.Get(static.DiDatabase).(database.Database),
		gl:  ctn.Get(static.DiGuildLog).(guildlog.Logger).Section("postban"),
		rep: ctn.Get(static.DiReport).(report.Provider),
		pmw: ctn.Get(static.DiPermissions).(permissions.Provider),
	}
}

func (t ListenerPostBan) Handler(s discordutil.ISession, e *discordgo.GuildBanAdd) {
	modlogChan, err := t.db.GetGuildModLog(e.GuildID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		t.error(e.GuildID, "failed getting modlog channel", err)
		return
	}
	if modlogChan == "" {
		return
	}

	time.Sleep(5 * time.Second)
	auditLog, err := s.GuildAuditLog(e.GuildID, "", "",
		int(discordgo.AuditLogActionMemberBanAdd), 10)
	if err != nil {
		t.error(e.GuildID, "failed getting guild audit log", err)
		return
	}

	var banEntry *discordgo.AuditLogEntry
	for _, logEntry := range auditLog.AuditLogEntries {
		if logEntry.TargetID == e.User.ID {
			banEntry = logEntry
			break
		}
	}
	if banEntry == nil {
		logrus.WithField("guildid", e.GuildID).
			Warn("could not find any ban entry in audit log")
		return
	}

	rep := models.Report{
		Type:       models.TypeBan,
		GuildID:    e.GuildID,
		ExecutorID: banEntry.UserID,
		VictimID:   e.User.ID,
		Msg:        banEntry.Reason,
		Anonymous:  true,
	}

	emb := rep.AsEmbed("")
	emb.Title = "User banned"
	emb.Description = "A user has just been banned."

	msg, err := s.ChannelMessageSendEmbed(modlogChan, emb)
	if err != nil {
		t.error(e.GuildID, "failed sending ban message", err)
	}

	_, err = t.ken.Components().Add(msg.ID, msg.ChannelID).
		Condition(func(ctx ken.ComponentContext) bool {
			ok, _, err := t.pmw.CheckPermissions(s, e.GuildID, ctx.User().ID, "sp.guild.mod.report")
			return ok && err == nil
		}).
		AddActionsRow(func(b ken.ComponentAssembler) {
			b.Add(discordgo.Button{
				CustomID: xid.New().String(),
				Label:    "Create Report in shinpuru",
				Style:    discordgo.PrimaryButton,
			}, func(ctx ken.ComponentContext) bool {
				reasonId := xid.New().String()
				attachmentId := xid.New().String()
				cModal, err := ctx.OpenModal("Create Ban Report Entry", "", func(b ken.ComponentAssembler) {
					b.AddActionsRow(func(b ken.ComponentAssembler) {
						b.Add(discordgo.TextInput{
							CustomID:  reasonId,
							Label:     "Reason",
							Style:     discordgo.TextInputParagraph,
							Value:     rep.Msg,
							Required:  true,
							MinLength: 3,
						}, nil)
					})
					b.AddActionsRow(func(b ken.ComponentAssembler) {
						b.Add(discordgo.TextInput{
							CustomID:    attachmentId,
							Label:       "Attachment URL",
							Style:       discordgo.TextInputShort,
							Placeholder: "A media URL attached to the report.",
						}, nil)
					})
				})
				if err != nil {
					return false
				}

				ctxModal := <-cModal
				if err = ctxModal.Defer(); err != nil {
					return false
				}

				rep.Msg = ctxModal.GetComponentByID(reasonId).GetValue()
				rep.AttachmentURL = ctxModal.GetComponentByID(attachmentId).GetValue()

				_, err = t.rep.PushReport(rep)
				if err != nil {
					return false
				}

				ctxModal.FollowUpEmbed(&discordgo.MessageEmbed{
					Description: "The report has been created.",
				}).DeleteAfter(8 * time.Second)
				s.ChannelMessageDelete(msg.ChannelID, msg.ID)

				return true
			})

			b.Add(discordgo.Button{
				CustomID: xid.New().String(),
				Label:    "No further action",
				Style:    discordgo.SecondaryButton,
			}, func(ctx ken.ComponentContext) bool {
				return true
			})
		}, true).Build()
	if err != nil {
		t.error(e.GuildID, "failed appending message components to message", err)
	}
}

func (t ListenerPostBan) error(guildID string, msg string, err error) {
	logrus.WithError(err).
		WithField("guild", guildID).
		Error(msg)
	t.gl.Errorf(guildID, "%s: %s", msg, err.Error())
}
