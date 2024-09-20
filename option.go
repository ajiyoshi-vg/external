package external

type option struct {
	chunkSize int
}

type Option func(*option)

func ChunkSize(size int) Option {
	return func(o *option) {
		o.chunkSize = size
	}
}
