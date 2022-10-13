package app

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader, isFullPath bool) (int, string) {
	var url string
	if isFullPath {
		url = ts.URL + path
	} else {
		url = path
	}

	req, err := http.NewRequest(method, url, body)
	require.NoError(t, err)

	transport := http.Transport{}
	resp, err := transport.RoundTrip(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp.StatusCode, string(respBody)
}

func TestRedirect(t *testing.T) {
	r := NewRouter()
	ts := httptest.NewServer(r)
	defer ts.Close()

	statusCode, body := testRequest(t, ts, "POST", "/", strings.NewReader("https://vk.com"), true)
	assert.Equal(t, http.StatusCreated, statusCode)

	statusCode, _ = testRequest(t, ts, "GET", body, nil, false)
	assert.Equal(t, http.StatusTemporaryRedirect, statusCode)
}
