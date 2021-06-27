package listeners

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/zekroTJA/shinpuru/internal/services/metrics"
)

type ListenerMetrics struct {
}

func NewListenerMetrics() *ListenerMetrics {
	return &ListenerMetrics{}
}

func (l *ListenerMetrics) Listener(s *discordgo.Session, e *discordgo.Event) {
	metrics.DiscordEventTriggers.
		With(prometheus.Labels{"event": strings.ToLower(e.Type)}).
		Add(1)
}
