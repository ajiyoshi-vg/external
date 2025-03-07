package codec

import (
	"io"
	"iter"
)

type Codec[T any] interface {
	Encode(seq iter.Seq[T], w io.Writer) error
	Decode(r io.ReadCloser) iter.Seq[T]
}
