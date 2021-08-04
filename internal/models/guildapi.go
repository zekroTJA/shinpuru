package models

type GuildAPISettings struct {
	Enabled        bool   `json:"enabled"`
	AllowedOrigins string `json:"allowed_origins"`
	Protected      bool   `json:"protected"`
	TokenHash      string `json:"token_hash,omitempty"`
}

func (g *GuildAPISettings) Hydrate() *GuildAPISettings {
	g.Protected = g.TokenHash != ""
	return g
}
