package external

import (
	"bufio"
	"encoding/json"
	"io"
	"iter"
	"os"
)

type Chunk[T any] struct {
	data    []T
	tmpFile *string
}

func NewChunk[T any](data []T) *Chunk[T] {
	return &Chunk[T]{data: data}
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
	name := tempFile.Name()
	x.tmpFile = &name
	return x.store(tempFile)
}

func (x *Chunk[T]) store(w io.WriteCloser) error {
	defer w.Close()
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
