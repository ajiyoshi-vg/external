package external

import (
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortInt(t *testing.T) {
	cases := map[string]struct {
		input  []int
		expect []int
	}{
		"sort int": {
			input:  []int{3, 1, 2},
			expect: []int{1, 2, 3},
		},
	}

	for title, c := range cases {
		t.Run(title, func(t *testing.T) {
			actual := Sort(slices.Values(c.input))
			assert.Equal(t, c.expect, slices.Collect(actual))
		})
	}
}

type item struct {
	Name   string
	Price  int
	secret string
}

func TestSortFunc(t *testing.T) {
	cases := map[string]struct {
		input    []item
		expect   []item
		copmpare func(item, item) int
	}{
		"sort item by name": {
			input: []item{
				{Name: "c", Price: 1},
				{Name: "a", Price: 3},
				{Name: "b", Price: 2},
			},
			expect: []item{
				{Name: "a", Price: 3},
				{Name: "b", Price: 2},
				{Name: "c", Price: 1},
			},
			copmpare: func(a, b item) int {
				return strings.Compare(a.Name, b.Name)
			},
		},
		"sort item by price": {
			input: []item{
				{Name: "c", Price: 1},
				{Name: "a", Price: 3},
				{Name: "b", Price: 2},
			},
			expect: []item{
				{Name: "c", Price: 1},
				{Name: "b", Price: 2},
				{Name: "a", Price: 3},
			},
			copmpare: func(a, b item) int {
				return a.Price - b.Price
			},
		},
		"forget fields that cannot be restored by json.Unmarshal": {
			input: []item{
				{Name: "c", Price: 1, secret: "secret"},
				{Name: "a", Price: 3, secret: "secret"},
				{Name: "b", Price: 2, secret: "secret"},
			},
			expect: []item{
				{Name: "a", Price: 3},
				{Name: "b", Price: 2},
				{Name: "c", Price: 1},
			},
			copmpare: func(a, b item) int {
				return strings.Compare(a.Name, b.Name)
			},
		},
	}

	for title, c := range cases {
		t.Run(title, func(t *testing.T) {
			actual := SortFunc(slices.Values(c.input), c.copmpare, ChunkSize(1))
			assert.Equal(t, c.expect, slices.Collect(actual))
		})
	}
}
