// Package embedbuilder provides a builder pattern
// to create discordgo message embeds.
package embedbuilder

import "github.com/bwmarrin/discordgo"

// EmbedBuilder provides a builder pattern to
// create a discordgo message embed.
type EmbedBuilder struct {
	emb *discordgo.MessageEmbed
}

// New returns a fresh EmbedBuilder.
func New() *EmbedBuilder {
	return &EmbedBuilder{&discordgo.MessageEmbed{}}
}

// WithAuthor sets an author to the embed.
func (b *EmbedBuilder) WithAuthor(name, url, iconUrl, proxyIconUrl string) *EmbedBuilder {
	b.emb.Author = &discordgo.MessageEmbedAuthor{
		URL:          url,
		Name:         name,
		IconURL:      iconUrl,
		ProxyIconURL: proxyIconUrl,
	}
	return b
}

// WithColor sets a color to the embed.
func (b *EmbedBuilder) WithColor(color int) *EmbedBuilder {
	b.emb.Color = color
	return b
}

// WithAuthor adds an author to the embed.
func (b *EmbedBuilder) WithDescription(description string) *EmbedBuilder {
	b.emb.Description = description
	return b
}

// AddField adds a field to the embed.
func (b *EmbedBuilder) AddField(name, value string, inline ...bool) *EmbedBuilder {
	if value == "" {
		value = "*nil*"
	}
	field := &discordgo.MessageEmbedField{
		Name:   name,
		Value:  value,
		Inline: len(inline) > 0 && inline[0],
	}
	if b.emb.Fields == nil {
		b.emb.Fields = []*discordgo.MessageEmbedField{field}
	} else {
		b.emb.Fields = append(b.emb.Fields, field)
	}
	return b
}

// AddInlineField adds an inline field to the embed.
func (b *EmbedBuilder) AddInlineField(name, value string) *EmbedBuilder {
	return b.AddField(name, value, true)
}

// WithFooter sets a footer to the embed.
func (b *EmbedBuilder) WithFooter(text, iconUrl, proxyIconUrl string) *EmbedBuilder {
	b.emb.Footer = &discordgo.MessageEmbedFooter{
		Text:         text,
		IconURL:      iconUrl,
		ProxyIconURL: proxyIconUrl,
	}
	return b
}

// WithImage sets an image to the embed.
func (b *EmbedBuilder) WithImage(url, proxyUrl string, width, height int) *EmbedBuilder {
	b.emb.Image = &discordgo.MessageEmbedImage{
		URL:      url,
		ProxyURL: proxyUrl,
		Width:    width,
		Height:   height,
	}
	return b
}

// WithProvider sets a provider to the embed.
func (b *EmbedBuilder) WithProvider(name, url string) *EmbedBuilder {
	b.emb.Provider = &discordgo.MessageEmbedProvider{
		Name: name,
		URL:  url,
	}
	return b
}

// WithThumbnail sets a thumbnail to the embed.
func (b *EmbedBuilder) WithThumbnail(url, proxyUrl string, width, height int) *EmbedBuilder {
	b.emb.Thumbnail = &discordgo.MessageEmbedThumbnail{
		URL:      url,
		ProxyURL: proxyUrl,
		Width:    width,
		Height:   height,
	}
	return b
}

// WithTimestamp sets a timestamp to the embed.
func (b *EmbedBuilder) WithTimestamp(timestamp string) *EmbedBuilder {
	b.emb.Timestamp = timestamp
	return b
}

// WithTitle sets a title to the embed.
func (b *EmbedBuilder) WithTitle(title string) *EmbedBuilder {
	b.emb.Title = title
	return b
}

// AsType sets the type to the embed.
func (b *EmbedBuilder) AsType(typ discordgo.EmbedType) *EmbedBuilder {
	b.emb.Type = typ
	return b
}

// WithURL sets the URL to the embed.
func (b *EmbedBuilder) WithURL(url string) *EmbedBuilder {
	b.emb.URL = url
	return b
}

// WithFooter sets a video to the embed.
func (b *EmbedBuilder) WithVideo(url string, width, height int) *EmbedBuilder {
	b.emb.Video = &discordgo.MessageEmbedVideo{
		URL:    url,
		Width:  width,
		Height: height,
	}
	return b
}

// Build returns the result embed.
func (b *EmbedBuilder) Build() *discordgo.MessageEmbed {
	return b.emb
}
