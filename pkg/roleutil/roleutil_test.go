package roleutil

import (
	"testing"

	"github.com/bwmarrin/discordgo"
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
