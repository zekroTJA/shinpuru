package report

import (
	"errors"
	"testing"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/snowflake"
	"github.com/sarulabs/di/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zekroTJA/shinpuru/internal/models"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/util/snowflakenodes"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/internal/util/testutil"
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

func (m *reportMock) Reset() {
	m.s.Calls = nil
	m.db.Calls = nil
	m.cfg.Calls = nil
	m.st.Calls = nil
	m.tp.Calls = nil
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
	t.tp.On("Now").Return(time.Time{})

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

func TestPushKick(t *testing.T) {
	m := getReportMock(func(m reportMock) {
		m.db.On("AddReport", mock.AnythingOfType("models.Report")).
			Return(nil)
		m.db.On("GetGuildModLog", mock.AnythingOfType("string")).
			Return("channel-modlog", nil)
		m.db.On("DeleteReport", mock.Anything).Return(nil)

		m.s.On("UserChannelCreate", mock.AnythingOfType("string")).
			Return(&discordgo.Channel{
				ID: "channel-id",
			}, nil)
		m.s.On("ChannelMessageSendEmbed", mock.AnythingOfType("string"), mock.AnythingOfType("*discordgo.MessageEmbed")).
			Return(nil, nil)

		m.st.On("Guild", mock.AnythingOfType("string"), mock.AnythingOfType("bool")).
			Return(&discordgo.Guild{
				ID: "guild-id",
				Roles: []*discordgo.Role{
					{ID: "role-admin", Position: 0, Permissions: 0x8},
					{ID: "role-0", Position: 0},
					{ID: "role-1", Position: 1},
				},
			}, nil)
	})

	s, err := New(m.ct)
	assert.Nil(t, err)

	// ----- Positive Test -----

	m.st.On("Member", "guild-id", "victim-id").
		Once().
		Return(&discordgo.Member{
			User: &discordgo.User{
				ID: "victim-id",
			},
			Roles: []string{"role-0"},
		}, nil)

	m.st.On("Member", "guild-id", "executor-id").
		Once().
		Return(&discordgo.Member{
			User: &discordgo.User{
				ID: "executor-id",
			},
			Roles: []string{"role-1"},
		}, nil)

	m.s.On("GuildMemberDeleteWithReason", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Once().
		Return(nil)

	rep := models.Report{
		ID:         snowflake.ParseInt64(1),
		Type:       69,
		VictimID:   "victim-id",
		ExecutorID: "executor-id",
		GuildID:    "guild-id",
		Msg:        "Some reason",
	}
	res, err := s.PushKick(rep)
	assert.Nil(t, err)
	assert.NotEqual(t, rep.ID, res.ID)
	assert.Equal(t, res.Type, models.TypeKick)
	rep.ID = res.ID
	rep.Type = res.Type
	assert.Equal(t, rep, res)
	m.db.AssertCalled(t, "AddReport", rep)
	m.s.AssertCalled(t, "GuildMemberDeleteWithReason", "guild-id", "victim-id", mock.AnythingOfType("string"))

	// ----- Negative Test: Victim and Reporter have same role -----

	m.Reset()

	m.st.On("Member", "guild-id", "victim-id").
		Once().
		Return(&discordgo.Member{
			User: &discordgo.User{
				ID: "victim-id",
			},
			Roles: []string{"role-1"},
		}, nil)

	m.st.On("Member", "guild-id", "executor-id").
		Once().
		Return(&discordgo.Member{
			User: &discordgo.User{
				ID: "executor-id",
			},
			Roles: []string{"role-1"},
		}, nil)

	rep = models.Report{
		ID:         snowflake.ParseInt64(1),
		Type:       69,
		VictimID:   "victim-id",
		ExecutorID: "executor-id",
		GuildID:    "guild-id",
		Msg:        "Some reason",
	}
	res, err = s.PushKick(rep)
	assert.EqualError(t, err, ErrRoleDiff.Error())
	m.db.AssertNotCalled(t, "AddReport", mock.Anything)
	m.s.AssertNotCalled(t, "GuildMemberDeleteWithReason", "guild-id", "victim-id", mock.AnythingOfType("string"))

	// ----- Negative Test: Victim has higer role than executor -----

	m.Reset()

	m.st.On("Member", "guild-id", "victim-id").
		Once().
		Return(&discordgo.Member{
			User: &discordgo.User{
				ID: "victim-id",
			},
			Roles: []string{"role-1"},
		}, nil)

	m.st.On("Member", "guild-id", "executor-id").
		Once().
		Return(&discordgo.Member{
			User: &discordgo.User{
				ID: "executor-id",
			},
			Roles: []string{"role-0"},
		}, nil)

	rep = models.Report{
		ID:         snowflake.ParseInt64(1),
		Type:       69,
		VictimID:   "victim-id",
		ExecutorID: "executor-id",
		GuildID:    "guild-id",
		Msg:        "Some reason",
	}
	res, err = s.PushKick(rep)
	assert.EqualError(t, err, ErrRoleDiff.Error())
	m.db.AssertNotCalled(t, "AddReport", mock.Anything)
	m.s.AssertNotCalled(t, "GuildMemberDeleteWithReason", "guild-id", "victim-id", mock.AnythingOfType("string"))

	// ----- Positive Test: Victim has higer role than executor but executor is Admin -----

	m.Reset()

	m.st.On("Member", "guild-id", "victim-id").
		Once().
		Return(&discordgo.Member{
			User: &discordgo.User{
				ID: "victim-id",
			},
			Roles: []string{"role-1"},
		}, nil)

	m.st.On("Member", "guild-id", "executor-id").
		Once().
		Return(&discordgo.Member{
			User: &discordgo.User{
				ID: "executor-id",
			},
			Roles: []string{"role-admin"},
		}, nil)

	m.s.On("GuildMemberDeleteWithReason", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Once().
		Return(nil)

	rep = models.Report{
		ID:         snowflake.ParseInt64(1),
		Type:       69,
		VictimID:   "victim-id",
		ExecutorID: "executor-id",
		GuildID:    "guild-id",
		Msg:        "Some reason",
	}
	res, err = s.PushKick(rep)
	assert.Nil(t, err)
	assert.NotEqual(t, rep.ID, res.ID)
	assert.Equal(t, res.Type, models.TypeKick)
	rep.ID = res.ID
	rep.Type = res.Type
	assert.Equal(t, rep, res)
	m.db.AssertCalled(t, "AddReport", rep)
	m.s.AssertCalled(t, "GuildMemberDeleteWithReason", "guild-id", "victim-id", mock.AnythingOfType("string"))

	// ----- Negative Test: Victim has left -----

	m.Reset()

	m.st.On("Member", "guild-id", "victim-id").
		Once().
		Return(nil, testutil.DiscordRestError(discordgo.ErrCodeUnknownMember))

	m.st.On("Member", "guild-id", "executor-id").
		Once().
		Return(&discordgo.Member{
			User: &discordgo.User{
				ID: "victim-id",
			},
			Roles: []string{"role-1"},
		}, nil)

	rep = models.Report{
		ID:         snowflake.ParseInt64(1),
		Type:       69,
		VictimID:   "victim-id",
		ExecutorID: "executor-id",
		GuildID:    "guild-id",
		Msg:        "Some reason",
	}
	res, err = s.PushKick(rep)
	assert.EqualError(t, err, ErrMemberHasLeft.Error())
	m.db.AssertNotCalled(t, "AddReport", mock.Anything)
	m.s.AssertNotCalled(t, "GuildMemberDeleteWithReason", "guild-id", "victim-id", mock.AnythingOfType("string"))
}

func TestPushBan(t *testing.T) {
	m := getReportMock(func(m reportMock) {
		m.db.On("AddReport", mock.AnythingOfType("models.Report")).
			Return(nil)
		m.db.On("GetGuildModLog", mock.AnythingOfType("string")).
			Return("channel-modlog", nil)
		m.db.On("DeleteReport", mock.Anything).Return(nil)

		m.s.On("UserChannelCreate", mock.AnythingOfType("string")).
			Return(&discordgo.Channel{
				ID: "channel-id",
			}, nil)
		m.s.On("ChannelMessageSendEmbed", mock.AnythingOfType("string"), mock.AnythingOfType("*discordgo.MessageEmbed")).
			Return(nil, nil)

		m.st.On("Guild", mock.AnythingOfType("string"), mock.AnythingOfType("bool")).
			Return(&discordgo.Guild{
				ID: "guild-id",
				Roles: []*discordgo.Role{
					{ID: "role-admin", Position: 0, Permissions: 0x8},
					{ID: "role-0", Position: 0},
					{ID: "role-1", Position: 1},
				},
			}, nil)
	})

	s, err := New(m.ct)
	assert.Nil(t, err)

	// ----- Positive Test -----

	m.st.On("Member", "guild-id", "victim-id").
		Once().
		Return(&discordgo.Member{
			User: &discordgo.User{
				ID: "victim-id",
			},
			Roles: []string{"role-0"},
		}, nil)

	m.st.On("Member", "guild-id", "executor-id").
		Once().
		Return(&discordgo.Member{
			User: &discordgo.User{
				ID: "executor-id",
			},
			Roles: []string{"role-1"},
		}, nil)

	m.s.On("GuildBanCreateWithReason", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int")).
		Once().
		Return(nil)

	rep := models.Report{
		ID:         snowflake.ParseInt64(1),
		Type:       69,
		VictimID:   "victim-id",
		ExecutorID: "executor-id",
		GuildID:    "guild-id",
		Msg:        "Some reason",
	}
	res, err := s.PushBan(rep)
	assert.Nil(t, err)
	assert.NotEqual(t, rep.ID, res.ID)
	assert.Equal(t, res.Type, models.TypeBan)
	rep.ID = res.ID
	rep.Type = res.Type
	assert.Equal(t, rep, res)
	m.db.AssertCalled(t, "AddReport", rep)
	m.s.AssertCalled(t, "GuildBanCreateWithReason", "guild-id", "victim-id", mock.AnythingOfType("string"), mock.AnythingOfType("int"))

	// ----- Positive Test with Timeout -----

	m.st.On("Member", "guild-id", "victim-id").
		Once().
		Return(&discordgo.Member{
			User: &discordgo.User{
				ID: "victim-id",
			},
			Roles: []string{"role-0"},
		}, nil)

	m.st.On("Member", "guild-id", "executor-id").
		Once().
		Return(&discordgo.Member{
			User: &discordgo.User{
				ID: "executor-id",
			},
			Roles: []string{"role-1"},
		}, nil)

	m.s.On("GuildBanCreateWithReason", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int")).
		Once().
		Return(nil)

	timeout := time.Time{}.Add(24 * time.Hour)
	rep = models.Report{
		ID:         snowflake.ParseInt64(1),
		Type:       69,
		VictimID:   "victim-id",
		ExecutorID: "executor-id",
		GuildID:    "guild-id",
		Msg:        "Some reason",
		Timeout:    &timeout,
	}
	res, err = s.PushBan(rep)
	assert.Nil(t, err)
	assert.NotEqual(t, rep.ID, res.ID)
	assert.Equal(t, res.Type, models.TypeBan)
	rep.ID = res.ID
	rep.Type = res.Type
	assert.Equal(t, rep, res)
	m.db.AssertCalled(t, "AddReport", rep)
	m.s.AssertCalled(t, "GuildBanCreateWithReason", "guild-id", "victim-id", mock.AnythingOfType("string"), mock.AnythingOfType("int"))

	// ----- Negative Test: Invalid Timeout -----

	m.Reset()

	timeout = time.Time{}.Add(-24 * time.Hour)
	rep = models.Report{
		ID:         snowflake.ParseInt64(1),
		Type:       69,
		VictimID:   "victim-id",
		ExecutorID: "executor-id",
		GuildID:    "guild-id",
		Msg:        "Some reason",
		Timeout:    &timeout,
	}
	res, err = s.PushBan(rep)
	assert.EqualError(t, err, ErrInvalidTimeout.Error())
	m.db.AssertNotCalled(t, "AddReport", mock.Anything)
	m.s.AssertNotCalled(t, "GuildBanCreateWithReason", "guild-id", "victim-id", mock.AnythingOfType("string"), mock.AnythingOfType("int"))

	// ----- Negative Test: Victim and Reporter have same role -----

	m.Reset()

	m.st.On("Member", "guild-id", "victim-id").
		Once().
		Return(&discordgo.Member{
			User: &discordgo.User{
				ID: "victim-id",
			},
			Roles: []string{"role-1"},
		}, nil)

	m.st.On("Member", "guild-id", "executor-id").
		Once().
		Return(&discordgo.Member{
			User: &discordgo.User{
				ID: "executor-id",
			},
			Roles: []string{"role-1"},
		}, nil)

	rep = models.Report{
		ID:         snowflake.ParseInt64(1),
		Type:       69,
		VictimID:   "victim-id",
		ExecutorID: "executor-id",
		GuildID:    "guild-id",
		Msg:        "Some reason",
	}
	res, err = s.PushBan(rep)
	assert.EqualError(t, err, ErrRoleDiff.Error())
	m.db.AssertNotCalled(t, "AddReport", mock.Anything)
	m.s.AssertNotCalled(t, "GuildBanCreateWithReason", "guild-id", "victim-id", mock.AnythingOfType("string"), mock.AnythingOfType("int"))

	// ----- Negative Test: Victim has higer role than executor -----

	m.Reset()

	m.st.On("Member", "guild-id", "victim-id").
		Once().
		Return(&discordgo.Member{
			User: &discordgo.User{
				ID: "victim-id",
			},
			Roles: []string{"role-1"},
		}, nil)

	m.st.On("Member", "guild-id", "executor-id").
		Once().
		Return(&discordgo.Member{
			User: &discordgo.User{
				ID: "executor-id",
			},
			Roles: []string{"role-0"},
		}, nil)

	rep = models.Report{
		ID:         snowflake.ParseInt64(1),
		Type:       69,
		VictimID:   "victim-id",
		ExecutorID: "executor-id",
		GuildID:    "guild-id",
		Msg:        "Some reason",
	}
	res, err = s.PushBan(rep)
	assert.EqualError(t, err, ErrRoleDiff.Error())
	m.db.AssertNotCalled(t, "AddReport", mock.Anything)
	m.s.AssertNotCalled(t, "GuildBanCreateWithReason", "guild-id", "victim-id", mock.AnythingOfType("string"), mock.AnythingOfType("int"))

	// ----- Positive Test: Victim has higer role than executor but executor is Admin -----

	m.Reset()

	m.st.On("Member", "guild-id", "victim-id").
		Once().
		Return(&discordgo.Member{
			User: &discordgo.User{
				ID: "victim-id",
			},
			Roles: []string{"role-1"},
		}, nil)

	m.st.On("Member", "guild-id", "executor-id").
		Once().
		Return(&discordgo.Member{
			User: &discordgo.User{
				ID: "executor-id",
			},
			Roles: []string{"role-admin"},
		}, nil)

	m.s.On("GuildBanCreateWithReason", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int")).
		Once().
		Return(nil)

	rep = models.Report{
		ID:         snowflake.ParseInt64(1),
		Type:       69,
		VictimID:   "victim-id",
		ExecutorID: "executor-id",
		GuildID:    "guild-id",
		Msg:        "Some reason",
	}
	res, err = s.PushBan(rep)
	assert.Nil(t, err)
	assert.NotEqual(t, rep.ID, res.ID)
	assert.Equal(t, res.Type, models.TypeBan)
	rep.ID = res.ID
	rep.Type = res.Type
	assert.Equal(t, rep, res)
	m.db.AssertCalled(t, "AddReport", rep)
	m.s.AssertCalled(t, "GuildBanCreateWithReason", "guild-id", "victim-id", mock.AnythingOfType("string"), mock.AnythingOfType("int"))

	// ----- Positive Test: Anonymous Report -----

	m.Reset()

	m.st.On("Member", "guild-id", "executor-id").
		Once().
		Return(&discordgo.Member{
			User: &discordgo.User{
				ID: "executor-id",
			},
			Roles: []string{"role-admin"},
		}, nil)

	m.s.On("GuildBanCreateWithReason", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int")).
		Once().
		Return(nil)

	rep = models.Report{
		ID:         snowflake.ParseInt64(1),
		Type:       69,
		VictimID:   "victim-id",
		ExecutorID: "executor-id",
		GuildID:    "guild-id",
		Msg:        "Some reason",
		Anonymous:  true,
	}
	res, err = s.PushBan(rep)
	assert.Nil(t, err)
	assert.NotEqual(t, rep.ID, res.ID)
	assert.Equal(t, res.Type, models.TypeBan)
	assert.True(t, res.Anonymous)
	rep.ID = res.ID
	rep.Type = res.Type
	rep.Anonymous = res.Anonymous
	assert.Equal(t, rep, res)
	m.db.AssertCalled(t, "AddReport", rep)
	m.s.AssertCalled(t, "GuildBanCreateWithReason", "guild-id", "victim-id", mock.AnythingOfType("string"), mock.AnythingOfType("int"))

	// ----- Positive Test: Implicitely Anonymous Report (See issue #378) -----

	m.Reset()

	m.st.On("Member", "guild-id", "executor-id").
		Once().
		Return(&discordgo.Member{
			User: &discordgo.User{
				ID: "executor-id",
			},
			Roles: []string{"role-admin"},
		}, nil)

	m.st.On("Member", "guild-id", "victim-id").
		Once().
		Return(nil, testutil.DiscordRestError(discordgo.ErrCodeUnknownMember))

	m.s.On("GuildBanCreateWithReason", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int")).
		Once().
		Return(nil)

	rep = models.Report{
		ID:         snowflake.ParseInt64(1),
		Type:       69,
		VictimID:   "victim-id",
		ExecutorID: "executor-id",
		GuildID:    "guild-id",
		Msg:        "Some reason",
	}
	res, err = s.PushBan(rep)
	assert.Nil(t, err)
	assert.NotEqual(t, rep.ID, res.ID)
	assert.Equal(t, res.Type, models.TypeBan)
	rep.ID = res.ID
	rep.Type = res.Type
	rep.Anonymous = res.Anonymous
	assert.Equal(t, rep, res)
	m.db.AssertCalled(t, "AddReport", rep)
	m.s.AssertCalled(t, "GuildBanCreateWithReason", "guild-id", "victim-id", mock.AnythingOfType("string"), mock.AnythingOfType("int"))

	// ----- Negative Test: Ban Process Failed -----

	m.Reset()

	m.db.On("DeleteReport").Once().Return(nil)

	m.st.On("Member", "guild-id", "victim-id").
		Once().
		Return(&discordgo.Member{
			User: &discordgo.User{
				ID: "victim-id",
			},
			Roles: []string{"role-0"},
		}, nil)

	m.st.On("Member", "guild-id", "executor-id").
		Once().
		Return(&discordgo.Member{
			User: &discordgo.User{
				ID: "executor-id",
			},
			Roles: []string{"role-1"},
		}, nil)

	m.s.On("GuildBanCreateWithReason", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int")).
		Once().
		Return(errors.New("test error"))

	rep = models.Report{
		ID:         snowflake.ParseInt64(1),
		Type:       69,
		VictimID:   "victim-id",
		ExecutorID: "executor-id",
		GuildID:    "guild-id",
		Msg:        "Some reason",
	}
	res, err = s.PushBan(rep)
	assert.EqualError(t, err, "test error")
	m.db.AssertCalled(t, "AddReport", mock.Anything)
	m.s.AssertCalled(t, "GuildBanCreateWithReason", "guild-id", "victim-id", mock.AnythingOfType("string"), mock.AnythingOfType("int"))
	m.db.AssertCalled(t, "DeleteReport", mock.Anything)
}

func TestPushMute(t *testing.T) {
	m := getReportMock(func(m reportMock) {
		m.db.On("AddReport", mock.AnythingOfType("models.Report")).
			Return(nil)
		m.db.On("GetGuildModLog", mock.AnythingOfType("string")).
			Return("channel-modlog", nil)
		m.db.On("DeleteReport", mock.Anything).Return(nil)

		m.s.On("UserChannelCreate", mock.AnythingOfType("string")).
			Return(&discordgo.Channel{
				ID: "channel-id",
			}, nil)
		m.s.On("ChannelMessageSendEmbed", mock.AnythingOfType("string"), mock.AnythingOfType("*discordgo.MessageEmbed")).
			Return(nil, nil)

		m.st.On("Guild", mock.AnythingOfType("string"), mock.AnythingOfType("bool")).
			Return(&discordgo.Guild{
				ID: "guild-id",
				Roles: []*discordgo.Role{
					{ID: "role-admin", Position: 0, Permissions: 0x8},
					{ID: "role-0", Position: 0},
					{ID: "role-1", Position: 1},
				},
			}, nil)
	})

	s, err := New(m.ct)
	assert.Nil(t, err)

	// ----- Positive Test -----

	m.st.On("Member", "guild-id", "victim-id").
		Once().
		Return(&discordgo.Member{
			User: &discordgo.User{
				ID: "victim-id",
			},
			Roles: []string{"role-0"},
		}, nil)

	m.st.On("Member", "guild-id", "executor-id").
		Once().
		Return(&discordgo.Member{
			User: &discordgo.User{
				ID: "executor-id",
			},
			Roles: []string{"role-1"},
		}, nil)

	m.s.On("GuildMemberTimeout", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("*time.Time")).
		Once().
		Return(nil)

	rep := models.Report{
		ID:         snowflake.ParseInt64(1),
		Type:       69,
		VictimID:   "victim-id",
		ExecutorID: "executor-id",
		GuildID:    "guild-id",
		Msg:        "Some reason",
	}
	res, err := s.PushMute(rep)
	assert.Nil(t, err)
	assert.NotEqual(t, rep.ID, res.ID)
	assert.Equal(t, res.Type, models.TypeMute)
	rep.ID = res.ID
	rep.Type = res.Type
	assert.Equal(t, rep, res)
	m.db.AssertCalled(t, "AddReport", rep)
	m.s.AssertCalled(t, "GuildMemberTimeout", "guild-id", "victim-id", mock.AnythingOfType("*time.Time"))

	// ----- Negative Test: Invalid Timeout -----

	m.Reset()

	timeout := time.Time{}.Add(-24 * time.Hour)
	rep = models.Report{
		ID:         snowflake.ParseInt64(1),
		Type:       69,
		VictimID:   "victim-id",
		ExecutorID: "executor-id",
		GuildID:    "guild-id",
		Msg:        "Some reason",
		Timeout:    &timeout,
	}
	res, err = s.PushMute(rep)
	assert.EqualError(t, err, ErrInvalidTimeout.Error())
	m.db.AssertNotCalled(t, "AddReport", mock.Anything)
	m.s.AssertNotCalled(t, "GuildMemberTimeout", "guild-id", "victim-id", mock.AnythingOfType("*time.Time"))

	// ----- Negative Test: Victim and Reporter have same role -----

	m.Reset()

	m.st.On("Member", "guild-id", "victim-id").
		Once().
		Return(&discordgo.Member{
			User: &discordgo.User{
				ID: "victim-id",
			},
			Roles: []string{"role-1"},
		}, nil)

	m.st.On("Member", "guild-id", "executor-id").
		Once().
		Return(&discordgo.Member{
			User: &discordgo.User{
				ID: "executor-id",
			},
			Roles: []string{"role-1"},
		}, nil)

	rep = models.Report{
		ID:         snowflake.ParseInt64(1),
		Type:       69,
		VictimID:   "victim-id",
		ExecutorID: "executor-id",
		GuildID:    "guild-id",
		Msg:        "Some reason",
	}
	res, err = s.PushMute(rep)
	assert.EqualError(t, err, ErrRoleDiff.Error())
	m.db.AssertNotCalled(t, "AddReport", mock.Anything)
	m.s.AssertNotCalled(t, "GuildMemberTimeout", "guild-id", "victim-id", mock.AnythingOfType("*time.Time"))

	// ----- Negative Test: Victim has higer role than executor -----

	m.Reset()

	m.st.On("Member", "guild-id", "victim-id").
		Once().
		Return(&discordgo.Member{
			User: &discordgo.User{
				ID: "victim-id",
			},
			Roles: []string{"role-1"},
		}, nil)

	m.st.On("Member", "guild-id", "executor-id").
		Once().
		Return(&discordgo.Member{
			User: &discordgo.User{
				ID: "executor-id",
			},
			Roles: []string{"role-0"},
		}, nil)

	rep = models.Report{
		ID:         snowflake.ParseInt64(1),
		Type:       69,
		VictimID:   "victim-id",
		ExecutorID: "executor-id",
		GuildID:    "guild-id",
		Msg:        "Some reason",
	}
	res, err = s.PushMute(rep)
	assert.EqualError(t, err, ErrRoleDiff.Error())
	m.db.AssertNotCalled(t, "AddReport", mock.Anything)
	m.s.AssertNotCalled(t, "GuildMemberTimeout", "guild-id", "victim-id", mock.AnythingOfType("*time.Time"))

	// ----- Positive Test: Victim has higer role than executor but executor is Admin -----

	m.Reset()

	m.st.On("Member", "guild-id", "victim-id").
		Once().
		Return(&discordgo.Member{
			User: &discordgo.User{
				ID: "victim-id",
			},
			Roles: []string{"role-1"},
		}, nil)

	m.st.On("Member", "guild-id", "executor-id").
		Once().
		Return(&discordgo.Member{
			User: &discordgo.User{
				ID: "executor-id",
			},
			Roles: []string{"role-admin"},
		}, nil)

	m.s.On("GuildMemberTimeout", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("*time.Time")).
		Once().
		Return(nil)

	rep = models.Report{
		ID:         snowflake.ParseInt64(1),
		Type:       69,
		VictimID:   "victim-id",
		ExecutorID: "executor-id",
		GuildID:    "guild-id",
		Msg:        "Some reason",
	}
	res, err = s.PushMute(rep)
	assert.Nil(t, err)
	assert.NotEqual(t, rep.ID, res.ID)
	assert.Equal(t, res.Type, models.TypeMute)
	rep.ID = res.ID
	rep.Type = res.Type
	assert.Equal(t, rep, res)
	m.db.AssertCalled(t, "AddReport", rep)
	m.s.AssertCalled(t, "GuildMemberTimeout", "guild-id", "victim-id", mock.AnythingOfType("*time.Time"))

	// ----- Positive Test: No Reason -----

	m.Reset()

	m.st.On("Member", "guild-id", "victim-id").
		Once().
		Return(&discordgo.Member{
			User: &discordgo.User{
				ID: "victim-id",
			},
			Roles: []string{"role-1"},
		}, nil)

	m.st.On("Member", "guild-id", "executor-id").
		Once().
		Return(&discordgo.Member{
			User: &discordgo.User{
				ID: "executor-id",
			},
			Roles: []string{"role-admin"},
		}, nil)

	m.s.On("GuildMemberTimeout", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("*time.Time")).
		Once().
		Return(nil)

	rep = models.Report{
		ID:         snowflake.ParseInt64(1),
		Type:       69,
		VictimID:   "victim-id",
		ExecutorID: "executor-id",
		GuildID:    "guild-id",
		Msg:        "",
	}
	res, err = s.PushMute(rep)
	assert.Nil(t, err)
	assert.NotEqual(t, rep.ID, res.ID)
	assert.Equal(t, res.Type, models.TypeMute)
	rep.ID = res.ID
	rep.Type = res.Type
	rep.Msg = "no reason specified"
	assert.Equal(t, rep, res)
	m.db.AssertCalled(t, "AddReport", rep)
	m.s.AssertCalled(t, "GuildMemberTimeout", "guild-id", "victim-id", testutil.Nil[time.Time]())
}

func TestRevokeMute(t *testing.T) {
	m := getReportMock(func(m reportMock) {
		m.db.On("GetGuildModLog", mock.AnythingOfType("string")).
			Return("channel-modlog", nil)
		m.db.On("DeleteReport", mock.Anything).Return(nil)

		m.s.On("UserChannelCreate", mock.AnythingOfType("string")).
			Return(&discordgo.Channel{
				ID: "channel-id",
			}, nil)
		m.s.On("ChannelMessageSendEmbed", mock.AnythingOfType("string"), mock.AnythingOfType("*discordgo.MessageEmbed")).
			Return(nil, nil)

		m.st.On("Guild", mock.AnythingOfType("string"), mock.AnythingOfType("bool")).
			Return(&discordgo.Guild{
				ID: "guild-id",
				Roles: []*discordgo.Role{
					{ID: "role-admin", Position: 0, Permissions: 0x8},
					{ID: "role-0", Position: 0},
					{ID: "role-1", Position: 1},
				},
			}, nil)
	})

	// ----- Positive Test -----

	m.st.On("Member", "guild-id", "victim-id").
		Once().
		Return(&discordgo.Member{
			User: &discordgo.User{
				ID: "victim-id",
			},
			Roles: []string{"role-0"},
		}, nil)

	m.st.On("Member", "guild-id", "executor-id").
		Once().
		Return(&discordgo.Member{
			User: &discordgo.User{
				ID: "executor-id",
			},
			Roles: []string{"role-1"},
		}, nil)

	m.s.On("GuildMemberTimeout", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("*time.Time")).
		Once().
		Return(nil)

	m.db.On("GetReportsFiltered", "guild-id", "victim-id", models.TypeMute, mock.AnythingOfType("int"), mock.AnythingOfType("int")).
		Once().
		Return([]models.Report{
			{
				ID:      snowflake.ParseInt64(123),
				Timeout: &time.Time{},
			},
		}, nil)

	m.db.On("ExpireReports", mock.AnythingOfType("string")).
		Once().
		Return(nil)

	s, err := New(m.ct)
	assert.Nil(t, err)

	emb, err := s.RevokeMute("guild-id", "executor-id", "victim-id", "")
	assert.Nil(t, err)
	m.s.AssertCalled(t, "GuildMemberTimeout", "guild-id", "victim-id", testutil.Nil[time.Time]())
	m.s.AssertCalled(t, "ChannelMessageSendEmbed", "channel-modlog", emb)
	m.db.AssertCalled(t, "ExpireReports", "123")

	// ----- Positive Test: Admin -----

	m.st.On("Member", "guild-id", "victim-id").
		Once().
		Return(&discordgo.Member{
			User: &discordgo.User{
				ID: "victim-id",
			},
			Roles: []string{"role-0"},
		}, nil)

	m.st.On("Member", "guild-id", "executor-id").
		Once().
		Return(&discordgo.Member{
			User: &discordgo.User{
				ID: "executor-id",
			},
			Roles: []string{"role-admin"},
		}, nil)

	m.s.On("GuildMemberTimeout", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("*time.Time")).
		Once().
		Return(nil)

	m.db.On("GetReportsFiltered", "guild-id", "victim-id", models.TypeMute, mock.AnythingOfType("int"), mock.AnythingOfType("int")).
		Once().
		Return([]models.Report{
			{
				ID:      snowflake.ParseInt64(123),
				Timeout: &time.Time{},
			},
		}, nil)

	m.db.On("ExpireReports", mock.AnythingOfType("string")).
		Once().
		Return(nil)

	s, err = New(m.ct)
	assert.Nil(t, err)

	emb, err = s.RevokeMute("guild-id", "executor-id", "victim-id", "")
	assert.Nil(t, err)
	m.s.AssertCalled(t, "GuildMemberTimeout", "guild-id", "victim-id", testutil.Nil[time.Time]())
	m.s.AssertCalled(t, "ChannelMessageSendEmbed", "channel-modlog", emb)
	m.db.AssertCalled(t, "ExpireReports", "123")

	// ----- Positive Test: No prior Reports -----

	m.Reset()

	m.st.On("Member", "guild-id", "victim-id").
		Once().
		Return(&discordgo.Member{
			User: &discordgo.User{
				ID: "victim-id",
			},
			Roles: []string{"role-0"},
		}, nil)

	m.st.On("Member", "guild-id", "executor-id").
		Once().
		Return(&discordgo.Member{
			User: &discordgo.User{
				ID: "executor-id",
			},
			Roles: []string{"role-1"},
		}, nil)

	m.s.On("GuildMemberTimeout", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("*time.Time")).
		Once().
		Return(nil)

	m.db.On("GetReportsFiltered", "guild-id", "victim-id", models.TypeMute, mock.AnythingOfType("int"), mock.AnythingOfType("int")).
		Once().
		Return(nil, nil)

	s, err = New(m.ct)
	assert.Nil(t, err)

	_, err = s.RevokeMute("guild-id", "executor-id", "victim-id", "")
	assert.Nil(t, err)
	m.s.AssertCalled(t, "GuildMemberTimeout", "guild-id", "victim-id", testutil.Nil[time.Time]())
	m.db.AssertNotCalled(t, "ExpireReports", "123")
}
