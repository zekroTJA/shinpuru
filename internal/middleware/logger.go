package middleware

import (
	"github.com/sirupsen/logrus"
	"github.com/zekroTJA/shireikan"
)

// LoggerMiddlewrae implements shireikan.Middleware to
// log executed commands.
type LoggerMiddlewrae struct{}

func (m *LoggerMiddlewrae) Handle(
	cmd shireikan.Command,
	ctx shireikan.Context,
	layer shireikan.MiddlewareLayer,
) (next bool, err error) {

	gid := "[dm]"
	if ctx.GetGuild() != nil {
		gid = ctx.GetGuild().ID
	}

	logrus.WithFields(logrus.Fields{
		"gid": gid,
		"uid": ctx.GetUser().ID,
		"cid": ctx.GetChannel().ID,
	}).Info("COMMANDS :: ", cmd.GetInvokes()[0], ctx.GetArgs())

	return true, nil
}

func (m *LoggerMiddlewrae) GetLayer() shireikan.MiddlewareLayer {
	return shireikan.LayerAfterCommand
}
