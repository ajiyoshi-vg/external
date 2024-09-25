package external

import (
	"iter"

	"golang.org/x/exp/constraints"
)

func Compare[T constraints.Ordered](a, b T) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

func Sort[T constraints.Ordered](seq iter.Seq[T], opt ...Option) iter.Seq[T] {
	return New(Compare[T], opt...).Sort(seq)
}

func SortFunc[T any](seq iter.Seq[T], cmp func(T, T) int, opt ...Option) iter.Seq[T] {
	return New(cmp, opt...).Sort(seq)
}

func nop[T any](func(T) bool) {}
