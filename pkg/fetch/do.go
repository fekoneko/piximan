package fetch

import (
	"fmt"
	"io"
	"net/http"

	"github.com/fekoneko/piximan/pkg/logext"
)

func Do(client http.Client, url string) ([]byte, error) {
	if request, err := newRequest(url); err == nil {
		logext.Fetch(url)
		return doWithRequest(client, request)
	} else {
		return nil, err
	}
}

func DoAuthorized(client http.Client, url string, sessionId string) ([]byte, error) {
	if request, err := newRequest(url); err == nil {
		request.Header.Add("Cookie", "PHPSESSID="+sessionId)
		logext.AuthorizedFetch(url)
		return doWithRequest(client, request)
	} else {
		return nil, err
	}
}

func newRequest(url string) (*http.Request, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", "Mozilla/5.0")
	request.Header.Add("Referer", "https://www.pixiv.net/")

	return request, nil
}

func doWithRequest(client http.Client, request *http.Request) ([]byte, error) {
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	// TODO: should suspend and retry later if network issues occured
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status code is: %v", response.Status)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
