package twitchnotify

import "time"

type Config struct {
	TimerDelay time.Duration `json:"timderdelay"`
}

func defaultConfig(configs []Config) Config {
	defaultConfig := Config{
		TimerDelay: 60 * time.Second,
	}

	if len(configs) > 0 {
		return configs[0]
	}

	return defaultConfig
}
