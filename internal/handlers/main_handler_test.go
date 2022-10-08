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

func TestMainHandlerPost(t *testing.T) {
	type wantPost struct {
		contentType string
		statusCode  int
	}
	tests := []struct {
		name     string
		request  string
		postBody string
		want     wantPost
	}{
		{
			name:     "[positive] query to /",
			request:  "/",
			postBody: "https://vk.com",
			want: wantPost{
				contentType: "text/plain; charset=utf-8",
				statusCode:  201,
			},
		},
		{
			name:     "[negative] query to /{anything}",
			request:  "/mjkjn",
			postBody: "https://vk.com",
			want: wantPost{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
			},
		},
		{
			name:     "[negative] incorrect body",
			request:  "/",
			postBody: "ht:/vk.om",
			want: wantPost{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.request, strings.NewReader(tt.postBody))
			w := httptest.NewRecorder()
			h := MainHandler()

			h.ServeHTTP(w, request)
			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			if result.StatusCode == 201 {
				b, err := io.ReadAll(result.Body)
				require.NoError(t, err)

				err = result.Body.Close()
				require.NoError(t, err)

				assert.True(t, utils.IsURL(string(b)))
			}
		})
	}
}

func TestMainHandlerGet(t *testing.T) {
	type wantGet struct {
		contentType string
		statusCode  int
		location    string
	}
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
			hPost := MainHandler()

			hPost.ServeHTTP(wPost, requestPost)
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
			hPost.ServeHTTP(wGet, requestGet)
			resultGet := wGet.Result()
			defer resultGet.Body.Close()

			assert.Equal(t, tt.want.statusCode, resultGet.StatusCode)
			assert.Equal(t, tt.want.contentType, resultGet.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.location, resultGet.Header.Get("Location"))
		})
	}
}
