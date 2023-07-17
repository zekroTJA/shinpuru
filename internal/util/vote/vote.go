package vote

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"strings"
	"time"

	"github.com/wcharczuk/go-chart/drawing"
	"github.com/zekroTJA/shinpuru/internal/services/timeprovider"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/discordutil"

	"github.com/bwmarrin/discordgo"
	"github.com/wcharczuk/go-chart"
)

// VoteState defines the lifecycle state of a Vote.
type VoteState int

const (
	VoteStateOpen VoteState = iota
	VoteStateClosed
	VoteStateClosedNC
	VoteStateExpired
)

// VotesRunning maps running vote IDs to
// their vote instances.
var VotesRunning = map[string]Vote{}

// VoteEmotes contains the emotes used to tick a vote.
var VoteEmotes = strings.Fields("\u0031\u20E3 \u0032\u20E3 \u0033\u20E3 \u0034\u20E3 \u0035\u20E3 \u0036\u20E3 \u0037\u20E3 \u0038\u20E3 \u0039\u20E3 \u0030\u20E3")

// Vote wraps the information and current
// state of a vote and its ticks.
type Vote struct {
	ID            string
	MsgID         string
	CreatorID     string
	GuildID       string
	ChannelID     string
	Description   string
	ImageURL      string
	Expires       time.Time
	Possibilities []string
	Ticks         map[string]*Tick
}

// Tick wraps a user ID and the index of
// the selection ticked.
type Tick struct {
	UserID string
	Tick   int
}

// Unmarshal tries to deserialize a raw data string
// to a Vote object. Errors occured during
// deserialization are returned as well.
func Unmarshal(data string) (Vote, error) {
	rawData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return Vote{}, err
	}

	var res Vote
	buffer := bytes.NewBuffer(rawData)
	gobdec := gob.NewDecoder(buffer)
	err = gobdec.Decode(&res)
	return res, err
}

// Marshal serializes the vote to a raw data string.
//
// The vote object is encoded to a byte array using
// the gob encoder and then encoded to a base64 string.
func (v *Vote) Marshal() (string, error) {
	var buffer bytes.Buffer
	gobenc := gob.NewEncoder(&buffer)
	err := gobenc.Encode(v)
	if err != nil {
		return "", err
	}
	gobres := buffer.Bytes()
	res := base64.StdEncoding.EncodeToString(gobres)
	return res, nil
}

// AsEmbed creates a discordgo.MessageEmbed from the
// vote. If voteState is passed, the state will be
// displayed as well. Otherwise, it will be assumed
// that the vote is open.
//
// If voteState is VoteStateClosed or VoteStateExpired,
// a pie chart will be generated representing the
// distribution of vote ticks and sent as image to
// the channel.
func (v *Vote) AsEmbed(s *discordgo.Session, voteState ...VoteState) (*discordgo.MessageEmbed, error) {
	state := VoteStateOpen
	if len(voteState) > 0 {
		state = voteState[0]
	}

	creator, err := s.User(v.CreatorID)
	if err != nil {
		return nil, err
	}
	title := "Open Vote"
	color := static.ColorEmbedDefault

	switch state {
	case VoteStateClosed, VoteStateClosedNC:
		title = "Vote closed"
		color = static.ColorEmbedOrange
	case VoteStateExpired:
		title = "Vote expired"
		color = static.ColorEmbedViolett
	}

	totalTicks := make(map[int]int)
	for _, t := range v.Ticks {
		if _, ok := totalTicks[t.Tick]; !ok {
			totalTicks[t.Tick] = 1
		} else {
			totalTicks[t.Tick]++
		}
	}

	description := v.Description + "\n\n"
	for i, p := range v.Possibilities {
		description += fmt.Sprintf("%s    %s  -  `%d`\n", VoteEmotes[i], p, totalTicks[i])
	}

	footerText := fmt.Sprintf("ID: %s", v.ID)
	if (v.Expires != time.Time{} && state == VoteStateOpen) {
		footerText = fmt.Sprintf("%s | Expires: %s", footerText, v.Expires.Format("01/02 15:04 MST"))
	}

	emb := &discordgo.MessageEmbed{
		Color:       color,
		Title:       title,
		Description: description,
		Author: &discordgo.MessageEmbedAuthor{
			IconURL: creator.AvatarURL("16x16"),
			Name:    creator.String(),
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: footerText,
		},
	}

	if len(totalTicks) > 0 && (state == VoteStateClosed || state == VoteStateExpired) {

		values := make([]chart.Value, len(v.Possibilities))

		for i, p := range v.Possibilities {
			values[i] = chart.Value{
				Value: float64(totalTicks[i]),
				Label: p,
			}
		}

		pie := chart.PieChart{
			Width:  512,
			Height: 512,
			Values: values,
			Background: chart.Style{
				FillColor: drawing.ColorTransparent,
			},
		}

		imgData := []byte{}
		buff := bytes.NewBuffer(imgData)
		err = pie.Render(chart.PNG, buff)
		if err != nil {
			return nil, err
		}

		_, err := s.ChannelMessageSendComplex(v.ChannelID, &discordgo.MessageSend{
			File: &discordgo.File{
				Name:   fmt.Sprintf("vote_chart_%s.png", v.ID),
				Reader: buff,
			},
			Reference: &discordgo.MessageReference{
				MessageID: v.MsgID,
				ChannelID: v.ChannelID,
				GuildID:   v.GuildID,
			},
		})
		if err != nil {
			return nil, err
		}
	}

	if v.ImageURL != "" {
		emb.Image = &discordgo.MessageEmbedImage{
			URL: v.ImageURL,
		}
	}

	return emb, nil
}

