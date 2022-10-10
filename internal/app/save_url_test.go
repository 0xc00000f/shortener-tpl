package app

import (
	"github.com/0xc00000f/shortener-tpl/internal/utils"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSaveURL(t *testing.T) {
	r := NewRouter()
	ts := httptest.NewServer(r)
	defer ts.Close()

	statusCode, body := testRequest(t, ts, "POST", "/", strings.NewReader("https://vk.com"), true)
	assert.Equal(t, http.StatusCreated, statusCode)
	assert.True(t, utils.IsURL(body))
}
