package scan

import (
	"bufio"
	"io"
	"iter"
	"log/slog"
	"time"
)

func Chunk[T any](seq iter.Seq[T], size int) iter.Seq[[]T] {
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
			buf := make([]byte, len(s.Bytes()))
			copy(buf, s.Bytes())
			if !yield(buf) {
				return
			}
		}
	}
}

func Uniq[T comparable](sorted iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		var last T
		for x := range sorted {
			if x != last {
				if !yield(x) {
					return
				}
				last = x
			}
		}
	}
}

func Prove[T any](name string, seq iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		last := time.Now()
		var total time.Duration
		i := 0

		defer func() {
			slog.Info("prove",
				"name", name,
				"total", total,
				"average", ave(total, i),
				"count", i)
		}()

		for x := range seq {
			total += time.Since(last)
			last = time.Now()
			i++
			if !yield(x) {
				return
			}
		}
	}
}

func ave(total time.Duration, count int) time.Duration {
	if count == 0 {
		return 0
	}
	return total / time.Duration(count)
}
