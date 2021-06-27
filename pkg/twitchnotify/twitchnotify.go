// Package twitchnotify provides functionalities
// to watch the state of twitch streams and
// notifying changes by polling the twitch REST
// API.
package twitchnotify

import (
	"errors"
	"fmt"
	"image"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/generaltso/vibrant"
	"github.com/zekroTJA/shinpuru/pkg/httpreq"
)

// UserIdent is the type of identificator for getting
// users via the Twitch API.
type UserIdent string

const (
	IdentID    UserIdent = "id"
	IdentLogin UserIdent = "login"

	maxUserCap = 1000

	// clockDuration = 30 * time.Second
	clockDuration = 60 * time.Second

	oAuth2Endpoint = "https://id.twitch.tv/oauth2/token"
	helixEndpoint  = "https://api.twitch.tv/helix"
)

var (
	ErrNotFound            = errors.New("not found")
	ErrInvalidResponseType = errors.New("invalid response type")
	ErrMaxUsersReached     = errors.New("max registered users reached")
)

// NotifyHandler describes a callback handler when a
// stream either goes online or offline passing the
// stream data as well as the user data of the streamer.
type NotifyHandler func(*Stream, *User)

// Credentials hold the client ID and client secret
// of a twitch API application.
type Credentials struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

// NotifyWorker provides general utilities to fetch
// watched online streamers and call notify handler
// callbacks when a stream goes online or offline.
type NotifyWorker struct {
	creds              *Credentials
	wentOnlineHandler  NotifyHandler
	wentOfflineHandler NotifyHandler

	mx        *sync.Mutex
	timer     *time.Ticker
	users     map[string]*User
	wereLive  []*Stream
	gameCache map[string]*Game

	bearerToken string
	bearerValid time.Time
}

// New initializes a new NotifyWorker instance and
// starts the worker timer loop.
func New(
	creds Credentials,
	wentOnlineHandler NotifyHandler,
	wentOfflineHandler NotifyHandler,
	config ...Config,
) (worker *NotifyWorker, err error) {
	conf := defaultConfig(config)

	worker = &NotifyWorker{
		creds:              &creds,
		wentOfflineHandler: wentOfflineHandler,
		wentOnlineHandler:  wentOnlineHandler,

		mx:        &sync.Mutex{},
		users:     make(map[string]*User),
		wereLive:  make([]*Stream, 0),
		gameCache: make(map[string]*Game),
	}

	if err = worker.getBearerToken(); err != nil {
		return
	}

	if conf.TimerDelay > 0 {
		worker.timer = time.NewTicker(conf.TimerDelay)

		go func() {
			for {
				<-worker.timer.C
				worker.Handle()
			}
		}()
	}

	return
}

// GetUser tries to fetch a user either by login or by ID, specified by
// typ which one of these methods is used.
// Returns the fetched user object and occured errors during fetch.
func (w *NotifyWorker) GetUser(identifyer string, typ UserIdent) (*User, error) {
	url := fmt.Sprintf("%s/users?%s=%s", helixEndpoint, typ, identifyer)

	data := new(usersDataWrapper)
	if err := w.doAuthenticatedGet(url, data); err != nil {
		return nil, err
	}

	if len(data.Data) < 1 || data.Data[0] == nil {
		return nil, ErrNotFound
	}

	return data.Data[0], nil
}

// AddUser adds the specified twitch User to the watch
// list. If maxUserCap is reached, an ErrMaxUsersreached
// error is returned.
func (w *NotifyWorker) AddUser(u *User) error {
	if len(w.users) >= maxUserCap {
		return ErrMaxUsersReached
	}

	w.users[u.ID] = u

	return nil
}

// GetEmbed assembles and returns an embed reference
// from the given Stream and User objects.
func GetEmbed(d *Stream, u *User) *discordgo.MessageEmbed {
	emb := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("%s just started streaming!", u.DisplayName),
		URL:         fmt.Sprintf("https://twitch.tv/%s", u.LoginName),
		Description: fmt.Sprintf("**%s**\n\nCurrent viewers: `%d`", d.Title, d.ViewerCount),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: u.AviURL,
		},
		Image: &discordgo.MessageEmbedImage{
			URL:    fmt.Sprintf("%s?rid=%d", d.ThumbnailURL, rand.Int()),
			Width:  1280,
			Height: 720,
		},
		Footer: &discordgo.MessageEmbedFooter{
			IconURL: strings.Replace(d.Game.IconURL, "{width}x{height}", "16x16", 1),
			Text:    "Playing " + d.Game.Name,
		},
	}

	if body, err := httpreq.GetFile(u.AviURL); err == nil {
		if imgData, _, err := image.Decode(body); err == nil {
			if palette, err := vibrant.NewPaletteFromImage(imgData); err == nil {
				for name, swatch := range palette.ExtractAwesome() {
					if name == "Vibrant" {
						emb.Color = int(swatch.Color)
					}
				}
			}
		}
	}

	return emb
}

