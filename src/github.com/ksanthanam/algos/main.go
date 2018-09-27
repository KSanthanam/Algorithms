package main

import (
	"flag"
	"fmt"
	"runtime"

	nqueens "github.com/ksanthanam/nqueens"
)

// var (
// 	size uint = 4
// )

func main() {
	runtime.GOMAXPROCS(50)
	// fmt.Println(runtime.GOMAXPROCS(-1))
	debugPtr := flag.Bool("debug", false, "Debug true/false")
	sizePtr := flag.Int("size", 10, "Debug true/false")
	levelPtr := flag.Int("level", 1, "Debug with lower levels")
	flag.Parse()
	nqueens.DEBUG = *debugPtr
	nqueens.LEVEL = *levelPtr
	size := uint(*sizePtr)
	fmt.Println(fmt.Sprintf("Running N Queen solution for size(%d) with DEBUG is %t and LEVEL is %d ", size, *debugPtr, *levelPtr))
	nqueens.NQueenSolutions(size)
}
