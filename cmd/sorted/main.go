package main

import (
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/ajiyoshi-vg/external/scan"
)

var opt struct {
	in string
}

func init() {
	flag.StringVar(&opt.in, "in", "", "input file(empty: stdin)")
	flag.Parse()
}

func main() {
	if err := run(); err != nil {
		slog.Error(err.Error())
	}
}

func run() error {
	if opt.in == "" {
		return sorted(os.Stdin)
	}
	f, err := os.Open(opt.in)
	if err != nil {
		return err
	}
	defer f.Close()
	return sorted(f)
}

func sorted(r io.Reader) error {
	prev := ""
	n := 0
	for line := range scan.Lines(r) {
		if prev > line {
			return fmt.Errorf("not sorted: [%s] > [%s]", prev, line)
		}
		prev = line
		n++
	}
	fmt.Fprintf(os.Stderr, "%d lines sorted\n", n)
	return nil
}
