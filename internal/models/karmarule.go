package models

import "github.com/bwmarrin/snowflake"

type KarmaAction string

const (
	KarmaActionToogleRole  KarmaAction = "TOGGLE_ROLE"
	KarmaActionKick        KarmaAction = "KICK"
	KarmaActionBan         KarmaAction = "BAN"
	KarmaActionSendMessage KarmaAction = "SEND_MESSAGE"
)

type KarmaTriggerType int

const (
	KarmaTriggerBelow KarmaTriggerType = iota
	KarmaTriggerAbove
)

type KarmaRule struct {
	ID       snowflake.ID     `json:"id"`
	GuildID  string           `json:"guildid"`
	Trigger  KarmaTriggerType `json:"trigger"`
	Value    int              `json:"value"`
	Argument string           `json:"argument"`
}
