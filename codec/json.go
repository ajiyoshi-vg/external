package codec

import (
	"encoding/json"
	"io"
	"iter"
	"log"
)

type JSON[T any] struct{}

var _ Codec[int] = (*JSON[int])(nil)

func (JSON[T]) Encode(seq iter.Seq[T], w io.Writer) error {
	enc := json.NewEncoder(w)
	for x := range seq {
		if err := enc.Encode(x); err != nil {
			return err
		}
	}
	return nil
}

func (JSON[T]) Decode(r io.ReadCloser) iter.Seq[T] {
	dec := json.NewDecoder(r)
	return func(yield func(T) bool) {
		defer r.Close()
		for dec.More() {
			var x T
			if err := dec.Decode(&x); err != nil {
				log.Printf("decode failed: %v", err)
				return
			}
			if !yield(x) {
				return
			}
		}
	}
}
