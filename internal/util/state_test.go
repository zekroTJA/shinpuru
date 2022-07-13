package util

import (
	"errors"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zekroTJA/shinpuru/mocks"
	"github.com/zekroTJA/shinpuru/pkg/multierror"
)

func TestUpdateGuildMemberStats(t *testing.T) {
	st := &mocks.IState{}
	s := &mocks.ISession{}
	var err error

	// ----------------

	st.On("Guilds").Once().Return(nil, errors.New("test error"))

	err = UpdateGuildMemberStats(st, s)
	assert.EqualError(t, err, "test error")

	// ----------------

	st.On("Guilds").Once().Return([]*discordgo.Guild{
		{ID: "guild-0"},
		{ID: "guild-1"},
		{ID: "guild-2"},
	}, nil)
	st.On("SetGuild", mock.Anything).Times(2).Return(nil)

	g0 := &discordgo.Guild{
		ID:          "guild-0",
		MemberCount: 69,
	}
	s.On("GuildWithCounts", "guild-0").Once().Return(g0, nil)
	g1 := &discordgo.Guild{
		ID:          "guild-1",
		MemberCount: 1337,
	}
	s.On("GuildWithCounts", "guild-1").Once().Return(g1, nil)
	s.On("GuildWithCounts", "guild-2").Once().Return(nil, errors.New("test error"))

	err = UpdateGuildMemberStats(st, s)
	st.AssertCalled(t, "SetGuild", g0)
	st.AssertCalled(t, "SetGuild", g1)
	mErr := multierror.New()
	mErr.Append(errors.New("test error"))
	assert.EqualError(t, err, mErr.Error())
}
