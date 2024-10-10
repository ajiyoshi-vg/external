package main

import (
	"bufio"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/google/uuid"
	"github.com/schollz/progressbar/v3"
)

var opt struct {
	num int64
}

func init() {
	flag.Int64Var(&opt.num, "n", 1000*1000, "number")
	flag.Parse()
}

func main() {
	if err := run(); err != nil {
		slog.Error(err.Error())
	}
}

func run() error {
	p := progressbar.Default(opt.num)
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	for range opt.num {
		fmt.Fprintln(w, uuid.New().String())
		_ = p.Add(1)
	}
	return nil
}
