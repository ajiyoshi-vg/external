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
			log.Println("got a chunk")
			xs = append(xs, x)
		}
		done <- NewChunks(xs)
	}()

	wg := sync.WaitGroup{}
	for buf := range Buffer(seq, s.opt.chunkSize) {
		wg.Add(1)
		go func(data []T) {
			defer wg.Done()
			defer log.Println("store chunk end")

			s.sort(data)
			chunk := NewChunk(data)
			if len(data) == s.opt.chunkSize {
				err := chunk.Store()
				if err != nil {
					log.Println(err)
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
				log.Println(err)
			}
		}()

		m := NewMerger(s.cmp, runs)
		merged, err := m.Merged()
		if err != nil {
			log.Println(err)
			return
		}

		yieldAll(merged, yield)
	}
}

func (s *Sorter[T]) sort(data []T) []T {
	log.Println("sort chunk start")
	defer log.Println("sort chunk end")
	sort.Slice(data, func(i, j int) bool {
		return s.cmp(data[i], data[j]) < 0
	})
	return data
}
