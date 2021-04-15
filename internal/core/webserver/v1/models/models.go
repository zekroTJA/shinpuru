package models

import "time"

type AccessTokenResponse struct {
	Token   string    `json:"token"`
	Expires time.Time `json:"expires"`
}
