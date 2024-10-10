package emit

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChan(t *testing.T) {
	cases := map[string]struct {
		input  []int
		expect []int
	}{
		"normal": {
			input:  []int{1, 2, 3},
			expect: []int{1, 2, 3},
		},
		"empty": {
			input:  nil,
			expect: []int{},
		},
	}

	for title, c := range cases {
		t.Run(title, func(t *testing.T) {
			seq := Chan(slices.Values(c.input))
			actual := make([]int, 0, len(c.input))
			for x := range seq {
				actual = append(actual, x)
			}
			assert.Equal(t, c.expect, actual)
		})
	}
}
