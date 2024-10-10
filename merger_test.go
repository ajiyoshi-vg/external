package external

import (
	"iter"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMerger(t *testing.T) {
	cases := map[string]struct {
		input  []iter.Seq[int]
		expect []int
	}{
		"normal": {
			input: []iter.Seq[int]{
				slices.Values([]int{1, 3, 5}),
				slices.Values([]int{2, 4, 6}),
			},
			expect: []int{1, 2, 3, 4, 5, 6},
		},
		"empty": {
			input:  []iter.Seq[int]{},
			expect: nil,
		},
		"nil": {
			input:  nil,
			expect: nil,
		},
	}

	for title, c := range cases {
		t.Run(title, func(t *testing.T) {
			seq := NewMerger(Compare[int]).Merge(c.input)
			actual := slices.Collect(seq)
			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestMerge(t *testing.T) {
	cases := map[string]struct {
		a, b   <-chan int
		expect []int
	}{
		"normal": {
			a:      produce(1, 3, 5),
			b:      produce(2, 4, 6),
			expect: []int{1, 2, 3, 4, 5, 6},
		},
		"empty": {
			a:      nil,
			b:      produce(1, 2, 3),
			expect: []int{1, 2, 3},
		},
		"empty2": {
			a:      produce(1, 2, 3),
			b:      nil,
			expect: []int{1, 2, 3},
		},
		"empty3": {
			a:      nil,
			b:      nil,
			expect: []int{},
		},
	}

	for title, c := range cases {
		t.Run(title, func(t *testing.T) {
			seq := Merge(c.a, c.b, Compare[int])
			actual := make([]int, 0, len(c.expect))
			for x := range seq {
				actual = append(actual, x)
			}
			assert.Equal(t, c.expect, actual)
		})
	}
}

func produce(xs ...int) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for _, x := range xs {
			ch <- x
		}
	}()
	return ch
}
