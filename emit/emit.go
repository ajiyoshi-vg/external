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

func Chan[T any](ch <-chan T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for x := range ch {
			if !yield(x) {
				return
			}
		}
	}
}
