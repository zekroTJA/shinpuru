package report

import (
	"time"

	"github.com/bwmarrin/snowflake"
)

type UnbanRequestState int

const (
	UnbanRequestStatePending UnbanRequestState = iota
	UnbanRequestStateDeclined
	UnbanRequestStateAccepted
)

type UnbanRequest struct {
	ID          snowflake.ID      `json:"id"`
	UserID      string            `json:"user_id"`
	GuildID     string            `json:"guild_id"`
	UserTag     string            `json:"user_tag"`
	Message     string            `json:"message"`
	Status      UnbanRequestState `json:"status"`
	ProcessedBy string            `json:"processed_by"`
	Processed   time.Time         `json:"processed"`
}
