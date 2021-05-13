package karma

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di/v2"
	"github.com/sirupsen/logrus"
	"github.com/zekroTJA/shinpuru/internal/models"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
	"github.com/zekroTJA/shinpuru/pkg/embedbuilder"
)

// Service provides functionalities to check karma state,
// karma blocklist and alter karma of a user.
type Service struct {
	s  *discordgo.Session
	db database.Database
}

// NewKarmaService initializes a new Service
// instance.
func NewKarmaService(container di.Container) (k *Service) {
	k = &Service{}

	k.s = container.Get(static.DiDiscordSession).(*discordgo.Session)
	k.db = container.Get(static.DiDatabase).(database.Database)

	return
}

// GetState returns the current enabled/disabled state
// of the karma system on the specified guild.
func (k *Service) GetState(guildID string) (ok bool, err error) {
	ok, err = k.db.GetKarmaState(guildID)
	if database.IsErrDatabaseNotFound(err) {
		err = nil
	}
	return
}

// IsBlockListed returns true if the passed user on the
// specified guild is blocked from gaining or giving
// karma.
func (k *Service) IsBlockListed(guildID, userID string) (isBlocklisted bool, err error) {
	isBlocklisted, err = k.db.IsKarmaBlockListed(guildID, userID)
	if database.IsErrDatabaseNotFound(err) {
		err = nil
	}
	return
}

// Update adds or removes karma of the given value of the
// specified user.
func (k *Service) Update(guildID, userID string, value int) (err error) {
	err = k.db.UpdateKarma(userID, guildID, value)

	rules, err := k.db.GetKarmaRules(guildID)
	if err != nil && !database.IsErrDatabaseNotFound(err) {
		return
	}
	if len(rules) > 0 {
		var valAfter int
		valAfter, err = k.db.GetKarma(userID, guildID)
		if err != nil {
			return
		}
		valBefore := valAfter - value

		for _, rule := range rules {
			if value > 0 {
				switch rule.Action {
				case models.KarmaActionToggleRole:
					if rule.Trigger == models.KarmaTriggerAbove && valBefore <= rule.Value && valAfter > rule.Value {
						if err = k.s.GuildMemberRoleAdd(guildID, userID, rule.Argument); err != nil {
							logrus.WithError(err).WithField("gid", guildID).WithField("uid", userID).Error("KARMA :: failed adding role")
						}
					} else if rule.Trigger == models.KarmaTriggerBelow && valBefore < rule.Value && valAfter >= rule.Value {
						if err = k.s.GuildMemberRoleRemove(guildID, userID, rule.Argument); err != nil {
							logrus.WithError(err).WithField("gid", guildID).WithField("uid", userID).Error("KARMA :: failed removing role")
						}
					}
				case models.KarmaActionSendMessage:
					if rule.Trigger == models.KarmaTriggerAbove && valBefore <= rule.Value && valAfter > rule.Value {
						k.trySendKarmaMessage(userID, guildID, true, rule.Value, rule.Argument)
					}
				case models.KarmaActionKick:
					if rule.Trigger == models.KarmaTriggerAbove && valBefore <= rule.Value && valAfter > rule.Value {
						k.tryKick(userID, guildID, true, rule.Value)
					}
				case models.KarmaActionBan:
					if rule.Trigger == models.KarmaTriggerAbove && valBefore <= rule.Value && valAfter > rule.Value {
						k.tryBan(userID, guildID, true, rule.Value)
					}
				}
			} else if value < 0 {
				switch rule.Action {
				case models.KarmaActionToggleRole:
					if rule.Trigger == models.KarmaTriggerAbove && valBefore > rule.Value && valAfter <= rule.Value {
						if err = k.s.GuildMemberRoleRemove(guildID, userID, rule.Argument); err != nil {
							logrus.WithError(err).WithField("gid", guildID).WithField("uid", userID).Error("KARMA :: failed removing role")
						}
					} else if rule.Trigger == models.KarmaTriggerBelow && valBefore >= rule.Value && valAfter < rule.Value {
						if err = k.s.GuildMemberRoleAdd(guildID, userID, rule.Argument); err != nil {
							logrus.WithError(err).WithField("gid", guildID).WithField("uid", userID).Error("KARMA :: failed adding role")
						}
					}
				case models.KarmaActionSendMessage:
					if rule.Trigger == models.KarmaTriggerBelow && valBefore >= rule.Value && valAfter < rule.Value {
						k.trySendKarmaMessage(userID, guildID, false, rule.Value, rule.Argument)
					}
				case models.KarmaActionKick:
					if rule.Trigger == models.KarmaTriggerBelow && valBefore >= rule.Value && valAfter < rule.Value {
						k.tryKick(userID, guildID, false, rule.Value)
					}
				case models.KarmaActionBan:
					if rule.Trigger == models.KarmaTriggerBelow && valBefore >= rule.Value && valAfter < rule.Value {
						k.tryBan(userID, guildID, false, rule.Value)
					}
				}
			}
		}
	}

	return
}

