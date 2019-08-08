package util

import (
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var ImgUrlRx = regexp.MustCompile(`https?:\/\/([\w-]+\.)+([\w-]+)(\/[\w-]+)*.*\.(png|jpg|jpeg|gif|ico|tiff|img|bmp|mp4)`)
var ImgUrlSRx = regexp.MustCompile(`^https?:\/\/([\w-]+\.)+([\w-]+)(\/[\w-]+)*.*\.(png|jpg|jpeg|gif|ico|tiff|img|bmp|mp4)$`)

func ExtractImageURLFromMessage(text string, attachments []*discordgo.MessageAttachment) (string, string) {
	var imgLink string

	if len(attachments) > 0 {
		imgLink = attachments[0].URL
	} else {
		rxResult := ImgUrlRx.FindString(text)
		if rxResult != "" {
			text = strings.Replace(text, rxResult, "", 1)
			imgLink = rxResult
		}
	}

	return text, imgLink
}
