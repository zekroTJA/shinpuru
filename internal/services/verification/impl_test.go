package verification

import (
	"errors"
	"testing"

	"github.com/sarulabs/di/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zekroTJA/shinpuru/internal/models"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/mocks"
)

type verificationMock struct {
	s   *mocks.ISession
	db  *mocks.Database
	cfg *mocks.ConfigProvider
	gl  *mocks.Logger

	ct di.Container
}

func getVerificationMock(prep ...func(m verificationMock)) verificationMock {
	var t verificationMock

	t.s = &mocks.ISession{}
	t.db = &mocks.Database{}
	t.cfg = &mocks.ConfigProvider{}
	t.gl = &mocks.Logger{}

	if len(prep) != 0 {
		prep[0](t)
	}

	t.gl.On("Errorf", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
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
			Return([]*models.VerificationQueueEntry{
				{GuildID: "guild-id", UserID: "user-0"},
				{GuildID: "guild-id", UserID: "user-1"},
				{GuildID: "guild-id", UserID: "user-2"},
			}, nil)
		m.db.On("RemoveVerificationQueue", mock.AnythingOfType("string"), mock.AnythingOfType("string")).
			Return(true, nil)

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
}
