package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"iter"
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
	i := 0
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	for x := range sortedLines(r) {
		fmt.Fprintln(out, x)
		i++
		if opt.verbose && i%(1*1000*1000) == 0 {
			log.Println(i)
		}
	}
	return nil
}

func sortedLines(r io.Reader) iter.Seq[string] {
	sorted := external.Sort(scan.Lines(r))
	if opt.verbose {
		return scan.Prove("sort", sorted)
	}
	return sorted
}
