package listeners

import (
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/snowflake"
	"github.com/sarulabs/di/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/mocks"
)

type antiraidHandlerMock struct {
	session *mocks.ISession
	db      *mocks.Database
	logger  *mocks.Logger
	state   *mocks.IState
	vs      *mocks.VerificationProvider

	ct di.Container

	getEvent func() *discordgo.GuildMemberAdd
}

func getAntiraidHandlerMock(prep ...func(t antiraidHandlerMock)) antiraidHandlerMock {
	t := antiraidHandlerMock{}

	t.session = &mocks.ISession{}
	t.db = &mocks.Database{}
	t.logger = &mocks.Logger{}
	t.state = &mocks.IState{}
	t.vs = &mocks.VerificationProvider{}

	if len(prep) != 0 {
		prep[0](t)
	}

	t.session.On("GuildEdit", mock.Anything, mock.Anything).Return(nil, nil)
	t.session.On("UserChannelCreate", mock.Anything).Return(&discordgo.Channel{ID: "test-ch"}, nil)
	t.session.On("ChannelMessageSendEmbed", mock.Anything, mock.Anything).Return(nil, nil)

	t.db.On("GetAntiraidState", mock.Anything).Return(true, nil)
	t.db.On("GetAntiraidRegeneration", mock.Anything).Return(30, nil)
	t.db.On("GetAntiraidBurst", mock.Anything).Return(3, nil)
	t.db.On("GetGuildModLog", mock.Anything).Return("", nil)
	t.db.On("GetAntiraidVerification", mock.Anything).Return(false, nil)
	t.db.On("AddToAntiraidJoinList", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	t.logger.On("Errorf", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
	t.logger.On("Section", mock.Anything).Return(t.logger)

	t.state.On("Guild", mock.Anything, mock.Anything).Return(&discordgo.Guild{
		OwnerID: "test-owner",
		Roles: []*discordgo.Role{
			{
				ID:          "test-admin-role",
				Permissions: 0x8,
			},
		},
	}, nil)
	t.state.On("Members", mock.Anything).Return([]*discordgo.Member{
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

	t.vs.On("SetEnabled", mock.Anything, mock.Anything).Return(nil)

	ct, _ := di.NewBuilder()
	ct.Add(di.Def{
		Name:  static.DiDatabase,
		Build: func(ctn di.Container) (interface{}, error) { return t.db, nil },
	})
	ct.Add(di.Def{
		Name:  static.DiGuildLog,
		Build: func(ctn di.Container) (interface{}, error) { return t.logger, nil },
	})
	ct.Add(di.Def{
		Name:  static.DiState,
		Build: func(ctn di.Container) (interface{}, error) { return t.state, nil },
	})
	ct.Add(di.Def{
		Name:  static.DiVerification,
		Build: func(ctn di.Container) (interface{}, error) { return t.vs, nil },
	})

	t.ct = ct.Build()

	snowflake.Epoch = 1420070400000
	node, _ := snowflake.NewNode(0)
	t.getEvent = func() *discordgo.GuildMemberAdd {
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

	return t
}

func TestHandleMemberAdd(t *testing.T) {
	const runs = 10

	collected := make([]string, 0, runs)
	m := getAntiraidHandlerMock(func(m antiraidHandlerMock) {
		m.db.On("GetAntiraidRegeneration", mock.Anything).Return(1, nil)
		m.db.On("GetAntiraidBurst", mock.Anything).Return(3, nil)
		m.db.On("AddToAntiraidJoinList", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("time.Time")).
			Run(func(args mock.Arguments) {
				assert.Equal(t, args[0], "test-guild")
				collected = append(collected, args[1].(string))
			}).
			Return(nil)
	})

	l := NewListenerAntiraid(m.ct)

	l.HandlerMemberAdd(m.session, m.getEvent()) // 2 tickets
	l.HandlerMemberAdd(m.session, m.getEvent()) // 1 ticket

	time.Sleep(1 * time.Second) // 2 tickets

	events := make([]string, 0, runs)
	for i := 0; i < runs; i++ {
		e := m.getEvent()
		if i > 1 {
			events = append(events, e.User.ID)
		}
		l.HandlerMemberAdd(m.session, e)
	}

	m.session.AssertNumberOfCalls(t, "ChannelMessageSendEmbed", 3)
	m.db.AssertNumberOfCalls(t, "AddToAntiraidJoinList", runs-3+1)

	assert.ElementsMatch(t, events, collected)
}

func TestHandleMemberAdd_Threadsafety(t *testing.T) {
	const runs, burst = 10000, 50

	m := getAntiraidHandlerMock(func(t antiraidHandlerMock) {
		t.db.On("GetAntiraidBurst", mock.Anything).Return(50, nil)
	})

	l := NewListenerAntiraid(m.ct)

	var wg sync.WaitGroup
	wg.Add(runs)

	rand.Seed(time.Now().Unix())
	for i := 0; i < runs; i++ {
		go func() {
			time.Sleep(time.Duration(rand.Intn(1000) * int(time.Millisecond)))
			l.HandlerMemberAdd(m.session, m.getEvent())
			wg.Done()
		}()
	}

	wg.Wait()

	m.session.AssertNumberOfCalls(t, "ChannelMessageSendEmbed", 3)
	m.db.AssertNumberOfCalls(t, "AddToAntiraidJoinList", runs-burst)
}
