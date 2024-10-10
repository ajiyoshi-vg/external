package emit

import (
	"iter"

	"github.com/ajiyoshi-vg/external/scan"
)

type bufferOption struct {
	chunkSize  int
	bufferSize int
}

func BufferSize(n int) func(*bufferOption) {
	return func(opt *bufferOption) {
		opt.bufferSize = n
	}
}
func ChunkSize(n int) func(*bufferOption) {
	return func(opt *bufferOption) {
		opt.chunkSize = n
	}
}

func Buffered[T any](
	seq iter.Seq[T],
	yield func(T) bool,
	opts ...func(*bufferOption),
) {
	opt := &bufferOption{
		chunkSize:  1000,
		bufferSize: 100,
	}
	for _, f := range opts {
		f(opt)
	}
	chunked := scan.Chunk(seq, opt.chunkSize)
	ch := NewChan(chunked, opt.bufferSize)
	Flatten(scan.Chan(ch), yield)
}
