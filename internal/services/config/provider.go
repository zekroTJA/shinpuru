package config

import (
	oldconfig "github.com/zekroTJA/shinpuru/internal/models"
)

type Provider interface {
	Config() *oldconfig.Config
	Parse() error
}
