package external

import (
	"bufio"
	"errors"
	"io"
	"iter"
	"os"
	"slices"

	"github.com/ajiyoshi-vg/external/codec"
)

type ChunkFile[T any] struct {
	data    []T
	tmpFile *string
	length  int
	codec   codec.Codec[T]
}

func NewChunkFile[T any](data []T) *ChunkFile[T] {
	return &ChunkFile[T]{
		data:   data,
		length: len(data),
		codec:  codec.JSON[T]{},
	}
}

func (x *ChunkFile[T]) Length() int {
	return x.length
}

func (x *ChunkFile[T]) Clean() error {
	if x.tmpFile != nil {
		return os.Remove(*x.tmpFile)
	}
	return nil
}

func (x *ChunkFile[T]) Store() error {
	tempFile, err := os.CreateTemp("", "run")
	if err != nil {
		return err
	}
	defer tempFile.Close()
	name := tempFile.Name()
	if err := x.store(tempFile); err != nil {
		return errors.Join(err, os.Remove(name))
	}
	x.tmpFile = &name
	x.data = nil
	return nil
}

func (x *ChunkFile[T]) store(w io.WriteCloser) error {
	buf := bufio.NewWriter(w)
	defer buf.Flush()
	return x.codec.Encode(slices.Values(x.data), buf)
}

func (x *ChunkFile[T]) Restore() (iter.Seq[T], error) {
	if x.data != nil {
		return slices.Values(x.data), nil
	}
	tempFile, err := os.Open(*x.tmpFile)
	if err != nil {
		return nil, err
	}
	return x.codec.Decode(tempFile), nil
}
