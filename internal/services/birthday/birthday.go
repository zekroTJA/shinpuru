package birthday

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di/v2"
	"github.com/zekroTJA/shinpuru/internal/models"
	"github.com/zekroTJA/shinpuru/internal/services/config"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/services/guildlog"
	"github.com/zekroTJA/shinpuru/internal/services/timeprovider"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"
	"github.com/zekroTJA/shinpuru/pkg/giphy"
	"github.com/zekrotja/dgrs"
	"github.com/zekrotja/rogu"
	"github.com/zekrotja/rogu/log"
	"github.com/zekrotja/sop"
)

type BirthdayService struct {
	gif     *giphy.Client
	db      database.Database
	st      *dgrs.State
	session *discordgo.Session
	gl      guildlog.Logger
	tp      timeprovider.Provider
	log     rogu.Logger
}

func New(ctn di.Container) *BirthdayService {
	b := &BirthdayService{
		db:      ctn.Get(static.DiDatabase).(database.Database),
		st:      ctn.Get(static.DiState).(*dgrs.State),
		session: ctn.Get(static.DiDiscordSession).(*discordgo.Session),
		gl:      ctn.Get(static.DiGuildLog).(guildlog.Logger).Section("birthday"),
		tp:      ctn.Get(static.DiTimeProvider).(timeprovider.Provider),
		log:     log.Tagged("Birthdays"),
	}

	cfg := ctn.Get(static.DiConfig).(config.Provider)
	if apiKey := cfg.Config().Giphy.APIKey; apiKey != "" {
		b.gif = giphy.New(apiKey, "v1")
	}

	return b
}

func (b *BirthdayService) Schedule() (err error) {
	bdays, err := b.db.GetBirthdays("")
	if err != nil {
		return
	}

	shardId, shardTotal := discordutil.GetShardOfSession(b.session)
	if shardTotal > 1 {
		bdays = sop.Slice(bdays).
			Filter(func(v models.Birthday, _ int) bool {
				id, err := discordutil.GetShardOfGuild(v.GuildID, shardTotal)
				return err == nil && id == shardId
			}).
			Unwrap()
	}

	bdayMap := sop.GroupE[models.Birthday](
		sop.Slice(bdays), func(v models.Birthday, i int) (string, models.Birthday) {
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

		gbds = gbds.Filter(isTodayFilter(b.tp.Now))
		if gbds.Len() == 0 {
			continue
		}

		bdayChan, err := b.db.GetGuildBirthdayChan(guild.ID)
		if err != nil && !database.IsErrDatabaseNotFound(err) {
			b.log.Error().Err(err).Field("gid", guild.ID).Msg("Failed getting birthday channel")
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

		gbds.Each(func(v models.Birthday, i int) {
			memb, err := b.st.Member(v.GuildID, v.UserID)
			if memb == nil || memb.User == nil {
				if err == nil || discordutil.IsErrCode(err, discordgo.ErrCodeUnknownMember) {
					err = b.db.DeleteBirthday(v.GuildID, v.UserID)
				}
			} else {
				err = b.sendMessage(memb, bdayChan, v)
			}
			if err != nil {
				b.log.Error().Err(err).Field("gid", guild.ID).Msg("Failed handling birthday")
				b.gl.Errorf(guild.ID, "Failed handling birthday: %s", err.Error())
			}
		})
	}

	return
}

func (b *BirthdayService) sendMessage(memb *discordgo.Member, chanID string, bd models.Birthday) (err error) {
	age := ""
	if bd.ShowYear {
		age = suffix(b.tp.Now().Year()-bd.Date.Year()) + " "
	}

	userMention := memb.Mention() + "'"
	if !strings.HasSuffix(strings.ToLower(memb.User.Username), "s") && !strings.HasSuffix(strings.ToLower(memb.User.Username), "z") {
		userMention += "s"
	}

	desc := fmt.Sprintf(
		"Today is %s %sbirthday!\n\nHappy birthday to you!  🥳 🎉 🎊",
		userMention, age)

	emb := &discordgo.MessageEmbed{
		Color:       static.ColorEmbedDefault,
		Description: desc,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL:    memb.User.AvatarURL(""),
			Width:  24,
			Height: 24,
		},
	}

	if gif := b.randomGif(); gif != nil {
		width, _ := strconv.Atoi(gif.Width)
		height, _ := strconv.Atoi(gif.Height)
		emb.Image = &discordgo.MessageEmbedImage{
			URL:    gif.Url,
			Width:  width,
			Height: height,
		}
	}

	_, err = b.session.ChannelMessageSendEmbed(chanID, emb)
	return
}

func (b *BirthdayService) randomGif() (img *giphy.Image) {
	if b.gif == nil {
		return
	}
	rng := rand.Intn(100)
	res, err := b.gif.Search("birthday", 1, rng, "pg")
	if err != nil {
		b.log.Error().Err(err).Msg("Failed searching for birthday gif")
		return
	}
	if len(res) != 0 {
		img = &res[0].Images.FixedWidth
	}
	return
}

func isTodayFilter(now func() time.Time) func(v models.Birthday, i int) bool {
	return func(v models.Birthday, i int) bool {
		m1, d1, h1 := format(now().UTC())
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
