package external

import (
	"errors"
	"iter"
	"runtime"
	"sort"

	"github.com/ajiyoshi-vg/external/scan"
	"golang.org/x/sync/errgroup"
)

type Splitter[T any] struct {
	cmp      func(T, T) int
	newChunk func([]T) ChunkStore[T]
	opt      option
}

type ChunkStore[T any] interface {
	Store() error
	Chunk[T]
}

func NewSplitter[T any](cmp func(T, T) int, opt ...Option) *Splitter[T] {
	newChunk := func(data []T) ChunkStore[T] { return NewChunkFile(data) }
	return NewSplitterWithChunkStore(cmp, newChunk, opt...)
}

func NewSplitterWithChunkStore[T any](
	cmp func(T, T) int,
	newChunk func([]T) ChunkStore[T],
	opt ...Option,
) *Splitter[T] {
	ret := &Splitter[T]{
		cmp:      cmp,
		newChunk: newChunk,
		opt: option{
			chunkSize: 1000 * 1000 * 3,
			limit:     runtime.NumCPU(),
		},
	}

	for _, f := range opt {
		f(&ret.opt)
	}

	return ret
}

func (s *Splitter[T]) Split(seq iter.Seq[T]) (*Chunks[T], error) {
	ch := make(chan Chunk[T])
	result := collectSlice(ch)

	g := errgroup.Group{}
	g.SetLimit(s.opt.limit)
	for data := range scan.Chunk(seq, s.opt.chunkSize) {
		g.Go(func() error {
			chunk, err := s.sortedChunk(data)
			if err != nil {
				return err
			}
			ch <- chunk
			return nil
		})
	}
	err := g.Wait()
	close(ch)

	ret := NewChunks(<-result)
	if err != nil {
		return nil, errors.Join(err, ret.Clean())
	}
	return ret, nil
}

func collectSlice[T any](ch <-chan T) <-chan []T {
	result := make(chan []T)
	go func() {
		defer close(result)
		xs := make([]T, 0, 10)
		for x := range ch {
			xs = append(xs, x)
		}
		result <- xs
	}()
	return result
}

func (s *Splitter[T]) sortedChunk(data []T) (Chunk[T], error) {
	s.sort(data)
	ret := s.newChunk(data)
	if len(data) == s.opt.chunkSize {
		if err := ret.Store(); err != nil {
			return nil, err
		}
	}
	return ret, nil

}

func (s *Splitter[T]) sort(data []T) {
	sort.Slice(data, func(i, j int) bool {
		return s.cmp(data[i], data[j]) < 0
	})
}
