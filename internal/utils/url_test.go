package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsURL(t *testing.T) {
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := IsURL(tt.url)
			assert.Equal(t, tt.isValid, v)
		})
	}
}
