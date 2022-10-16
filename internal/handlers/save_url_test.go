package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/0xc00000f/shortener-tpl/internal/api"
	"github.com/0xc00000f/shortener-tpl/internal/logic"
	"github.com/0xc00000f/shortener-tpl/internal/storage"
	"github.com/0xc00000f/shortener-tpl/internal/utils"

	"github.com/stretchr/testify/assert"
)

func TestSaveURL(t *testing.T) {
	storage := storage.NewMemoryStorage()
	sa := api.NewShortenerAPI(logic.NewURLEncoder(
		logic.SetStorage(storage),
		logic.SetLength(7),
	))
	apiInstance := NewRouter(sa)
	ts := httptest.NewServer(apiInstance)
	defer ts.Close()

	statusCode, body := testRequest(t, ts, "POST", "/", strings.NewReader("https://vk.com"), true)
	assert.Equal(t, http.StatusCreated, statusCode)
	assert.True(t, utils.IsURL(body))
}

func TestSaveURLJson(t *testing.T) {
	shortLength := 7

	storage := storage.NewMemoryStorage()
	sa := api.NewShortenerAPI(logic.NewURLEncoder(
		logic.SetStorage(storage),
		logic.SetLength(shortLength),
	))
	apiInstance := NewRouter(sa)
	ts := httptest.NewServer(apiInstance)
	defer ts.Close()

	statusCode, body := testRequest(t, ts, "POST", "/api/shorten",
		strings.NewReader(fmt.Sprintf("{\"%v\":\"%v\"}", "url", "https://vk.com")), true)
	assert.Equal(t, http.StatusCreated, statusCode)

	resp := ShortResponse{}
	err := json.Unmarshal([]byte(body), &resp)
	require.NoError(t, err)
	assert.True(t, utils.IsURL(resp.Result))
}
