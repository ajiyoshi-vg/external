package scan

import (
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChunk(t *testing.T) {
	cases := map[string]struct {
		input  []int
		expect [][]int
		size   int
	}{
		"chunk by 3": {
			input:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			expect: [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
			size:   3,
		},
		"chunk by 2": {
			input:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			expect: [][]int{{1, 2}, {3, 4}, {5, 6}, {7, 8}, {9}},
			size:   2,
		},
		"empty": {
			input:  []int{},
			expect: nil,
			size:   3,
		},
	}

	for title, c := range cases {
		t.Run(title, func(t *testing.T) {
			seq := slices.Values(c.input)
			actual := slices.Collect(Chunk(seq, c.size))
			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestLines(t *testing.T) {
	cases := map[string]struct {
		input  string
		expect []string
	}{
		"empty": {
			input:  "",
			expect: nil,
		},
		"one line": {
			input:  "hello\n",
			expect: []string{"hello"},
		},
		"two lines": {
			input:  "hello\nworld\n",
			expect: []string{"hello", "world"},
		},
	}

	for title, c := range cases {
		t.Run(title, func(t *testing.T) {
			seq := Lines(strings.NewReader(c.input))
			actual := slices.Collect(seq)
			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestByteLines(t *testing.T) {
	cases := map[string]struct {
		input  string
		expect [][]byte
	}{
		"empty": {
			input:  "",
			expect: nil,
		},
		"one line": {
			input:  "hello\n",
			expect: [][]byte{[]byte("hello")},
		},
		"two lines": {
			input:  "hello\nworld\n",
			expect: [][]byte{[]byte("hello"), []byte("world")},
		},
	}

	for title, c := range cases {
		t.Run(title, func(t *testing.T) {
			seq := ByteLines(strings.NewReader(c.input))
			actual := slices.Collect(seq)
			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestUniq(t *testing.T) {
	cases := map[string]struct {
		input  []int
		expect []int
	}{
		"empty": {
			input:  nil,
			expect: nil,
		},
		"one element": {
			input:  []int{1},
			expect: []int{1},
		},
		"two elements": {
			input:  []int{1, 1},
			expect: []int{1},
		},
		"three elements": {
			input:  []int{1, 2, 2},
			expect: []int{1, 2},
		},
	}

	for title, c := range cases {
		t.Run(title, func(t *testing.T) {
			seq := slices.Values(c.input)
			actual := slices.Collect(Uniq(seq))
			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestChan(t *testing.T) {
	cases := map[string]struct {
		input  <-chan int
		expect []int
	}{
		"normal": {
			input:  produce(1, 2, 3),
			expect: []int{1, 2, 3},
		},
	}

	for title, c := range cases {
		t.Run(title, func(t *testing.T) {
			actual := slices.Collect(Chan(c.input))
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
