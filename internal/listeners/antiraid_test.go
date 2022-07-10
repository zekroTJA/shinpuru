package listeners

import (
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/snowflake"
	"github.com/sarulabs/di/v2"
	"github.com/stretchr/testify/mock"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/mocks"
)

func TestHandleMemberAdd_Threadsafety(t *testing.T) {
	sMock := &mocks.ISession{}
	sMock.On("GuildEdit", mock.Anything, mock.Anything).Return(nil, nil)
	sMock.On("UserChannelCreate", mock.Anything).Return(&discordgo.Channel{ID: "test-ch"}, nil)
	sMock.On("ChannelMessageSendEmbed", mock.Anything, mock.Anything).Return(nil, nil)

	dbMock := &mocks.Database{}
	dbMock.On("GetAntiraidState", mock.Anything).Return(true, nil)
	dbMock.On("GetAntiraidRegeneration", mock.Anything).Return(30, nil)
	dbMock.On("GetAntiraidBurst", mock.Anything).Return(50, nil)
	dbMock.On("GetGuildModLog", mock.Anything).Return("", nil)
	dbMock.On("GetAntiraidVerification", mock.Anything).Return(false, nil)
	dbMock.On("AddToAntiraidJoinList", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	glMock := &mocks.Logger{}
	glMock.On("Errorf", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
	glMock.On("Section", mock.Anything).Return(glMock)

	stMock := &mocks.IState{}
	stMock.On("Guild", mock.Anything, mock.Anything).Return(&discordgo.Guild{
		OwnerID: "test-owner",
		Roles: []*discordgo.Role{
			{
				ID:          "test-admin-role",
				Permissions: 0x8,
			},
		},
	}, nil)
	stMock.On("Members", mock.Anything).Return([]*discordgo.Member{
		{User: &discordgo.User{
			ID: "test-owner",
		}},
		{User: &discordgo.User{
			ID: "test-admin-1",
		}, Roles: []string{"test-admin-role"}},
		{User: &discordgo.User{
			ID: "test-admin-2",
		}, Roles: []string{"test-admin-role"}},
	}, nil)

	vsMock := &mocks.VerificationProvider{}
	vsMock.On("SetEnabled", mock.Anything, mock.Anything).Return(nil)

	ct, _ := di.NewBuilder()
	ct.Add(di.Def{
		Name:  static.DiDatabase,
		Build: func(ctn di.Container) (interface{}, error) { return dbMock, nil },
	})
	ct.Add(di.Def{
		Name:  static.DiGuildLog,
		Build: func(ctn di.Container) (interface{}, error) { return glMock, nil },
	})
	ct.Add(di.Def{
		Name:  static.DiState,
		Build: func(ctn di.Container) (interface{}, error) { return stMock, nil },
	})
	ct.Add(di.Def{
		Name:  static.DiVerification,
		Build: func(ctn di.Container) (interface{}, error) { return vsMock, nil },
	})

	l := NewListenerAntiraid(ct.Build())

	snowflake.Epoch = 1420070400000
	node, _ := snowflake.NewNode(0)
	getEvent := func() *discordgo.GuildMemberAdd {
		id := node.Generate().String()
		return &discordgo.GuildMemberAdd{
			Member: &discordgo.Member{
				GuildID: "test-guild",
				User: &discordgo.User{
					ID: id,
				},
			},
		}
	}

	const runs = 10000

	var wg sync.WaitGroup
	wg.Add(runs)

	rand.Seed(time.Now().Unix())
	for i := 0; i < runs; i++ {
		go func() {
			time.Sleep(time.Duration(rand.Intn(1000) * int(time.Millisecond)))
			l.HandlerMemberAdd(sMock, getEvent())
			wg.Done()
		}()
	}

	wg.Wait()

	sMock.AssertNumberOfCalls(t, "ChannelMessageSendEmbed", 3)
	dbMock.AssertNumberOfCalls(t, "AddToAntiraidJoinList", runs-50)
}
