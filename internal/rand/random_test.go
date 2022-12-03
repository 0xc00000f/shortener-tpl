package rand

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	rand := New(true)
	tests := []struct {
		name            string
		letters         int
		predictableText string
	}{
		{
			name:            "6 letters random string",
			letters:         6,
			predictableText: "BpLnfg",
		},
		{
			name:            "72 letters random string",
			letters:         72,
			predictableText: "Dsc2WD8F2qNfHK5a84jjJkwzDkh9h2fhfUVuS9jZ8uVbhV3vC5AWX39IVUWSP2NcHciWvqZT",
		},
		{
			name:            "0 letters random string",
			letters:         0,
			predictableText: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			randString := rand.String(tt.letters)
			assert.Equal(t, tt.letters, len(randString))
			assert.Equal(t, tt.predictableText, randString)
		})
	}
}
