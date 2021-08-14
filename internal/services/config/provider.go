package config

import (
	oldconfig "github.com/zekroTJA/shinpuru/internal/config"
)

type Provider interface {
	Config() *oldconfig.Config
	Parse() error
}
