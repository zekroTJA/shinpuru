package inits

import (
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/listeners"
	"github.com/zekroTJA/shinpuru/internal/services/config"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/snowflakenodes"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
	"github.com/zekrotja/dgrs"
	"github.com/zekrotja/rogu/log"
)

func InitDiscordBotSession(container di.Container) (release func()) {
	release = func() {}

	log := log.Tagged("Discord")
	log.Info().Msg("Initializing bot session ...")

	err := snowflakenodes.Setup()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed setting up snowflake nodes")
	}

	session := container.Get(static.DiDiscordSession).(*discordgo.Session)
	cfg := container.Get(static.DiConfig).(config.Provider)

	session.Token = "Bot " + cfg.Config().Discord.Token
	session.Identify.Intents = discordgo.MakeIntent(static.Intents)
	session.StateEnabled = false

	if shardCfg := cfg.Config().Discord.Sharding; shardCfg.Total > 1 {
		st := container.Get(static.DiState).(*dgrs.State)

		var id int
		if shardCfg.AutoID {
			d := time.Duration(rand.Int63n(int64(5 * time.Second)))
			log.Info().
				Field("d", d.Round(time.Millisecond)).
				Msg("Sleeping before retrieving shard ID")
			time.Sleep(d)
			if id, err = st.ReserveShard(shardCfg.Pool); err != nil {
				log.Fatal().Err(err).Msg("Failed receiving alive shards from state")
			}
			release = func() {
				log.Info().Field("id", id).Msg("Releasing shard ID")
				if err = st.ReleaseShard(shardCfg.Pool, id); err != nil {
					log.Error().Err(err).Msg("Failed releasing shard ID")
				}
			}
		} else {
			id = shardCfg.ID
			if id < 0 || id >= shardCfg.Total {
				log.Fatal().Msgf("Shard ID must be in range [0, %d)", shardCfg.Total)
			}
		}

		log.Info().
			Field("id", id).
			Field("total", shardCfg.Total).
			Msg("Running in sharded mode")

		session.Identify.Shard = &[2]int{id, shardCfg.Total}
	}

	listenerInviteBlock := listeners.NewListenerInviteBlock(container)
	listenerGhostPing := listeners.NewListenerGhostPing(container)
	listenerColors := listeners.NewColorListener(container)

	listenerJDoodle, err := listeners.NewListenerJdoodle(container)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed setting up code execution listener")
	}

	listenerStarboard := listeners.NewListenerStarboard(container)
	listenerVerification := listeners.NewListenerVerifications(container)
	listenerAutoVoice := listeners.NewListenerAutoVoice(container)
	listenerGuilds := listeners.NewListenerGuildAdd(container)
	listenerRoleSelects := listeners.NewListenerRoleselect(container)
	listenerStatus := listeners.NewListenerStatus()

	session.AddHandler(listeners.NewListenerReady(container).Handler)
	session.AddHandler(listeners.NewListenerMemberAdd(container).Handler)
	session.AddHandler(listeners.NewListenerMemberRemove(container).Handler)
	session.AddHandler(listeners.NewListenerVote(container).Handler)
	session.AddHandler(listeners.NewListenerChannelCreate(container).Handler)
	session.AddHandler(listeners.NewListenerVoiceUpdate(container).Handler)
	session.AddHandler(discordutil.WrapHandler(listeners.NewListenerKarma(container).Handler))
	session.AddHandler(discordutil.WrapHandler(listeners.NewListenerAntiraid(container).HandlerMemberAdd))
	session.AddHandler(listeners.NewListenerBotMention(container).Listener)
	session.AddHandler(listeners.NewListenerDMSync(container).Handler)
	session.AddHandler(discordutil.WrapHandler(listeners.NewListenerPostBan(container).Handler))

	session.AddHandler(listenerGhostPing.HandlerMessageCreate)
	session.AddHandler(listenerGhostPing.HandlerMessageDelete)
	session.AddHandler(discordutil.WrapHandler(listenerInviteBlock.HandlerMessageSend))
	session.AddHandler(discordutil.WrapHandler(listenerInviteBlock.HandlerMessageEdit))

	session.AddHandler(listenerJDoodle.HandlerMessageCreate)
	session.AddHandler(listenerJDoodle.HandlerMessageUpdate)
	session.AddHandler(listenerJDoodle.HandlerReactionAdd)

	session.AddHandler(listenerColors.HandlerMessageCreate)
	session.AddHandler(listenerColors.HandlerMessageEdit)
	session.AddHandler(listenerColors.HandlerMessageReaction)

	session.AddHandler(listenerStarboard.ListenerReactionAdd)
	session.AddHandler(listenerStarboard.ListenerReactionRemove)

	session.AddHandler(listenerVerification.HandlerMemberAdd)
	session.AddHandler(listenerVerification.HandlerMemberRemove)

	session.AddHandler(listenerAutoVoice.HandlerVoiceUpdate)
	session.AddHandler(listenerAutoVoice.HandlerChannelDelete)

	session.AddHandler(listenerGuilds.HandlerReady)
	session.AddHandler(listenerGuilds.HandlerCreate)

	session.AddHandler(discordutil.WrapHandler(listenerRoleSelects.HandlerMessageBulkDelete))
	session.AddHandler(discordutil.WrapHandler(listenerRoleSelects.HandlerMessageDelete))
	session.AddHandler(discordutil.WrapHandler(listenerRoleSelects.Ready))

	session.AddHandler(listenerStatus.ListenerConnect)
	session.AddHandler(listenerStatus.ListenerDisconnect)

	session.AddHandler(func(s *discordgo.Session, e *discordgo.MessageCreate) {
		atomic.AddUint64(&util.StatsMessagesAnalysed, 1)
	})

	if cfg.Config().Metrics.Enable {
		session.AddHandler(listeners.NewListenerMetrics().Listener)
	}

	err = session.Open()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed connecting Discord bot session")
	}

	return
}
