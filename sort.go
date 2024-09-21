package external

import (
	"iter"

	"golang.org/x/exp/constraints"
)

func compare[T constraints.Ordered](a, b T) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

func Sort[T constraints.Ordered](seq iter.Seq[T], opt ...Option) iter.Seq[T] {
	return New(compare[T], opt...).Sort(seq)
}
