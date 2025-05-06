package fetch

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/fekoneko/piximan/pkg/logext"
)

const BUFFER_SIZE = 4096

// TODO: add some delay to avoid flooding
//       - how much delay is needed?
//       - should we do separate dalays for pixiv.net and pximg.net?
//       - should the delay be different for authorized requests?

func Do(client http.Client, url string, onProgress func(int, int)) ([]byte, error) {
	request, err := newRequest(url)
	if err != nil {
		return nil, err
	}

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
	for i := 0; i < 10000; i += rand.Intn(100) {
		onProgress(min(i, 10000), 10000)
		time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
	}
	return nil, fmt.Errorf("test error")

	// response, err := client.Do(request)
	// if err != nil {
	// 	return nil, err
	// }
	// defer response.Body.Close()

	// // TODO: should suspend and retry later if network issues occured
	// if response.StatusCode != http.StatusOK {
	// 	return nil, fmt.Errorf("response status code is: %v", response.Status)
	// }

	// if response.ContentLength <= 0 {
	// 	body, err := io.ReadAll(response.Body)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	onProgress(len(body), len(body))

	// 	return body, nil
	// } else {
	// 	body := make([]byte, 0, response.ContentLength)
	// 	buffer := make([]byte, BUFFER_SIZE)

	// 	for {
	// 		readLength, err := response.Body.Read(buffer)
	// 		if err == io.EOF {
	// 			return body, nil
	// 		} else if err != nil {
	// 			return nil, err
	// 		} else {
	// 			body = append(body, buffer[:readLength]...)
	// 			onProgress(len(body), int(response.ContentLength))
	// 		}
	// 	}
	// }
}
