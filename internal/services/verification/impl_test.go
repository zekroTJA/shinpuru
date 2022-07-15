package verification

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zekroTJA/shinpuru/internal/models"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/internal/util/testutil"
	"github.com/zekroTJA/shinpuru/mocks"
)

type verificationMock struct {
	s   *mocks.ISession
	db  *mocks.Database
	cfg *mocks.ConfigProvider
	gl  *mocks.Logger
	tp  *mocks.TimeProvider

	ct di.Container
}

func getVerificationMock(prep ...func(m verificationMock)) verificationMock {
	var t verificationMock

	t.s = &mocks.ISession{}
	t.db = &mocks.Database{}
	t.cfg = &mocks.ConfigProvider{}
	t.gl = &mocks.Logger{}
	t.tp = &mocks.TimeProvider{}

	if len(prep) != 0 {
		prep[0](t)
	}

	t.gl.On("Errorf", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	t.gl.On("Errorf", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	t.gl.On("Section", mock.Anything).Return(t.gl)

	ct, _ := di.NewBuilder()
	ct.Add(
		di.Def{
			Name:  static.DiDiscordSession,
			Build: func(ctn di.Container) (interface{}, error) { return t.s, nil },
		},
		di.Def{
			Name:  static.DiDatabase,
			Build: func(ctn di.Container) (interface{}, error) { return t.db, nil },
		},
		di.Def{
			Name:  static.DiConfig,
			Build: func(ctn di.Container) (interface{}, error) { return t.cfg, nil },
		},
		di.Def{
			Name:  static.DiGuildLog,
			Build: func(ctn di.Container) (interface{}, error) { return t.gl, nil },
		},
		di.Def{
			Name:  static.DiTimeProvider,
			Build: func(ctn di.Container) (interface{}, error) { return t.tp, nil },
		},
	)

	t.ct = ct.Build()

	return t
}

func TestGetEnabled(t *testing.T) {
	m := getVerificationMock(func(m verificationMock) {
		m.db.On("GetGuildVerificationRequired", "guild-enabled").Return(true, nil)
		m.db.On("GetGuildVerificationRequired", "guild-disabled").Return(false, nil)
		m.db.On("GetGuildVerificationRequired", "guild-error").Return(true, errors.New("test error"))
	})

	p := New(m.ct)

	ok, err := p.GetEnabled("guild-enabled")
	assert.Nil(t, err)
	assert.True(t, ok)

	ok, err = p.GetEnabled("guild-disabled")
	assert.Nil(t, err)
	assert.False(t, ok)

	ok, err = p.GetEnabled("guild-error")
	assert.EqualError(t, err, "test error")
	assert.False(t, ok)
}

func TestSetEnabled(t *testing.T) {
	m := getVerificationMock(func(m verificationMock) {
		m.db.On("SetGuildVerificationRequired", mock.AnythingOfType("string"), mock.AnythingOfType("bool")).
			Return(nil)
		m.db.On("GetVerificationQueue", mock.AnythingOfType("string"), mock.AnythingOfType("string")).
			Return([]models.VerificationQueueEntry{
				{GuildID: "guild-id", UserID: "user-0"},
				{GuildID: "guild-id", UserID: "user-1"},
				{GuildID: "guild-id", UserID: "user-2"},
				{GuildID: "guild-id", UserID: "user-left"},
			}, nil)
		m.db.On("RemoveVerificationQueue", mock.AnythingOfType("string"), mock.AnythingOfType("string")).
			Return(true, nil)

		m.s.On("GuildMemberTimeout", mock.AnythingOfType("string"), "user-left", mock.Anything).
			Return(testutil.DiscordRestError(discordgo.ErrCodeUnknownMember))
		m.s.On("GuildMemberTimeout", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("*time.Time")).
			Run(func(args mock.Arguments) {
				// This needs to be checked separately because
				// m.s.AssertCalled(t, "GuildMemberTimeout", "guild-id", "user-2", nil)
				// does not seem to work.
				assert.Nil(t, args[2], nil)
			}).
			Return(nil)
	})

	p := New(m.ct)

	err := p.SetEnabled("guild-id", true)
	assert.Nil(t, err)
	m.db.AssertCalled(t, "SetGuildVerificationRequired", "guild-id", true)

	err = p.SetEnabled("guild-id", false)
	assert.Nil(t, err)
	m.db.AssertCalled(t, "SetGuildVerificationRequired", "guild-id", false)
	m.s.AssertCalled(t, "GuildMemberTimeout", "guild-id", "user-0", mock.AnythingOfType("*time.Time"))
	m.s.AssertCalled(t, "GuildMemberTimeout", "guild-id", "user-1", mock.AnythingOfType("*time.Time"))
	m.s.AssertCalled(t, "GuildMemberTimeout", "guild-id", "user-2", mock.AnythingOfType("*time.Time"))
	m.db.AssertCalled(t, "RemoveVerificationQueue", "guild-id", "user-0")
	m.db.AssertCalled(t, "RemoveVerificationQueue", "guild-id", "user-1")
	m.db.AssertCalled(t, "RemoveVerificationQueue", "guild-id", "user-2")
	m.db.AssertCalled(t, "RemoveVerificationQueue", "guild-id", "user-left")
}

func TestIsVerified(t *testing.T) {
	m := getVerificationMock(func(m verificationMock) {
		m.db.On("GetUserVerified", "user-error").
			Return(false, errors.New("test error"))
		m.db.On("GetUserVerified", "user-notlisted").
			Return(false, database.ErrDatabaseNotFound)
		m.db.On("GetUserVerified", "user-nonverified").
			Return(false, nil)
		m.db.On("GetUserVerified", "user-verified").
			Return(true, nil)
	})

	p := New(m.ct)

	ok, err := p.IsVerified("user-error")
	assert.EqualError(t, err, "test error")
	assert.False(t, ok)

	ok, err = p.IsVerified("user-notlisted")
	assert.Nil(t, err)
	assert.False(t, ok)

	ok, err = p.IsVerified("user-nonverified")
	assert.Nil(t, err)
	assert.False(t, ok)

	ok, err = p.IsVerified("user-verified")
	assert.Nil(t, err)
	assert.True(t, ok)
}

func TestEnqueueVerification(t *testing.T) {
	cfg := &models.Config{}
	cfg.WebServer.PublicAddr = "publicaddr"

	m := getVerificationMock(func(m verificationMock) {
		m.db.On("GetUserVerified", "user-verified").
			Return(true, nil)
		m.db.On("GetUserVerified", mock.AnythingOfType("string")).
			Return(false, nil)
		m.db.On("AddVerificationQueue", mock.Anything).
			Return(nil)
		m.db.On("GetGuildJoinMsg", mock.AnythingOfType("string")).
			Return("joinmsg-chan", nil)

		m.s.On("UserChannelCreate", mock.AnythingOfType("string")).
			Return(&discordgo.Channel{ID: "channel-id"}, nil)
		m.s.On("ChannelMessageSendEmbed", mock.AnythingOfType("string"), mock.AnythingOfType("*discordgo.MessageEmbed")).
			Return(nil, nil)
		m.s.On("ChannelMessageSendComplex", mock.AnythingOfType("string"), mock.AnythingOfType("*discordgo.MessageSend")).
			Return(nil, nil)

		m.cfg.On("Config").Return(cfg)

		m.tp.On("Now").Return(time.Time{})
	})

	p := New(m.ct)

	timeoutTime := time.Time{}.Add(timeout)

	// ----- Non Verified User ------

	m.s.On("GuildMemberTimeout", "guild-id", "user-nonverified", mock.Anything).Once().Return(nil).
		Run(func(args mock.Arguments) {
			assert.Equal(t, *args[2].(*time.Time), timeoutTime)
		})

	err := p.EnqueueVerification(discordgo.Member{
		GuildID: "guild-id",
		User: &discordgo.User{
			ID: "user-nonverified",
		},
	})

	assert.Nil(t, err)
	m.db.AssertCalled(t, "AddVerificationQueue", models.VerificationQueueEntry{
		GuildID:   "guild-id",
		UserID:    "user-nonverified",
		Timestamp: time.Time{},
	})
	m.s.AssertCalled(t, "GuildMemberTimeout", "guild-id", "user-nonverified", mock.Anything)
	m.s.AssertCalled(t, "UserChannelCreate", "user-nonverified")
	m.s.AssertCalled(t, "ChannelMessageSendEmbed", "channel-id", mock.Anything)

	// ----- Verified User ------

	err = p.EnqueueVerification(discordgo.Member{
		GuildID: "guild-id",
		User: &discordgo.User{
			ID: "user-verified",
		},
	})

	assert.Nil(t, err)
	m.db.AssertNotCalled(t, "AddVerificationQueue", models.VerificationQueueEntry{
		GuildID:   "guild-id",
		UserID:    "user-verified",
		Timestamp: time.Time{},
	})
	m.s.AssertNotCalled(t, "GuildMemberTimeout", "guild-id", "user-verified", mock.Anything)
	m.s.AssertNotCalled(t, "GuildMemberTimeout", "guild-id", "user-verified", mock.Anything)
	m.s.AssertNotCalled(t, "UserChannelCreate", "user-verified")

	// ----- Bot User ------

	err = p.EnqueueVerification(discordgo.Member{
		GuildID: "guild-id",
		User: &discordgo.User{
			ID:  "user-bot",
			Bot: true,
		},
	})

	assert.Nil(t, err)
	m.db.AssertNotCalled(t, "AddVerificationQueue", models.VerificationQueueEntry{
		GuildID:   "guild-id",
		UserID:    "user-bot",
		Timestamp: time.Time{},
	})
	m.s.AssertNotCalled(t, "GuildMemberTimeout", "guild-id", "user-bot", mock.Anything)
	m.s.AssertNotCalled(t, "GuildMemberTimeout", "guild-id", "user-bot", mock.Anything)
	m.s.AssertNotCalled(t, "UserChannelCreate", "user-bot")

	// ----- Edge Case: User is nil for some reason ------

	err = p.EnqueueVerification(discordgo.Member{
		GuildID: "guild-id",
	})

	assert.NotNil(t, err)
}

func TestVerify(t *testing.T) {
	m := getVerificationMock(func(m verificationMock) {
		m.db.On("SetUserVerified", mock.AnythingOfType("string"), true).
			Return(nil)
		m.db.On("GetVerificationQueue", mock.AnythingOfType("string"), mock.AnythingOfType("string")).
			Return([]models.VerificationQueueEntry{
				{GuildID: "guild-id", UserID: "user-0"},
				{GuildID: "guild-id", UserID: "user-1"},
				{GuildID: "guild-id", UserID: "user-2"},
				{GuildID: "guild-id", UserID: "user-left"},
			}, nil)
		m.db.On("RemoveVerificationQueue", mock.AnythingOfType("string"), mock.AnythingOfType("string")).
			Return(true, nil)

		m.s.On("GuildMemberTimeout", mock.AnythingOfType("string"), "user-left", mock.Anything).
			Return(&discordgo.RESTError{Message: &discordgo.APIErrorMessage{Code: discordgo.ErrCodeUnknownMember}})
		m.s.On("GuildMemberTimeout", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("*time.Time")).
			Run(func(args mock.Arguments) {
				// This needs to be checked separately because
				// m.s.AssertCalled(t, "GuildMemberTimeout", "guild-id", "user-2", nil)
				// does not seem to work.
				assert.Nil(t, args[2], nil)
			}).
			Return(nil)
	})

	p := New(m.ct)

	err := p.Verify("user-id")
	assert.Nil(t, err)
}

func TestKickRoutine(t *testing.T) {
	testError := errors.New("test error")

	m := getVerificationMock(func(m verificationMock) {
		m.db.On("GetVerificationQueue", mock.AnythingOfType("string"), mock.AnythingOfType("string")).
			Return([]models.VerificationQueueEntry{
				{GuildID: "guild-id", UserID: "user-left", Timestamp: time.Time{}.Add(-timeout - 1*time.Minute)},
				{GuildID: "guild-id", UserID: "user-error", Timestamp: time.Time{}.Add(-timeout - 1*time.Minute)},
				{GuildID: "guild-id", UserID: "user-0", Timestamp: time.Time{}.Add(-timeout - 1*time.Minute)},
				{GuildID: "guild-id", UserID: "user-1", Timestamp: time.Time{}.Add(-timeout)},
				{GuildID: "guild-id", UserID: "user-2", Timestamp: time.Time{}.Add(-1 * time.Minute)},
			}, nil)
		m.db.On("RemoveVerificationQueue", mock.AnythingOfType("string"), mock.AnythingOfType("string")).
			Return(true, nil)

		m.s.On("GuildMemberTimeout", mock.AnythingOfType("string"), "user-left", mock.Anything).
			Return(&discordgo.RESTError{
				Message: &discordgo.APIErrorMessage{
					Code: discordgo.ErrCodeUnknownMember,
				},
				Response:     &http.Response{},
				ResponseBody: []byte{},
			})
		m.s.On("GuildMemberTimeout", mock.AnythingOfType("string"), "user-error", mock.Anything).
			Return(testError)
		m.s.On("GuildMemberTimeout", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("*time.Time")).
			Run(func(args mock.Arguments) {
				// This needs to be checked separately because
				// m.s.AssertCalled(t, "GuildMemberTimeout", "guild-id", "user-2", nil)
				// does not seem to work.
				assert.Nil(t, args[2], nil)
			}).
			Return(nil)
		m.s.On("GuildMemberDelete", mock.AnythingOfType("string"), mock.AnythingOfType("string")).
			Return(nil)

		m.tp.On("Now").Return(time.Time{})
	})

	p := New(m.ct)

	p.KickRoutine()

	m.s.AssertCalled(t, "GuildMemberTimeout", "guild-id", "user-0", mock.Anything)
	m.s.AssertCalled(t, "GuildMemberDelete", "guild-id", "user-0")
	m.db.AssertCalled(t, "RemoveVerificationQueue", "guild-id", "user-0")
	m.gl.AssertNotCalled(t, "Errorf", mock.AnythingOfType("string"), mock.Anything)

	m.s.AssertCalled(t, "GuildMemberTimeout", "guild-id", "user-1", mock.Anything)
	m.s.AssertCalled(t, "GuildMemberDelete", "guild-id", "user-1")
	m.db.AssertCalled(t, "RemoveVerificationQueue", "guild-id", "user-1")
	m.gl.AssertNotCalled(t, "Errorf", mock.AnythingOfType("string"), mock.Anything)

	m.s.AssertNotCalled(t, "GuildMemberTimeout", "guild-id", "user-2", mock.Anything)
	m.s.AssertNotCalled(t, "GuildMemberDelete", "guild-id", "user-2")
	m.db.AssertNotCalled(t, "RemoveVerificationQueue", "guild-id", "user-2")
	m.gl.AssertNotCalled(t, "Errorf", mock.AnythingOfType("string"), mock.Anything)

	m.s.AssertCalled(t, "GuildMemberTimeout", "guild-id", "user-left", mock.Anything)
	m.s.AssertNotCalled(t, "GuildMemberDelete", "guild-id", "user-left")
	m.db.AssertCalled(t, "RemoveVerificationQueue", "guild-id", "user-left")
	m.gl.AssertNotCalled(t, "Errorf", mock.AnythingOfType("string"), mock.Anything)

	m.s.AssertCalled(t, "GuildMemberTimeout", "guild-id", "user-error", mock.Anything)
	m.s.AssertCalled(t, "GuildMemberDelete", "guild-id", "user-error")
	m.db.AssertCalled(t, "RemoveVerificationQueue", "guild-id", "user-error")
	m.gl.AssertCalled(t, "Errorf", "guild-id", mock.AnythingOfType("string"), testError.Error())
}
