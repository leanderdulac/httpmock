package httpmock

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {

	ms := New(t, true, true)
	require.NotNil(t, ms.server, "server must be instantiated")

}

func TestGet(t *testing.T) {
	ms := New(t, true, true)
	ms.AddGetHandler("/test", func(w http.ResponseWriter, r *http.Request) {
		ServerDumpRequest(r, true)
	})

	resp, err := http.Get(ms.GetURL("/test"))
	require.NoError(t, err, "server are not fine :(")
	require.NoError(t, ClientDumpResponse(resp, true), "response are not ok :(")

	ms.Close()
}

func TestAddGetHandler(t *testing.T) {
	ms := New(t, false, false)
	ms.AddGetHandler("/test", func(w http.ResponseWriter, r *http.Request) {
		ServerDumpRequest(r, true)
	})

	resp, err := http.Get(ms.GetURL("/test"))
	require.NoError(t, err, "server are not fine :(")
	require.NoError(t, ClientDumpResponse(resp, true), "response are not ok :(")

	ms.Close()
}

func TestAddPostHandler(t *testing.T) {
	ms := New(t, false, false)
	ms.AddPostHandler("/test", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()

		require.NoError(t, err, "error not expected here")
		require.Equal(t, "test123", string(body), "unexpected body :/")
	})

	resp, err := http.Post(ms.GetURL("/test"), "text/plain", bytes.NewBufferString("test123"))
	require.NoError(t, err, "server are not fine :(")
	require.NoError(t, ClientDumpResponse(resp, true), "response are not ok :(")

	ms.Close()
}
