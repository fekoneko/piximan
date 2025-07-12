package dto

import "github.com/fekoneko/piximan/internal/imageext"

type Page struct {
	Urls struct {
		Thumb     *string `json:"thumb"`
		ThumbMini *string `json:"thumb_mini"`
		Small     *string `json:"small"`
		Regular   *string `json:"regular"`
		Original  *string `json:"original"`
	} `json:"urls"`
}

func (p *Page) FromDto(size imageext.Size) *string {
	switch size {
	case imageext.SizeThumbnail:
		if p.Urls.Thumb != nil {
			return p.Urls.Thumb
		} else {
			return p.Urls.ThumbMini
		}
	case imageext.SizeSmall:
		return p.Urls.Small
	case imageext.SizeMedium:
		return p.Urls.Regular
	case imageext.SizeOriginal:
		return p.Urls.Original
	default:
		return nil
	}
}
