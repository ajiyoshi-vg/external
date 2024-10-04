package external

type option struct {
	chunkSize int
	limit     int
}

type Option func(*option)

func ChunkSize(size int) Option {
	return func(o *option) {
		o.chunkSize = size
	}
}

func Limit(n int) Option {
	return func(o *option) {
		o.limit = n
	}
}
