package filter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilter_Tokenize(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name       string
		expression string
		tokens     []token
		err        error
	}{
		{
			name:       "empty",
			expression: "",
			tokens:     nil,
			err:        nil,
		},
		{
			name:       "field eq 'string'",
			expression: "field eq 'string'",
			tokens: []token{
				{FIELD, "field"},
				{EQ, "eq"},
				{STRING, "string"},
			},
			err: nil,
		},
	}

	for _, tc := range tt {
		tc := tc // pin

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tokens, err := tokenize(tc.expression)
			if err != nil {
				assert.ErrorIs(t, err, tc.err)

				return
			}

			assert.NoError(t, tc.err, "expected error, got none")
			assert.Equal(t, tc.tokens, tokens)
		})
	}
}
