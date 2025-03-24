package downloader

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/fekoneko/piximan/pkg/collection/work"
)

func (d *Downloader) fetchWork(id uint64) (*work.Work, error) {
	url := fmt.Sprintf("https://www.pixiv.net/ajax/illust/%v", id)
	request, _ := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Add("User-Agent", "Mozilla/5.0")
	response, err := d.client.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status code is: %v", response.Status)
	}

	work, err := workFromResponse(response)
	if err != nil {
		return nil, err
	}

	return work, nil
}

func workFromResponse(response *http.Response) (*work.Work, error) {
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var unmarshalled ApiResponse[work.ApiDto]
	if err := json.Unmarshal(body, &unmarshalled); err != nil {
		return nil, err
	}

	return work.FromApiDto(&unmarshalled.Body, time.Now()), nil
}