// AsField creates a discordgo.MessageEmbedField from
// the vote information.
func (v *Vote) AsField() *discordgo.MessageEmbedField {
	shortenedDescription := v.Description
	if len(shortenedDescription) > 200 {
		shortenedDescription = shortenedDescription[200:] + "..."
	}

	expiresTxt := "never"
	if (v.Expires != time.Time{}) {
		expiresTxt = v.Expires.Format("01/02 15:04 MST")
	}

	return &discordgo.MessageEmbedField{
		Name: "VID: " + v.ID,
		Value: fmt.Sprintf("**Description:** %s\n**Expires:** %s\n`%d votes`\n[*jump to msg*](%s)",
			shortenedDescription, expiresTxt, len(v.Ticks), discordutil.GetMessageLink(&discordgo.Message{
				ID:        v.MsgID,
				ChannelID: v.ChannelID,
			}, v.GuildID)),
	}
}

// AddReactions adds the reactions to the votes message
// for each selection possibility.
//
// Vote emotes are used from VoteEmotes.
func (v *Vote) AddReactions(s *discordgo.Session) error {
	for i := 0; i < len(v.Possibilities); i++ {
		err := s.MessageReactionAdd(v.ChannelID, v.MsgID, VoteEmotes[i])
		if err != nil {
			return err
		}
	}
	return nil
}

// Tick sets the tick for the specified user to the vote.
func (v *Vote) Tick(s *discordgo.Session, userID string, tick int) (err error) {
	if userID, err = HashUserID(userID, []byte(v.ID)); err != nil {
		return
	}

	if t, ok := v.Ticks[userID]; ok {
		t.Tick = tick
	} else {
		v.Ticks[userID] = &Tick{
			UserID: userID,
			Tick:   tick,
		}
	}

	emb, err := v.AsEmbed(s)
	if err != nil {
		return
	}

	_, err = s.ChannelMessageEditEmbed(v.ChannelID, v.MsgID, emb)
	return
}

// SetExpire sets the expiration for a vote.
func (v *Vote) SetExpire(s *discordgo.Session, d time.Duration, tp timeprovider.Provider) error {
	v.Expires = tp.Now().Add(d)

	emb, err := v.AsEmbed(s)
	if err != nil {
		return err
	}
	_, err = s.ChannelMessageEditEmbed(v.ChannelID, v.MsgID, emb)

	return err
}

// Close closes the vote and removes it
// from the VotesRunning map.
func (v *Vote) Close(s *discordgo.Session, voteState VoteState) error {
	delete(VotesRunning, v.ID)
	emb, err := v.AsEmbed(s, voteState)
	if err != nil {
		return err
	}
	_, err = s.ChannelMessageEditEmbed(v.ChannelID, v.MsgID, emb)
	if err != nil {
		return err
	}
	err = s.MessageReactionsRemoveAll(v.ChannelID, v.MsgID)
	return err
}
