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

type ChunkFile[T any] struct {
	data    []T
	tmpFile *string
	length  int
}

func NewChunk[T any](data []T) *ChunkFile[T] {
	return &ChunkFile[T]{
		data:   data,
		length: len(data),
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
	enc := json.NewEncoder(buf)
	for _, v := range x.data {
		if err := enc.Encode(v); err != nil {
			return err
		}
	}
	return nil
}

func (x *ChunkFile[T]) Restore() (iter.Seq[T], error) {
	if x.data != nil {
		return slices.Values(x.data), nil
	}
	tempFile, err := os.Open(*x.tmpFile)
	if err != nil {
		return nil, err
	}
	return x.restore(tempFile), nil
}

func (x *ChunkFile[T]) restore(r io.ReadCloser) iter.Seq[T] {
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
