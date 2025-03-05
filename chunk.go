package external

import (
	"errors"
	"iter"
)

type Chunks[T any] struct {
	chunk []Chunk[T]
}

type Chunk[T any] interface {
	Length() int
	Restore() (iter.Seq[T], error)
	Clean() error
}

var _ Chunk[int] = (*ChunkFile[int])(nil)

func NewChunks[T any](cs []Chunk[T]) *Chunks[T] {
	return &Chunks[T]{chunk: cs}
}

func (x *Chunks[T]) Length() int {
	ret := 0
	for _, c := range x.chunk {
		ret += c.Length()
	}
	return ret
}

func (x *Chunks[T]) Iters() ([]iter.Seq[T], error) {
	ret := make([]iter.Seq[T], 0, len(x.chunk))
	for _, c := range x.chunk {
		iter, err := c.Restore()
		if err != nil {
			return nil, err
		}
		ret = append(ret, iter)
	}
	return ret, nil
}

func (x *Chunks[T]) Clean() error {
	var ret error
	for _, c := range x.chunk {
		if err := c.Clean(); err != nil {
			ret = errors.Join(ret, err)
		}
	}
	return ret
}
