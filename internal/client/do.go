package client

import (
	"fmt"
	"io"
	"net/http"

	"github.com/fekoneko/piximan/internal/logext"
)

const BUFFER_SIZE = 4096

func (c *Client) Do(url string, onProgress func(int, int)) ([]byte, http.Header, error) {
	request, err := newRequest(url)
	if err != nil {
		return nil, nil, err
	}
	c.start(request)
	defer c.done(request)

	removeBar, updateBar := logext.Request(url)
	defer removeBar()

	return c.doWithRequest(request, func(current int, total int) {
		updateBar(current, total)
		if onProgress != nil {
			onProgress(current, total)
		}
	})
}

func (c *Client) DoAuthorized(
	url string, onProgress func(int, int),
) ([]byte, http.Header, error) {
	sessionId, authorized := c.sessionId()
	if !authorized {
		return nil, nil, fmt.Errorf("authorization is required")
	}

	request, err := newRequest(url)
	if err != nil {
		return nil, nil, err
	}
	c.start(request)
	defer c.done(request)

	request.Header.Add("Cookie", "PHPSESSID="+sessionId)
	removeBar, updateBar := logext.AuthorizedRequest(url)
	defer removeBar()

	return c.doWithRequest(request, func(current int, total int) {
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

func (c *Client) doWithRequest(
	request *http.Request, onProgress func(int, int),
) ([]byte, http.Header, error) {
	response, err := c.client().Do(request)
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

func (c *Client) start(request *http.Request) {
	if request == nil || request.URL == nil {
		return
	}

	switch request.URL.Host {
	case "i.pximg.net":
		c.pximgRequestGroup.Start()
	default:
		c.defaultRequestGroup.Start()
	}
}

func (c *Client) done(request *http.Request) {
	switch request.URL.Host {
	case "i.pximg.net":
		c.pximgRequestGroup.Done()
	default:
		c.defaultRequestGroup.Done()
	}
}
