package listeners

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/ratelimit"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/timedmap"
)

const (
	reactionsAddKarma    = "👍👌⭐✔"
	reactionsRemoveKarma = "👎❌"

	lifetimePerMessage = 24 * time.Hour

	rateLimiterTokens   = 5
	rateLimiterRestore  = time.Hour / rateLimiterTokens
	lifetimeRateLimiter = rateLimiterRestore * rateLimiterTokens
)

const (
	typNull   = 0
	typAdd    = 1
	typRemove = -1
)

type ListenerKarma struct {
	db    database.Database
	cache *timedmap.TimedMap

	msgsApplied timedmap.Section
	limiters    timedmap.Section
}

func NewListenerKarma(db database.Database) *ListenerKarma {
	cache := timedmap.New(5 * time.Minute)
	return &ListenerKarma{
		db:    db,
		cache: cache,

		msgsApplied: cache.Section(0),
		limiters:    cache.Section(1),
	}
}

func (l *ListenerKarma) Handler(s *discordgo.Session, e *discordgo.MessageReactionAdd) {
	if e.UserID == s.State.User.ID {
		return
	}

	var typ int
	if strings.Contains(reactionsAddKarma, e.MessageReaction.Emoji.Name) {
		typ = typAdd
	} else if strings.Contains(reactionsRemoveKarma, e.MessageReaction.Emoji.Name) {
		typ = typRemove
	}

	if typ == typNull {
		return
	}

	user, err := s.User(e.UserID)
	if err != nil {
		util.Log.Errorf("failed getting user %s: %s", e.UserID, err.Error())
		return
	}

	if user.Bot {
		return
	}

	if l.isMessageAlreadyApplied(e.UserID, e.MessageID) {
		return
	}

	if !l.rateLimiterTake(e.UserID, e.GuildID) {
		// TODO: Send message that karma credits are exceeded
		return
	}

	msg, err := s.State.Message(e.ChannelID, e.MessageID)
	if err != nil {
		if msg, err = s.ChannelMessage(e.ChannelID, e.MessageID); err != nil {
			util.Log.Errorf("failed getting message %s: %s", e.MessageID, err.Error())
			return
		}
	}

	if msg.Author.Bot || msg.Author.ID == e.UserID {
		return
	}

	if err = l.db.UpdateKarma(msg.Author.ID, e.GuildID, typ); err != nil {
		util.Log.Errorf("failed updating karma: %s", err.Error())
		return
	}

	l.applyMessage(e.UserID, e.MessageID)
}

func (l *ListenerKarma) isMessageAlreadyApplied(userID, msgID string) bool {
	key := fmt.Sprintf("%s:%s", userID, msgID)
	return l.msgsApplied.Contains(key)
}

func (l *ListenerKarma) applyMessage(userID, msgID string) {
	key := fmt.Sprintf("%s:%s", userID, msgID)
	l.msgsApplied.Set(key, true, lifetimePerMessage)
}

func (l *ListenerKarma) rateLimiterTake(userID, guildID string) bool {
	key := fmt.Sprintf("%s:%s", userID, guildID)

	limiter, ok := l.limiters.GetValue(key).(*ratelimit.Limiter)

	if !ok || limiter == nil {
		limiter = ratelimit.NewLimiter(rateLimiterRestore, rateLimiterTokens)
		l.limiters.Set(key, limiter, lifetimeRateLimiter)
	}

	expires, err := l.limiters.GetExpires(key)
	if err != nil {
		expires = time.Now()
	}

	refresh := lifetimeRateLimiter - expires.Sub(time.Now())
	l.limiters.Refresh(key, refresh)

	return limiter.Allow()
}
