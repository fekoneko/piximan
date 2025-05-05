package fetch

import (
	"fmt"
	"io"
	"net/http"

	"github.com/fekoneko/piximan/pkg/logext"
)

const BUFFER_SIZE = 4096

// TODO: add some delay to avoid flooding
//       - how much delay is needed?
//       - should we do separate dalays for pixiv.net and pximg.net?
//       - should the delay be different for authorized requests?

func Do(client http.Client, url string) ([]byte, error) {
	if request, err := newRequest(url); err == nil {
		logext.Request(url)
		return doWithRequest(client, request, func(readLength int64, totalLength int64) {
			fmt.Printf("progress %v/%v\n", readLength, totalLength)
		})
	} else {
		return nil, err
	}
}

func DoAuthorized(client http.Client, url string, sessionId string) ([]byte, error) {
	if request, err := newRequest(url); err == nil {
		request.Header.Add("Cookie", "PHPSESSID="+sessionId)
		logext.AuthorizedRequest(url)
		return doWithRequest(client, request, func(readLength int64, totalLength int64) {
			fmt.Printf("progress %v/%v\n", readLength, totalLength)
		})
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

func doWithRequest(client http.Client, request *http.Request, onProgress func(int64, int64)) ([]byte, error) {
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
		length := int64(len(body))
		onProgress(length, length)

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
				length := int64(len(body))
				onProgress(length, response.ContentLength)
			}
		}
	}
}
