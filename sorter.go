package external

import (
	"iter"
	"log"
	"slices"
	"sort"
	"strings"
)

type Sorter[T any] struct {
	cmp func(T, T) int
	opt option
}

func SortString(seq iter.Seq[string], opt ...Option) iter.Seq[string] {
	return New(strings.Compare, opt...).Sort(seq)
}

func New[T any](cmp func(T, T) int, opt ...Option) *Sorter[T] {
	ret := &Sorter[T]{
		cmp: cmp,
		opt: option{
			chunkSize: 1000 * 1000,
		},
	}

	for _, f := range opt {
		f(&ret.opt)
	}

	return ret
}

func (s *Sorter[T]) Sort(seq iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		runs := make([]*run[T], 0, 10)
		defer func() {
			for _, run := range runs {
				if err := run.Clean(); err != nil {
					log.Println(err)
				}
			}
		}()

		chunk := make([]T, 0, s.opt.chunkSize)
		for x := range seq {
			chunk = append(chunk, x)
			if len(chunk) == s.opt.chunkSize {
				log.Println("sort chunk start")
				s.sort(chunk)
				log.Println("sort chunk end")
				run := NewRun(chunk)
				if err := run.Store(); err != nil {
					return
				}
				runs = append(runs, run)
				chunk = make([]T, 0, s.opt.chunkSize)
			}
		}
		s.sort(chunk)

		if len(runs) == 0 {
			for _, x := range chunk {
				if !yield(x) {
					return
				}
			}
			return
		}

		// merge runs
		ss := make([]iter.Seq[T], 0, len(runs)+1)
		for _, run := range runs {
			seq, err := run.Restore()
			if err != nil {
				return
			}
			ss = append(ss, seq)
		}
		ss = append(ss, slices.Values(chunk))
		yieldAll(s.MergeAll(ss), yield)
	}
}

func (s *Sorter[T]) MergeAll(xs []iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		if len(xs) == 0 {
			return
		}
		if len(xs) == 1 {
			yieldAll(xs[0], yield)
			return
		}
		if len(xs) == 2 {
			yieldAll(s.Merge(xs[0], xs[1]), yield)
			return
		}
		a := s.MergeAll(xs[:len(xs)/2])
		b := s.MergeAll(xs[len(xs)/2:])
		yieldAll(s.Merge(a, b), yield)
	}
}

func (s *Sorter[T]) Merge(a, b iter.Seq[T]) iter.Seq[T] {
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

			if s.cmp(x, y) < 0 {
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

func (s *Sorter[T]) sort(data []T) []T {
	sort.Slice(data, func(i, j int) bool {
		return s.cmp(data[i], data[j]) < 0
	})
	return data
}
