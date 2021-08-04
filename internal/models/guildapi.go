package models

type GuildAPISettings struct {
	Enabled        bool   `json:"enabled"`
	AllowedOrigins string `json:"allowed_origins"`
	Protected      bool   `json:"protected"`
	TokenHash      string `json:"-"`
}

func (g *GuildAPISettings) Hydrate() *GuildAPISettings {
	g.Protected = g.TokenHash != ""
	return g
}
