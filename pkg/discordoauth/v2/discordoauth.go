// package discordoauth provides fasthttp handlers
// to authenticate with via the Discord OAuth2
// endpoint. v2 swapped support for fasthttp-router
// to support for fiber handlers.
package discordoauth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gofiber/fiber/v2"
)

const (
	endpointOauth = "https://discord.com/api/oauth2/token"
	endpointMe    = "https://discord.com/api/users/@me"
)

// OnErrorFunc is the function to be used to handle errors during
// authentication.
type OnErrorFunc func(ctx *fiber.Ctx, status int, msg string) error

// OnSuccessFuc is the func to be used to handle the successful
// authentication.
type OnSuccessFuc func(ctx *fiber.Ctx, userID string) error

// DiscordOAuth provides http handlers for
// authenticating a discord User by your Discord
// OAuth application.
type DiscordOAuth struct {
	clientID     string
	clientSecret string
	redirectURI  string

	onError   OnErrorFunc
	onSuccess OnSuccessFuc
}

type oAuthTokenResponse struct {
	Error       string `json:"error"`
	AccessToken string `json:"access_token"`
}

type getUserMeResponse struct {
	Error string `json:"error"`
	ID    string `json:"id"`
}

// NewDiscordOAuth returns a new instance of DiscordOAuth.
func NewDiscordOAuth(clientID, clientSecret, redirectURI string, onError OnErrorFunc, onSuccess OnSuccessFuc) *DiscordOAuth {
	if onError == nil {
		onError = func(ctx *fiber.Ctx, status int, msg string) error { return nil }
	}
	if onSuccess == nil {
		onSuccess = func(ctx *fiber.Ctx, userID string) error { return nil }
	}

	return &DiscordOAuth{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURI:  redirectURI,

		onError:   onError,
		onSuccess: onSuccess,
	}
}

// HandlerInit returns a redirect response to the OAuth Apps
// authentication page.
func (d *DiscordOAuth) HandlerInit(ctx *fiber.Ctx) error {
	uri := fmt.Sprintf("https://discord.com/api/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=identify",
		d.clientID, url.QueryEscape(d.redirectURI))
	return ctx.Redirect(uri, fiber.StatusTemporaryRedirect)
}

// HandlerCallback will be requested by discordapp.com on successful
// app authentication. This handler will check the validity of the passed
// authorization code by getting a bearer token and trying to get self
// user data by requesting them using the bearer token.
// If this fails, onError will be called. Else, onSuccess will be
// called passing the userID of the user authenticated.
func (d *DiscordOAuth) HandlerCallback(ctx *fiber.Ctx) error {
	code := ctx.Query("code")

	// 1. Request getting bearer token by app auth code

	data := map[string][]string{
		"client_id":     {d.clientID},
		"client_secret": {d.clientSecret},
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {d.redirectURI},
		"scope":         {"identify"},
	}

	values := url.Values(data)
	req, err := http.NewRequest("POST", endpointOauth,
		bytes.NewBuffer([]byte(values.Encode())))
	if err != nil {
		return d.onError(ctx, http.StatusInternalServerError, "failed creating request: "+err.Error())
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return d.onError(ctx, http.StatusInternalServerError, "failed executing request: "+err.Error())
	}

	if res.StatusCode >= 300 {
		return d.onError(ctx, http.StatusUnauthorized, "")
	}

	resAuthBody := new(oAuthTokenResponse)
	err = parseJSONBody(res.Body, resAuthBody)
	if err != nil {
		return d.onError(ctx, http.StatusInternalServerError, "failed parsing Discord API response: "+err.Error())
	}

	if resAuthBody.Error != "" || resAuthBody.AccessToken == "" {
		return d.onError(ctx, http.StatusUnauthorized, "")
	}

	// 2. Request getting user ID

	req, err = http.NewRequest("GET", endpointMe, nil)
	if err != nil {
		return d.onError(ctx, http.StatusInternalServerError, "failed creating request: "+err.Error())
	}

	req.Header.Set("Authorization", "Bearer "+resAuthBody.AccessToken)

	res, err = http.DefaultClient.Do(req)
	if err != nil {
		return d.onError(ctx, http.StatusInternalServerError, "failed executing request: "+err.Error())
	}

	if res.StatusCode >= 300 {
		return d.onError(ctx, http.StatusUnauthorized, "")
	}

	resGetMe := new(getUserMeResponse)
	err = parseJSONBody(res.Body, resGetMe)
	if err != nil {
		return d.onError(ctx, http.StatusInternalServerError, "failed parsing Discord API response: "+err.Error())
	}

	if resGetMe.Error != "" || resGetMe.ID == "" {
		return d.onError(ctx, http.StatusUnauthorized, "")
	}

	return d.onSuccess(ctx, resGetMe.ID)
}

func parseJSONBody(body io.ReadCloser, v interface{}) error {
	dec := json.NewDecoder(body)
	return dec.Decode(v)
}
