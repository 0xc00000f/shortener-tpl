package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandStringRunes(t *testing.T) {
	tests := []struct {
		name    string
		letters int
	}{
		{
			name:    "6 letters random string",
			letters: 6,
		},
		{
			name:    "72 letters random string",
			letters: 72,
		},
		{
			name:    "0 letters random string",
			letters: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			randString := RandStringRunes(tt.letters)
			assert.Equal(t, len(randString), tt.letters)
		})
	}
}
