package external

import (
	"iter"
	"log"
)

type Merger[T any] struct {
	cmp  func(T, T) int
	runs []*Chunk[T]
}

func NewMerger[T any](cmp func(T, T) int, runs []*Chunk[T]) *Merger[T] {
	return &Merger[T]{cmp: cmp, runs: runs}
}

func Merge[T any](a, b iter.Seq[T], cmp func(T, T) int) iter.Seq[T] {
	return func(yield func(T) bool) {
		log.Println("merge start")
		defer log.Println("merge end")

		nextA, stopA := iter.Pull(a)
		defer stopA()
		nextB, stopB := iter.Pull(b)
		defer stopB()

		x, okA := nextA()
		y, okB := nextB()
		for okA || okB {
			if !okA {
				if !yield(y) {
					return
				}
				yieldAllNext(nextB, yield)
				return
			}
			if !okB {
				if !yield(x) {
					return
				}
				yieldAllNext(nextA, yield)
				return
			}

			if cmp(x, y) < 0 {
				if !yield(x) {
					return
				}
				x, okA = nextA()
			} else {
				if !yield(y) {
					return
				}
				y, okB = nextB()
			}
		}
	}
}
