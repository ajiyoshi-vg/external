package external

import (
	"errors"
	"iter"
	"sort"
	"sync"

	"github.com/ajiyoshi-vg/external/emit"
	"github.com/ajiyoshi-vg/external/scan"
)

type Sorter[T any] struct {
	cmp func(T, T) int
	opt option
	err error
}

func New[T any](cmp func(T, T) int, opt ...Option) *Sorter[T] {
	ret := &Sorter[T]{
		cmp: cmp,
		opt: option{
			chunkSize: 1000 * 1000 * 3,
		},
	}

	for _, f := range opt {
		f(&ret.opt)
	}

	return ret
}

func (s *Sorter[T]) Split(seq iter.Seq[T]) *Chunks[T] {
	done := make(chan *Chunks[T])
	defer close(done)

	ch := make(chan *Chunk[T])

	go func() {
		xs := make([]*Chunk[T], 0, 10)
		for x := range ch {
			xs = append(xs, x)
		}
		done <- NewChunks(xs)
	}()

	wg := sync.WaitGroup{}
	for buf := range scan.Chunk(seq, s.opt.chunkSize) {
		wg.Add(1)
		go func(data []T) {
			defer wg.Done()

			s.sort(data)
			chunk := NewChunk(data)
			if len(data) == s.opt.chunkSize {
				err := chunk.Store()
				if err != nil {
					s.err = errors.Join(s.err, err)
					return
				}
			}
			ch <- chunk
		}(buf)
	}
	wg.Wait()
	close(ch)

	return <-done
}

func (s *Sorter[T]) Sort(seq iter.Seq[T]) iter.Seq[T] {
	runs := s.Split(seq)
	return func(yield func(T) bool) {
		defer func() {
			if err := runs.Clean(); err != nil {
				s.err = errors.Join(s.err, err)
			}
		}()

		m := NewMerger(s.cmp, runs)
		merged, err := m.Merged()
		if err != nil {
			s.err = errors.Join(s.err, err)
			return
		}

		emit.All(merged, yield)
	}
}

func (s *Sorter[T]) Err() error {
	return s.err
}

func (s *Sorter[T]) sort(data []T) []T {
	sort.Slice(data, func(i, j int) bool {
		return s.cmp(data[i], data[j]) < 0
	})
	return data
}
