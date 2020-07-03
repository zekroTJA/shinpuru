package imgstore

import (
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	ImgUrlSRx = regexp.MustCompile(`^https?:\/\/([\w-]+\.)+([\w-]+)(\/[\w-]+)*.*\.(png|jpg|jpeg|gif|ico|tiff|img|bmp|mp4)$`)

	imgUrlRx = regexp.MustCompile(`https?:\/\/([\w-]+\.)+([\w-]+)(\/[\w-]+)*.*\.(png|jpg|jpeg|gif|ico|tiff|img|bmp|mp4)`)
)

// ExtractFromMessage tries to extract an image URL from the passed
// text or message attachments and returns the text of the message
// excluding the matched link and the image link.
func ExtractFromMessage(text string, attachments []*discordgo.MessageAttachment) (resText, imgLink string) {
	resText = text

	if len(attachments) > 0 {
		imgLink = attachments[0].URL
	} else {
		rxResult := imgUrlRx.FindString(text)
		if rxResult != "" {
			resText = strings.Replace(text, rxResult, "", 1)
			imgLink = rxResult
		}
	}

	return resText, imgLink
}
