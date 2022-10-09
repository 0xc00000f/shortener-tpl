package handlers

import (
	"github.com/0xc00000f/shortener-tpl/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRedirect(t *testing.T) {
	type wantGet struct {
		contentType string
		statusCode  int
		location    string
	}
	//storage := storage.NewStorage()
	tests := []struct {
		name        string
		requestPost string
		postBody    string
		want        wantGet
	}{
		{
			name:        "[positive] query to /",
			requestPost: "http://localhost:8080/",
			postBody:    "https://vk.com",
			want: wantGet{
				contentType: "text/plain; charset=utf-8",
				statusCode:  307,
				location:    "https://vk.com",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// prepare data
			requestPost := httptest.NewRequest(http.MethodPost, tt.requestPost, strings.NewReader(tt.postBody))
			wPost := httptest.NewRecorder()
			SaveURL(wPost, requestPost)

			resultPost := wPost.Result()
			defer resultPost.Body.Close()

			assert.Equal(t, resultPost.StatusCode, 201)

			b, err := io.ReadAll(resultPost.Body)
			require.NoError(t, err)
			err = resultPost.Body.Close()
			require.NoError(t, err)

			assert.True(t, utils.IsURL(string(b)))

			uri := string(b)

			// test
			log.Print("URI:", uri)
			requestGet := httptest.NewRequest(http.MethodGet, uri, nil)
			wGet := httptest.NewRecorder()
			Redirect(wGet, requestGet)

			resultGet := wGet.Result()
			defer resultGet.Body.Close()

			assert.Equal(t, tt.want.statusCode, resultGet.StatusCode)
			assert.Equal(t, tt.want.contentType, resultGet.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.location, resultGet.Header.Get("Location"))
		})
	}
}
