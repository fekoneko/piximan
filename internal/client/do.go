package client

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/fekoneko/piximan/internal/logger"
)

const BUFFER_SIZE = 4096

func (c *Client) Do(url string, onProgress func(int, int)) ([]byte, http.Header, error) {
	request, err := newRequest(url)
	if err != nil {
		return nil, nil, err
	}

	return c.doWithRequest(request, logger.Request, onProgress)
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
	request.Header.Add("Cookie", "PHPSESSID="+sessionId)

	return c.doWithRequest(request, logger.AuthorizedRequest, onProgress)
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
	request *http.Request, log func(url string) (func(), func(int, int)), onProgress func(int, int),
) ([]byte, http.Header, error) {
	c.startRequest(request)
	defer c.requestDone(request)

	retryDelay := time.Duration(0)
	for {
		removeBar, updateBar := log(request.URL.String())
		body, headers, err := c.tryRequest(request, func(current int, total int) {
			updateBar(current, total)
			if onProgress != nil {
				onProgress(current, total)
			}
		})
		removeBar()

		if err == nil {
			return body, headers, nil
		} else if errors.Is(err, StatusError{}) {
			return nil, nil, err
		}

		logger.Warning("request failed: %v (retrying in %v)",
			err, retryDelay,
		)
		time.Sleep(retryDelay)
		retryDelay = retryDelay*2 + time.Second*10
	}
}

func (c *Client) tryRequest(request *http.Request, onProgress func(int, int)) ([]byte, http.Header, error) {
	response, err := c.client().Do(request)
	if err != nil {
		return nil, nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, response.Header, StatusError{
			error: fmt.Errorf("response status code is: %v", response.Status),
			code:  response.StatusCode,
		}
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

func (c *Client) startRequest(request *http.Request) {
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

func (c *Client) requestDone(request *http.Request) {
	switch request.URL.Host {
	case "i.pximg.net":
		c.pximgRequestGroup.Done()
	default:
		c.defaultRequestGroup.Done()
	}
}
