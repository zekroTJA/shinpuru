package discordoauth

import (
	"time"

	"github.com/zekrotja/jwt"
)

type stateClaims struct {
	jwt.PublicClaims

	Payload map[string]string `json:"pld,omitempty"`
}

func (d *DiscordOAuth) getHandler() jwt.Handler[stateClaims] {
	return jwt.NewHandler[stateClaims](jwt.NewHmacSha512(d.stateSigningKey))
}

func (d *DiscordOAuth) encodeAndSignWithPayload(payload map[string]string) (string, error) {
	now := time.Now()
	var c stateClaims
	c.SetExpTime(now.Add(5 * time.Minute))
	c.SetIat(now)
	c.SetNbfTime(now)
	c.Iss = "discordOauthValidator"

	if len(payload) != 0 {
		c.Payload = payload
	}

	return d.getHandler().EncodeAndSign(c)
}

func (d *DiscordOAuth) decodeAndValidate(token string) (stateClaims, error) {
	return d.getHandler().DecodeAndValidate(token)
}
