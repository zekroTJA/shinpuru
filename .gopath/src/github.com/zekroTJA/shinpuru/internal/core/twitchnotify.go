package core

import (
	"errors"
	"fmt"
	"image"
	"strings"
	"time"

	"github.com/generaltso/vibrant"
	"github.com/zekroTJA/shinpuru/internal/util"

	"github.com/bwmarrin/discordgo"
)

const clockDuration = 60 * time.Second

const (
	TwitchNotifyIdentLogin = "login"
	TwitchNotifyIdentID    = "id"
)

type TwitchNotifyData struct {
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

	Game *TwitchNotifyGame
}

type TwitchNotifyUser struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	LoginName   string `json:"login"`
	AviURL      string `json:"profile_image_url"`
}

type TwitchNotifyGame struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	IconURL string `json:"box_art_url"`
}

type TwitchNotifyHandler func(*TwitchNotifyData, *TwitchNotifyUser)

type TwitchNotifyWorker struct {
	timer              *time.Ticker
	users              map[string]*TwitchNotifyUser
	clientID           string
	pastResponses      []*TwitchNotifyData
	wentOnlineHandler  TwitchNotifyHandler
	wentOfflineHandler TwitchNotifyHandler
	gameIDCache        map[string]*TwitchNotifyGame
}

type TwitchNotifyDBEntry struct {
	GuildID      string
	ChannelID    string
	TwitchUserID string
}

func (g *TwitchNotifyGame) formatIconURL(res string) {
	g.IconURL = strings.Replace(g.IconURL, "{width}x{height}", res, 1)
}

func NewTwitchNotifyWorker(clientID string, wentOnlineHandler TwitchNotifyHandler, wentOfflineHandler TwitchNotifyHandler) *TwitchNotifyWorker {
	worker := &TwitchNotifyWorker{
		users:              make(map[string]*TwitchNotifyUser),
		clientID:           clientID,
		wentOnlineHandler:  wentOnlineHandler,
		wentOfflineHandler: wentOfflineHandler,
		gameIDCache:        make(map[string]*TwitchNotifyGame),
	}

	timer := time.NewTicker(clockDuration)

	go func() {
		for {
			<-timer.C
			worker.handler()
		}
	}()

	worker.timer = timer

	return worker
}

func (w *TwitchNotifyWorker) handler() error {
	if len(w.users) < 1 {
		return nil
	}
	urlParam := ""
	for uID := range w.users {
		urlParam += "&user_id=" + uID
	}

	res, err := HTTPRequest("GET", "https://api.twitch.tv/helix/streams?"+urlParam[1:], map[string]string{
		"Client-ID": w.clientID,
	}, nil)

	if err != nil {
		return err
	}

	var data struct {
		Data []*TwitchNotifyData `json:"data"`
	}
	err = res.ParseJSONBody(&data)
	if err != nil {
		return err
	}

	for _, cData := range data.Data {
		var isStillOffline bool
		for _, pData := range w.pastResponses {
			if cData.ID == pData.ID {
				isStillOffline = true
			}
		}

		if !isStillOffline {
			user, _ := w.users[cData.UserID]
			if game, ok := w.gameIDCache[cData.GameID]; !ok {
				res, err := HTTPRequest("GET", "https://api.twitch.tv/helix/games?id="+cData.GameID, map[string]string{
					"Client-ID": w.clientID,
				}, nil)
				if err == nil {
					var body struct {
						Data []*TwitchNotifyGame `json:"data"`
					}
					if res.ParseJSONBody(&body) == nil && &body != nil && len(body.Data) > 0 {
						game = body.Data[0]
						game.formatIconURL("50x70")
						w.gameIDCache[cData.GameID] = game
						cData.Game = game
					}
				} else {
					util.Log.Error("failed requesting game name: ", err)
				}
			} else {
				cData.Game = game
			}

			if cData.GameID == "" {
				cData.Game = &TwitchNotifyGame{
					Name: "game not found",
				}
			}

			cData.ThumbnailURL = strings.Replace(cData.ThumbnailURL, "{width}x{height}", "1280x720", 1)
			w.wentOnlineHandler(cData, user)
		}
	}

	for _, pData := range w.pastResponses {
		var isStillOnline bool
		for _, cData := range data.Data {
			if pData.ID == cData.ID {
				isStillOnline = true
			}
		}

		if !isStillOnline {
			if w.wentOfflineHandler != nil {
				w.wentOfflineHandler(pData, nil)
			}
		}
	}

	w.pastResponses = make([]*TwitchNotifyData, len(data.Data))
	copy(w.pastResponses, data.Data)

	return nil
}

func (w *TwitchNotifyWorker) GetUser(identifyer, identType string) (*TwitchNotifyUser, error) {
	res, err := HTTPRequest("GET", "https://api.twitch.tv/helix/users?"+identType+"="+identifyer, map[string]string{
		"Client-ID": w.clientID,
	}, nil)

	if err != nil {
		return nil, err
	}

	var data struct {
		Data []*TwitchNotifyUser `json:"data"`
	}

	err = res.ParseJSONBody(&data)
	if err != nil {
		return nil, err
	}

	if len(data.Data) < 1 || data.Data[0].ID == "" {
		return nil, errors.New("not found")
	}

	return data.Data[0], nil
}

func (w *TwitchNotifyWorker) AddUser(u *TwitchNotifyUser) error {
	if len(w.users) >= 1000 {
		return errors.New("max reached")
	}
	w.users[u.ID] = u
	return nil
}

func TwitchNotifyGetEmbed(d *TwitchNotifyData, u *TwitchNotifyUser) *discordgo.MessageEmbed {
	emb := &discordgo.MessageEmbed{
		Title:       u.DisplayName + " just started streaming!",
		URL:         "https://twitch.tv/" + u.LoginName,
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
			IconURL: d.Game.IconURL,
			Text:    "Playing " + d.Game.Name,
		},
	}

	if body, err := HTTPGetFile(u.AviURL); err == nil {
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
