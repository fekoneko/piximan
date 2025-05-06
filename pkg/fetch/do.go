package fetch

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/fekoneko/piximan/pkg/logext"
)

const BUFFER_SIZE = 4096

func Do(client http.Client, url string, onProgress func(int, int)) ([]byte, error) {
	request, err := newRequest(url)
	if err != nil {
		return nil, err
	}
	delayRequest(request)

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
	client http.Client, url string, sessionId string, onProgress func(int, int),
) ([]byte, error) {
	request, err := newRequest(url)
	if err != nil {
		return nil, err
	}
	delayRequest(request)

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
	client http.Client, request *http.Request, onProgress func(int, int),
) ([]byte, error) {
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// TODO: should suspend and retry later if network issues occured
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status code is: %v", response.Status)
	}

	if response.ContentLength <= 0 {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		onProgress(len(body), len(body))

		return body, nil
	} else {
		body := make([]byte, 0, response.ContentLength)
		buffer := make([]byte, BUFFER_SIZE)

		for {
			readLength, err := response.Body.Read(buffer)
			if err == io.EOF {
				return body, nil
			} else if err != nil {
				return nil, err
			} else {
				body = append(body, buffer[:readLength]...)
				onProgress(len(body), int(response.ContentLength))
			}
		}
	}
}

var pixivMutex = sync.Mutex{}
var pximgMutex = sync.Mutex{}
var otherMutex = sync.Mutex{}

const pixivDelay = time.Second * 3
const pximgDelay = time.Second * 1
const otherDelay = time.Second * 3

var prevPixivTime time.Time
var prevPximgTime time.Time
var prevOtherTime time.Time

func delayRequest(request *http.Request) {
	if request == nil || request.URL == nil {
		return
	}

	switch request.URL.Host {
	case "www.pixiv.net":
		pixivMutex.Lock()
		duration := time.Until(prevPixivTime.Add(pixivDelay))
		time.Sleep(duration)
		prevPixivTime = time.Now()
		pixivMutex.Unlock()

	case "i.pximg.net":
		// TODO: download in batches of 5 for i.pximg.net
		pximgMutex.Lock()
		duration := time.Until(prevPximgTime.Add(pximgDelay))
		time.Sleep(duration)
		prevPximgTime = time.Now()
		pximgMutex.Unlock()

	default:
		otherMutex.Lock()
		duration := time.Until(prevOtherTime.Add(otherDelay))
		time.Sleep(duration)
		prevOtherTime = time.Now()
		otherMutex.Unlock()
	}
}
