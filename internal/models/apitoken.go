package models

import "time"

type APITokenEntry struct {
	UserID     string    `json:"userid"`
	Salt       string    `json:"salt"`
	Created    time.Time `json:"created"`
	Expires    time.Time `json:"expires"`
	LastAccess time.Time `json:"lastaccess"`
	Hits       int       `json:"hits"`
}
