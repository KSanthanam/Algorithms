package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"

	nqueens "github.com/ksanthanam/nqueens"
)

// var (
// 	size uint = 4
// )

func main() {
	// runtime.GOMAXPROCS(15)
	// fmt.Println(runtime.GOMAXPROCS(-1))
	debugPtr := flag.Bool("debug", false, "Debug true/false")
	sizePtr := flag.Int("size", 25, "Debug true/false")
	levelPtr := flag.Int("level", 1, "Debug with lower levels")
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")
	flag.Parse()
	nqueens.DEBUG = *debugPtr
	nqueens.LEVEL = *levelPtr
	size := uint(*sizePtr)

	fmt.Println("CPU Profile", *cpuprofile)
	if *cpuprofile != "" {
		fmt.Println("Createing ", *cpuprofile)
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	fmt.Println(fmt.Sprintf("Running N Queen solution for size(%d) with DEBUG is %t and LEVEL is %d ", size, *debugPtr, *levelPtr))
	nqueens.NQueenSolutions(size)
}
