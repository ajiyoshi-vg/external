package external

import (
	"iter"

	"github.com/ajiyoshi-vg/external/emit"
	"github.com/ajiyoshi-vg/external/scan"
)

type Merger[T any] struct {
	cmp func(T, T) int
}

func NewMerger[T any](cmp func(T, T) int) *Merger[T] {
	return &Merger[T]{cmp: cmp}
}

func (m *Merger[T]) Merge(xs []iter.Seq[T]) iter.Seq[T] {
	return scan.Chan(m.merge(xs))
}

func (m *Merger[T]) merge(xs []iter.Seq[T]) <-chan T {
	if len(xs) == 0 {
		return emit.Chan[T](nop)
	}
	if len(xs) == 1 {
		return emit.Chan(xs[0])
	}
	a := m.merge(xs[:len(xs)/2])
	b := m.merge(xs[len(xs)/2:])
	return Merge(a, b, m.cmp)
}

func Merge[T any](a, b <-chan T, cmp func(T, T) int) <-chan T {
	ret := make(chan T, 100)

	yieldAll := func(ch <-chan T) {
		if ch == nil {
			return
		}
		for x := range ch {
			ret <- x
		}
	}

	go func() {
		defer close(ret)

		if a == nil {
			yieldAll(b)
			return
		}
		if b == nil {
			yieldAll(a)
			return
		}

		nextA, nextB, okA, okB := both(a, b)

		for okA && okB {
			if cmp(nextA, nextB) < 0 {
				ret <- nextA
				nextA, okA = <-a
			} else {
				ret <- nextB
				nextB, okB = <-b
			}
		}

		if okB {
			ret <- nextB
			yieldAll(b)
		}
		if okA {
			ret <- nextA
			yieldAll(a)
		}

	}()
	return ret
}

func both[T any](a, b <-chan T) (T, T, bool, bool) {
	select {
	case nextA, okA := <-a:
		nextB, okB := <-b
		return nextA, nextB, okA, okB
	case nextB, okB := <-b:
		nextA, okA := <-a
		return nextA, nextB, okA, okB
	}
}
