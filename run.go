package external

import (
	"bufio"
	"encoding/json"
	"io"
	"iter"
	"os"
)

type run[T any] struct {
	data    []T
	tmpFile *string
}

func NewRun[T any](data []T) *run[T] {
	return &run[T]{data: data}
}

func (x *run[T]) Clean() error {
	if x.tmpFile != nil {
		return os.Remove(*x.tmpFile)
	}
	return nil
}

func (x *run[T]) Store() error {
	tempFile, err := os.CreateTemp("", "run")
	if err != nil {
		return err
	}
	defer tempFile.Close()
	name := tempFile.Name()
	x.tmpFile = &name
	return x.store(tempFile)
}

func (x *run[T]) store(w io.Writer) error {
	bf := bufio.NewWriter(w)
	defer bf.Flush()
	for _, v := range x.data {
		buf, err := json.Marshal(v)
		if err != nil {
			return err
		}
		if _, err := bf.Write(buf); err != nil {
			return err
		}
		if _, err := bf.WriteRune('\n'); err != nil {
			return err
		}
	}
	return nil
}

func (x *run[T]) Restore() (iter.Seq[T], error) {
	tempFile, err := os.Open(*x.tmpFile)
	if err != nil {
		return nil, err
	}
	return x.restore(tempFile), nil
}

func (x *run[T]) restore(r io.Reader) iter.Seq[T] {
	return func(yield func(T) bool) {
		for line := range Lines(r) {
			var v T
			if err := json.Unmarshal([]byte(line), &v); err != nil {
				return
			}
			if !yield(v) {
				return
			}
		}
	}
}
