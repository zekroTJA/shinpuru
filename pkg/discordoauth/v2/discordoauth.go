// package discordoauth provides fasthttp handlers
// to authenticate with via the Discord OAuth2
// endpoint. v2 swapped support for fasthttp-router
// to support for fiber handlers.
package discordoauth

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"github.com/zekrotja/jwt"
)

const (
	endpointOauth = "https://discord.com/api/oauth2/token"
	endpointMe    = "https://discord.com/api/users/@me"
)

type SuccessResult struct {
	UserID string
	State  map[string]string
}

// OnErrorFunc is the function to be used to handle errors during
// authentication.
type OnErrorFunc func(ctx *fiber.Ctx, status int, msg string) error

// OnSuccessFuc is the func to be used to handle the successful
// authentication.
type OnSuccessFuc func(ctx *fiber.Ctx, res SuccessResult) error

// DiscordOAuth provides http handlers for
// authenticating a discord User by your Discord
// OAuth application.
type DiscordOAuth struct {
	clientID        string
	clientSecret    string
	redirectURI     string
	stateSigningKey []byte

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
func NewDiscordOAuth(clientID, clientSecret, redirectURI string, onError OnErrorFunc, onSuccess OnSuccessFuc) (*DiscordOAuth, error) {
	if onError == nil {
		onError = func(ctx *fiber.Ctx, status int, msg string) error { return nil }
	}
	if onSuccess == nil {
		onSuccess = func(ctx *fiber.Ctx, res SuccessResult) error { return nil }
	}

	signingKey := make([]byte, 128)
	_, err := rand.Read(signingKey)
	if err != nil {
		return nil, err
	}

	return &DiscordOAuth{
		clientID:        clientID,
		clientSecret:    clientSecret,
		redirectURI:     redirectURI,
		stateSigningKey: signingKey,

		onError:   onError,
		onSuccess: onSuccess,
	}, nil
}

// HandlerInit returns a redirect response to the OAuth Apps
// authentication page.
func (d *DiscordOAuth) HandlerInit(ctx *fiber.Ctx) error {
	return d.HandlerInitWithState(ctx, nil)
}

// HandlerInitWithState returns a redirect response to the OAuth Apps
// authentication page with a passed state map which is then packed
// into the result passed to the onSuccess handler.
func (d *DiscordOAuth) HandlerInitWithState(ctx *fiber.Ctx, state map[string]string) error {
	stateToken, err := d.encodeAndSignWithPayload(state)
	if err != nil {
		return err
	}

	uri := fmt.Sprintf("https://discord.com/api/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=identify&state=%s",
		d.clientID, url.QueryEscape(d.redirectURI), stateToken)
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
	state := ctx.Query("state")

	if state == "" {
		return d.onError(ctx, fasthttp.StatusUnauthorized, "no state value returned")
	}

	claims, err := d.decodeAndValidate(state)
	if err != nil {
		switch err {
		case jwt.ErrInvalidSignature:
			return d.onError(ctx, fasthttp.StatusUnauthorized, "invalid state signature")
		case jwt.ErrNotValidYet, jwt.ErrTokenExpired, jwt.ErrInvalidTokenFormat:
			return d.onError(ctx, fasthttp.StatusBadRequest, "invalid state format: "+err.Error())
		default:
			return d.onError(ctx, fasthttp.StatusInternalServerError, "state validation failed: "+err.Error())
		}
	}

	// 1. Request getting bearer token by app auth code

	data := map[string][]string{
		"client_id":     {d.clientID},
		"client_secret": {d.clientSecret},
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {d.redirectURI},
		"scope":         {"identify"},
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)

	values := url.Values(data)

	req.Header.SetMethod("POST")
	req.SetRequestURI(endpointOauth)
	req.SetBody([]byte(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if err := fasthttp.Do(req, res); err != nil {
		return d.onError(ctx, fasthttp.StatusInternalServerError, "failed executing request: "+err.Error())
	}

	if res.StatusCode() >= 300 {
		return d.onError(ctx, fasthttp.StatusUnauthorized, "invalid auth code")
	}

	resAuthBody := new(oAuthTokenResponse)
	err = parseJSONBody(res.Body(), resAuthBody)
	if err != nil {
		return d.onError(ctx, fasthttp.StatusInternalServerError, "failed parsing Discord API response: "+err.Error())
	}

	if resAuthBody.Error != "" || resAuthBody.AccessToken == "" {
		return d.onError(ctx, fasthttp.StatusUnauthorized, "empty auth response")
	}

	// 2. Request getting user ID

	req.Header.Reset()
	req.ResetBody()
	req.Header.SetMethod("GET")
	req.SetRequestURI(endpointMe)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", resAuthBody.AccessToken))

	if err = fasthttp.Do(req, res); err != nil {
		return d.onError(ctx, fasthttp.StatusInternalServerError, "failed executing request: "+err.Error())
	}

	if res.StatusCode() >= 300 {
		return d.onError(ctx, fasthttp.StatusUnauthorized, "user request failed")
	}

	resGetMe := new(getUserMeResponse)
	err = parseJSONBody(res.Body(), resGetMe)
	if err != nil {
		return d.onError(ctx, fasthttp.StatusInternalServerError, "failed parsing Discord API response: "+err.Error())
	}

	if resGetMe.Error != "" || resGetMe.ID == "" {
		return d.onError(ctx, fasthttp.StatusUnauthorized, "empty user response")
	}

	return d.onSuccess(ctx, SuccessResult{
		UserID: resGetMe.ID,
		State:  claims.Payload,
	})
}

func parseJSONBody(body []byte, v interface{}) error {
	return json.Unmarshal(body, v)
}
