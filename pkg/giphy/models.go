package giphy

import "fmt"

type response[T any] struct {
	Data T `json:"data"`
}

type Image struct {
	Url    string `json:"url"`
	Width  string `json:"width"`
	Height string `json:"height"`
	Size   string `json:"size"`
}

type Images struct {
	Still480W              Image `json:"480w_still"`
	Downsized              Image `json:"downsized"`
	DownsizedLarge         Image `json:"downsized_large"`
	DownsizedMedium        Image `json:"downsized_medium"`
	DownsizedSmall         Image `json:"downsized_small"`
	DownsizedStill         Image `json:"downsized_still"`
	FixedHeight            Image `json:"fixed_height"`
	FixedHeightDownsampled Image `json:"fixed_height_downsampled"`
	FixedHeightSmall       Image `json:"fixed_height_small"`
	FixedHeightSmallStill  Image `json:"fixed_height_small_still"`
	FixedHeightStill       Image `json:"fixed_height_still"`
	FixedWidth             Image `json:"fixed_width"`
	FixedWidthDownsampled  Image `json:"fixed_width_downsampled"`
	FixedWidthSmall        Image `json:"fixed_width_small"`
	FixedWidthSmallStill   Image `json:"fixed_width_small_still"`
	FixedWidthStill        Image `json:"fixed_width_still"`
	Hd                     Image `json:"hd"`
	Looping                Image `json:"looping"`
	Original               Image `json:"original"`
	OriginalMp4            Image `json:"original_mp4"`
	OriginalStill          Image `json:"original_still"`
	Preview                Image `json:"preview"`
	PreviewGif             Image `json:"preview_gif"`
	PreviewWebp            Image `json:"preview_webp"`
}

type Gif struct {
	Type     string `json:"type"`
	Id       string `json:"id"`
	Slug     string `json:"slug"`
	Url      string `json:"url"`
	BitlyUrl string `json:"bitly_url"`
	EmbedUrl string `json:"embed_url"`
	Username string `json:"username"`
	Source   string `json:"source"`
	Rating   string `json:"rating"`
	Title    string `json:"string"`
	Images   Images `json:"images"`
}

type Error struct {
	Message string `json:"message"`
	Code    int
}

func (e Error) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}