// getBearerToken tries to authenticate with the configured
// twitch app credentials and retrieves a bearer token which
// is then used for further request authentication.
func (w *NotifyWorker) getBearerToken() error {
	url := fmt.Sprintf("%s?client_id=%s&client_secret=%s&grant_type=client_credentials",
		oAuth2Endpoint, w.creds.ClientID, w.creds.ClientSecret)

	res, err := httpreq.Post(url, nil, nil)
	if err != nil {
		return err
	}

	var token bearerTokenResponse
	if err = res.JSON(&token); err != nil {
		return err
	}

	w.bearerToken = token.AccessToken
	w.bearerValid = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)

	return nil
}

// doAuthenticatedGet executes a GET request to the twitch API
// using the retrieved bearer token for authentication.
// If the token is unset or has expired, a new token will be
// retrieved and the request fill be executed afterwards.
// The request result will be put in the passed data reference
// and errors occured are returned.
func (w *NotifyWorker) doAuthenticatedGet(url string, data interface{}) (err error) {
	if w.bearerToken == "" || time.Now().After(w.bearerValid) {
		if err = w.getBearerToken(); err != nil {
			return
		}
	}

	res, err := httpreq.Get(url, map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", w.bearerToken),
		"Client-ID":     w.creds.ClientID,
	})
	if err != nil {
		return
	}

	err = res.JSON(data)

	return
}

func (w *NotifyWorker) getStreams() ([]*Stream, error) {
	userIDs := make([]string, len(w.users))
	var i int
	for k := range w.users {
		userIDs[i] = k
		i++
	}

	url := fmt.Sprintf("%s/streams?user_id=%s",
		helixEndpoint, strings.Join(userIDs, "&user_id="))

	data := new(streamsDataWrapper)
	if err := w.doAuthenticatedGet(url, data); err != nil {
		return nil, err
	}

	streams := make([]*Stream, len(data.Data))
	for i, s := range data.Data {
		streams[i] = s
	}

	return streams, nil
}

// getGame tries to fetch a twitch Game type by the given
// gameID. The fetched game is returned as well as occured
// errors during fetch.
func (w *NotifyWorker) getGame(gameID string) (*Game, error) {
	game, ok := w.gameCache[gameID]
	if ok {
		return game, nil
	}

	url := fmt.Sprintf("%s/games?id=%s", helixEndpoint, gameID)
	data := new(gamesDataWrapper)
	if err := w.doAuthenticatedGet(url, data); err != nil {
		return nil, err
	}

	if len(data.Data) < 1 || data.Data[0] == nil {
		return nil, ErrNotFound
	}

	game = data.Data[0]
	w.gameCache[gameID] = game

	return game, nil
}

// Handle is the callback function executed on ech timer tick.
func (w *NotifyWorker) Handle() error {
	if len(w.users) < 1 {
		return nil
	}

	// Request watched streams which are currently live.
	streams, err := w.getStreams()
	if err != nil {
		return err
	}

	// Execute wentOnlineHandler for each stream which
	// is now live and was not live in the request before.
	for _, stream := range streams {
		var wasOnline bool

		w.mx.Lock()
		for _, nd := range w.wereLive {
			if nd.ID == stream.ID {
				wasOnline = true
				break
			}
		}
		w.mx.Unlock()

		if wasOnline {
			continue
		}

		game, err := w.getGame(stream.GameID)
		if err != nil {
			return err
		}

		stream.Game = game
		stream.ThumbnailURL = strings.Replace(stream.ThumbnailURL, "{width}x{height}", "1280x720", 1)

		user := w.users[stream.UserID]

		w.wentOnlineHandler(stream, user)
	}

	// Execute wentOfflineHandler for each stream which was
	// online in the request before and is now offline.
	w.mx.Lock()
	defer w.mx.Unlock()
	for _, nd := range w.wereLive {
		var stillOnline bool

		for _, stream := range streams {
			if stream.ID == nd.ID {
				stillOnline = true
				break
			}
		}

		if !stillOnline {
			user := w.users[nd.UserID]
			w.wentOfflineHandler(nd, user)
		}
	}

	// Update the last request state.
	w.wereLive = make([]*Stream, len(streams))
	copy(w.wereLive, streams)

	return nil
}
