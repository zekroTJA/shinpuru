package inits

import (
	"github.com/sarulabs/di/v2"
	"github.com/sirupsen/logrus"
	"github.com/zekroTJA/shinpuru/internal/services/config"
	"github.com/zekroTJA/shinpuru/internal/services/webserver/auth"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/discordoauth/v2"
)

func InitDiscordOAuth(container di.Container) *discordoauth.DiscordOAuth {
	cfg := container.Get(static.DiConfig).(config.Provider)
	oauthHandler := container.Get(static.DiOAuthHandler).(auth.RequestHandler)

	doa, err := discordoauth.NewDiscordOAuth(
		cfg.Config().Discord.ClientID,
		cfg.Config().Discord.ClientSecret,
		cfg.Config().WebServer.PublicAddr+static.EndpointAuthCB,
		oauthHandler.LoginFailedHandler,
		oauthHandler.LoginSuccessHandler,
	)

	if err != nil {
		logrus.WithError(err).Fatal("Discord OAuth initialization failed")
	}

	return doa
}
