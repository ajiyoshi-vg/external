package emit

import "iter"

func All[T any](seq iter.Seq[T], yield func(T) bool) {
	next, stop := iter.Pull(seq)
	defer stop()
	After(next, yield)
}

func After[T any](next func() (T, bool), yield func(T) bool) {
	for x, ok := next(); ok; x, ok = next() {
		if !yield(x) {
			return
		}
	}
}
