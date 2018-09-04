package httpmock

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWaitForURLs(t *testing.T) {

	urlList := []string{
		"http://localhost:8080",
		"http://localhost",
		"http://127.0.0.1:8080",
		"http://127.0.0.1",
		"http://google.com:80",
		"http://google.com",
	}

	t.Run("str2url", func(t *testing.T) {
		urls, err := str2url(urlList)
		require.NoError(t, err)
		require.NotNil(t, urls)
		require.Equal(t, len(urlList), len(urls))
	})

	t.Run("Sleep", func(t *testing.T) {

		err := WaitForURLs(urlList, 2)
		require.Error(t, err)
	})
}
