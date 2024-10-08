# `external` go library

External is a library that takes an almost arbitrary iter.Seq[T] using Go1.23's range over func API and returns a sorted iter.Seq[T]. The difference from sort.Slice is that it can be used for data sets that are too large to fit into memory. external.Sort splits the massive input into small chunks, sorts each one using sort.Slice, writes them to temporary files, and merges them to return the final sorted iter.Seq[T]. By using json.Encoder for writing to temporary files and json.Decoder for reading back, it is possible to handle arbitrary iter.Seq[T]. In other words, types that can be encoded/decoded with the standard json package can be sorted with external.Sort/SortFunc. Simple integers or strings meet this condition, and any struct with public fields containing integers, strings, or slices thereof will also meet this condition without modifications. Structures with floating-point numbers and private fields may lose information during writing/reading. To avoid this, implement json.Marshaler/json.Unmarshaler.

## Usage Example

```go
import (
	"github.com/ajiyoshi-vg/external"
	"github.com/ajiyoshi-vg/external/scan"
)

func sort(r io.Reader) error {
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
    // Sort scanned strings line by line
    // Write sorted lines to standard output
	for x := range external.Sort(scan.Lines(r)) {
		fmt.Fprintln(out, x)
	}
	return nil
}
```

## Performance

Here are the results of a benchmark sorting 20 million lines of UUIDs. On a MacBook Pro 2023 with an Apple M3, it was faster than the sort command without any options. I assume the sort command uses a similar strategy to external.Sort, so I don't know why this is the case. Since the sort command does not write out in JSON, ultimately, the sort command should be faster. It might be achieved by adjusting the options. It seems like the sort command is not using multiple CPUs, but according to the man page, the default value for the parallel option is the same as the number of CPUs.


```bash
$ make bench sort_command
time cat bench.dat | go run cmd/sort/main.go | go run cmd/sorted/main.go

real    0m10.192s
user    0m20.503s
sys     0m3.466s
sort --version
2.3-Apple (165.100.8)
time cat bench.dat | sort | go run cmd/sorted/main.go

real    1m51.083s
user    1m49.487s
sys     0m1.324s
```

## License

MIT
