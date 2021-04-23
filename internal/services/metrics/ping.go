package metrics

import (
	"math"
	"time"

	"github.com/go-ping/ping"
)

const discordAPIendpoint = "gateway.discord.gg"

func startPingWatcher(interval time.Duration) (pinger *ping.Pinger, err error) {
	pinger, err = ping.NewPinger(discordAPIendpoint)
	if err != nil {
		return
	}

	pinger.SetPrivileged(true)
	pinger.RecordRtts = false
	pinger.Interval = interval
	pinger.Timeout = time.Duration(math.MaxInt64)

	pinger.OnRecv = func(p *ping.Packet) {
		DiscordGatewayPing.Set(float64(p.Rtt.Microseconds()) / 1000)
	}

	go pinger.Run()

	return
}
