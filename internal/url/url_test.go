package url_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/0xc00000f/shortener-tpl/internal/url"
)

func TestValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		url     string
		isValid bool
	}{
		{
			name:    "not valid #1",
			url:     "http:::/not.valid/a//a??a?b=&&c#hi",
			isValid: false,
		},
		{
			name:    "not valid #2",
			url:     "http//google.com",
			isValid: false,
		},
		{
			name:    "not valid #3",
			url:     "google.com",
			isValid: false,
		},
		{
			name:    "not valid #4",
			url:     "/foo/bar",
			isValid: false,
		},
		{
			name:    "not valid #5",
			url:     "http://",
			isValid: false,
		},
		{
			name:    "valid #1",
			url:     "http://google.com",
			isValid: true,
		},
		{
			name:    "valid #2",
			url:     "https://ya.ru",
			isValid: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			v := url.Valid(tt.url)
			assert.Equal(t, tt.isValid, v)
		})
	}
}
