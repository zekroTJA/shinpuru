package util

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	rxNumber = regexp.MustCompile(`^\d+$`)
)

func IsNumber(str string) bool {
	return rxNumber.MatchString(str)
}

func EnsureNotEmpty(str, def string) string {
	if str == "" {
		return def
	}
	return str
}

func ByteCountFormatter(bc uint64) string {
	f1k := float64(1024)
	if bc < 1024 {
		return fmt.Sprintf("%d B", bc)
	}
	if bc < 1024*1024 {
		return fmt.Sprintf("%.3f kiB", float64(bc)/f1k)
	}
	if bc < 1024*1024*1024 {
		return fmt.Sprintf("%.3f MiB", float64(bc)/f1k/f1k)
	}
	if bc < 1024*1024*1024*1024 {
		return fmt.Sprintf("%.3f GiB", float64(bc)/f1k/f1k/f1k)
	}
	return fmt.Sprintf("%.3f TiB", float64(bc)/f1k/f1k/f1k/f1k)
}

func BoolAsString(cond bool, ifTrue, ifFalse string) string {
	if cond {
		return ifTrue
	}
	return ifFalse
}

func IndexOfStrArray(str string, arr []string) int {
	for i, v := range arr {
		if v == str {
			return i
		}
	}
	return -1
}

func StringArrayContains(str string, arr []string) bool {
	return IndexOfStrArray(str, arr) > -1
}

func GetMessageLink(msg *discordgo.Message, guildID string) string {
	return fmt.Sprintf("https://discordapp.com/channels/%s/%s/%s", guildID, msg.ChannelID, msg.ID)
}

func GetDiscordSnowflakeCreationTime(snowflake string) (time.Time, error) {
	sfI, err := strconv.ParseInt(snowflake, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	timestamp := (sfI >> 22) + 1420070400000
	return time.Unix(timestamp/1000, timestamp), nil
}

func IsAdmin(g *discordgo.Guild, m *discordgo.Member) bool {
	if m == nil || g == nil {
		return false
	}

	for _, r := range g.Roles {
		if r.Permissions&0x8 != 0 {
			for _, mrID := range m.Roles {
				if r.ID == mrID {
					return true
				}
			}
		}
	}

	return false
}

// TODO: Deprecated
func DeleteMessageLater(s *discordgo.Session, msg *discordgo.Message, duration time.Duration) {
	if msg == nil {
		return
	}
	time.AfterFunc(duration, func() {
		s.ChannelMessageDelete(msg.ChannelID, msg.ID)
	})
}
