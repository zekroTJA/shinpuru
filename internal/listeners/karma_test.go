package listeners

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di/v2"
	"github.com/stretchr/testify/mock"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/mocks"
)

const karmaUp, karmaDown = "👍", "👎"

type karmaHandlerMock struct {
	session *mocks.ISession
	db      *mocks.Database
	logger  *mocks.Logger
	karma   *mocks.KarmaProvider
	st      *mocks.IState

	ct di.Container
}

func getKarmaHandlerMock(prep ...func(t karmaHandlerMock)) karmaHandlerMock {
	var t karmaHandlerMock

	t.session = &mocks.ISession{}
	t.db = &mocks.Database{}
	t.logger = &mocks.Logger{}
	t.karma = &mocks.KarmaProvider{}
	t.st = &mocks.IState{}

	if len(prep) != 0 {
		prep[0](t)
	}

	t.logger.On("Errorf", mock.Anything, mock.Anything, mock.Anything).Return()
	t.logger.On("Errorf", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
	t.logger.On("Section", mock.Anything).Return(t.logger)

	t.db.On("GetKarmaTokens", mock.Anything).Return(3, nil)
	t.db.On("GetKarmaEmotes", mock.Anything).Return(karmaUp, karmaDown, nil)

	t.st.On("SelfUser").Return(&discordgo.User{ID: "self-id"}, nil)
	t.st.On("Message", "channel-id", "message-self").Return(&discordgo.Message{
		ID:        "message-self",
		ChannelID: "channel-id",
		GuildID:   "guild-id",
		Author: &discordgo.User{
			ID: "user-id",
		},
	}, nil)
	t.st.On("Message", "channel-id", "message-bot").Return(&discordgo.Message{
		ID:        "message-bot",
		ChannelID: "channel-id",
		GuildID:   "guild-id",
		Author: &discordgo.User{
			ID:  "user-bot",
			Bot: true,
		},
	}, nil)
	t.st.On("Message", "channel-id", "message-blocked").Return(&discordgo.Message{
		ID:        "message-bot",
		ChannelID: "channel-id",
		GuildID:   "guild-id",
		Author: &discordgo.User{
			ID: "user-blocked",
		},
	}, nil)
	t.st.On("Message", "channel-id", mock.Anything).Return(&discordgo.Message{
		ID:        "message-id-" + strconv.Itoa(rand.Int()),
		ChannelID: "channel-id",
		GuildID:   "guild-id",
		Author: &discordgo.User{
			ID: "author-id",
		},
	}, nil)

	t.karma.On("GetState", "guild-enabled").Return(true, nil)
	t.karma.On("GetState", "guild-disabled").Return(false, nil)
	t.karma.On("IsBlockListed", mock.Anything, "user-blocked").Return(true, nil)
	t.karma.On("IsBlockListed", mock.Anything, mock.Anything).Return(false, nil)
	t.karma.On("Update", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	t.session.On("User", "user-bot").Return(&discordgo.User{ID: "user-bot", Bot: true}, nil)
	t.session.On("User", "user-id").Return(&discordgo.User{ID: "user-id"}, nil)
	t.session.On("User", "user-id-2").Return(&discordgo.User{ID: "user-id-2"}, nil)
	t.session.On("User", "user-id-3").Return(&discordgo.User{ID: "user-id-3"}, nil)
	t.session.On("UserChannelCreate", "user-id").Return(&discordgo.Channel{ID: "user-id"}, nil)
	t.session.On("ChannelMessageSendEmbed", "user-id", mock.Anything).Return(nil, nil)
	t.session.On("ChannelMessage", "channel-id", "message-uncached").Return(&discordgo.Message{
		ID:        "message-id",
		ChannelID: "channel-id",
		GuildID:   "guild-id",
		Author: &discordgo.User{
			ID:  "author-id",
			Bot: true,
		},
	}, nil)

	ct, _ := di.NewBuilder()
	ct.Add(
		di.Def{
			Name:  static.DiDatabase,
			Build: func(ctn di.Container) (interface{}, error) { return t.db, nil },
		},
		di.Def{
			Name:  static.DiGuildLog,
			Build: func(ctn di.Container) (interface{}, error) { return t.logger, nil },
		},
		di.Def{
			Name:  static.DiKarma,
			Build: func(ctn di.Container) (interface{}, error) { return t.karma, nil },
		},
		di.Def{
			Name:  static.DiState,
			Build: func(ctn di.Container) (interface{}, error) { return t.st, nil },
		},
	)

	t.ct = ct.Build()

	return t
}

func TestKarmaHandler(t *testing.T) {
	m := getKarmaHandlerMock()

	l := NewListenerKarma(m.ct)

	// Bot User
	l.Handler(m.session, &discordgo.MessageReactionAdd{
		MessageReaction: &discordgo.MessageReaction{
			UserID:    "user-bot",
			MessageID: "message-id-1",
			ChannelID: "channel-id",
			GuildID:   "guild-enabled",
			Emoji: discordgo.Emoji{
				Name: karmaUp,
			},
		},
	})

	m.karma.AssertNotCalled(t, "Update", "guild-enabled", "author-id", "user-bot", 1)

	// Self User
	l.Handler(m.session, &discordgo.MessageReactionAdd{
		MessageReaction: &discordgo.MessageReaction{
			UserID:    "user-id",
			MessageID: "message-self",
			ChannelID: "channel-id",
			GuildID:   "guild-enabled",
			Emoji: discordgo.Emoji{
				Name: karmaUp,
			},
		},
	})

	m.karma.AssertNotCalled(t, "Update", "guild-enabled", "user-id", "user-id", 1)

	// Blocked Sender User
	l.Handler(m.session, &discordgo.MessageReactionAdd{
		MessageReaction: &discordgo.MessageReaction{
			UserID:    "user-blocked",
			MessageID: "message-id-1",
			ChannelID: "channel-id",
			GuildID:   "guild-enabled",
			Emoji: discordgo.Emoji{
				Name: karmaUp,
			},
		},
	})

	m.karma.AssertNotCalled(t, "Update", "guild-enabled", "author-id", "user-blocked", 1)

	// Blocked Receiver User
	l.Handler(m.session, &discordgo.MessageReactionAdd{
		MessageReaction: &discordgo.MessageReaction{
			UserID:    "user-id",
			MessageID: "message-blocked",
			ChannelID: "channel-id",
			GuildID:   "guild-enabled",
			Emoji: discordgo.Emoji{
				Name: karmaUp,
			},
		},
	})

	m.karma.AssertNotCalled(t, "Update", "guild-enabled", "user-blocked", "user-id", 1)

	// Disabled guild
	l.Handler(m.session, &discordgo.MessageReactionAdd{
		MessageReaction: &discordgo.MessageReaction{
			UserID:    "user-id",
			MessageID: "message-id-1",
			ChannelID: "channel-id",
			GuildID:   "guild-disabled",
			Emoji: discordgo.Emoji{
				Name: karmaUp,
			},
		},
	})

	m.karma.AssertNotCalled(t, "Update", "guild-disabled", "author-id", "user-id", 1)

	// Only apply to message once
	for i := 0; i < 3; i++ {
		l.Handler(m.session, &discordgo.MessageReactionAdd{
			MessageReaction: &discordgo.MessageReaction{
				UserID:    "user-id",
				MessageID: "message-id-1",
				ChannelID: "channel-id",
				GuildID:   "guild-enabled",
				Emoji: discordgo.Emoji{
					Name: karmaUp,
				},
			},
		})
	}

	l.Handler(m.session, &discordgo.MessageReactionAdd{
		MessageReaction: &discordgo.MessageReaction{
			UserID:    "user-id",
			MessageID: "message-id-2",
			ChannelID: "channel-id",
			GuildID:   "guild-enabled",
			Emoji: discordgo.Emoji{
				Name: karmaUp,
			},
		},
	})

	l.Handler(m.session, &discordgo.MessageReactionAdd{
		MessageReaction: &discordgo.MessageReaction{
			UserID:    "user-id",
			MessageID: "message-id-3",
			ChannelID: "channel-id",
			GuildID:   "guild-enabled",
			Emoji: discordgo.Emoji{
				Name: karmaDown,
			},
		},
	})

	// Tokens exceeded
	l.Handler(m.session, &discordgo.MessageReactionAdd{
		MessageReaction: &discordgo.MessageReaction{
			UserID:    "user-id",
			MessageID: "message-id-4",
			ChannelID: "channel-id",
			GuildID:   "guild-enabled",
			Emoji: discordgo.Emoji{
				Name: karmaDown,
			},
		},
	})

	m.karma.AssertCalled(t, "Update", "guild-enabled", "author-id", "user-id", 1)
	m.karma.AssertCalled(t, "Update", "guild-enabled", "author-id", "user-id", -1)

	m.karma.AssertNumberOfCalls(t, "Update", 3)
}
