package models

import (
	"errors"

	"github.com/bwmarrin/snowflake"
)

type KarmaAction string

const (
	KarmaActionToogleRole  KarmaAction = "TOGGLE_ROLE"
	KarmaActionKick        KarmaAction = "KICK"
	KarmaActionBan         KarmaAction = "BAN"
	KarmaActionSendMessage KarmaAction = "SEND_MESSAGE"
)

func (a KarmaAction) Validate() bool {
	switch a {
	case KarmaActionToogleRole, KarmaActionKick, KarmaActionBan, KarmaActionSendMessage:
		return true
	default:
		return false
	}
}

type KarmaTriggerType int

const (
	KarmaTriggerBelow KarmaTriggerType = iota
	KarmaTriggerAbove

	karmaTriggerMax
)

func (tt KarmaTriggerType) Validate() bool {
	return tt >= 0 && tt < karmaTriggerMax
}

type KarmaRule struct {
	ID       snowflake.ID     `json:"id"`
	GuildID  string           `json:"guildid"`
	Trigger  KarmaTriggerType `json:"trigger"`
	Value    int              `json:"value"`
	Action   KarmaAction      `json:"action"`
	Argument string           `json:"argument"`
}

func (r *KarmaRule) Validate() error {
	if !r.Trigger.Validate() {
		return errors.New("invalid value for trigger")
	}
	if !r.Action.Validate() {
		return errors.New("invalid value for action")
	}
	return nil
}
