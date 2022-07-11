package listeners

import (
	"strconv"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di/v2"
	"github.com/stretchr/testify/mock"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/mocks"
)

type inviteBlockMock struct {
	session *mocks.ISession
	db      *mocks.Database
	logger  *mocks.Logger
	pmw     *mocks.PermissionsProvider

	ct di.Container
}

func getInviteBlockMock(f ...func(m inviteBlockMock)) inviteBlockMock {
	var t inviteBlockMock

	t.session = &mocks.ISession{}
	t.db = &mocks.Database{}
	t.logger = &mocks.Logger{}
	t.pmw = &mocks.PermissionsProvider{}

	t.logger.On("Errorf", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
	t.logger.On("Section", mock.Anything).Return(t.logger)

	if len(f) != 0 {
		f[0](t)
	}

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
			Name:  static.DiPermissions,
			Build: func(ctn di.Container) (interface{}, error) { return t.pmw, nil },
		},
	)

	t.ct = ct.Build()

	return t
}

func TestHandlerMessageSend(t *testing.T) {
	m := getInviteBlockMock(func(t inviteBlockMock) {
		t.db.On("GetGuildInviteBlock", "guild-id").Return("true", nil)
		t.db.On("GetGuildInviteBlock", "guild-id-disabled").Return("", nil)

		t.pmw.On("CheckPermissions", mock.Anything, mock.Anything, "user-allowed", "!sp.guild.mod.inviteblock.send").
			Return(true, false, nil)
		t.pmw.On("CheckPermissions", mock.Anything, mock.Anything, "user-notallowed", "!sp.guild.mod.inviteblock.send").
			Return(false, false, nil)

		t.session.On("GuildInvites", mock.Anything).Return([]*discordgo.Invite{
			{Code: "poggers"},
		}, nil)
		t.session.On("UserChannelCreate", mock.Anything).Return(&discordgo.Channel{ID: mock.Anything}, nil)
		t.session.On("ChannelMessageSendEmbed", mock.Anything, mock.Anything).Return(nil, nil)
		t.session.On("ChannelMessageDelete", "channel-id", mock.Anything).Return(nil)
	})

	l := NewListenerInviteBlock(m.ct)

	i := -1
	getMessage := func(user, content string, guildid ...string) *discordgo.MessageCreate {
		gid := "guild-id"
		if len(guildid) != 0 {
			gid = guildid[0]
		}
		i++
		return &discordgo.MessageCreate{
			Message: &discordgo.Message{
				Content:   content,
				ID:        strconv.Itoa(i),
				ChannelID: "channel-id",
				GuildID:   gid,
				Author: &discordgo.User{
					ID: user,
				},
			},
		}
	}

	positives := []*discordgo.MessageCreate{
		getMessage(
			"user-notallowed",
			"https://discord.gg/5nPSMXn6"),
		getMessage(
			"user-notallowed",
			"https://discord.com/invite/5nPSMXn6"),
		getMessage(
			"user-notallowed",
			"https://discordapp.com/invite/5nPSMXn6"),
		getMessage(
			"user-notallowed",
			"discord.gg/5nPSMXn6"),
		getMessage(
			"user-notallowed",
			"discord.com/invite/5nPSMXn6"),
		getMessage(
			"user-notallowed",
			"discordapp.com/invite/5nPSMXn6"),
		getMessage(
			"user-notallowed",
			"test content https://discord.gg/5nPSMXn6 askjhd aksdh aksjdh asda"),
		getMessage(
			"user-notallowed",
			"contenthttps://discord.gg/5nPSMXn6asdasd"),
		getMessage(
			"user-notallowed",
			"https://discord.zekro.de"),
		getMessage(
			"user-notallowed",
			"discord.zekro.de"),
	}

	negatives := []*discordgo.MessageCreate{
		getMessage(
			"user-notallowed",
			"some message"),
		getMessage(
			"user-notallowed",
			"https://zekro.de"),
		getMessage(
			"user-notallowed",
			"test content discord.gg/poggers"),
		getMessage(
			"user-allowed",
			"test content discord.gg/5nPSMXn6"),
		getMessage(
			"user-allowed",
			"some message"),
		getMessage(
			"user-notallowed",
			"https://discord.gg/5nPSMXn6",
			"guild-id-disabled"),
		getMessage(
			"user-notallowed",
			"some message",
			"guild-id-disabled"),
	}

	for _, e := range append(positives, negatives...) {
		l.HandlerMessageSend(m.session, e)
	}

	m.session.AssertNumberOfCalls(t, "ChannelMessageDelete", len(positives))
	m.session.AssertNumberOfCalls(t, "UserChannelCreate", len(positives))

	for _, e := range positives {
		m.session.AssertCalled(t, "ChannelMessageDelete", e.ChannelID, e.ID)
	}

	for _, e := range negatives {
		m.session.AssertNotCalled(t, "ChannelMessageDelete", e.ChannelID, e.ID)
	}
}
