package downloader

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/fekoneko/piximan/pkg/collection/work"
	"github.com/fekoneko/piximan/pkg/downloader/dto"
)

func (d *Downloader) fetchWork(id uint64) (*work.Work, error) {
	url := fmt.Sprintf("https://www.pixiv.net/ajax/illust/%v", id)
	response, err := d.fetch(url)
	if err != nil {
		return nil, err
	}

	work, err := workFromResponse(response)
	if err != nil {
		return nil, err
	}

	return work, nil
}

func workFromResponse(body []byte) (*work.Work, error) {
	var unmarshalled dto.Response[dto.Work]
	if err := json.Unmarshal(body, &unmarshalled); err != nil {
		return nil, err
	}

	return unmarshalled.Body.ToWork(time.Now()), nil
}
