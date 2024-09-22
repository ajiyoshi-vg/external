package external

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplit(t *testing.T) {
	cases := map[string]struct {
		input  []int
		expect int
		split  *Splitter[int]
	}{
		"split int": {
			input:  []int{3, 1, 2},
			expect: 3,
			split:  NewSplitter(compare[int]),
		},
		"split int(chunked)": {
			input:  []int{3, 1, 2},
			expect: 3,
			split:  NewSplitter(compare[int], ChunkSize(2)),
		},
	}

	for title, c := range cases {
		t.Run(title, func(t *testing.T) {
			actual, err := c.split.Split(slices.Values(c.input))
			assert.NoError(t, err)
			assert.Equal(t, c.expect, actual.Length())
		})
	}
}
