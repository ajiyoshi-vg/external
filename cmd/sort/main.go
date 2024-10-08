package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/ajiyoshi-vg/external"
	"github.com/ajiyoshi-vg/external/scan"
)

var opt struct {
	in      string
	verbose bool
}

func init() {
	flag.StringVar(&opt.in, "in", "", "input file(empty: stdin)")
	flag.BoolVar(&opt.verbose, "v", false, "verbose")
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
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	for i, x := range scan.WithIndex(external.Sort(scan.Lines(r))) {
		fmt.Fprintln(out, x)
		if opt.verbose && i%(1*1000*1000) == 0 {
			log.Println(i)
		}
	}
	return nil
}
