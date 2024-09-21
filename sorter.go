package external

import (
	"iter"
	"log"
	"sort"
	"strings"
	"sync"
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
			chunkSize: 1000 * 1000 * 5,
		},
	}

	for _, f := range opt {
		f(&ret.opt)
	}

	return ret
}

func (s *Sorter[T]) Split(seq iter.Seq[T]) []*Chunk[T] {
	done := make(chan []*Chunk[T])
	defer close(done)

	ch := make(chan *Chunk[T])

	go func() {
		ret := make([]*Chunk[T], 0, 10)
		for run := range ch {
			ret = append(ret, run)
		}
		done <- ret
	}()

	wg := sync.WaitGroup{}
	for buf := range Buffer(seq, s.opt.chunkSize) {
		wg.Add(1)
		go func(data []T) {
			defer wg.Done()
			defer log.Println("store chunk end")

			s.sort(data)
			run := NewChunk(data)
			if err := run.Store(); err != nil {
				log.Println(err)
				return
			}
			ch <- run
		}(buf)
	}
	wg.Wait()
	close(ch)

	return <-done
}

func (s *Sorter[T]) Sort(seq iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		runs := s.Split(seq)

		defer func() {
			for _, run := range runs {
				if err := run.Clean(); err != nil {
					log.Println(err)
				}
			}
		}()

		// merge runs
		ss := make([]iter.Seq[T], 0, len(runs))
		for _, run := range runs {
			seq, err := run.Restore()
			if err != nil {
				return
			}
			ss = append(ss, seq)
		}
		i := 0
		for x := range s.Merge(ss) {
			if !yield(x) {
				return
			}
			i++
			if i%1000000 == 0 {
				log.Println(i)
			}
		}
	}
}

func (s *Sorter[T]) Merge(xs []iter.Seq[T]) iter.Seq[T] {
	if len(xs) == 0 {
		return nil
	}
	if len(xs) == 1 {
		return xs[0]
	}
	a := s.Merge(xs[:len(xs)/2])
	b := s.Merge(xs[len(xs)/2:])
	return Merge(a, b, s.cmp)
}

func (s *Sorter[T]) sort(data []T) []T {
	log.Println("sort chunk start")
	defer log.Println("sort chunk end")
	sort.Slice(data, func(i, j int) bool {
		return s.cmp(data[i], data[j]) < 0
	})
	return data
}
