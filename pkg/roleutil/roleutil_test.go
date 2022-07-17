package roleutil

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

func TestSortRoles(t *testing.T) {
	makeRoles := func() []*discordgo.Role {
		return []*discordgo.Role{
			{Position: 0},
			{Position: 4},
			{Position: 1},
			{Position: 5},
			{Position: 3},
			{Position: 2},
		}
	}

	roles := makeRoles()
	SortRoles(roles, false)

	for i, r := range roles {
		if i != r.Position {
			t.Errorf("role pos was %d (exp: %d)", r.Position, i)
		}
	}

	roles = makeRoles()
	SortRoles(roles, true)

	for i, r := range roles {
		exp := len(roles) - i - 1
		if exp != r.Position {
			t.Errorf("role pos was %d (exp: %d)", r.Position, exp)
		}
	}
}

func TestPositionDiff(t *testing.T) {
	guild := &discordgo.Guild{
		Roles: []*discordgo.Role{
			{ID: "0", Position: 0},
			{ID: "1", Position: 1},
			{ID: "2", Position: 2},
			{ID: "3", Position: 3},
			{ID: "4", Position: 4},
		},
	}

	m1 := &discordgo.Member{}
	m2 := &discordgo.Member{}

	m1.Roles = []string{"4"}
	m2.Roles = []string{"1"}

	diff := PositionDiff(m1, m2, guild)
	assert.Equal(t, 3, diff)

	diff = PositionDiff(m2, m1, guild)
	assert.Equal(t, -3, diff)

	m1.Roles = []string{"2", "3", "4"}
	m2.Roles = []string{"1", "0"}

	diff = PositionDiff(m1, m2, guild)
	assert.Equal(t, 3, diff)

	diff = PositionDiff(m2, m1, guild)
	assert.Equal(t, -3, diff)

	m1.Roles = []string{"2", "3", "4"}
	m2.Roles = []string{"1", "0", "4"}

	diff = PositionDiff(m1, m2, guild)
	assert.Equal(t, 0, diff)

	diff = PositionDiff(m2, m1, guild)
	assert.Equal(t, 0, diff)

	m1.Roles = []string{}
	m2.Roles = []string{"1", "0", "4"}

	diff = PositionDiff(m1, m2, guild)
	assert.Equal(t, -5, diff)

	diff = PositionDiff(m2, m1, guild)
	assert.Equal(t, 5, diff)

	m1.Roles = nil
	m2.Roles = []string{"1", "0", "4"}

	diff = PositionDiff(m1, m2, guild)
	assert.Equal(t, -5, diff)

	diff = PositionDiff(m2, m1, guild)
	assert.Equal(t, 5, diff)

	m1.Roles = nil
	m2.Roles = nil

	diff = PositionDiff(m1, m2, guild)
	assert.Equal(t, 0, diff)

	diff = PositionDiff(m2, m1, guild)
	assert.Equal(t, 0, diff)

	// ----------------------------------

	guild.Roles = nil

	m1.Roles = []string{"2", "3"}
	m2.Roles = []string{"1", "0"}

	diff = PositionDiff(m1, m2, guild)
	assert.Equal(t, 0, diff)

	diff = PositionDiff(m2, m1, guild)
	assert.Equal(t, 0, diff)
}
