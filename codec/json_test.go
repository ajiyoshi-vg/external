package codec

import (
	"bytes"
	"io"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSON(t *testing.T) {
	cases := map[string]struct {
		input  []int
		expect []int
	}{
		"normal": {
			input:  []int{1, 2, 3},
			expect: []int{1, 2, 3},
		},
	}

	for title, c := range cases {
		t.Run(title, func(t *testing.T) {
			codec := &JSON[int]{}
			seq := slices.Values(c.input)
			var b bytes.Buffer
			err := codec.Encode(seq, &b)
			assert.NoError(t, err)

			r := io.NopCloser(&b)
			actual := slices.Collect(codec.Decode(r))
			assert.Equal(t, c.expect, actual)
		})
	}
}