// CheckAndUpdate is shorthand for GetState, IsBlockListed
// and Update in one single pipe.
func (k *Service) CheckAndUpdate(guildID string, object *discordgo.User, value int) (ok bool, err error) {
	if object.Bot {
		return
	}

	enabled, err := k.GetState(guildID)
	if !enabled || err != nil {
		return
	}

	isBLocklisted, err := k.IsBlockListed(guildID, object.ID)
	if isBLocklisted || err != nil {
		return
	}

	err = k.Update(object.ID, guildID, value)
	ok = err == nil
	return
}

func (k *Service) trySendKarmaMessage(userID, guildID string, added bool, value int, content string) {
	ch, err := k.s.UserChannelCreate(userID)
	if err != nil {
		logrus.WithError(err).WithField("uid", userID).WithField("gid", guildID).Error("KARMA :: failed opening dm channel")
		return
	}

	guild, err := discordutil.GetGuild(k.s, guildID)
	if err != nil {
		logrus.WithError(err).WithField("uid", userID).WithField("gid", guildID).Error("KARMA :: failed getting guild details")
		return
	}

	pre := "dropped below"
	if added {
		pre = "rised above"
	}

	emb := embedbuilder.New().
		WithDescription(fmt.Sprintf("Your karma %s %d points on guild %s.\n\n%s",
			pre, value, guild.Name, content))

	if added {
		emb.WithColor(static.ColorEmbedGreen)
	} else {
		emb.WithColor(static.ColorEmbedOrange)
	}

	_, err = k.s.ChannelMessageSendEmbed(ch.ID, emb.Build())
	if err != nil {
		logrus.WithError(err).WithField("uid", userID).WithField("gid", guildID).Error("KARMA :: failed sending dm")
	}
}

func (k *Service) tryKick(userID, guildID string, added bool, value int) {
	k.trySendKarmaMessage(userID, guildID, added, value,
		"Because of that, you have been automatically kicked from the guild.")

	pre := "dropped below"
	if added {
		pre = "rised above"
	}
	reason := fmt.Sprintf("Karma %s %d points", pre, value)

	if err := k.s.GuildMemberDeleteWithReason(guildID, userID, reason); err != nil {
		logrus.WithError(err).WithField("uid", userID).WithField("gid", guildID).Error("KARMA :: failed kicking member")
	}
}

func (k *Service) tryBan(userID, guildID string, added bool, value int) {
	k.trySendKarmaMessage(userID, guildID, added, value,
		"Because of that, you have been automatically banned from the guild.")

	pre := "dropped below"
	if added {
		pre = "rised above"
	}
	reason := fmt.Sprintf("Karma %s %d points", pre, value)

	if err := k.s.GuildBanCreateWithReason(guildID, userID, reason, 7); err != nil {
		logrus.WithError(err).WithField("uid", userID).WithField("gid", guildID).Error("KARMA :: failed banning member")
	}
}
