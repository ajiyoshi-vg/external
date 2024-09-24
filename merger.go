package external

import (
	"iter"

	"github.com/ajiyoshi-vg/external/emit"
)

type Merger[T any] struct {
	cmp func(T, T) int
}

func NewMerger[T any](cmp func(T, T) int) *Merger[T] {
	return &Merger[T]{cmp: cmp}
}

func (m *Merger[T]) Merge(xs []iter.Seq[T]) iter.Seq[T] {
	if len(xs) == 0 {
		return func(yield func(T) bool) {}
	}
	if len(xs) == 1 {
		return xs[0]
	}
	a := m.Merge(xs[:len(xs)/2])
	b := m.Merge(xs[len(xs)/2:])
	return Merge(a, b, m.cmp)
}

func Merge[T any](a, b iter.Seq[T], cmp func(T, T) int) iter.Seq[T] {
	return func(yield func(T) bool) {
		nextA, stopA := iter.Pull(a)
		defer stopA()
		nextB, stopB := iter.Pull(b)
		defer stopB()

		a, okA := nextA()
		b, okB := nextB()
		for okA || okB {
			if !okA {
				if !yield(b) {
					return
				}
				emit.Pull(nextB, yield)
				return
			}
			if !okB {
				if !yield(a) {
					return
				}
				emit.Pull(nextA, yield)
				return
			}

			if cmp(a, b) < 0 {
				if !yield(a) {
					return
				}
				a, okA = nextA()
			} else {
				if !yield(b) {
					return
				}
				b, okB = nextB()
			}
		}
	}
}
