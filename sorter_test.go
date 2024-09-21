package external

import (
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSorter(t *testing.T) {
	cases := map[string]struct {
		input  []string
		expect []string
		sorter *Sorter[string]
	}{
		"sort string": {
			input:  []string{"c", "a", "b"},
			expect: []string{"a", "b", "c"},
			sorter: New(strings.Compare),
		},
		"chunked": {
			input: []string{
				"z", "y", "x", "w", "v", "u", "t", "s", "r", "q",
			},
			expect: []string{
				"q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
			},
			sorter: New(strings.Compare, ChunkSize(2)),
		},
	}

	for title, c := range cases {
		t.Run(title, func(t *testing.T) {
			seq := slices.Values(c.input)
			actual := slices.Collect(c.sorter.Sort(seq))
			assert.Equal(t, c.expect, actual)
		})
	}
}
