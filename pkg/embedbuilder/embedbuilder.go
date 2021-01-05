package embedbuilder

import "github.com/bwmarrin/discordgo"

type EmbedBuilder struct {
	emb *discordgo.MessageEmbed
}

func New() *EmbedBuilder {
	return &EmbedBuilder{&discordgo.MessageEmbed{}}
}

func (b *EmbedBuilder) WithAuthor(name, url, iconUrl, proxyIconUrl string) *EmbedBuilder {
	b.emb.Author = &discordgo.MessageEmbedAuthor{
		URL:          url,
		Name:         name,
		IconURL:      iconUrl,
		ProxyIconURL: proxyIconUrl,
	}
	return b
}

func (b *EmbedBuilder) WithColor(color int) *EmbedBuilder {
	b.emb.Color = color
	return b
}

func (b *EmbedBuilder) WithDescription(description string) *EmbedBuilder {
	b.emb.Description = description
	return b
}

func (b *EmbedBuilder) AddField(name, value string, inline bool) *EmbedBuilder {
	if value == "" {
		value = "*nil*"
	}
	field := &discordgo.MessageEmbedField{
		Name:   name,
		Value:  value,
		Inline: inline,
	}
	if b.emb.Fields == nil {
		b.emb.Fields = []*discordgo.MessageEmbedField{field}
	} else {
		b.emb.Fields = append(b.emb.Fields, field)
	}
	return b
}

func (b *EmbedBuilder) WithFooter(text, iconUrl, proxyIconUrl string) *EmbedBuilder {
	b.emb.Footer = &discordgo.MessageEmbedFooter{
		Text:         text,
		IconURL:      iconUrl,
		ProxyIconURL: proxyIconUrl,
	}
	return b
}

func (b *EmbedBuilder) WithImage(url, proxyUrl string, width, height int) *EmbedBuilder {
	b.emb.Image = &discordgo.MessageEmbedImage{
		URL:      url,
		ProxyURL: proxyUrl,
		Width:    width,
		Height:   height,
	}
	return b
}

func (b *EmbedBuilder) WithProvider(name, url string) *EmbedBuilder {
	b.emb.Provider = &discordgo.MessageEmbedProvider{
		Name: name,
		URL:  url,
	}
	return b
}

func (b *EmbedBuilder) WithThumbnail(url, proxyUrl string, width, height int) *EmbedBuilder {
	b.emb.Thumbnail = &discordgo.MessageEmbedThumbnail{
		URL:      url,
		ProxyURL: proxyUrl,
		Width:    width,
		Height:   height,
	}
	return b
}

func (b *EmbedBuilder) WithTimestamp(timestamp string) *EmbedBuilder {
	b.emb.Timestamp = timestamp
	return b
}

func (b *EmbedBuilder) WithTitle(title string) *EmbedBuilder {
	b.emb.Title = title
	return b
}

func (b *EmbedBuilder) AsType(typ discordgo.EmbedType) *EmbedBuilder {
	b.emb.Type = typ
	return b
}

func (b *EmbedBuilder) WithURL(url string) *EmbedBuilder {
	b.emb.URL = url
	return b
}

func (b *EmbedBuilder) WithVideo(url string, width, height int) *EmbedBuilder {
	b.emb.Video = &discordgo.MessageEmbedVideo{
		URL:    url,
		Width:  width,
		Height: height,
	}
	return b
}

func (b *EmbedBuilder) Build() *discordgo.MessageEmbed {
	return b.emb
}
