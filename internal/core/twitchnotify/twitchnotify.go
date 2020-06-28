package twitchnotify

import (
	"errors"
	"fmt"
	"image"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/generaltso/vibrant"
	"github.com/zekroTJA/shinpuru/internal/core/config"
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

// NotifyWorker provides general utilities to fetch
// watched online streamers and call notify handler
// callbacks when a stream goes online or offline.
type NotifyWorker struct {
	conf               *config.TwitchApp
	wentOnlineHandler  NotifyHandler
	wentOfflineHandler NotifyHandler

	timer     *time.Ticker
	users     map[string]*User
	wereLive  []*Stream
	gameCache map[string]*Game

	bearerToken string
	bearerValid time.Time
}

// Stream wraps information about a twitch stream.
type Stream struct {
	ID           string   `json:"id"`
	UserID       string   `json:"user_id"`
	UserName     string   `json:"user_name"`
	GameID       string   `json:"game_id"`
	CommunityIDs []string `json:"community_ids"`
	Type         string   `json:"type"`
	Title        string   `json:"title"`
	ViewerCount  int      `json:"viewer_count"`
	StartedAt    string   `json:"started_at"`
	Language     string   `json:"language"`
	ThumbnailURL string   `json:"thumbnail_url"`

	Game *Game
}

// User wraps information about a twitch streamer.
type User struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	LoginName   string `json:"login"`
	Description string `json:"description"`
	AviURL      string `json:"profile_image_url"`
}

// Game wraps information about a twitch game type.
type Game struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	IconURL string `json:"box_art_url"`
}

// DBEntry specifies a database entry for tracking
// twitch users.
type DBEntry struct {
	GuildID      string
	ChannelID    string
	TwitchUserID string
}

type usersDataWrapper struct {
	Data []*User `json:"data"`
}

type streamsDataWrapper struct {
	Data []*Stream `json:"data"`
}

type gamesDataWrapper struct {
	Data []*Game `json:"data"`
}

type bearerTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// New initializes a new NotifyWorker instance and
// starts the worker timer loop.
func New(conf *config.TwitchApp, wentOnlineHandler NotifyHandler, wentOfflineHandler NotifyHandler) (*NotifyWorker, error) {
	worker := &NotifyWorker{
		conf:               conf,
		wentOfflineHandler: wentOfflineHandler,
		wentOnlineHandler:  wentOnlineHandler,

		users:     make(map[string]*User),
		wereLive:  make([]*Stream, 0),
		gameCache: make(map[string]*Game),
	}

	if err := worker.getBearerToken(); err != nil {
		return nil, err
	}

	worker.timer = time.NewTicker(clockDuration)

	go func() {
		for {
			<-worker.timer.C
			worker.handler()
		}
	}()

	return worker, nil
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
			URL:    d.ThumbnailURL,
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
		oAuth2Endpoint, w.conf.ClientID, w.conf.ClientSecret)

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
		"Client-ID":     w.conf.ClientID,
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

// handler is the callback function executed on ech timer tick.
func (w *NotifyWorker) handler() error {
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

	streamIterator1:
		for _, nd := range w.wereLive {
			if nd.ID == stream.ID {
				wasOnline = true
				break streamIterator1
			}
		}

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
	for _, nd := range w.wereLive {
		var stillOnline bool

	streamIterator2:
		for _, stream := range streams {
			if stream.ID == nd.ID {
				stillOnline = true
				break streamIterator2
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
