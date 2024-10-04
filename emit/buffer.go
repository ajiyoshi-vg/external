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

func Buffered[T any](seq iter.Seq[T], yield func(T) bool, opts ...func(*bufferOption)) {
	opt := bufferOption{
		chunkSize:  1000,
		bufferSize: 5 * 1000,
	}
	for _, f := range opts {
		f(&opt)
	}
	ch := make(chan []T, opt.bufferSize)
	go func() {
		for xs := range scan.Chunk(seq, opt.chunkSize) {
			ch <- xs
		}
		close(ch)
	}()
	for xs := range ch {
		for _, x := range xs {
			if !yield(x) {
				return
			}
		}
	}
}
