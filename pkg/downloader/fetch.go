package downloader

import (
	"fmt"
	"net/http"
)

func (d *Downloader) fetch(url string) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", "Mozilla/5.0")
	response, err := d.client.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status code is: %v", response.Status)
	}
	return response, nil
}
