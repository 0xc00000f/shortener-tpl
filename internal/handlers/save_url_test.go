package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/0xc00000f/shortener-tpl/internal/api"
	"github.com/0xc00000f/shortener-tpl/internal/logic"
	"github.com/0xc00000f/shortener-tpl/internal/storage"
	"github.com/0xc00000f/shortener-tpl/internal/url"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveURL(t *testing.T) {
	storage := storage.NewMemoryStorage()
	logic := logic.NewURLEncoder(
		logic.SetStorage(storage),
		logic.SetLength(7),
	)
	sa := api.NewShortenerAPI(api.SetLogic(logic))
	apiInstance := NewRouter(sa)
	ts := httptest.NewServer(apiInstance)
	defer ts.Close()

	statusCode, body := testRequest(t, ts, "POST", "/", strings.NewReader("https://vk.com"), true)
	assert.Equal(t, http.StatusCreated, statusCode)
	assert.True(t, url.Valid(body))
}

func TestSaveURLJson(t *testing.T) {
	shortLength := 7

	storage := storage.NewMemoryStorage()
	logic := logic.NewURLEncoder(
		logic.SetStorage(storage),
		logic.SetLength(shortLength),
	)
	sa := api.NewShortenerAPI(api.SetLogic(logic))
	apiInstance := NewRouter(sa)
	ts := httptest.NewServer(apiInstance)
	defer ts.Close()

	statusCode, body := testRequest(t, ts, "POST", "/api/shorten",
		strings.NewReader(fmt.Sprintf("{\"%v\":\"%v\"}", "url", "https://vk.com")), true)
	assert.Equal(t, http.StatusCreated, statusCode)

	resp := ShortResponse{}
	err := json.Unmarshal([]byte(body), &resp)
	require.NoError(t, err)
	assert.True(t, url.Valid(resp.Result))
}
