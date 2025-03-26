package downloader

import (
	"encoding/json"
	"fmt"

	"github.com/fekoneko/piximan/pkg/downloader/dto"
	"github.com/fekoneko/piximan/pkg/encode"
)

func (d *Downloader) fetchFramesData(id uint64) (string, []encode.Frame, error) {
	url := fmt.Sprintf("https://www.pixiv.net/ajax/illust/%v/ugoira_meta", id)
	response, err := d.fetch(url)
	if err != nil {
		return "", nil, err
	}

	url, frames, err := dataFromResponse(response)
	if err != nil {
		return "", nil, err
	}

	return url, frames, nil
}

func dataFromResponse(body []byte) (string, []encode.Frame, error) {
	var unmarshalled dto.Response[dto.FramesData]
	if err := json.Unmarshal(body, &unmarshalled); err != nil {
		return "", nil, err
	}

	return unmarshalled.Body.Src, unmarshalled.Body.Frames, nil
}
