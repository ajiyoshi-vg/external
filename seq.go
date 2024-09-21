package external

import (
	"bufio"
	"io"
	"iter"
)

func yieldAllNext[T any](next func() (T, bool), yield func(T) bool) {
	for x, ok := next(); ok; x, ok = next() {
		if !yield(x) {
			return
		}
	}
}

func yieldAll[T any](seq iter.Seq[T], yield func(T) bool) {
	next, stop := iter.Pull(seq)
	defer stop()
	yieldAllNext(next, yield)
}

func Buffer[T any](seq iter.Seq[T], size int) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		buf := make([]T, 0, size)
		for x := range seq {
			buf = append(buf, x)
			if len(buf) == size {
				if !yield(buf) {
					return
				}
				buf = make([]T, 0, size)
			}
		}
		if len(buf) > 0 {
			yield(buf)
		}
	}
}

func Lines(r io.Reader) iter.Seq[string] {
	s := bufio.NewScanner(r)
	return func(yield func(string) bool) {
		for s.Scan() {
			if !yield(s.Text()) {
				return
			}
		}
	}
}

func ByteLines(r io.Reader) iter.Seq[[]byte] {
	s := bufio.NewScanner(r)
	return func(yield func([]byte) bool) {
		for s.Scan() {
			if !yield(s.Bytes()) {
				return
			}
		}
	}
}
