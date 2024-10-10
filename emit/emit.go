package emit

import "iter"

func All[T any](seq iter.Seq[T], yield func(T) bool) {
	for x := range seq {
		if !yield(x) {
			return
		}
	}
}

func Pull[T any](next func() (T, bool), yield func(T) bool) {
	for x, ok := next(); ok; x, ok = next() {
		if !yield(x) {
			return
		}
	}
}

func Then[T any](seq iter.Seq[T], f func()) iter.Seq[T] {
	return func(yield func(T) bool) {
		defer f()
		All(seq, yield)
	}
}

func Flatten[T any](seq iter.Seq[[]T], yield func(T) bool) {
	for xs := range seq {
		for _, x := range xs {
			if !yield(x) {
				return
			}
		}
	}
}

func Chan[T any](seq iter.Seq[T]) <-chan T {
	return NewChan(seq, 0)
}

func NewChan[T any](seq iter.Seq[T], n int) <-chan T {
	ret := make(chan T, n)
	go func() {
		defer close(ret)
		for x := range seq {
			ret <- x
		}
	}()
	return ret
}
