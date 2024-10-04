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
	cmp func(T, T) int
	opt option
}

func NewSplitter[T any](cmp func(T, T) int, opt ...Option) *Splitter[T] {
	ret := &Splitter[T]{
		cmp: cmp,
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

	g := errgroup.Group{}
	g.SetLimit(s.opt.limit)
	for data := range scan.Chunk(seq, s.opt.chunkSize) {
		g.Go(func() error {
			chunk, err := s.chunk(data)
			if err != nil {
				return err
			}
			ch <- chunk
			return nil
		})
	}
	err := g.Wait()
	close(ch)

	ret := <-done
	if err != nil {
		return nil, errors.Join(err, ret.Clean())
	}
	return ret, nil
}

func (s *Splitter[T]) chunk(data []T) (*Chunk[T], error) {
	s.sort(data)
	ret := NewChunk(data)
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
