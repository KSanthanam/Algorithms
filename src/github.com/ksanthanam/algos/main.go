package main

import (
	"flag"
	"strconv"

	nqueens "github.com/ksanthanam/nqueens"
)

var (
	size uint = 10
)

func main() {
	flag.Parse()
	sizeStr := flag.Arg(0)
	if sizeInt, err := strconv.Atoi(sizeStr); err == nil {
		size = uint(sizeInt)
	}
	nqueens.NQueenSolutions(size)
}
