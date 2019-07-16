package inits

import (
	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/snowflake"
	"github.com/zekroTJA/shinpuru/internal/commands"
	"github.com/zekroTJA/shinpuru/internal/core"
	"github.com/zekroTJA/shinpuru/internal/listeners"
	"github.com/zekroTJA/shinpuru/internal/util"
)

func InitDiscordBotSession(session *discordgo.Session, config *core.Config, database core.Database, cmdHandler *commands.CmdHandler, lct *core.LCTimer) {
	snowflake.Epoch = util.DefEpoche
	err := util.SetupSnowflakeNodes()
	if err != nil {
		util.Log.Fatal("Failed setting up snowflake nodes: ", err)
	}

	session.Token = "Bot " + config.Discord.Token

	listenerInviteBlock := listeners.NewListenerInviteBlock(database, cmdHandler)
	listenerGhostPing := listeners.NewListenerGhostPing(database, cmdHandler)

	session.AddHandler(listeners.NewListenerReady(config, database, lct).Handler)
	session.AddHandler(listeners.NewListenerCmd(config, database, cmdHandler).Handler)
	session.AddHandler(listeners.NewListenerGuildJoin(config).Handler)
	session.AddHandler(listeners.NewListenerMemberAdd(database).Handler)
	session.AddHandler(listeners.NewListenerMemberRemove(database).Handler)
	session.AddHandler(listeners.NewListenerVote(database).Handler)
	session.AddHandler(listeners.NewListenerChannelCreate(database).Handler)
	session.AddHandler(listeners.NewListenerVoiceUpdate(database).Handler)
	session.AddHandler(listeners.NewListenerJdoodle(database).Handler)

	session.AddHandler(listenerGhostPing.HandlerMessageCreate)
	session.AddHandler(listenerGhostPing.HandlerMessageDelete)
	session.AddHandler(listenerInviteBlock.HandlerMessageSend)
	session.AddHandler(listenerInviteBlock.HandlerMessageEdit)

	err = session.Open()
	if err != nil {
		util.Log.Fatal("Failed connecting Discord bot session:", err)
	}
}
