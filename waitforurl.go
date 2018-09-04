package httpmock

import (
	"errors"
	"net"
	"net/url"
	"time"
)

// ErrWaitForURLs come up when could not wait for url(s)
var ErrWaitForURLs = errors.New("could not wait for url")

func str2url(urlList []string) ([]*url.URL, error) {

	URLs := []*url.URL{}

	for _, rawURL := range urlList {
		URL, err := url.Parse(rawURL)
		if err != nil {
			return nil, err
		}
		URLs = append(URLs, URL)
	}

	return URLs, nil
}

func dial(URLs []*url.URL) bool {

	success := 0

	for _, URL := range URLs {
		conn, err := net.Dial("tcp", URL.Host)
		if err != nil {
			continue
		}
		success++
		conn.Close()
	}
	return len(URLs) == success
}

// WaitForURLs sleep until url list come up
func WaitForURLs(urlList []string, limit int) error {

	URLs, err := str2url(urlList)
	if err != nil {
		return err
	}

	count := 0
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for _ = range ticker.C {

		if count >= limit {
			return ErrWaitForURLs
		}
		count++

		if dial(URLs) {
			return nil
		}
	}
	return ErrWaitForURLs
}
