package external

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"iter"
	"os"
	"slices"
)

type Chunks[T any] struct {
	chunk []*Chunk[T]
}

func NewChunks[T any](chunk []*Chunk[T]) *Chunks[T] {
	return &Chunks[T]{chunk: chunk}
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

type Chunk[T any] struct {
	data    []T
	tmpFile *string
	length  int
}

func NewChunk[T any](data []T) *Chunk[T] {
	return &Chunk[T]{
		data:   data,
		length: len(data),
	}
}

func (x *Chunk[T]) Length() int {
	return x.length
}

func (x *Chunk[T]) Clean() error {
	if x.tmpFile != nil {
		return os.Remove(*x.tmpFile)
	}
	return nil
}

func (x *Chunk[T]) Store() error {
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

func (x *Chunk[T]) store(w io.WriteCloser) error {
	buf := bufio.NewWriter(w)
	defer buf.Flush()
	enc := json.NewEncoder(buf)
	for _, v := range x.data {
		if err := enc.Encode(v); err != nil {
			return err
		}
	}
	return nil
}

func (x *Chunk[T]) Restore() (iter.Seq[T], error) {
	if x.data != nil {
		return slices.Values(x.data), nil
	}
	tempFile, err := os.Open(*x.tmpFile)
	if err != nil {
		return nil, err
	}
	return x.restore(tempFile), nil
}

func (x *Chunk[T]) restore(r io.ReadCloser) iter.Seq[T] {
	dec := json.NewDecoder(bufio.NewReader(r))
	return func(yield func(T) bool) {
		defer r.Close()
		for dec.More() {
			var v T
			if err := dec.Decode(&v); err != nil {
				return
			}
			if !yield(v) {
				return
			}
		}
	}
}
