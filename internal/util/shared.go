package util

import (
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func ExtractImageURLFromMessage(text string, attachments []*discordgo.MessageAttachment) (string, string) {
	var imgLink string

	if len(attachments) > 0 {
		imgLink = attachments[0].URL
	} else {
		imgRx := regexp.MustCompile(`https?:\/\/([\w-]+\.)+([\w-]+)(\/[\w-]+)*.*\.(png|jpg|jpeg|gif|ico|tiff|img|bmp)`)
		rxResult := imgRx.FindString(text)
		if rxResult != "" {
			text = strings.Replace(text, rxResult, "", 1)
			imgLink = rxResult
		}
	}

	return text, imgLink
}
