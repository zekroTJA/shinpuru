package karma

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/models"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/guildlog"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/embedbuilder"
	"github.com/zekrotja/dgrs"
	"github.com/zekrotja/rogu"
	"github.com/zekrotja/rogu/log"
)

// Service provides functionalities to check karma state,
// karma blocklist and alter karma of a user.
type Service struct {
	s   *discordgo.Session
	db  database.Database
	gl  guildlog.Logger
	st  *dgrs.State
	log *rogu.Logger
}

var _ Provider = (*Service)(nil)

// NewKarmaService initializes a new Service
// instance.
func NewKarmaService(container di.Container) (k *Service) {
	k = &Service{}

	k.s = container.Get(static.DiDiscordSession).(*discordgo.Session)
	k.db = container.Get(static.DiDatabase).(database.Database)
	k.gl = container.Get(static.DiGuildLog).(guildlog.Logger).Section("karma")
	k.st = container.Get(static.DiState).(*dgrs.State)
	k.log = log.Tagged("Karma")

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
func (k *Service) Update(guildID, userID, executorID string, value int) (err error) {
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
							k.log.Error().Err(err).Fields("gid", guildID, "uid", userID).Msg("Failed adding role")
							k.gl.Errorf(guildID, "Failed adding role to user (%s): %s", userID, err.Error())
						}
					} else if rule.Trigger == models.KarmaTriggerBelow && valBefore < rule.Value && valAfter >= rule.Value {
						if err = k.s.GuildMemberRoleRemove(guildID, userID, rule.Argument); err != nil {
							k.log.Error().Err(err).Fields("gid", guildID, "uid", userID).Msg("Failed removing role")
							k.gl.Errorf(guildID, "Failed removing role to user (%s): %s", userID, err.Error())
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
							k.log.Error().Err(err).Fields("gid", guildID, "uid", userID).Msg("Failed removing role")
							k.gl.Errorf(guildID, "Failed removing role to user (%s): %s", userID, err.Error())
						}
					} else if rule.Trigger == models.KarmaTriggerBelow && valBefore >= rule.Value && valAfter < rule.Value {
						if err = k.s.GuildMemberRoleAdd(guildID, userID, rule.Argument); err != nil {
							k.log.Error().Err(err).Fields("gid", guildID, "uid", userID).Msg("Failed adding role")
							k.gl.Errorf(guildID, "Failed adding role to user (%s): %s", userID, err.Error())
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

	if value < 0 {
		err = k.ApplyPenalty(guildID, executorID)
	}

	return
}

func (k *Service) ApplyPenalty(guildID, userID string) (err error) {
	if userID == "" {
		return
	}

	enabled, err := k.db.GetKarmaPenalty(guildID)
	if database.IsErrDatabaseNotFound(err) {
		err = nil
	}
	if err != nil || !enabled {
		return
	}
	fmt.Println(enabled, guildID, userID)

	err = k.Update(guildID, userID, "", -1)
	return
}

// CheckAndUpdate is shorthand for GetState, IsBlockListed
// and Update in one single pipe.
func (k *Service) CheckAndUpdate(guildID, executorID string, object *discordgo.User, value int) (ok bool, err error) {
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

	err = k.Update(guildID, object.ID, executorID, value)
	ok = err == nil
	return
}

func (k *Service) trySendKarmaMessage(userID, guildID string, added bool, value int, content string) {
	ch, err := k.s.UserChannelCreate(userID)
	if err != nil {
		k.log.Error().Err(err).Fields("uid", userID, "gid", guildID).Msg("Failed opening dm channel")
		return
	}

	guild, err := k.st.Guild(guildID, true)
	if err != nil {
		k.log.Error().Err(err).Fields("uid", userID, "gid", guildID).Msg("Failed getting guild details")
		k.gl.Errorf(guildID, "Failed getting guild details: %s", err.Error())
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
		k.log.Error().Err(err).Fields("uid", userID, "gid", guildID).Msg("Failed sending dm")
		k.gl.Errorf(guildID, "Failed sending DM to user (%s): %s", userID, err.Error())
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
		k.log.Error().Err(err).Fields("uid", userID, "gid", guildID).Msg("Failed kicking member")
		k.gl.Errorf(guildID, "Failed kicking user (%s): %s", userID, err.Error())
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
		k.log.Error().Err(err).Fields("uid", userID, "gid", guildID).Msg("Failed banning member")
		k.gl.Errorf(guildID, "Failed banning user (%s): %s", userID, err.Error())
	}
}
