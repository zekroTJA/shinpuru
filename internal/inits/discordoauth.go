package inits

import (
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/core/config"
	"github.com/zekroTJA/shinpuru/internal/core/webserver/auth/oauth"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	discordoauth "github.com/zekroTJA/shinpuru/pkg/discordoauth/v2"
)

func InitDiscordOAuth(container di.Container) *discordoauth.DiscordOAuth {
	cfg := container.Get(static.DiConfig).(*config.Config)
	oauthHandler := container.Get(static.DiOAuthHandler).(oauth.OAuthHandler)

	return discordoauth.NewDiscordOAuth(
		cfg.Discord.ClientID,
		cfg.Discord.ClientSecret,
		cfg.WebServer.PublicAddr+static.EndpointAuthCB,
		oauthHandler.LoginFailedHandler,
		oauthHandler.LoginSuccessHandler,
	)
}
