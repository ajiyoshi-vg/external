package external

import (
	"errors"
	"iter"

	"github.com/ajiyoshi-vg/external/emit"
)

type Sorter[T any] struct {
	cmp func(T, T) int
	opt []Option
	err error
}

func New[T any](cmp func(T, T) int, opt ...Option) *Sorter[T] {
	ret := &Sorter[T]{
		cmp: cmp,
		opt: opt,
	}

	return ret
}

func (x *Sorter[T]) Sort(seq iter.Seq[T]) iter.Seq[T] {
	x.err = nil
	s := NewSplitter(x.cmp, x.opt...)
	m := NewMerger(x.cmp)

	chunks, err := s.Split(seq)
	if err != nil {
		x.catch(err)
		return nop
	}

	cs, err := chunks.Iters()
	if err != nil {
		x.catch(err)
		return nop
	}

	return func(yield func(T) bool) {
		defer func() {
			x.catch(chunks.Clean())
		}()

		emit.All(m.Merge(cs), yield)
	}
}

func (s *Sorter[T]) Err() error {
	return s.err
}

func (s *Sorter[T]) catch(err error) {
	s.err = errors.Join(s.err, err)
}
