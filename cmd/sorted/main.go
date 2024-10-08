package main

import (
	"bufio"
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
	br := bufio.NewReader(r)
	prev := ""
	for line := range scan.Lines(br) {
		if prev > line {
			return fmt.Errorf("not sorted: [%s] > [%s]", prev, line)
		}
		prev = line
	}
	return nil
}
