package twitchnotify

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
