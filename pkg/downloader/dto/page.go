package dto

type Page struct {
	Urls struct {
		ThumbMini string `json:"thumb_mini"`
		Small     string `json:"small"`
		Regular   string `json:"regular"`
		Original  string `json:"original"`
	} `json:"urls"`
}

func (p *Page) FromDto() *[4]string {
	return &[4]string{
		p.Urls.ThumbMini,
		p.Urls.Small,
		p.Urls.Regular,
		p.Urls.Original,
	}
}
