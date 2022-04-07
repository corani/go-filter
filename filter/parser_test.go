package filter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilter_ParseExpression(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name       string
		expression string
		nodes      *node
		err        error
	}{
		{
			name:       "empty",
			expression: "",
			nodes: &node{
				oper:   BOOL,
				asBool: true,
			},
			err: nil,
		},
		{
			name:       "field eq 'string'",
			expression: "field eq 'string'",
			nodes: &node{
				oper: EQ,
				children: []*node{
					{oper: FIELD, asField: "field"},
					{oper: STRING, asString: "string"},
				},
			},
			err: nil,
		},
	}

	for _, tc := range tt {
		tc := tc // pin

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			nodes, err := parseExpression(tc.expression)
			if err != nil {
				if !assert.ErrorIs(t, err, tc.err) {
					if v, err := tokenize(tc.expression); err == nil {
						t.Logf("tokens: %v", v)
					}
				}

				return
			}

			assert.NoError(t, tc.err, "expected error, got none")
			assert.Equal(t, tc.nodes, nodes)
		})
	}
}
