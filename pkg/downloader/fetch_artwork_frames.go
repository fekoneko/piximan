package downloader

import (
	"encoding/json"
	"fmt"

	"github.com/fekoneko/piximan/pkg/downloader/dto"
	"github.com/fekoneko/piximan/pkg/encode"
)

func (d *Downloader) fetchArtworkFrames(id uint64) (string, []encode.Frame, error) {
	url := fmt.Sprintf("https://www.pixiv.net/ajax/illust/%v/ugoira_meta", id)
	body, err := d.fetch(url)
	if err != nil {
		return "", nil, err
	}

	var unmarshalled dto.Response[dto.FramesData]
	if err := json.Unmarshal(body, &unmarshalled); err != nil {
		return "", nil, err
	}

	url, frames := unmarshalled.Body.FromDto()

	return url, frames, nil
}
