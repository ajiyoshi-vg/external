package main

import (
	"bufio"
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
	i := 0
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	for x := range external.SortString(external.Lines(r)) {
		fmt.Fprintln(out, x)
		i++
		if i%(1*1000*1000) == 0 {
			log.Println(i)
		}
	}
	return nil
}
