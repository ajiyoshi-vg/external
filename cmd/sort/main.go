package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/ajiyoshi-vg/external"
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
		log.Fatal(err)
	}
}

func run() error {
	if opt.in == "" {
		return sort(os.Stdin)
	}
	f, err := os.Open(opt.in)
	if err != nil {
		return err
	}
	defer f.Close()
	return sort(f)
}

func sort(r io.Reader) error {
	for x := range external.SortString(external.Lines(r)) {
		fmt.Println(x)
	}
	return nil
}
