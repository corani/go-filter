package filter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilter_Apply(t *testing.T) {
	t.Parallel()

	type export struct {
		Type  string
		Value int
	}

	tt := []struct {
		name       string
		expression string
		data       []export
		err        error
		exp        []interface{}
	}{
		{
			name:       "empty",
			expression: "",
			data: []export{
				{"export", 33},
				{"delta", 42},
			},
			exp: []interface{}{
				export{"export", 33},
				export{"delta", 42},
			},
		},
		{
			name:       "simple eq",
			expression: "Type eq 'export'",
			data: []export{
				{"export", 33},
				{"delta", 42},
			},
			exp: []interface{}{
				export{"export", 33},
			},
		},
		{
			name:       "simple lt",
			expression: "Value lt 42",
			data: []export{
				{"export", 33},
				{"delta", 42},
			},
			exp: []interface{}{
				export{"export", 33},
			},
		},
	}

	for _, tc := range tt {
		tc := tc // pin

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			f := MustParse(tc.expression)

			res, err := f.Apply(tc.data)
			if err != nil {
				assert.ErrorIs(t, err, tc.err)

				return
			}

			assert.NoError(t, tc.err, "expected error, got none")
			assert.Equal(t, tc.exp, res)
		})
	}
}
