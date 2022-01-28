package birthday

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di/v2"
	"github.com/sirupsen/logrus"
	"github.com/zekroTJA/shinpuru/internal/models"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/guildlog"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekrotja/dgrs"
	"github.com/zekrotja/sop"
)

type BirthdayService struct {
	db      database.Database
	st      *dgrs.State
	session *discordgo.Session
	gl      guildlog.Logger
}

func New(ctn di.Container) *BirthdayService {
	return &BirthdayService{
		db:      ctn.Get(static.DiDatabase).(database.Database),
		st:      ctn.Get(static.DiState).(*dgrs.State),
		session: ctn.Get(static.DiDiscordSession).(*discordgo.Session),
		gl:      ctn.Get(static.DiGuildLog).(guildlog.Logger).Section("birthday"),
	}
}

func (b *BirthdayService) Schedule() (err error) {
	bdays, err := b.db.GetBirthdays("")
	if err != nil {
		return
	}

	bdayMap := sop.GroupE[*models.Birthday](
		sop.Slice(bdays), func(v *models.Birthday, i int) (string, *models.Birthday) {
			return v.GuildID, v
		},
	)

	if len(bdayMap) == 0 {
		return
	}

	guilds, err := b.st.Guilds()
	if err != nil {
		return
	}

	for _, guild := range guilds {
		gbds, ok := bdayMap[guild.ID]
		if !ok || gbds.Len() == 0 {
			continue
		}
		bdayChan, err := b.db.GetGuildBirthdayChan(guild.ID)
		if err != nil && !database.IsErrDatabaseNotFound(err) {
			logrus.WithError(err).WithField("gid", guild.ID).Error("failed getting birthday channel")
			b.gl.Errorf(guild.ID, "Failed getting birthday channel: %s", err.Error())
			continue
		}
		if bdayChan == "" {
			continue
		}

		ch, _ := b.st.Channel(bdayChan)
		if ch == nil {
			b.gl.Warnf(guild.ID, "Birthday channel has been disabled because it could not be found on the guild")
			b.db.SetGuildBirthdayChan(guild.ID, "")
			continue
		}

		gbds.
			Filter(isTodayFilter()).
			Each(func(v *models.Birthday, i int) {
				err := b.sendMessage(bdayChan, v)
				if err != nil {
					logrus.WithError(err).WithField("gid", guild.ID).Error("failed sending birthday message")
					b.gl.Errorf(guild.ID, "Failed sending birthday message: %s", err.Error())
				}
			})
	}

	return
}

func (b *BirthdayService) sendMessage(chanID string, bd *models.Birthday) (err error) {
	user, err := b.st.User(bd.UserID)
	if err != nil {
		return
	}

	age := ""
	if bd.ShowYear {
		age = suffix(time.Now().Year()-bd.Date.Year()) + " "
	}

	userMention := user.Mention() + "'"
	if !strings.HasPrefix(strings.ToLower(user.Username), "s") && !strings.HasPrefix(strings.ToLower(user.Username), "z") {
		userMention += "s"
	}

	desc := fmt.Sprintf(
		"Today is %s %sbirthday!\n\nHappy birthday to you!  🥳 🎉 🎊",
		userMention, age)

	_, err = b.session.ChannelMessageSendEmbed(chanID, &discordgo.MessageEmbed{
		Color:       static.ColorEmbedDefault,
		Description: desc,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL:    user.AvatarURL(""),
			Width:  24,
			Height: 24,
		},
	})
	return
}

func isTodayFilter() func(v *models.Birthday, i int) bool {
	now := time.Now().UTC()
	return func(v *models.Birthday, i int) bool {
		m1, d1, h1 := format(now)
		m2, d2, h2 := format(v.Date)
		return m1 == m2 && d1 == d2 && h1 == h2
	}
}

func format(date time.Time) (m time.Month, d, h int) {
	_, m, d = date.Date()
	h, _, _ = date.Clock()
	return
}

func suffix(i int) string {
	v := strconv.Itoa(i)
	suffix := "th"
	switch i % 10 {
	case 1:
		suffix = "st"
	case 2:
		suffix = "nd"
	case 3:
		suffix = "rd"
	}
	return v + suffix
}
