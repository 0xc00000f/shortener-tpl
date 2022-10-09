package handlers

import (
	"github.com/0xc00000f/shortener-tpl/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSaveURL(t *testing.T) {
	type wantPost struct {
		contentType string
		statusCode  int
	}
	//storage := storage.NewStorage()
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
			SaveURL(w, request)
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
