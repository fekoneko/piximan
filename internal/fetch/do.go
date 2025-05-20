package fetch

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/fekoneko/piximan/internal/logext"
)

const BUFFER_SIZE = 4096
const PXIMG_PENDING_LIMIT = 5

func Do(client *http.Client, url string, onProgress func(int, int)) ([]byte, error) {
	request, err := newRequest(url)
	if err != nil {
		return nil, err
	}
	lock(request)
	defer unlock(request)

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
) ([]byte, error) {
	request, err := newRequest(url)
	if err != nil {
		return nil, err
	}
	lock(request)
	defer unlock(request)

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

const PXIMG_DELAY = time.Second * 1
const DEFAULT_DELAY = time.Second * 2

var numPximgPending = 0
var pximgCond = sync.NewCond(&sync.Mutex{})
var defaultMutex = sync.Mutex{}
var prevDefaultTime time.Time

func lock(request *http.Request) {
	if request == nil || request.URL == nil {
		return
	}

	switch request.URL.Host {
	case "i.pximg.net":
		pximgCond.L.Lock()
		for numPximgPending >= PXIMG_PENDING_LIMIT {
			pximgCond.Wait()
		}
		numPximgPending++
		pximgCond.L.Unlock()

	default:
		defaultMutex.Lock()
		duration := time.Until(prevDefaultTime.Add(DEFAULT_DELAY))
		time.Sleep(duration)
	}
}

func unlock(request *http.Request) {
	switch request.URL.Host {
	case "i.pximg.net":
		pximgCond.L.Lock()
		numPximgPending--
		pximgCond.Broadcast()
		pximgCond.L.Unlock()

	default:
		prevDefaultTime = time.Now()
		defaultMutex.Unlock()
	}
}
