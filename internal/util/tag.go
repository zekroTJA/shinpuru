package util

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/bwmarrin/snowflake"
)

type Tag struct {
	ID        snowflake.ID
	Ident     string
	CreatorID string
	GuildID   string
	Content   string
	Created   time.Time
	LastEdit  time.Time
}

type author struct {
	nameTag  string
	imageURL string
}

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
		Color:       ColorEmbedDefault,
		Footer: &discordgo.MessageEmbedFooter{
			Text: footer,
		},
	}
}

func (t *Tag) AsEntry(s *discordgo.Session) string {
	author := t.formattedAuthor(s)

	return fmt.Sprintf("**%s** by %s [`%s`]", t.Ident, author.nameTag, t.ID)
}

func (t *Tag) RawContent() string {
	return fmt.Sprintf("```md\n%s\n```", t.Content)
}
