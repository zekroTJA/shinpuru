package metrics

import (
	"time"

	"github.com/go-ping/ping"
)

const discordAPIendpoint = "gateway.discord.gg"

// PingWatcher detects the average round trip time
// to the discord API gateway endpoint in the given
// interval and saves it.
type PingWatcher struct {
	pinger *ping.Pinger
	ticker *time.Ticker

	LastRead  *ping.Statistics
	OnElapsed func(*ping.Statistics, error)
}

// NewPingWatcher intializes a new PingWatcher instance
// and starts the watch timer with the given interval.
func NewPingWatcher(interval time.Duration) (pw *PingWatcher, err error) {
	pw = new(PingWatcher)

	pw.pinger, err = ping.NewPinger(discordAPIendpoint)
	if err != nil {
		return
	}

	pw.pinger.SetPrivileged(true)
	pw.pinger.Count = 3

	pw.ticker = time.NewTicker(interval)
	go pw.tickerWorker()

	return
}

func (pw *PingWatcher) tickerWorker() {
	for {
		go pw.recordPing()
		<-pw.ticker.C
	}
}

func (pw *PingWatcher) recordPing() {
	err := pw.pinger.Run()
	pw.LastRead = pw.pinger.Statistics()
	if pw.OnElapsed != nil {
		pw.OnElapsed(pw.LastRead, err)
	}
}
