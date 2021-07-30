package fetch

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/dgrs"
)

type DgrsDataOutlet struct {
	state      *dgrs.State
	forceFetch bool
}

var _ DataOutlet = (*DgrsDataOutlet)(nil)

func WrapDrgs(state *dgrs.State, forceFetch ...bool) DgrsDataOutlet {
	ff := len(forceFetch) > 0 && forceFetch[0]
	return DgrsDataOutlet{state, ff}
}

func (o DgrsDataOutlet) GuildRoles(guildID string) ([]*discordgo.Role, error) {
	return o.state.Roles(guildID, o.forceFetch)
}

func (o DgrsDataOutlet) GuildMembers(guildID string, _ string, _ int) (st []*discordgo.Member, err error) {
	return o.state.Members(guildID, o.forceFetch)
}

func (o DgrsDataOutlet) GuildChannels(guildID string) (st []*discordgo.Channel, err error) {
	return o.state.Channels(guildID, o.forceFetch)
}
