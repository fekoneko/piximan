package fetch

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/fekoneko/piximan/internal/logext"
	"github.com/fekoneko/piximan/internal/syncext"
)

const BUFFER_SIZE = 4096

func Do(client *http.Client, url string, onProgress func(int, int)) ([]byte, http.Header, error) {
	request, err := newRequest(url)
	if err != nil {
		return nil, nil, err
	}
	start(request)
	defer done(request)

	removeBar, updateBar := logext.Request(url)
	defer removeBar()

	return doWithRequest(client, request, func(current int, total int) {
		updateBar(current, total)
		if onProgress != nil {
			onProgress(current, total)
		}
	})
}

func DoAuthorized(
	client *http.Client, url string, sessionId string, onProgress func(int, int),
) ([]byte, http.Header, error) {
	request, err := newRequest(url)
	if err != nil {
		return nil, nil, err
	}
	start(request)
	defer done(request)

	request.Header.Add("Cookie", "PHPSESSID="+sessionId)
	removeBar, updateBar := logext.AuthorizedRequest(url)
	defer removeBar()

	return doWithRequest(client, request, func(current int, total int) {
		updateBar(current, total)
		if onProgress != nil {
			onProgress(current, total)
		}
	})
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

func doWithRequest(
	client *http.Client, request *http.Request, onProgress func(int, int),
) ([]byte, http.Header, error) {
	response, err := client.Do(request)
	if err != nil { // TODO: should suspend and retry later if network issues occured
		return nil, response.Header, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, response.Header, fmt.Errorf("response status code is: %v", response.Status)
	}

	if response.ContentLength <= 0 {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, response.Header, err
		}
		onProgress(len(body), len(body))

		return body, response.Header, nil
	} else {
		body := make([]byte, 0, response.ContentLength)
		buffer := make([]byte, BUFFER_SIZE)

		for {
			readLength, err := response.Body.Read(buffer)
			if err == io.EOF {
				return body, response.Header, nil
			} else if err != nil {
				return nil, response.Header, err
			} else {
				body = append(body, buffer[:readLength]...)
				onProgress(len(body), int(response.ContentLength))
			}
		}
	}
}

var piximgRequestGroup = syncext.NewRequestGroup(5, time.Second*1)
var defaultRequestGroup = syncext.NewRequestGroup(1, time.Second*2)

func start(request *http.Request) {
	if request == nil || request.URL == nil {
		return
	}

	switch request.URL.Host {
	case "i.pximg.net":
		piximgRequestGroup.Start()
	default:
		defaultRequestGroup.Start()
	}
}

func done(request *http.Request) {
	switch request.URL.Host {
	case "i.pximg.net":
		piximgRequestGroup.Done()
	default:
		defaultRequestGroup.Done()
	}
}
