package tag

import (
	"fmt"
	"time"

	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekrotja/discordgo"

	"github.com/bwmarrin/snowflake"
)

// Tag wraps a chat tag object.
type Tag struct {
	ID        snowflake.ID
	Ident     string
	CreatorID string
	GuildID   string
	Content   string
	Created   time.Time
	LastEdit  time.Time
}

// author wraps a name tag and avatar imageURL
// of an author user.
type author struct {
	nameTag  string
	imageURL string
}

// AsEmbed creates a discordgo.MessageEmbed from
// the tag information.
func (t *Tag) AsEmbed(s *discordgo.Session) *discordgo.MessageEmbed {
	footer := ""

	author := t.formattedAuthor(s)

	if t.Created == t.LastEdit {
		footer = fmt.Sprintf("Created %s",
			t.Created.Format(time.RFC822))
	} else {
		footer = fmt.Sprintf("Created %s | Last edit %s",
			t.Created.Format(time.RFC822), t.LastEdit.Format(time.RFC822))
	}

	return &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    author.nameTag,
			IconURL: author.imageURL,
		},
		Description: t.Content,
		Color:       static.ColorEmbedDefault,
		Footer: &discordgo.MessageEmbedFooter{
			Text: footer,
		},
	}
}

// AsEntry returns a single formatted string line
// to represent a tag.
func (t *Tag) AsEntry(s *discordgo.Session) string {
	author := t.formattedAuthor(s)

	return fmt.Sprintf("**%s** by %s [`%s`]", t.Ident, author.nameTag, t.ID)
}

// RawContent returns the content of the tags body
// in a markdown code embed for a discord message.
func (t *Tag) RawContent() string {
	return fmt.Sprintf("```md\n%s\n```", t.Content)
}

// formattedAuthor returns an author object from the
// CreatorID of the tag.
func (t *Tag) formattedAuthor(s *discordgo.Session) *author {
	authorF := new(author)
	author, err := s.GuildMember(t.GuildID, t.CreatorID)
	if err == nil && author != nil {
		authorF.nameTag = fmt.Sprintf("%s#%s", author.User.Username, author.User.Discriminator)
		authorF.imageURL = author.User.AvatarURL("")
	} else {
		authorF.nameTag = fmt.Sprintf("<not on guild> (%s)", t.CreatorID)
	}

	return authorF
}
