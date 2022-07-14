package report

import (
	"errors"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/snowflake"
	"github.com/sarulabs/di/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zekroTJA/shinpuru/internal/models"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/util/snowflakenodes"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/mocks"
)

func init() {
	snowflakenodes.Setup()
}

type reportMock struct {
	s   *mocks.ISession
	db  *mocks.Database
	cfg *mocks.ConfigProvider
	st  *mocks.IState
	tp  *mocks.TimeProvider

	ct di.Container
}

func getReportMock(prep ...func(m reportMock)) reportMock {
	var t reportMock

	t.s = &mocks.ISession{}
	t.db = &mocks.Database{}
	t.cfg = &mocks.ConfigProvider{}
	t.st = &mocks.IState{}
	t.tp = &mocks.TimeProvider{}

	if len(prep) != 0 {
		prep[0](t)
	}

	t.cfg.On("Config").Return(&models.Config{})

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
			Name:  static.DiTimeProvider,
			Build: func(ctn di.Container) (interface{}, error) { return t.tp, nil },
		},
		di.Def{
			Name:  static.DiState,
			Build: func(ctn di.Container) (interface{}, error) { return t.st, nil },
		},
	)

	t.ct = ct.Build()

	return t
}

func TestPushReport(t *testing.T) {
	m := getReportMock(func(m reportMock) {
		m.db.On("AddReport", mock.AnythingOfType("models.Report")).
			Return(nil)
		m.db.On("GetGuildModLog", "guild-nomodlog-1").
			Return("", database.ErrDatabaseNotFound)
		m.db.On("GetGuildModLog", "guild-nomodlog-2").
			Return("", nil)
		m.db.On("GetGuildModLog", mock.AnythingOfType("string")).
			Return("channel-modlog", nil)

		m.s.On("UserChannelCreate", "victim-nodm-1").
			Return(nil, errors.New("test error"))
		m.s.On("UserChannelCreate", "victim-nodm-2").
			Return(nil, nil)
		m.s.On("UserChannelCreate", mock.AnythingOfType("string")).
			Return(&discordgo.Channel{
				ID: "channel-id",
			}, nil)
		m.s.On("ChannelMessageSendEmbed", mock.AnythingOfType("string"), mock.AnythingOfType("*discordgo.MessageEmbed")).
			Return(nil, nil)
	})

	s, err := New(m.ct)
	assert.Nil(t, err)

	// ----- Report Warn Victom with DM and Modlog -----

	rep := models.Report{
		ID:         snowflake.ParseInt64(1), // Hard set ID; expect getting overwritten
		Type:       models.TypeWarn,
		GuildID:    "guild-id",
		VictimID:   "victim-id",
		ExecutorID: "exec-id",
		Msg:        "Some message",
	}
	res, err := s.PushReport(rep)

	assert.Nil(t, err)
	assert.NotEqual(t, rep.ID, res.ID)
	rep.ID = res.ID
	assert.Equal(t, rep, res)
	m.db.AssertCalled(t, "AddReport", rep)
	m.s.AssertCalled(t, "ChannelMessageSendEmbed", "channel-modlog", rep.AsEmbed(""))
	m.s.AssertCalled(t, "UserChannelCreate", "victim-id")
	m.s.AssertCalled(t, "ChannelMessageSendEmbed", "channel-id", rep.AsEmbed(""))

	// ----- Report Warn Victom with DM and NO Modlog -----

	m.s.Calls = nil

	rep = models.Report{
		ID:         snowflake.ParseInt64(1),
		Type:       models.TypeWarn,
		GuildID:    "guild-nomodlog-1",
		VictimID:   "victim-id",
		ExecutorID: "exec-id",
		Msg:        "Some message",
	}
	res, err = s.PushReport(rep)

	assert.Nil(t, err)
	assert.NotEqual(t, rep.ID, res.ID)
	rep.ID = res.ID
	assert.Equal(t, rep, res)
	m.db.AssertCalled(t, "AddReport", rep)
	m.s.AssertNotCalled(t, "ChannelMessageSendEmbed", "channel-modlog", mock.Anything)
	m.s.AssertNotCalled(t, "ChannelMessageSendEmbed", "", mock.Anything)
	m.s.AssertCalled(t, "UserChannelCreate", "victim-id")
	m.s.AssertCalled(t, "ChannelMessageSendEmbed", "channel-id", rep.AsEmbed(""))

	m.s.Calls = nil

	rep = models.Report{
		ID:         snowflake.ParseInt64(1),
		Type:       models.TypeWarn,
		GuildID:    "guild-nomodlog-2",
		VictimID:   "victim-id",
		ExecutorID: "exec-id",
		Msg:        "Some message",
	}
	res, err = s.PushReport(rep)

	assert.Nil(t, err)
	assert.NotEqual(t, rep.ID, res.ID)
	rep.ID = res.ID
	assert.Equal(t, rep, res)
	m.db.AssertCalled(t, "AddReport", rep)
	m.s.AssertNotCalled(t, "ChannelMessageSendEmbed", "channel-modlog", mock.Anything)
	m.s.AssertNotCalled(t, "ChannelMessageSendEmbed", "", mock.Anything)
	m.s.AssertCalled(t, "UserChannelCreate", "victim-id")
	m.s.AssertCalled(t, "ChannelMessageSendEmbed", "channel-id", rep.AsEmbed(""))

	// ----- Report Warn Victom with NO DM and Modlog -----

	m.s.Calls = nil

	rep = models.Report{
		ID:         snowflake.ParseInt64(1),
		Type:       models.TypeWarn,
		GuildID:    "guild-1",
		VictimID:   "victim-nodm-1",
		ExecutorID: "exec-id",
		Msg:        "Some message",
	}
	res, err = s.PushReport(rep)

	assert.Nil(t, err)
	assert.NotEqual(t, rep.ID, res.ID)
	rep.ID = res.ID
	assert.Equal(t, rep, res)
	m.db.AssertCalled(t, "AddReport", rep)
	m.s.AssertCalled(t, "ChannelMessageSendEmbed", "channel-modlog", mock.Anything)
	m.s.AssertCalled(t, "UserChannelCreate", "victim-nodm-1")
	m.s.AssertNotCalled(t, "ChannelMessageSendEmbed", "channel-id", mock.Anything)

	m.s.Calls = nil

	rep = models.Report{
		ID:         snowflake.ParseInt64(1),
		Type:       models.TypeWarn,
		GuildID:    "guild-1",
		VictimID:   "victim-nodm-2",
		ExecutorID: "exec-id",
		Msg:        "Some message",
	}
	res, err = s.PushReport(rep)

	assert.Nil(t, err)
	assert.NotEqual(t, rep.ID, res.ID)
	rep.ID = res.ID
	assert.Equal(t, rep, res)
	m.db.AssertCalled(t, "AddReport", rep)
	m.s.AssertCalled(t, "ChannelMessageSendEmbed", "channel-modlog", mock.Anything)
	m.s.AssertCalled(t, "UserChannelCreate", "victim-nodm-2")
	m.s.AssertNotCalled(t, "ChannelMessageSendEmbed", "channel-id", mock.Anything)
}
