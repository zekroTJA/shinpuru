package models

import (
	"errors"

	"github.com/bwmarrin/snowflake"
	"github.com/zekroTJA/shinpuru/pkg/checksum"
)

type KarmaAction string

const (
	KarmaActionToggleRole  KarmaAction = "TOGGLE_ROLE"
	KarmaActionKick        KarmaAction = "KICK"
	KarmaActionBan         KarmaAction = "BAN"
	KarmaActionSendMessage KarmaAction = "SEND_MESSAGE"
)

func (a KarmaAction) Validate() bool {
	switch a {
	case KarmaActionToggleRole, KarmaActionKick, KarmaActionBan, KarmaActionSendMessage:
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
	Checksum string           `json:"-"`
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

func (r *KarmaRule) CalculateChecksum() string {
	cop := *r
	cop.ID = 0
	r.Checksum = checksum.Must(checksum.SumMd5(&cop))
	return r.Checksum
}
