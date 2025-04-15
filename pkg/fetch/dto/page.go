package dto

type Page struct {
	Urls struct {
		Thumb     *string `json:"thumb"`
		ThumbMini *string `json:"thumb_mini"`
		Small     *string `json:"small"`
		Regular   *string `json:"regular"`
		Original  *string `json:"original"`
	} `json:"urls"`
}

func (p *Page) FromDto() *[4]string {
	thumbnailUrl := p.Urls.Thumb
	if thumbnailUrl == nil {
		thumbnailUrl = p.Urls.ThumbMini
	}

	if thumbnailUrl == nil || p.Urls.Small == nil || p.Urls.Regular == nil || p.Urls.Original == nil {
		return nil
	}

	return &[4]string{
		*thumbnailUrl,
		*p.Urls.Small,
		*p.Urls.Regular,
		*p.Urls.Original,
	}
}
