package versioncheck

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLatestVersion(t *testing.T) {
	p := NewGitHubProvider("zekroTJA", "shinpuru")
	_, err := p.GetLatestVersion()
	assert.Nil(t, err)
}
