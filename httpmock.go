package httpmock

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"testing"

	"github.com/fatih/color"
	"github.com/stretchr/testify/require"
)

// MockServer wrapper for httptest.Server
type MockServer struct {
	server *httptest.Server
	t      *testing.T
	get    map[string]GetHandler
	post   map[string]PostHandler
	put    map[string]PutHandler
}

// GetHandler Server GET handler
type GetHandler func(w http.ResponseWriter, r *http.Request)

// PostHandler Server POST handler
type PostHandler func(w http.ResponseWriter, r *http.Request)

// PutHandler Server PUT handler
type PutHandler func(w http.ResponseWriter, r *http.Request)

func dumpServerRequest(r *http.Request, printBody bool) error {
	dump, err := httputil.DumpRequest(r, printBody)
	if err != nil {
		return err
	}
	begin := color.GreenString("server.request.begin")
	end := color.GreenString("server.request.end")
	fmt.Printf("[%s]\n%s\n[%s]\n", begin, dump, end)

	return nil
}

// ClientDumpResponse dump response from server
func ClientDumpResponse(r *http.Response, printBody bool) error {
	dump, err := httputil.DumpResponse(r, printBody)
	if err != nil {
		return err
	}
	begin := color.BlueString("client.response.begin")
	end := color.BlueString("client.response.end")
	fmt.Printf("[%s]\n%s\n[%s]\n", begin, dump, end)

	return nil
}

// ServerDumpRequest dump request from client
func ServerDumpRequest(r *http.Request, printBody bool) error {
	return dumpServerRequest(r, printBody)
}

// New create a MockServer instance
func New(t *testing.T, dumpRequest bool, dumpBody bool) MockServer {

	ms := MockServer{
		t:    t,
		get:  map[string]GetHandler{},
		post: map[string]PostHandler{},
		put:  map[string]PutHandler{},
	}

	ms.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if true == dumpRequest {
			if err := dumpServerRequest(r, dumpBody); err != nil {
				http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			}
		}

		switch r.Method {
		case "GET":
			callback, ok := ms.get[r.URL.Path]
			require.Equal(ms.t, ok, true, fmt.Sprintf("could not found GET handler for: %s", r.URL.Path))
			callback(w, r)
		case "POST":
			callback, ok := ms.post[r.URL.Path]
			require.Equal(ms.t, ok, true, fmt.Sprintf("could not found POST handler for: %s", r.URL.Path))
			callback(w, r)
		case "PUT":
			callback, ok := ms.put[r.URL.Path]
			require.Equal(ms.t, ok, true, fmt.Sprintf("could not find PUT handler for: %s", r.URL.Path))
			callback(w, r)
		default:
			require.Fail(t, fmt.Sprintf("http method not found: %s", r.Method))
		}
	}))

	return ms
}

// GetURL return URL for given endpoint
func (ms *MockServer) GetURL(endpoint string) string {
	return fmt.Sprintf("%s%s", ms.server.URL, endpoint)
}

// Close close httptest.Server
func (ms *MockServer) Close() {
	ms.server.Close()
}

// AddGetHandler add GET handler to specific path
func (ms *MockServer) AddGetHandler(path string, handler GetHandler) {
	ms.get[path] = handler
}

// AddPostHandler add POST handler to specific path
func (ms *MockServer) AddPostHandler(path string, handler PostHandler) {
	ms.post[path] = handler
}

// AddPutHandler add PUT handler to specific path
func (ms *MockServer) AddPutHandler(path string, handler PutHandler) {
	ms.put[path] = handler
}
